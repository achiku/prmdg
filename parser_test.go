package main

import (
	"bytes"
	"fmt"
	"go/format"
	"testing"

	"github.com/achiku/varfmt"
	"github.com/kylelemons/godebug/pretty"
	hschema "github.com/lestrrat/go-jshschema"
	schema "github.com/lestrrat/go-jsschema"
	jsval "github.com/lestrrat/go-jsval"
)

func TestSchemaDefitionsPractive(t *testing.T) {
	sc, err := schema.ReadFile("./doc/schema/schema.json")
	if err != nil {
		t.Fatal(err)
	}

	for structName, df := range sc.Definitions {
		t.Logf("%s", varfmt.PublicVarName(structName))
		t.Logf("%s", df.BaseURL())
		for n, tp := range df.Definitions {
			if tp.IsResolved() {
				t.Logf("  %s: %s(%s) res: %t",
					varfmt.PublicVarName(n), tp.Type, tp.Format, df.IsPropRequired(n))
			} else {
				a, err := tp.Resolve(nil)
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("  %s: %s(%s) res: %t",
					varfmt.PublicVarName(n), a.Type, a.Format, df.IsPropRequired(n))
			}
		}
	}
}

func TestSchemaEndpointPractice(t *testing.T) {
	sc, err := schema.ReadFile("./doc/schema/schema.json")
	if err != nil {
		t.Fatal(err)
	}

	for structName, df := range sc.Definitions {
		t.Logf("%s", structName)
		hsc := hschema.New()
		if err := hsc.Extract(df.Extras); err != nil {
			t.Fatal(err)
		}
		for _, e := range hsc.Links {
			t.Logf(" %s: %s", e.Method, e.Href)
			t.Logf("   response: %s", e.Rel)
			if e.Schema != nil {
				t.Log("   request:")
				for name, props := range e.Schema.Properties {
					sh, err := resolveSchema(props, sc)
					if err != nil {
						t.Fatal(err)
					}
					t.Logf("       %s %v, %t", name, sh.Type, df.IsPropRequired(name))
				}
			}
		}
	}
}

func testNewParser(t *testing.T) *Parser {
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
	pretty.Print(res)
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

func TestParseValidator(t *testing.T) {
	parser := testNewParser(t)
	vl, err := parser.ParseValidators()
	if err != nil {
		t.Fatal(err)
	}
	g := jsval.NewGenerator()
	var src bytes.Buffer
	fmt.Fprintln(&src, "import \"github.com/lestrrat/go-jsval\"")
	g.Process(&src, vl...)
	b, err := format.Source(src.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", b)
}

func TestParseActionLargeJSON(t *testing.T) {
	sc, err := schema.ReadFile("./doc/large-example.json")
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
			// resp, err := format.Source([]byte(a.ResponseStruct()))
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// t.Logf("%s", resp)
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
