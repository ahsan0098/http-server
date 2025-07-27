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
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// HelloHdlr := handlers.NewHello(l)
	ProductHdlr := handlers.NewProducts(l)
	
	mux := http.NewServeMux()
	mux.Handle("/products/", ProductHdlr)
	// mux.Handle("/", HelloHdlr)
	// mux.Handle("/home", handlers.HomeHdlr())


	// creating own server
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
