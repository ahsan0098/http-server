// @title CoreNet Product API
// @version 1.0
// @description This is a simple API for managing products.
// @host localhost:9090

package main

import (
	"context"
	_ "corenethttp/docs"
	"corenethttp/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	ProductHdlr := handlers.ProductController(l)

	mux := chi.NewRouter()

	mux.Route("/products", func(r chi.Router) {

		r.Use(ProductHdlr.Validator)

		r.Post("/", ProductHdlr.CreateProduct)

		r.Put("/{id}", ProductHdlr.ProductUpdate)
	})

	mux.Get("/products", ProductHdlr.GetProducts)
	mux.Delete("/products/{id}", ProductHdlr.ProductDelete)
	mux.Get("/swagger/*", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:         "127.0.0.1:9090",
		Handler:      mux,
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
