package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Devashish08/post-comments-service/internal/model"
	"github.com/Devashish08/post-comments-service/internal/store"
	"github.com/go-chi/chi"
)

type PostHandler struct {
	Store store.Store
}

func NewPostHandler(s store.Store) *PostHandler {
	return &PostHandler{
		Store: s,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if post.Content == "" {
		http.Error(w, "Post content cannot be empty", http.StatusBadRequest)
		return
	}

	if err := h.Store.CreatePost(&post); err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Store.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postId")

	post, err := h.Store.GetPost(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}
