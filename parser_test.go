package main

import (
	"go/format"
	"testing"

	schema "github.com/lestrrat/go-jsschema"
)

func testNewParser(t *testing.T) *Parser {
	// sc, err := schema.ReadFile("./doc/large-example.json")
	sc, err := schema.ReadFile("./doc/schema/schema.json")
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
	// pretty.Print(res)
	// log.Printf("%v", res)
	for key, r := range res {
		t.Logf("%s/%s", key, r.Name)
		t.Logf("%s", r.Struct(FormatOption{}))
		// for _, prop := range r.Properties {
		// 	t.Logf("  %s %s: %s:%s %v",
		// 		prop.Name, prop.Types, prop.SecondTypes, prop.Reference, prop.SecondReference)
		// }
	}
	// t.Logf("%v", res)
}

func TestParseValidators(t *testing.T) {
	parser := testNewParser(t)
	vals, err := parser.ParseValidators()
	if err != nil {
		t.Fatal(err)
	}
	// for _, v := range vals {
	// 	t.Logf("%s", v.RegexpConst())
	// 	t.Logf("%s", v.RegexpVar())
	// 	t.Logf("%s", v.ValidatorFunc())
	// }
	ss, err := format.Source(vals.Render())
	if err != nil {
		t.Errorf("%s", vals.Render())
		t.Fatal(err)
	}
	t.Logf("%s", ss)
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
	// pretty.Print(res)
	for key, actions := range res {
		t.Log(key)
		for _, action := range actions {
			t.Logf("  %s: %s", action.Method, action.Href)
		}
	}
}

func TestParseActionLargeJSON(t *testing.T) {
	sc, err := schema.ReadFile("./doc/schema/schema.json")
	if err != nil {
		t.Fatal(err)
	}
	parser := &Parser{
		schema:  sc,
		pkgName: "model",
	}

	res, err := parser.ParseResources()
	if err != nil {
		t.Fatal(err)
	}
	act, err := parser.ParseActions(res)
	if err != nil {
		t.Fatal(err)
	}
	for _, ac := range act {
		for _, a := range ac {
			resp, err := format.Source(a.ResponseStruct(FormatOption{}))
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%s", resp)
			req, err := format.Source(a.RequestStruct(FormatOption{}))
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%s", req)
		}
	}
}
