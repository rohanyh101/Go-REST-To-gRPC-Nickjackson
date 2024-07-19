package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/roh4nyh/rest2micro/handlers"
)

func main() {
	l := log.New(os.Stdout, "hello-world", log.LstdFlags)
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)
	ph := handlers.NewProducts(l)

	mux := http.NewServeMux()
	mux.Handle("GET /hi", hh)
	mux.Handle("GET /bye", gh)
	mux.HandleFunc("GET /products", ph.GetProducts)
	mux.HandleFunc("POST /products", ph.AddProduct)
	mux.HandleFunc("PUT /products/{product_id}", ph.UpdateProduct)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "could not read request body", http.StatusInternalServerError)
			return
		}

		log.Printf("received request: %s", b)
		json.NewEncoder(w).Encode(string(b))
	})

	s := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("server is up and running on port :8080")
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("could not start server: %v", err)
		}
	}()

	// the graceful shutdown is necessary to ensure that the server is not terminated abruptly
	// if any request is still being processed like database transaction, etc.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <-sigChan
	log.Println("received terminate, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.Shutdown(tc)

	// if err := http.ListenAndServe(":8080", mux); err != nil {
	// 	log.Fatalf("could not start server: %v", err)
	// }
}
