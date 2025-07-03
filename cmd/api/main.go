// Package main provides the entry point for the post-comments REST API service.
// This service provides endpoints for managing blog posts and their comments,
// with support for markdown processing and in-memory storage.
//
// API Endpoints:
//
//	GET /health                     - Health check endpoint
//	POST /posts                     - Create a new post
//	GET /posts                      - Get all posts
//	GET /posts/{postId}             - Get a specific post with comments
//	POST /posts/{postId}/comments   - Create a comment on a post
//	GET /posts/{postId}/comments    - Get all comments for a post
//
// The service runs on port 8080 by default and uses an in-memory store
// for data persistence during the application lifecycle.
package main

import (
	"log"
	"net/http"

	"github.com/Devashish08/post-comments-service/internal/handler"
	"github.com/Devashish08/post-comments-service/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// main initializes and starts the HTTP server with all configured routes and middleware.
// The server uses Chi router with logging and recovery middleware for production readiness.
func main() {
	// Initialize dependencies
	db := store.NewInMemoryStore()
	postHandler := handler.NewPostHandler(db)
	commentHandler := handler.NewCommentHandler(db)

	// Setup router with middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"service is up running"}`))
	})

	// Post and comment routes
	r.Route("/posts", func(r chi.Router) {
		r.Post("/", postHandler.CreatePost) // POST /posts
		r.Get("/", postHandler.GetAllPosts) // GET /posts
		r.Route("/{postId}", func(r chi.Router) {
			r.Get("/", postHandler.GetPostByID)                   // GET /posts/{postId}
			r.Post("/comments", commentHandler.CreateComment)     // POST /posts/{postId}/comments
			r.Get("/comments", commentHandler.GetCommentsForPost) // GET /posts/{postId}/comments
		})
	})

	// Start server
	const port = "8080"
	log.Printf("Server is running on port %s...", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
