package main

import (
	"context"
	"corenethttp/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	ProductHdlr := handlers.NewProducts(l)

	ms := mux.NewRouter()

	ms.HandleFunc("/products", ProductHdlr.GetProducts).Methods("GET")

	
	postSub := ms.Methods("POST").Subrouter()
	postSub.Use(ProductHdlr.Validator)
	postSub.HandleFunc("/products", ProductHdlr.CreateProduct)

	
	putSub := ms.Methods("PUT").Subrouter()
	putSub.Use(ProductHdlr.Validator)
	putSub.HandleFunc("/products/{id:[0-9]+}", ProductHdlr.ProductUpdate)

	
	server := &http.Server{
		Addr:         "127.0.0.1:9090",
		Handler:      ms,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	SigChan := make(chan os.Signal, 1)
	signal.Notify(SigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Server running on http://127.0.0.1:9090")
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sig := <-SigChan

	l.Println("Received termination,  Graceful shutdown", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

	log.Println("Server shutdown complete.")
}
