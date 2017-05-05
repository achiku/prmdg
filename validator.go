package main

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/achiku/varfmt"
)

// Validator validator
type Validator struct {
	Name         string
	RegexpString string
}

// RegexpConstName const name
func (val Validator) RegexpConstName() string {
	return val.Name + "RegexString"
}

// RegexpVarName var regexp name
func (val Validator) RegexpVarName() string {
	return val.Name + "Regex"
}

// ValidateFuncName validator func name
func (val Validator) ValidateFuncName() string {
	return varfmt.PublicVarName(val.Name) + "Validator"
}

// RegexpConst const def
func (val Validator) RegexpConst() string {
	return fmt.Sprintf("%s = \"%s\"", val.RegexpConstName(), val.RegexpString)
}

// RegexpVar var regexp def
func (val Validator) RegexpVar() string {
	return fmt.Sprintf("%s = regexp.MustCompile(%s)", val.RegexpVarName(), val.RegexpConstName())
}

// Validator validator
func (val Validator) Validator() string {
	tmpl, _ := template.New("").Parse(`
	func {{ .ValidateFuncName }}(fl validator.FieldLevel) bool {
		return {{ .RegexpVarName }}.MatchString(fl.Field().String())
	}`)
	var src bytes.Buffer
	// ignore errors since it always succeeds
	tmpl.Execute(&src, val)
	return src.String()
}
