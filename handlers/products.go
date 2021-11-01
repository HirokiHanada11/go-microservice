package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/HirokiHanada11/go-microservices/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// json.Marshal returns json encoded struct.
// lp is a list of products and it has method ToJSON,
// which writes encoded json to specified io writer, response writer in this case.
// this method is slightly faster as it does not have to allocate any memory for the json encoded data

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	// d, err := json.Marshal(lp)
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")

	// recieve the validated data from the middleware using context value
	// .() casts data into the format
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&prod)
}

func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	// retrieve the URL parameter from the mux subrouter for put
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}

	p.l.Println("Handle PUT product", id)

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

// an empty struct to be used as a key for passing data from middleware through contex
// http request generates a context autmatically and can be accessed with r.Context()
type KeyProduct struct{}

// middleware function that takes in http handler and returns a new one
func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// send the data through context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
