package handlers

import (
	"net/http"

	"github.com/HirokiHanada11/go-microservices/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//  200: productsResponse

// GetProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	// d, err := json.Marshal(lp)
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// *** DEV NOTES ***
// json.Marshal returns json encoded struct.
// lp is a list of products and it has method ToJSON,
// which writes encoded json to specified io writer, response writer in this case.
// this method is slightly faster as it does not have to allocate any memory for the json encoded data
