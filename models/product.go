package models

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	CreateAt  string `json:"-"`
	UpdatedAt string `json:"-"`
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
		CreateAt:  time.Now().UTC().String(),
		UpdatedAt: time.Now().UTC().String(),
	},
	{
		ID:        2,
		Name:      "Product 2",
		Price:     220,
		CreateAt:  time.Now().UTC().String(),
		UpdatedAt: time.Now().UTC().String(),
	},
}
