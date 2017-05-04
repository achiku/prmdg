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

func (pr *Property) refToStructName() string {
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
			t = fmt.Sprintf("[]%s", varfmt.PublicVarName(pr.refToStructName()))
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
		t = fmt.Sprintf("*%s", varfmt.PublicVarName(pr.refToStructName()))
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
