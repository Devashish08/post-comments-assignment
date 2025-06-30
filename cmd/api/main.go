package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"service is up running"}`))
	})

	const port = "8080"
	log.Printf("Server is running on port %s...", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
