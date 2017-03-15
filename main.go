package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
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

	structCmd    = app.Command("struct", "generate struct file")
	validatorCmd = app.Command("validator", "generate validator file")
)

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch cmd {
	case structCmd.FullCommand():
		if err := generateStructFile(pkg, *fp, op); err != nil {
			app.Errorf("failed to generate struct file: %s", err)
		}
	case validatorCmd.FullCommand():
		if err := generateValidatorFile(pkg, *fp, op); err != nil {
			app.Errorf("failed to generate validator file: %s", err)
		}
	}
}

func generateStructFile(pkg *string, fp string, op *string) error {
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
	sort.Strings(resKeys)
	for _, k := range resKeys {
		res := resources[k]
		ss, err := format.Source(res.Struct())
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
			req, err := format.Source(action.RequestStruct())
			if err != nil {
				return errors.Wrapf(err, "failed to format request struct: %s, %s", k, action.Href)
			}
			src = append(src, req...)

			resp, err := format.Source(action.ResponseStruct())
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
	} else {
		out = os.Stdout
	}
	if _, err := out.Write(src); err != nil {
		return err
	}
	return nil
}

func generateValidatorFile(pkg *string, fp string, op *string) error {
	sc, err := schema.ReadFile(fp)
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", fp)
	}
	parser := NewParser(sc, *pkg)
	validators, err := parser.ParseValidators()
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
	} else {
		out = os.Stdout
	}
	if _, err := out.Write(src.Bytes()); err != nil {
		return err
	}
	return nil
}
