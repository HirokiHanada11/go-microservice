package handlers

import (
	"net/http"

	"github.com/HirokiHanada11/go-microservices/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//  200: productsResponse

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")

	// recieve the validated data from the middleware using context value
	// .() casts data into the format
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&prod)
}
