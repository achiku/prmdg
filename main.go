package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"

	schema "github.com/lestrrat-go/jsschema"
	"github.com/lestrrat-go/jsval"
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
		"jsval", "generate validator file using github.com/lestrrat-go/go-jsval")
	validatorCmd = app.Command(
		"validator", "generate validator file using github.com/go-playground/validator")

	scValidator = structCmd.Flag("validate-tag", "add `validate` tag to struct").Bool()
	scUseTitle  = structCmd.Flag("use-title", "use title tag in request/response struct name").Bool()
	scNullable  = structCmd.Flag("nullable", "use github.com/guregu/null for null value").Bool()
)

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	var (
		in  io.Reader
		out io.Writer
		err error
	)
	if *op != "" {
		out, err = os.Create(*op)
		if err != nil {
			app.Errorf("failed to create output file %s: %s", *op, err)
		}
	} else {
		out = os.Stdout
	}
	in, err = os.Open(*fp)
	if err != nil {
		app.Errorf("failed to open input file %s: %s", *fp, err)
	}

	switch cmd {
	case structCmd.FullCommand():
		if err := generateStructFile(pkg, in, out, *scValidator, *scUseTitle, *scNullable); err != nil {
			app.Errorf("failed to generate struct file: %s", err)
		}
	case jsValCmd.FullCommand():
		if err := generateJsValValidatorFile(pkg, in, out); err != nil {
			app.Errorf("failed to generate jsval validator file: %s", err)
		}
	case validatorCmd.FullCommand():
		if err := generateValidatorFile(pkg, in, out); err != nil {
			app.Errorf("failed to generate validator file: %s", err)
		}
	}

	if *op != "" {
		params := []string{"-w", *op}
		if err := exec.Command("goimports", params...).Run(); err != nil {
			app.Errorf("failed to goimports: %s", err)
		}
	}
}

func generateValidatorFile(pkg *string, fp io.Reader, op io.Writer) error {
	sc, err := schema.Read(fp)
	if err != nil {
		log.Printf("%s", err)
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

	if _, err := op.Write(src); err != nil {
		return err
	}
	return nil
}

func generateStructFile(pkg *string, fp io.Reader, op io.Writer, val, useTitle, nullable bool) error {
	sc, err := schema.Read(fp)
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

	if _, err := op.Write(src); err != nil {
		return err
	}
	return nil
}

func generateJsValValidatorFile(pkg *string, fp io.Reader, op io.Writer) error {
	sc, err := schema.Read(fp)
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
	if err := generator.Process(&src, validators...); err != nil {
		return err
	}

	if _, err := op.Write(src.Bytes()); err != nil {
		return err
	}
	return nil
}
