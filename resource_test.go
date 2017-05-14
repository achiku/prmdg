package main

import (
	"go/format"
	"testing"

	schema "github.com/lestrrat/go-jsschema"
)

func TestRefToStructName(t *testing.T) {
	cases := []struct {
		Prop Property
	}{
		{
			Prop: Property{
				Name:      "coupon",
				Types:     []schema.PrimitiveType{schema.ObjectType},
				Format:    "",
				Reference: "#/definitions/coupon",
				Required:  true,
			},
		},
		{
			Prop: Property{
				Name:      "couponType",
				Types:     []schema.PrimitiveType{schema.ObjectType},
				Format:    "",
				Reference: "#/definitions/coupon/definitions/types",
				Required:  true,
			},
		},
	}

	for _, c := range cases {
		t.Logf(c.Prop.refToStructName())
	}
}

func TestSimplePropertyField(t *testing.T) {
	cases := []struct {
		Prop     Property
		Expected string
	}{
		{
			Prop: Property{
				Name:     "name",
				Types:    []schema.PrimitiveType{schema.StringType},
				Format:   "",
				Required: true,
			},
			Expected: "Name string `json:\"name\" schema:\"name\"`",
		},
		{
			Prop: Property{
				Name:     "name",
				Types:    []schema.PrimitiveType{schema.StringType},
				Format:   "",
				Required: false,
			},
			Expected: "Name string `json:\"name,omitempty\" schema:\"name\"`",
		},
		{
			Prop: Property{
				Name:     "id",
				Types:    []schema.PrimitiveType{schema.IntegerType},
				Format:   "",
				Required: true,
			},
			Expected: "ID int64 `json:\"id\" schema:\"id\"`",
		},
		{
			Prop: Property{
				Name:     "createdAt",
				Types:    []schema.PrimitiveType{schema.StringType},
				Format:   "date-time",
				Required: true,
			},
			Expected: "CreatedAt time.Time `json:\"createdAt\" schema:\"createdAt\"`",
		},
	}

	for _, c := range cases {
		str := c.Prop.Field(FormatOption{
			Schema:    true,
			UseTitle:  false,
			Validator: false,
		})
		if string(str) != c.Expected {
			t.Errorf("want %s got %s", c.Expected, str)
		}
	}
}

func TestResourceStruct(t *testing.T) {
	res := Resource{
		Name:  "task",
		Title: "Task resource",
		Properties: []*Property{
			&Property{
				Name:     "id",
				Types:    []schema.PrimitiveType{schema.IntegerType},
				Format:   "",
				Required: true,
			},
			&Property{
				Name:     "name",
				Types:    []schema.PrimitiveType{schema.StringType},
				Format:   "",
				Required: true,
			},
			&Property{
				Name:     "createdAt",
				Types:    []schema.PrimitiveType{schema.StringType},
				Format:   "date-time",
				Required: true,
			},
			&Property{
				Name:     "completedAt",
				Types:    []schema.PrimitiveType{schema.StringType},
				Format:   "date-time",
				Required: false,
			},
		},
	}

	b := res.Struct(FormatOption{})
	ss, err := format.Source(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", ss)
}
