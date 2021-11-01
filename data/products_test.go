package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "Hellloo",
		Price: 1.02,
		SKU:   "abc-hcih-hogao",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
