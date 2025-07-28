package handlers

import (
	"context"
	"corenethttp/models"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Products struct {
	lg *log.Logger
}

func ProductController(l *log.Logger) *Products {
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

	prod, _ := r.Context().Value(KeyProduct{}).(*models.Product)

	models.AddProduct(prod)
}

func (p *Products) ProductUpdate(w http.ResponseWriter, r *http.Request) {
	p.lg.Print("Products Updating")

	prod, _ := r.Context().Value(KeyProduct{}).(*models.Product)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		prod := &models.Product{}
		err := prod.FromJson(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while parsing request body : %v", err), 400)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
