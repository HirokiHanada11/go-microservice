package data

import (
	"encoding/json"
	"io"
	"time"
)

/*
	`json:""` is a struct tag which specifies the format of the struct when encoded to json
	Product defines the structure for an API product
*/
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product //array of products

/*
	encodes products struct to json
*/
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w) //creates new encoder with io writer input
	return e.Encode(p)
}

/*
	it is cleaner to write CRUD methods with the data
*/
func GetProducts() Products {
	return productList
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Flothy mily coffe",
		Price:       4.55,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       2.50,
		SKU:         "efg456",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
