package models

import (
	"testing"
)

func TestValidation(t *testing.T) {
	p := &Product{}
	p.Name = "Testing Product"
	p.Price = 12.0
	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
