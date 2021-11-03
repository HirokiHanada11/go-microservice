package handlers

import (
	"net/http"
	"strconv"

	"github.com/HirokiHanada11/go-microservices/data"
	"github.com/gorilla/mux"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//  200: productsResponse

func (p Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	// retrieve the URL parameter from the mux subrouter for put
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Println("Handle PUT product", id)

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
