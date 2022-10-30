package main

import (
	"testing"

	schema "github.com/lestrrat-go/jsschema"
)

func testNewParser(t *testing.T) *Parser {
	sc, err := schema.ReadFile("./_example/doc/schema/schema.json")
	if err != nil {
		t.Fatal(err)
	}
	return &Parser{
		schema:  sc,
		pkgName: "model",
	}
}

func TestParseResources(t *testing.T) {
	parser := testNewParser(t)
	res, err := parser.ParseResources()
	if err != nil {
		t.Fatal(err)
	}
	for key, r := range res {
		t.Logf("%s/%s", key, r.Name)
	}
}

func TestParseJsValValidators(t *testing.T) {
	parser := testNewParser(t)
	vals, err := parser.ParseJsValValidators()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range vals {
		t.Logf("%s", v.Name)
	}
}

func TestParseValidators(t *testing.T) {
	parser := testNewParser(t)
	vals, err := parser.ParseValidators()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range vals {
		t.Logf("%s", v.Name)
	}
}

func TestParseActions(t *testing.T) {
	parser := testNewParser(t)
	r, err := parser.ParseResources()
	if err != nil {
		t.Fatal(err)
	}

	res, err := parser.ParseActions(r)
	if err != nil {
		t.Fatal(err)
	}
	for key, actions := range res {
		t.Log(key)
		for _, action := range actions {
			t.Logf("  %s: %s", action.Method, action.Href)
		}
	}
}
