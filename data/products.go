package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

/*
	`json:""` is a struct tag which specifies the format of the struct when encoded to json
	Product defines the structure for an API product
*/

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU) // defines a custom validation tag
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	//  sku is of format abc-abc-sdfsf
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1) //-1 extracts all matches

	return len(matches) == 1
}

type Products []*Product //array of products

// decodes json to product struct
// takes io reader as a parameter (Response Body) and returns an error
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// encodes products struct to json
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w) //creates new encoder with io writer input
	return e.Encode(p)
}

// 	it is cleaner to write CRUD methods with the data
func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	p.ID = id
	productList[pos] = p

	return nil
}

// DeleteProduct deletes a product from the database
func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
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
