package main

import (
	"net/url"
	"sort"
	"strings"

	"github.com/achiku/varfmt"
	"github.com/lestrrat-go/go-jsval/builder"
	hschema "github.com/lestrrat-go/jshschema"
	schema "github.com/lestrrat-go/jsschema"
	"github.com/lestrrat-go/jsval"
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

func sortProperties(props []*Property) []*Property {
	pMap := make(map[string]*Property)
	for _, p := range props {
		pMap[p.Name] = p
	}
	var names []string
	for n := range pMap {
		names = append(names, n)
	}
	sort.Strings(names)

	var sorted []*Property
	for _, n := range names {
		sorted = append(sorted, pMap[n])
	}
	return sorted
}

func sortActions(acs []Action) []Action {
	aMap := make(map[string]Action)
	for _, a := range acs {
		aMap[a.Method+a.Href] = a
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

// NewProperty new property
func NewProperty(name string, tp *schema.Schema, df *schema.Schema, root *schema.Schema) (*Property, error) {
	// save reference before resolving ref
	ref := tp.Reference
	fieldSchema, err := resolveSchema(tp, root)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve, %s", name)
	}
	fld := &Property{
		Name:      name,
		Format:    string(fieldSchema.Format),
		Types:     fieldSchema.Type,
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
			return nil, errors.Errorf("array type has to have an item: %s", name)
		}
		item := fieldSchema.Items.Schemas[0]
		resolvedItem, err := resolveSchema(item, root)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve: %s", name)
		}
		switch {
		case isMainResource(item.Reference) && resolvedItem.Type.Contains(schema.ObjectType):
			// reference to main resource object
			// log.Printf("ref to main resource: %s: %s", name, item.Reference)
			fld.SecondTypes = []schema.PrimitiveType{schema.ObjectType}
			fld.SecondReference = item.Reference
		case item.Properties == nil && fieldSchema.Properties != nil:
			// field schema already has properties = inline object
			// log.Printf("inline obj: %s: %v", name, fieldSchema.Properties)
			var inlineFields []*Property
			for k, prop := range fieldSchema.Properties {
				f, err := NewProperty(k, prop, df, root)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to perse inline object: %s", k)
				}
				inlineFields = append(inlineFields, f)
			}
			fld.InlineProperties = inlineFields
			fld.SecondTypes = []schema.PrimitiveType{schema.ObjectType}
		case item.Reference == "" && item.Properties == nil:
			// no reference, no item properties = primitive type
			// log.Printf("primitive type: %s %s", name, item.Type)
			fld.SecondTypes = item.Type
		case item.Reference != "" && !resolvedItem.Type.Contains(schema.ObjectType):
			// reference to primitive = resolved primitive type
			// log.Printf("resolved primitive type: %s %s", name, resolvedItem.Type)
			fld.SecondTypes = resolvedItem.Type
		case item.Reference == "" && item.Properties != nil:
			// no reference, item properties = inline object
			// parse properties, and recursively create inline fields
			// log.Printf("resolved inline obj: %s: %v", name, item.Properties)
			var inlineFields []*Property
			for k, prop := range item.Properties {
				f, err := NewProperty(k, prop, df, root)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to perse inline object: %s", k)
				}
				inlineFields = append(inlineFields, f)
			}
			fld.InlineProperties = inlineFields
			fld.SecondTypes = []schema.PrimitiveType{schema.ObjectType}
		case !isMainResource(item.Reference):
			// log.Printf("resolved inline obj: %s: %v", name, resolvedItem.Properties)
			var inlineFields []*Property
			for k, prop := range resolvedItem.Properties {
				f, err := NewProperty(k, prop, df, root)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to perse inline object: %s", k)
				}
				inlineFields = append(inlineFields, f)
			}
			fld.InlineProperties = inlineFields
			fld.SecondTypes = []schema.PrimitiveType{schema.ObjectType}
		}
		fld.PropType = PropTypeArray
	case fieldSchema.Type.Contains(schema.ObjectType):
		// if this field is a object
		switch {
		case fieldSchema.Reference == "" && fieldSchema.Properties != nil:
			// inline object without definitions
			var inlineFields []*Property
			for k, prop := range fieldSchema.Properties {
				f, err := NewProperty(k, prop, df, root)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to perse inline object: %s", k)
				}
				inlineFields = append(inlineFields, f)
			}
			fld.InlineProperties = inlineFields
		}
		fld.PropType = PropTypeObject
	default:
		// if this field is a scalar
		fld.PropType = PropTypeScalar
	}
	return fld, nil
}

// ParseValidators parse validator
func (p *Parser) ParseValidators() (Validators, error) {
	vals := make(Validators)
	for _, df := range p.schema.Definitions {
		for name, tp := range df.Properties {
			fs, err := resolveSchema(tp, p.schema)
			if err != nil {
				return nil, err
			}
			if fs.Pattern != nil &&
				!fs.Type.Contains(schema.ObjectType) && !fs.Type.Contains(schema.ArrayType) {
				v := Validator{
					Name:         name,
					RegexpString: fs.Pattern.String(),
				}
				vals[name] = v
			}
		}
	}
	return vals, nil
}

// ParseResources parse plain resource
func (p *Parser) ParseResources() (map[string]Resource, error) {
	res := make(map[string]Resource)
	// parse resource itself
	for id, df := range p.schema.Definitions {
		rs := Resource{
			Name:      id,
			Title:     df.Title,
			Schema:    df,
			IsPrimary: true,
		}
		// parse resource field
		var flds []*Property
		for name, tp := range df.Properties {
			fld, err := NewProperty(name, tp, df, p.schema)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse %s", id)
			}
			fld.InlineProperties = sortProperties(fld.InlineProperties)
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
			var encoding string
			if e.EncType == "" {
				encoding = "application/json"
			} else {
				encoding = e.EncType
			}
			ep := Action{
				Encoding: encoding,
				Href:     href,
				Method:   e.Method,
				Title:    e.Title,
				Rel:      e.Rel,
			}
			// parse request if exists
			if e.Schema != nil {
				var flds []*Property
				for name, tp := range e.Schema.Properties {
					fld, err := NewProperty(name, tp, df, p.schema)
					if err != nil {
						return nil, errors.Wrapf(err, "failed to parse %s", id)
					}
					flds = append(flds, fld)
				}
				ep.Request = &Resource{
					Name:       id,
					Properties: sortProperties(flds),
					Title:      e.Schema.Title,
					IsPrimary:  false,
				}
			}
			// parse response if exists
			if e.TargetSchema != nil {
				// http://json-schema.org/latest/json-schema-hypermedia.html#rfc.section.5.4
				switch {
				case e.TargetSchema.Reference == "":
					var flds []*Property
					for name, tp := range e.TargetSchema.Properties {
						fld, err := NewProperty(name, tp, df, p.schema)
						if err != nil {
							return nil, errors.Wrapf(err, "failed to parse %s", id)
						}
						flds = append(flds, fld)
					}
					ep.Response = &Resource{
						Name:       id,
						Properties: sortProperties(flds),
						Title:      e.TargetSchema.Title,
						Schema:     e.TargetSchema,
						IsPrimary:  false,
					}
				case e.TargetSchema.Reference != "" && IsRefToMainResource(e.TargetSchema.Reference):
					ep.Response = &Resource{
						Name:      id,
						Title:     e.TargetSchema.Title,
						Schema:    e.TargetSchema,
						IsPrimary: false,
					}
				case e.TargetSchema.Reference != "" && !IsRefToMainResource(e.TargetSchema.Reference):
					fld, err := NewProperty(e.TargetSchema.ID, e.TargetSchema, df, p.schema)
					if err != nil {
						return nil, errors.Wrapf(err, "failed to parse %s", id)
					}
					ep.Response = &Resource{
						Name:       id,
						Properties: sortProperties(fld.InlineProperties),
						Title:      e.TargetSchema.Title,
						Schema:     e.TargetSchema,
						IsPrimary:  false,
					}
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

// ParseJsValValidators parse validator
func (p *Parser) ParseJsValValidators() ([]*jsval.JSVal, error) {
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

func isMainResource(ref string) bool {
	if ref == "" {
		return false
	}
	tmp := strings.Replace(ref, "#/definitions/", "", 1)
	return !strings.Contains(tmp, "/")
}
