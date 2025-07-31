package models

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
)

// Product defines the structure for product API responses
// @Description Product represents an item in the store.
// @Tags model
type Product struct {
	ID int `json:"id" example:"1"`

	// Name of the product
	// required: true
	Name string `json:"name" validate:"required" example:"Apple Watch"`

	// Price of the product in cents
	// required: true
	Price int `json:"price" validate:"required,numeric" example:"299"`

	// Creation timestamp (not exposed in response)
	CreatedAt string `json:"-"`

	// Last updated timestamp (not exposed in response)
	UpdatedAt string `json:"-"`
}

func (p *Product) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return validate.Struct(p)
}

type Products []*Product

func GetProducts() Products {
	return productsList
}

func AddProduct(p *Product) {
	productsList = append(productsList, p)
}

func UpdateProduct(prod *Product, id int) error {

	i, err := findProduct(id)
	if err != nil {
		return err
	}

	productsList[i] = prod
	return nil
}

func DeleteProduct(id int) error {

	i, err := findProduct(id)
	if err != nil {
		return err
	}

	productsList = append(productsList[:i], productsList[i+1:]...)

	return nil
}

var ProductNotFound = fmt.Errorf("Product Not Found")

func findProduct(id int) (int, error) {
	for i, prod := range productsList {
		if prod.ID == id {
			return i, nil
		}
	}
	return -1, ProductNotFound
}

func (p *Products) ToJson(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(p)
}

func (p *Product) FromJson(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(p)
}

var productsList = Products{
	{
		ID:        1,
		Name:      "Product 1",
		Price:     120,
		CreatedAt: time.Now().UTC().String(),
		UpdatedAt: time.Now().UTC().String(),
	},
	{
		ID:        2,
		Name:      "Product 2",
		Price:     220,
		CreatedAt: time.Now().UTC().String(),
		UpdatedAt: time.Now().UTC().String(),
	},
}
