package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Goodbye struct {
	l *log.Logger
}

func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

type Response struct {
	Name string `json:"name"validate:"required,min=3,max=100"`
}

func (h *Goodbye) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Goodbye")

	// b, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, "could not read request body", http.StatusInternalServerError)
	// 	return
	// }

	var resp Response
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		http.Error(w, "could not decode request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(resp); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %v", err), http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf("Goodbye %s", resp.Name)))
}
