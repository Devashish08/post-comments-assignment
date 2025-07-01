package main

import (
	"log"
	"net/http"

	"github.com/Devashish08/post-comments-service/internal/handler"
	"github.com/Devashish08/post-comments-service/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	db := store.NewInMemoryStore()
	postHandler := handler.NewPostHandler(db)
	commentHandler := handler.NewCommentHandler(db)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"service is up running"}`))
	})

	r.Route("/posts", func(r chi.Router) {
		r.Post("/", postHandler.CreatePost) // POST /posts
		r.Get("/", postHandler.GetAllPosts) // GET /posts
		r.Route("/{postId}", func(r chi.Router) {
			r.Get("/", postHandler.GetPostByID)
			r.Post("/comments", commentHandler.CreateComment)
			r.Get("/comments", commentHandler.GetCommentsForPost)
		})
	})

	const port = "8080"
	log.Printf("Server is running on port %s...", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
