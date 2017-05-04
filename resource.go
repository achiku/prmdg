package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/achiku/varfmt"
	schema "github.com/lestrrat/go-jsschema"
)

// Resource plain resource
type Resource struct {
	Name       string
	Title      string
	Schema     *schema.Schema
	Properties []*Property
}

// FormatOption output struct format option
type FormatOption struct {
	Validator bool
	Schema    bool
	UseTitle  bool
}

// Struct returns struct go representation of resource
func (rs *Resource) Struct(op FormatOption) []byte {
	name := varfmt.PublicVarName(strings.Replace(rs.Name, "-", "_", -1))
	var src bytes.Buffer
	fmt.Fprintf(&src, "// %s struct for %s resource\n", name, rs.Name)
	fmt.Fprintf(&src, "type %s struct {\n", name)
	for _, p := range rs.Properties {
		fmt.Fprintf(&src, "%s\n", p.Field(op))
	}
	fmt.Fprint(&src, "}\n\n")
	return src.Bytes()
}

// Property resource properties
type Property struct {
	Name             string
	Format           string
	Types            schema.PrimitiveTypes
	SecondTypes      schema.PrimitiveTypes
	PropType         PropType
	Required         bool
	Reference        string
	SecondReference  string
	InlineProperties []*Property
	Pattern          *regexp.Regexp
	Schema           *schema.Schema
}

func (pr *Property) refToStructName() string {
	// FIXME: naieve
	var ref string
	if pr.SecondReference != "" {
		ref = pr.SecondReference
	} else {
		ref = pr.Reference
	}
	return strings.Replace(strings.Replace(ref, "#/definitions/", "", 1), "-", "_", -1)
}

func (pr *Property) inlineOjbect(op FormatOption) string {
	var inline bytes.Buffer
	fmt.Fprint(&inline, "struct{\n")
	for _, p := range pr.InlineProperties {
		fmt.Fprintf(&inline, "%s\n", p.Field(op))
	}
	fmt.Fprint(&inline, "} ")
	return inline.String()
}

func (pr *Property) inlineListOjbect(op FormatOption) string {
	var inline bytes.Buffer
	fmt.Fprint(&inline, "[]struct{\n")
	for _, p := range pr.InlineProperties {
		fmt.Fprintf(&inline, "%s\n", p.Field(op))
	}
	fmt.Fprint(&inline, "} ")
	return inline.String()
}

// Field returns go struct field representation of property
func (pr *Property) Field(op FormatOption) []byte {
	fieldName := varfmt.PublicVarName(strings.Replace(pr.Name, "-", "_", -1))
	// FIXME: need to support multiple types including 'null'
	// https://github.com/interagent/prmd/blob/master/docs/schemata.md#definitions
	var (
		t     string
		empty string
	)
	switch {
	case pr.PropType == PropTypeScalar:
		t = convertScalarProp(pr.Types, pr.Format)
	case pr.PropType == PropTypeArray:
		if len(pr.SecondTypes) == 1 && pr.SecondTypes.Contains(schema.ObjectType) {
			t = fmt.Sprintf("[]%s", varfmt.PublicVarName(pr.refToStructName()))
		} else if len(pr.InlineProperties) != 0 {
			t = pr.inlineListOjbect(op)
		} else {
			t = fmt.Sprintf("[]%s", convertScalarProp(pr.SecondTypes, pr.Format))
		}
	case pr.PropType == PropTypeObject && pr.refToStructName() != "":
		// reference to object
		t = fmt.Sprintf("*%s", varfmt.PublicVarName(pr.refToStructName()))
	case pr.PropType == PropTypeObject && pr.refToStructName() == "":
		// inline object
		t = pr.inlineOjbect(op)
	}
	if !pr.Required {
		empty = ",omitempty"
	}

	var src bytes.Buffer
	if op.Schema {
		fmt.Fprintf(&src, "%s %s `json:\"%s%s\" schema:\"%s\"`", fieldName, t, pr.Name, empty, pr.Name)
	} else {
		fmt.Fprintf(&src, "%s %s `json:\"%s%s\"`", fieldName, t, pr.Name, empty)
	}
	return src.Bytes()
}

func convertScalarProp(t schema.PrimitiveTypes, format string) string {
	switch {
	case t.Contains(schema.NumberType):
		return "float64"
	case t.Contains(schema.IntegerType):
		return "int64"
	case t.Contains(schema.BooleanType):
		return "bool"
	case t.Contains(schema.StringType):
		if format == "date-time" {
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
	Title    string
	Request  *Resource
	Response *Resource
}

// RequestStruct request struct
func (a *Action) RequestStruct(op FormatOption) []byte {
	if a.Request == nil {
		return []byte("")
	}
	var n string
	if op.UseTitle {
		n = a.Title
	} else {
		n = a.Rel
	}
	name := varfmt.PublicVarName(
		strings.Replace(
			strings.Replace(a.Response.Name+strings.Title(n), "-", "_", -1)+"Request", " ", "_", -1))

	var src bytes.Buffer
	fmt.Fprintf(&src, "// %s struct for %s\n", name, a.Request.Name)
	fmt.Fprintf(&src, "// %s: %s\n", a.Method, a.Href)
	fmt.Fprintf(&src, "type %s struct {\n", name)
	for _, p := range a.Request.Properties {
		fmt.Fprintf(&src, "%s\n", p.Field(op))
	}
	fmt.Fprint(&src, "}\n\n")
	return src.Bytes()
}

// ResponseStruct response struct
func (a *Action) ResponseStruct(op FormatOption) []byte {
	if a.Response == nil {
		return []byte("")
	}
	var n string
	if op.UseTitle {
		n = a.Title
	} else {
		n = a.Rel
	}
	name := varfmt.PublicVarName(
		strings.Replace(
			strings.Replace(a.Response.Name+strings.Title(n), "-", "_", -1)+"Response", " ", "_", -1))
	orgName := varfmt.PublicVarName(strings.Replace(a.Response.Name, "-", "_", -1))

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
