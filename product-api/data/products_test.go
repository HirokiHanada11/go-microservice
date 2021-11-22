package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "Hellloo",
		Price: 1.02,
		SKU:   "abc-hcih-hogao",
	}

	v := NewValidation()
	err := v.Validate(p)

	if err != nil {
		t.Fatal(err)
	}
}
