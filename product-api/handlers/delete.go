package handlers

import (
	"net/http"
	"strconv"

	"github.com/HirokiHanada11/go-microservices/product-api/data"
	"github.com/gorilla/mux"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Returns no content
// responses:
//  201: noContentResponse

// DeleteProduct deletes product from a database
func (p Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	// retrieve the URL parameter from the mux subrouter for put
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Debug("Deleting product", "id", id)

	err := p.productDB.DeleteProduct(id)

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
