package handlers

import (
	"log"
	"net/http"

	"github.com/HirokiHanada11/go-microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	//handle an update

	//catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

/*
	json.Marshal returns json encoded struct.
	lp is a list of products and it has method ToJSON,
	which writes encoded json to specified io writer, response writer in this case.
	this method is slightly faster as it does not have to allocate any memory for the json encoded data
*/
func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	// d, err := json.Marshal(lp)
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
