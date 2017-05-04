package main

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/achiku/varfmt"
	hschema "github.com/lestrrat/go-jshschema"
	schema "github.com/lestrrat/go-jsschema"
	jsval "github.com/lestrrat/go-jsval"
	"github.com/lestrrat/go-jsval/builder"
	"github.com/pkg/errors"
)

// PropType proper type
type PropType int

// Property types
const (
	PropTypeScalar PropType = iota
	PropTypeArray
	PropTypeObject
)

// Parser convertor
type Parser struct {
	schema  *schema.Schema
	pkgName string
}

// NewParser creates parser
func NewParser(sh *schema.Schema, pkgName string) *Parser {
	return &Parser{
		schema:  sh,
		pkgName: pkgName,
	}
}

// Resource plain resource
type Resource struct {
	Name       string
	Title      string
	Schema     *schema.Schema
	Properties []Property
}

// Struct returns struct go representation of resource
func (rs *Resource) Struct() []byte {
	name := varfmt.PublicVarName(rs.Name)
	var src bytes.Buffer
	fmt.Fprintf(&src, "// %s struct for %s resource\n", name, rs.Name)
	fmt.Fprintf(&src, "type %s struct {\n", name)
	for _, p := range rs.Properties {
		fmt.Fprintf(&src, "%s\n", p.Field())
	}
	fmt.Fprint(&src, "}\n\n")
	return src.Bytes()
}

// Property resource properties
type Property struct {
	Name             string
	Format           string
	Types            []string
	SecondTypes      []string
	PropType         PropType
	Required         bool
	Reference        string
	SecondReference  string
	Pattern          *regexp.Regexp
	Schema           *schema.Schema
	InlineProperties []Property
}

func (pr *Property) refToKey() string {
	var ref string
	if pr.SecondReference != "" {
		ref = pr.SecondReference
	} else {
		ref = pr.Reference
	}
	return strings.Replace(ref, "#/definitions/", "", 1)
}

// Field returns go struct field representation of property
func (pr *Property) Field() []byte {
	structName := varfmt.PublicVarName(strings.Replace(pr.Name, "-", "_", -1))
	// FIXME: need to support multiple types including 'null'
	// https://github.com/interagent/prmd/blob/master/docs/schemata.md#definitions
	var (
		t     string
		empty string
	)
	switch {
	case pr.PropType == PropTypeScalar && len(pr.Types) == 1:
		t = convertScalarProp(pr.Types[0], pr.Format)
	case pr.PropType == PropTypeArray:
		if len(pr.SecondTypes) == 1 && pr.SecondTypes[0] == "object" {
			t = fmt.Sprintf("[]%s", varfmt.PublicVarName(pr.refToKey()))
		} else if len(pr.InlineProperties) != 0 {
			var inline bytes.Buffer
			fmt.Fprint(&inline, "[]struct{\n")
			for _, p := range pr.InlineProperties {
				fmt.Fprintf(&inline, "%s\n", p.Field())
			}
			fmt.Fprint(&inline, "} ")
			t = inline.String()
		} else {
			t = fmt.Sprintf("[]%s", convertScalarProp(pr.SecondTypes[0], pr.Format))
		}
	case pr.PropType == PropTypeObject:
		t = fmt.Sprintf("*%s", varfmt.PublicVarName(pr.refToKey()))
	}
	if !pr.Required {
		empty = ",omitempty"
	}

	var src bytes.Buffer
	fmt.Fprintf(&src, "%s %s `json:\"%s%s\" schema:\"%s\"`", structName, t, pr.Name, empty, pr.Name)
	return src.Bytes()
}

func convertScalarProp(t, f string) string {
	switch t {
	case "number":
		return "float64"
	case "integer":
		return "int64"
	case "boolean":
		return "bool"
	case "string":
		if f == "date-time" {
			return "time.Time"
		}
		return "string"
	default:
		return ""
	}
}

// Action endpoint
type Action struct {
	Href     string
	Method   string
	Rel      string
	Request  *Resource
	Response *Resource
}

// RequestStruct request struct
func (a *Action) RequestStruct() []byte {
	if a.Request == nil {
		return []byte("")
	}
	name := varfmt.PublicVarName(
		strings.Replace(a.Response.Name+strings.Title(a.Rel), "-", "_", -1) + "Request")

	var src bytes.Buffer
	fmt.Fprintf(&src, "// %s struct for %s\n", name, a.Request.Name)
	fmt.Fprintf(&src, "// %s: %s\n", a.Method, a.Href)
	fmt.Fprintf(&src, "type %s struct {\n", name)
	for _, p := range a.Request.Properties {
		fmt.Fprintf(&src, "%s\n", p.Field())
	}
	fmt.Fprint(&src, "}\n\n")
	return src.Bytes()
}

// ResponseStruct response struct
func (a *Action) ResponseStruct() []byte {
	if a.Response == nil {
		return []byte("")
	}
	name := varfmt.PublicVarName(
		strings.Replace(a.Response.Name+strings.Title(a.Rel), "-", "_", -1) + "Response")
	orgName := varfmt.PublicVarName(a.Response.Name)

	var src bytes.Buffer
	fmt.Fprintf(&src, "// %s struct for %s\n", name, a.Response.Name)
	fmt.Fprintf(&src, "// %s: %s\n", a.Method, a.Href)
	if a.Rel == "instances" {
		fmt.Fprintf(&src, "type %s []%s\n", name, orgName)
		return src.Bytes()
	}
	fmt.Fprintf(&src, "type %s %s\n\n", name, orgName)
	return src.Bytes()
}

func resolveSchema(sch *schema.Schema, root *schema.Schema) (*schema.Schema, error) {
	if sch.IsResolved() {
		return sch, nil
	}
	sh, err := sch.Resolve(root)
	if err != nil {
		return nil, err
	}
	// FIXME: recursively resolving schema. may need to limit # of recursion.
	return resolveSchema(sh, root)
}

func typesToStrings(types schema.PrimitiveTypes) []string {
	var vals []string
	for _, tt := range types {
		vals = append(vals, tt.String())
	}
	return vals
}

func sortProperties(props []Property) []Property {
	pMap := make(map[string]Property)
	for _, p := range props {
		pMap[p.Name] = p
	}
	var names []string
	for n := range pMap {
		names = append(names, n)
	}
	sort.Strings(names)

	var sorted []Property
	for _, n := range names {
		sorted = append(sorted, pMap[n])
	}
	return sorted
}

func sortActions(acs []Action) []Action {
	aMap := make(map[string]Action)
	for _, a := range acs {
		aMap[a.Href] = a
	}
	var refs []string
	for r := range aMap {
		refs = append(refs, r)
	}
	sort.Strings(refs)
	var sorted []Action
	for _, r := range refs {
		sorted = append(sorted, aMap[r])
	}
	return sorted
}

func sortValidator(vals []*jsval.JSVal) []*jsval.JSVal {
	vMap := make(map[string]*jsval.JSVal)
	for _, v := range vals {
		vMap[v.Name] = v
	}
	var names []string
	for n := range vMap {
		names = append(names, n)
	}
	sort.Strings(names)
	var sorted []*jsval.JSVal
	for _, n := range names {
		sorted = append(sorted, vMap[n])
	}
	return sorted
}

// ParseResources parse plain resource
func (p *Parser) ParseResources() (map[string]Resource, error) {
	res := make(map[string]Resource)
	// parse resource itself
	for id, df := range p.schema.Definitions {
		rs := Resource{
			Name:   id,
			Title:  df.Title,
			Schema: df,
		}
		// parse resource field
		var flds []Property
		for name, tp := range df.Properties {
			// save reference before resolving ref
			ref := tp.Reference
			fieldSchema, err := resolveSchema(tp, p.schema)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to resolve, %s:%s", id, name)
			}
			fld := Property{
				Name:      name,
				Format:    string(fieldSchema.Format),
				Types:     typesToStrings(fieldSchema.Type),
				Required:  df.IsPropRequired(name),
				Pattern:   fieldSchema.Pattern,
				Reference: ref,
				Schema:    fieldSchema,
			}
			switch {
			case fieldSchema.Type.Contains(schema.ArrayType):
				// if this field is an array
				// currently this tool supports only one itme per array field
				if len(fieldSchema.Items.Schemas) != 1 {
					return nil, errors.Errorf("array type has to have an item: %s, %s", id, name)
				}
				item := fieldSchema.Items.Schemas[0]
				tmpItem, err := resolveSchema(item, p.schema)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to resolve, %s:%s", id, name)
				}
				switch {
				case item.Reference == "" && item.Properties == nil:
					// no reference, no item properties = primitive type
					fld.SecondTypes = typesToStrings(item.Type)
					// log.Printf("no ref, no prop: %s: %s", name, item.Reference)
				case item.Reference != "" && !tmpItem.Type.Contains(schema.ObjectType):
					// reference to primitive
					fld.SecondTypes = typesToStrings(tmpItem.Type)
					// log.Printf("ref to primitive: %s: %s", name, item.Reference)
				case item.Reference == "" && item.Properties != nil:
					// no reference, item properties = inline object
					// parse properties and create inline fields
					var inlineFields []Property
					for k, v := range item.Properties {
						f := Property{
							Name:      k,
							Format:    string(v.Format),
							Pattern:   v.Pattern,
							Reference: v.Reference,
							Types:     typesToStrings(v.Type),
							Schema:    v,
							PropType:  PropTypeScalar,
						}
						inlineFields = append(inlineFields, f)
					}
					fld.InlineProperties = inlineFields
					// log.Printf("no ref, inline prop: %s: %s", name, item.Reference)
				case item.Reference != "" && tmpItem.Type.Contains(schema.ObjectType):
					// reference to object
					fld.SecondTypes = []string{"object"}
					fld.SecondReference = item.Reference
					// log.Printf("ref to obj: %s: %s", name, item.Reference)
				}
				fld.PropType = PropTypeArray
			case fieldSchema.Type.Contains(schema.ObjectType):
				// if this field is a object
				fld.PropType = PropTypeObject
				fld.SecondTypes = []string{name}
			default:
				// if this field is a scalar
				fld.PropType = PropTypeScalar
			}
			flds = append(flds, fld)
		}
		rs.Properties = sortProperties(flds)
		res[id] = rs
	}
	return res, nil
}

// ParseActions parse endpoints
func (p *Parser) ParseActions(res map[string]Resource) (map[string][]Action, error) {
	eptsMap := make(map[string][]Action)
	for id, df := range p.schema.Definitions {
		// use json hyper schema to parse links
		hsc := hschema.New()
		if err := hsc.Extract(df.Extras); err != nil {
			return nil, errors.Wrapf(err, "failed to extract links for (%s)", id)
		}
		// parse endpoints
		var eps []Action
		for _, e := range hsc.Links {
			href, err := url.QueryUnescape(e.Href)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to unescape %s", e.Href)
			}
			ep := Action{
				Href:   href,
				Method: e.Method,
				Rel:    e.Rel,
			}
			// parse request if exists
			if e.Schema != nil {
				var flds []Property
				for name, props := range e.Schema.Properties {
					ref := props.Reference
					sh, err := resolveSchema(props, p.schema)
					if err != nil {
						return nil, errors.Wrapf(err, "failed to resolve, %s:%s", id, name)
					}
					fld := Property{
						Name:      name,
						Types:     typesToStrings(sh.Type),
						Format:    string(sh.Format),
						Required:  e.Schema.IsPropRequired(name),
						Pattern:   sh.Pattern,
						Reference: ref,
					}
					flds = append(flds, fld)
				}
				ep.Request = &Resource{
					Name:       id,
					Properties: sortProperties(flds),
					Title:      e.Schema.Title,
				}
			}
			// parse response if exists
			if e.TargetSchema != nil {
				// http://json-schema.org/latest/json-schema-hypermedia.html#rfc.section.5.4
				var flds []Property
				for name, props := range e.TargetSchema.Properties {
					sh, err := resolveSchema(props, p.schema)
					if err != nil {
						return nil, errors.Wrapf(err, "failed to resolve, %s:%s", id, name)
					}
					fld := Property{
						Name:     name,
						Types:    typesToStrings(sh.Type),
						Format:   string(sh.Format),
						Required: e.TargetSchema.IsPropRequired(name),
						Pattern:  sh.Pattern,
					}
					flds = append(flds, fld)
				}
				ep.Response = &Resource{
					Name:       id,
					Properties: sortProperties(flds),
					Title:      e.TargetSchema.Title,
				}
			} else {
				// if targetSchema is not set, use default resource for this link
				resp, ok := res[id]
				if !ok {
					return nil, errors.Errorf("resource not found: %s", id)
				}
				ep.Response = &resp
			}
			eps = append(eps, ep)
		}
		eptsMap[id] = sortActions(eps)
	}
	return eptsMap, nil
}

// ParseValidators parse validator
func (p *Parser) ParseValidators() ([]*jsval.JSVal, error) {
	var validators []*jsval.JSVal
	for id, df := range p.schema.Definitions {
		// use json hyper schema to parse links
		hsc := hschema.New()
		if err := hsc.Extract(df.Extras); err != nil {
			return nil, errors.Wrapf(err, "failed to extract links for (%s)", id)
		}

		for _, e := range hsc.Links {
			var v *jsval.JSVal
			if e.Schema == nil {
				v = jsval.New()
				v.SetRoot(jsval.Any())
			} else {
				sh, err := resolveSchema(e.Schema, p.schema)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to resolve, %s", id)
				}
				b := builder.New()
				v, err = b.BuildWithCtx(sh, p.schema)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to build validator: %s", id)
				}
			}
			v.Name = varfmt.PublicVarName(
				strings.Replace(id+strings.Title(e.Rel), "-", "_", -1) + "Validator")
			validators = append(validators, v)
		}
	}
	return sortValidator(validators), nil
}
