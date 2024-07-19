package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/roh4nyh/rest2micro/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
		return
	}

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Products")

	product_id := r.PathValue("product_id")
	id, _ := strconv.Atoi(product_id)

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "unable to unmarshal json", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Product updated successfully"))
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()

	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "unable to marshal json", http.StatusInternalServerError)
		return
	}

	// d, err := json.Marshal(lp)
	// if err != nil {
	// 	http.Error(w, "unable to marshal json", http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.Write(d)
}

// func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
