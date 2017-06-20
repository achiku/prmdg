package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerateStructFile(t *testing.T) {
	pkg := "taskyapi"
	op := ioutil.Discard
	cases := []struct {
		Validator bool
		UseTitle  bool
		Nullable  bool
	}{
		{Validator: false, UseTitle: false, Nullable: false},
		{Validator: true, UseTitle: false, Nullable: false},
		{Validator: true, UseTitle: true, Nullable: true},
	}
	for _, c := range cases {
		fp, err := os.Open("./example/doc/schema/schema.json")
		if err != nil {
			t.Fatal(err)
		}
		if err := generateStructFile(&pkg, fp, op, c.Validator, c.UseTitle, c.Nullable); err != nil {
			t.Fatal(err)
		}
		fp.Close()
	}
}

func TestGenerateJsValValidatorFile(t *testing.T) {
	pkg := "taskyapi"
	fp, err := os.Open("./example/doc/schema/schema.json")
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()
	op := ioutil.Discard
	if err := generateJsValValidatorFile(&pkg, fp, op); err != nil {
		t.Fatal(err)
	}
}

func TestGenerateValidatorFile(t *testing.T) {
	pkg := "taskyapi"
	fp, err := os.Open("./example/doc/schema/schema.json")
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()
	op := ioutil.Discard
	if err := generateValidatorFile(&pkg, fp, op); err != nil {
		t.Fatal(err)
	}
}
