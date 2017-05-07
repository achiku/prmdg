package main

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/achiku/varfmt"
)

// Validators validators
type Validators map[string]Validator

// Render rendor validators
func (vs Validators) Render() []byte {
	var src bytes.Buffer
	// constants
	fmt.Fprint(&src, "// regexp string constants\n")
	fmt.Fprint(&src, "const (\n")
	for _, v := range vs {
		fmt.Fprintf(&src, "%s\n", v.RegexpConst())
	}
	fmt.Fprint(&src, ")\n")

	// vars
	fmt.Fprint(&src, "// regexp objects\n")
	fmt.Fprint(&src, "var (\n")
	for _, v := range vs {
		fmt.Fprintf(&src, "%s\n", v.RegexpVar())
	}
	fmt.Fprint(&src, ")\n")

	fmt.Fprint(&src, "// use a single instance of Validate, it caches struct info\n")
	fmt.Fprint(&src, "var validate *validator.Validate\n")

	// function definitions
	for _, v := range vs {
		fmt.Fprintf(&src, "%s\n", v.ValidatorFunc())
	}

	// register validation functions
	fmt.Fprint(&src, "func init() {\n")
	fmt.Fprint(&src, "validate = validator.New()\n")
	for _, v := range vs {
		fmt.Fprintf(&src, "%s\n", v.RegisterFunc())
	}
	fmt.Fprint(&src, "}\n")
	return src.Bytes()
}

// Validator validator
type Validator struct {
	Name         string
	RegexpString string
}

// RegexpConstName const name
func (val Validator) RegexpConstName() string {
	return varfmt.PublicVarName(val.Name) + "RegexString"
}

// RegexpVarName var regexp name
func (val Validator) RegexpVarName() string {
	return varfmt.PublicVarName(val.Name) + "Regex"
}

// ValidateFuncName validator func name
func (val Validator) ValidateFuncName() string {
	return varfmt.PublicVarName(val.Name) + "Validator"
}

// RegexpConst const def
func (val Validator) RegexpConst() string {
	return fmt.Sprintf("%s = `%s`", val.RegexpConstName(), val.RegexpString)
}

// RegexpVar var regexp def
func (val Validator) RegexpVar() string {
	return fmt.Sprintf("%s = regexp.MustCompile(%s)", val.RegexpVarName(), val.RegexpConstName())
}

// ValidatorFunc validator
func (val Validator) ValidatorFunc() string {
	// ignore errors since it always succeeds
	tmpl, _ := template.New("").Parse(`
	// {{ .ValidateFuncName }} for validation
	func {{ .ValidateFuncName }}(fl validator.FieldLevel) bool {
		return {{ .RegexpVarName }}.MatchString(fl.Field().String())
	}`)
	var src bytes.Buffer
	tmpl.Execute(&src, val)
	return src.String()
}

// RegisterFunc register validator
func (val Validator) RegisterFunc() string {
	// ignore errors since it always succeeds
	tmpl, _ := template.New("").Parse(`
	if err := validate.RegisterValidation("{{ .ValidateFuncName }}", {{ .ValidateFuncName }}); err != nil {
		log.Fatal(err)
	}`)
	var src bytes.Buffer
	tmpl.Execute(&src, val)
	return src.String()
}
