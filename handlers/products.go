package handlers

import (
	"context"
	"corenethttp/models"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	lg *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.lg.Print("Products Fetching")

	prods := models.GetProducts()

	err := prods.ToJson(w)
	if err != nil {
		http.Error(w, "Error While Encoding", http.StatusInternalServerError)
	}
}

func (p *Products) CreateProduct(w http.ResponseWriter, r *http.Request) {
	p.lg.Print("Products Creating")

	prod, ok := r.Context().Value(KeyProduct{}).(*models.Product)
	if !ok {
		http.Error(w, "Missing product in context", http.StatusInternalServerError)
		return
	}

	models.AddProduct(prod)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Product added successfully")
}

func (p *Products) ProductUpdate(w http.ResponseWriter, r *http.Request) {
	p.lg.Print("Products Updating")

	prod, ok := r.Context().Value(KeyProduct{}).(*models.Product)
	if !ok {
		http.Error(w, "Missing product in context", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, fmt.Sprintf("Error while converting request parameter: %v", err), http.StatusBadRequest)
		return
	}
	err = models.UpdateProduct(prod, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Product updated successfully")
}

type KeyProduct struct{}

func (p *Products) Validator(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			prod := &models.Product{}
			err := prod.FromJson(r.Body)

			if err != nil {
				http.Error(w, fmt.Sprintf("Error while parsing request body : %v", err), 400)
				return
			}

			ctx := context.WithValue(r.Context(), KeyProduct{}, prod)

			// Pass the new context with product to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
