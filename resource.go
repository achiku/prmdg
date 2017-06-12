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
	IsPrimary  bool
}

// FormatOption output struct format option
type FormatOption struct {
	Validator bool
	Schema    bool
	UseTitle  bool
	UseNull   bool
}

// Struct returns struct go representation of resource
func (rs *Resource) Struct(op FormatOption) []byte {
	name := varfmt.PublicVarName(normalize(rs.Name))
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

func normalize(n string) string {
	return strings.Replace(
		strings.Replace(n, "-", "_", -1), " ", "_", -1)
}

// IsRefToMainResource is ref
func IsRefToMainResource(ref string) bool {
	if ref == "" {
		return false
	}
	tmp := strings.Replace(ref, "#/definitions/", "", 1)
	return !strings.Contains(tmp, "/")
}

// IsRefToMainResource check if first class resource
func (pr *Property) IsRefToMainResource() bool {
	var ref string
	if pr.SecondReference != "" {
		ref = pr.SecondReference
	} else {
		ref = pr.Reference
	}
	if ref == "" {
		return false
	}
	tmp := strings.Replace(ref, "#/definitions/", "", 1)
	return !strings.Contains(tmp, "/")
}

func (pr *Property) refToStructName() string {
	// FIXME: too naieve. use js-pointer.
	var ref string
	if pr.SecondReference != "" {
		ref = pr.SecondReference
	} else {
		ref = pr.Reference
	}
	return normalize(strings.Replace(ref, "#/definitions/", "", 1))
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
	fieldName := varfmt.PublicVarName(normalize(pr.Name))
	var (
		t     string
		empty string
	)
	switch {
	case pr.PropType == PropTypeScalar:
		t = pr.ScalarType(op)
	case pr.PropType == PropTypeArray:
		if len(pr.InlineProperties) == 0 && pr.IsRefToMainResource() && pr.SecondTypes.Contains(schema.ObjectType) {
			// referecnce to main resource object
			t = fmt.Sprintf("[]%s", varfmt.PublicVarName(normalize(pr.refToStructName())))
		} else if len(pr.InlineProperties) != 0 {
			// inline list object
			t = pr.inlineListOjbect(op)
		} else {
			// an array of primitive types
			t = fmt.Sprintf("[]%s", pr.ScalarType(op))
		}
	case pr.Types.Contains(schema.ObjectType) && pr.IsRefToMainResource():
		// reference to main resource object
		t = fmt.Sprintf("*%s", varfmt.PublicVarName(normalize(pr.refToStructName())))
	case pr.Types.Contains(schema.ObjectType) && !pr.IsRefToMainResource():
		// inline object
		t = pr.inlineOjbect(op)
	}
	if !pr.Required {
		empty = ",omitempty"
	}

	var src bytes.Buffer
	fmt.Fprintf(&src, "%s %s `json:\"%s%s\"", fieldName, t, pr.Name, empty)
	if op.Schema {
		fmt.Fprintf(&src, " schema:\"%s\"", pr.Name)
	}

	if op.Validator {
		if pr.Required && pr.Pattern == nil {
			fmt.Fprint(&src, " validate:\"required\"")
		} else if pr.Required && pr.Pattern != nil {
			fmt.Fprintf(&src, " validate:\"required,%s\"", fieldName+"Validator")
		}
	}
	fmt.Fprint(&src, "`")
	return src.Bytes()
}

// ScalarType returns go scalar type
func (pr *Property) ScalarType(op FormatOption) string {
	var types schema.PrimitiveTypes
	if pr.Types.Contains(schema.ArrayType) {
		types = pr.SecondTypes
	} else {
		types = pr.Types
	}
	if op.UseNull {
		if types.Contains(schema.NullType) || !pr.Required {
			switch {
			case types.Contains(schema.NumberType):
				return "null.Float"
			case types.Contains(schema.IntegerType):
				return "null.Int"
			case types.Contains(schema.BooleanType):
				return "bool"
			case types.Contains(schema.StringType):
				if pr.Format == "date-time" {
					return "null.Time"
				}
				return "null.String"
			default:
				return ""
			}
		} else {
			switch {
			case types.Contains(schema.NumberType):
				return "float64"
			case types.Contains(schema.IntegerType):
				return "int64"
			case types.Contains(schema.BooleanType):
				return "bool"
			case types.Contains(schema.StringType):
				if pr.Format == "date-time" {
					return "time.Time"
				}
				return "string"
			default:
				return ""
			}
		}
	}
	switch {
	case types.Contains(schema.NumberType):
		return "float64"
	case types.Contains(schema.IntegerType):
		return "int64"
	case types.Contains(schema.BooleanType):
		return "bool"
	case types.Contains(schema.StringType):
		if pr.Format == "date-time" {
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
		normalize(a.Response.Name+strings.Title(n)) + "Request")

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
		normalize(a.Response.Name+strings.Title(n)) + "Response")
	orgName := varfmt.PublicVarName(normalize(a.Response.Name))

	var src bytes.Buffer
	fmt.Fprintf(&src, "// %s struct for %s\n", name, a.Response.Name)
	fmt.Fprintf(&src, "// %s: %s\n", a.Method, a.Href)
	if a.Rel == "instances" {
		fmt.Fprintf(&src, "type %s []%s\n", name, orgName)
		return src.Bytes()
	}
	switch {
	case a.Response.IsPrimary:
		fmt.Fprintf(&src, "type %s %s\n\n", name, orgName)
	case a.Response.Schema != nil && IsRefToMainResource(a.Response.Schema.Reference):
		fmt.Fprintf(&src, "type %s %s\n\n", name, orgName)
	case a.Response.Schema != nil && a.Response.Schema.Reference == "" && len(a.Response.Properties) != 0:
		fmt.Fprintf(&src, "type %s struct {\n", name)
		for _, p := range a.Response.Properties {
			fmt.Fprintf(&src, "%s\n", p.Field(op))
		}
		fmt.Fprint(&src, "}\n\n")
	case a.Response.Schema != nil && !IsRefToMainResource(a.Response.Schema.Reference):
		pr := a.Response.Properties[0]
		fmt.Fprintf(&src, "type %s struct {\n", name)
		for _, p := range pr.InlineProperties {
			fmt.Fprintf(&src, "%s\n", p.Field(op))
		}
		fmt.Fprint(&src, "}\n\n")
	}
	return src.Bytes()
}
