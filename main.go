package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HirokiHanada11/go-microservices/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	ph := handlers.NewProducts(l)

	/*
		servemux stands for Http request multiplexer
		it matches the URL of each coming request against a list of registered patterns
		and calls the handler for the pattern that most closely matches the URL
	*/

	// creating server (serve mux) with gorilla mux library
	// mux simplifies the process of taking parameters from URL
	sm := mux.NewRouter()

	// defining method specific subrouters
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts) // regex can be used directly inside of the URL, and it autmatically does the matching
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"} //this package creates API Documentation UI
	sh := middleware.Redoc(ops, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./"))) //serves files inside the directory matching URL path

	//Defining server struct
	s := &http.Server{
		Addr:         ":9090",           //address of the server to run on
		Handler:      sm,                //specify handler
		IdleTimeout:  120 * time.Second, //timeout for the tcp connections to stay idle
		ReadTimeout:  1 * time.Second,   //max read time
		WriteTimeout: 1 * time.Second,   //max write time
	}

	/*
		starting a webserver with a goroutine so that it can be gracefully shutdown
	*/
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	/*
		sigchan is a channel that accepts termination signals for the program.
		once the channel receives the signal, it logs statement
	*/
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	/*
		contexts are used to control the progress of the http requests being handled in a goroutine
		all contexts have parent and children relationships and once the parent contex is canceled,
		all the child contexts are also canceled
		contex.timeout gives the timeout to the contex and once it runs out, it shutsdown the goroutine
	*/
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
	defer cancel()
}
