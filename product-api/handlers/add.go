package handlers

import (
	"net/http"

	"github.com/HirokiHanada11/go-microservices/product-api/data"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new products
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")

	// recieve the validated data from the middleware using context value
	// .() casts data into the format
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&prod)
}
