package main

import "testing"

func TestGenerateStructFile(t *testing.T) {
	pkg := "taskyapi"
	fp := "./example/doc/schema/schema.json"
	op := ""
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
		if err := generateStructFile(&pkg, fp, &op, c.Validator, c.UseTitle, c.Nullable); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGenerateJsValValidatorFile(t *testing.T) {
	pkg := "taskyapi"
	fp := "./example/doc/schema/schema.json"
	op := ""
	if err := generateJsValValidatorFile(&pkg, fp, &op); err != nil {
		t.Fatal(err)
	}
}

func TestGenerateValidatorFile(t *testing.T) {
	pkg := "taskyapi"
	fp := "./example/doc/schema/schema.json"
	op := ""
	if err := generateValidatorFile(&pkg, fp, &op); err != nil {
		t.Fatal(err)
	}
}
