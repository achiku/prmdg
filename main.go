package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"sort"

	schema "github.com/lestrrat/go-jsschema"
	jsval "github.com/lestrrat/go-jsval"
	"github.com/pkg/errors"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "0.0.1"
)

var (
	app = kingpin.New("prmdg", "prmd generated JSON Hyper Schema to Go")
	pkg = app.Flag("package", "package name for Go file").Default("main").Short('p').String()
	fp  = app.Flag("file", "path JSON Schema").Required().Short('f').String()
	op  = app.Flag("output", "path to Go output file").Short('o').String()

	structCmd = app.Command("struct", "generate struct file")
	jsValCmd  = app.Command(
		"jsval", "generate validator file using github.com/lestrrat/go-jsval")
	validatorCmd = app.Command(
		"validator", "generate validator file using github.com/go-playground/validator")

	scValidator = structCmd.Flag("validate-tag", "add `validate` tag to struct").Bool()
	scUseTitle  = structCmd.Flag("use-title", "use title tag in request/response struct name").Bool()
	scNullable  = structCmd.Flag("nullable", "use github.com/guregu/null for null value").Bool()
)

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch cmd {
	case structCmd.FullCommand():
		if err := generateStructFile(pkg, *fp, op, *scValidator, *scUseTitle, *scNullable); err != nil {
			app.Errorf("failed to generate struct file: %s", err)
		}
	case jsValCmd.FullCommand():
		if err := generateJsValValidatorFile(pkg, *fp, op); err != nil {
			app.Errorf("failed to generate jsval validator file: %s", err)
		}
	case validatorCmd.FullCommand():
		if err := generateValidatorFile(pkg, *fp, op); err != nil {
			app.Errorf("failed to generate validator file: %s", err)
		}
	}
}

func generateValidatorFile(pkg *string, fp string, op *string) error {
	sc, err := schema.ReadFile(fp)
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", fp)
	}
	parser := NewParser(sc, *pkg)
	vals, err := parser.ParseValidators()
	if err != nil {
		return err
	}
	var src []byte
	src = append(src, []byte(fmt.Sprintf("package %s\n\n", *pkg))...)
	ss, err := format.Source(vals.Render())
	if err != nil {
		return err
	}
	src = append(src, ss...)
	var out *os.File
	if *op != "" {
		out, err = os.Create(*op)
		if err != nil {
			return errors.Wrapf(err, "failed to create %s", *op)
		}
		defer out.Close()
		if _, err := out.Write(src); err != nil {
			return err
		}
		params := []string{"-w", *op}
		if err := exec.Command("goimports", params...).Run(); err != nil {
			return err
		}
	} else {
		if _, err := os.Stdout.Write(src); err != nil {
			return err
		}
	}
	return nil
}

func generateStructFile(pkg *string, fp string, op *string, val, useTitle, nullable bool) error {
	sc, err := schema.ReadFile(fp)
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", fp)
	}
	parser := NewParser(sc, *pkg)
	resources, err := parser.ParseResources()
	if err != nil {
		return err
	}
	links, err := parser.ParseActions(resources)
	if err != nil {
		return err
	}

	var src []byte
	src = append(src, []byte(fmt.Sprintf("package %s\n\n", *pkg))...)

	var resKeys []string
	for key := range resources {
		resKeys = append(resKeys, key)
	}
	stOpt := FormatOption{
		Validator: val,
		Schema:    false,
		UseTitle:  useTitle,
		UseNull:   nullable,
	}
	sort.Strings(resKeys)
	for _, k := range resKeys {
		res := resources[k]
		ss, err := format.Source(res.Struct(stOpt))
		if err != nil {
			return errors.Wrapf(err, "failed to format resource: %s: %s", res.Name, res.Title)
		}
		src = append(src, ss...)
	}

	var linkKeys []string
	for key := range links {
		linkKeys = append(linkKeys, key)
	}
	sort.Strings(linkKeys)
	for _, k := range linkKeys {
		actions := links[k]
		for _, action := range actions {
			var reqOpt FormatOption
			switch {
			case action.Method == "GET" || action.Encoding == "application/x-www-form-urlencoded":
				reqOpt = FormatOption{
					Validator: val,
					Schema:    true,
					UseTitle:  useTitle,
					UseNull:   nullable,
				}
			default:
				reqOpt = FormatOption{
					Validator: val,
					Schema:    false,
					UseTitle:  useTitle,
					UseNull:   nullable,
				}
			}
			req, err := format.Source(action.RequestStruct(reqOpt))
			if err != nil {
				return errors.Wrapf(err, "failed to format request struct: %s, %s", k, action.Href)
			}
			src = append(src, req...)
			resp, err := format.Source(action.ResponseStruct(reqOpt))
			if err != nil {
				return errors.Wrapf(err, "failed to format response struct: %s, %s", k, action.Href)
			}
			src = append(src, resp...)
		}
	}

	var out *os.File
	if *op != "" {
		out, err = os.Create(*op)
		if err != nil {
			return errors.Wrapf(err, "failed to create %s", *op)
		}
		defer out.Close()
		if _, err := out.Write(src); err != nil {
			return err
		}
		params := []string{"-w", *op}
		if err := exec.Command("goimports", params...).Run(); err != nil {
			return err
		}
	} else {
		if _, err := os.Stdout.Write(src); err != nil {
			return err
		}
	}
	return nil
}

func generateJsValValidatorFile(pkg *string, fp string, op *string) error {
	sc, err := schema.ReadFile(fp)
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", fp)
	}
	parser := NewParser(sc, *pkg)
	validators, err := parser.ParseJsValValidators()
	if err != nil {
		return err
	}
	generator := jsval.NewGenerator()
	var src bytes.Buffer
	fmt.Fprintf(&src, "package %s\n", *pkg)
	fmt.Fprint(&src, "import \"github.com/lestrrat/go-jsval\"\n")
	if err := generator.Process(&src, validators...); err != nil {
		return err
	}

	var out *os.File
	if *op != "" {
		out, err = os.Create(*op)
		if err != nil {
			return errors.Wrapf(err, "failed to create %s", *op)
		}
		defer out.Close()
		if _, err := out.Write(src.Bytes()); err != nil {
			return err
		}
		params := []string{"-w", *op}
		if err := exec.Command("goimports", params...).Run(); err != nil {
			return err
		}
	} else {
		if _, err := os.Stdout.Write(src.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
