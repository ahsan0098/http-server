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

// GetProducts godoc
// @Summary Get all products
// @Description Get all the products from the store
// @Tags products
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Product
// @Router /products [get]
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.lg.Print("Products Fetching")

	prods := models.GetProducts()

	err := prods.ToJson(w)
	if err != nil {
		http.Error(w, "Error While Encoding", http.StatusInternalServerError)
	}
}

// CreateProduct godoc
// @Summary Add a new product
// @Description Create a new product in the store
// @Tags products
// @Accept  json
// @Produce  json
// @Param   product  body  models.Product  true  "Product to create"
// @Success 201 {object} models.Product
// @Router /products [post]
func (p *Products) CreateProduct(w http.ResponseWriter, r *http.Request) {
	p.lg.Print("Products Creating")

	prod, _ := r.Context().Value(KeyProduct{}).(*models.Product)

	models.AddProduct(prod)
}

// ProductUpdate godoc
// @Summary Update an existing product
// @Description Update product by ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Product Data"
// @Success 200 {object} models.Product
// @Router /products/{id} [put]
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

// ProductDelete godoc
// @Summary      Delete product
// @Description  Delete a product using its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path int  true  "Product ID"
// @Success      200  {string}  string  "Product deleted successfully"
// @Failure      400  {string}  string  "Invalid ID"
// @Failure      404  {string}  string  "Product not found"
// @Router       /products/{id} [delete]
func (p *Products) ProductDelete(w http.ResponseWriter, r *http.Request) {
	p.lg.Print("Products Deleting")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Product deleted successfully")
}

type KeyProduct struct{}

func (p *Products) Validator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		prod := &models.Product{}

		if err := prod.FromJson(r.Body); err != nil {
			http.Error(w, fmt.Sprintf("Error while parsing request body : %v", err), 400)
			return
		}

		if err := prod.Validate(); err != nil {
			http.Error(w, fmt.Sprintf("Validation Errors : %v", err), 400)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
