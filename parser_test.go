package main

import (
	"go/format"
	"testing"

	schema "github.com/lestrrat/go-jsschema"
)

func testNewParser(t *testing.T) *Parser {
	// sc, err := schema.ReadFile("./doc/schema/schema.json")
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
		t.Logf("%s", r.Struct())
		// for _, prop := range r.Properties {
		// 	t.Logf("  %s %s: %s:%s %v",
		// 		prop.Name, prop.Types, prop.SecondTypes, prop.Reference, prop.SecondReference)
		// }
	}
	// t.Logf("%v", res)
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
			resp, err := format.Source(a.ResponseStruct())
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%s", resp)
			req, err := format.Source(a.RequestStruct())
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%s", req)
		}
	}
}

func TestPropertyField(t *testing.T) {
	cases := []struct {
		Prop     Property
		Expected string
	}{
		{
			Prop: Property{
				Name:     "name",
				Types:    []string{"string"},
				Format:   "",
				Required: true,
			},
			Expected: "Name string `json:\"name\" schema:\"name\"`",
		},
		{
			Prop: Property{
				Name:     "name",
				Types:    []string{"string"},
				Format:   "",
				Required: false,
			},
			Expected: "Name string `json:\"name,omitempty\" schema:\"name\"`",
		},
		{
			Prop: Property{
				Name:     "id",
				Types:    []string{"integer"},
				Format:   "",
				Required: true,
			},
			Expected: "ID int64 `json:\"id\" schema:\"id\"`",
		},
		{
			Prop: Property{
				Name:     "createdAt",
				Types:    []string{"string"},
				Format:   "date-time",
				Required: true,
			},
			Expected: "CreatedAt time.Time `json:\"createdAt\" schema:\"createdAt\"`",
		},
	}

	for _, c := range cases {
		str := c.Prop.Field()
		if string(str) != c.Expected {
			t.Errorf("want %s got %s", c.Expected, str)
		}
	}
}

func TestResourceStruct(t *testing.T) {
	res := Resource{
		Name:  "task",
		Title: "Task resource",
		Properties: []Property{
			Property{
				Name:     "id",
				Types:    []string{"integer"},
				Format:   "",
				Required: true,
			},
			Property{
				Name:     "name",
				Types:    []string{"string"},
				Format:   "",
				Required: true,
			},
			Property{
				Name:     "createdAt",
				Types:    []string{"string"},
				Format:   "date-time",
				Required: true,
			},
			Property{
				Name:     "completedAt",
				Types:    []string{"string"},
				Format:   "date-time",
				Required: false,
			},
		},
	}

	b := res.Struct()
	ss, err := format.Source(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", ss)
}
