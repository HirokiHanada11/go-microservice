package handlers

import (
	"net/http"
	"strconv"

	"github.com/HirokiHanada11/go-microservices/product-api/data"
	"github.com/gorilla/mux"
)

// swagger:route PUT /products/{id} products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// GetProducts returns the products from the data store

func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	// retrieve the URL parameter from the mux subrouter for put
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}

	p.l.Debug("Updating product", "id", id)

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = p.productDB.UpdateProduct(prod)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("Unable to delete product", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Error("Unable to delete product", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}
}
