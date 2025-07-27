package handlers

import (
	"corenethttp/models"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	lg *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		p.GetProducts(w, r)
		return
	}

	if r.Method == http.MethodPost {
		p.CreateProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {

		reg := regexp.MustCompile(`/products/(\d+)$`)
		data := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(data) != 1 {
			http.Error(w, "No Url Parameter Passed", 400)
			return

		}

		if len(data[0]) != 2 {
			http.Error(w, "Wrong Url Parameters Passed", 400)
			return

		}

		id, err := strconv.Atoi(data[0][1])
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid Non integer number Passed : %s", data[0][1]), 400)
			return
		}
		p.ProductUpdate(w, r, id)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.lg.Print("Products Fetching")

	prods := models.GetProducts()

	// //simple marshling
	// m, _ := json.Marshal(prods)

	// using Encoder that does not stores but directly passs to io.reader

	err := prods.ToJson(w)
	if err != nil {
		http.Error(w, "Error While Encoding", http.StatusInternalServerError)
	}
}

func (p *Products) CreateProduct(w http.ResponseWriter, r *http.Request) {
	p.lg.Print("Products Creating")

	prod := &models.Product{}
	err := prod.FromJson(r.Body)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error while parsing request body : %v", err), 400)
	}

	models.AddProduct(prod)
}

func (p *Products) ProductUpdate(w http.ResponseWriter, r *http.Request, id int) {
	p.lg.Print("Products Updating")

	prod := &models.Product{}
	err := prod.FromJson(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while parsing request body: %v", err), http.StatusBadRequest)
		return
	}

	prod.ID = id

	err = models.UpdateProduct(prod, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Product updated successfully")
}
