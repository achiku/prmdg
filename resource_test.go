package main

import (
	"go/format"
	"testing"
)

func TestSimplePropertyField(t *testing.T) {
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
		Properties: []*Property{
			&Property{
				Name:     "id",
				Types:    []string{"integer"},
				Format:   "",
				Required: true,
			},
			&Property{
				Name:     "name",
				Types:    []string{"string"},
				Format:   "",
				Required: true,
			},
			&Property{
				Name:     "createdAt",
				Types:    []string{"string"},
				Format:   "date-time",
				Required: true,
			},
			&Property{
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
