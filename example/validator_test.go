package exampleapi

import "testing"

func TestValidator(t *testing.T) {
	aa := "aaa"
	if err := TaskInstancesValidator.Validate(aa); err != nil {
		t.Fatal(err)
	}
}
