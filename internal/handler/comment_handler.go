package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Devashish08/post-comments-service/internal/model"
	"github.com/Devashish08/post-comments-service/internal/store"
	"github.com/go-chi/chi"
)

type CommentHandler struct {
	Store store.Store
}

func NewCommentHandler(s store.Store) *CommentHandler {
	return &CommentHandler{Store: s}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postId")

	var comment model.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)

		return
	}

	if comment.Content == "" {
		http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		return
	}

	comment.PostID = postID

	if err := h.Store.CreateComment(&comment); err != nil {
		if err.Error() == "post with id "+postID+" not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) GetCommentsForPost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postId")
	if _, err := h.Store.GetPost(postID); err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	comments, err := h.Store.GetCommentsByPostID(postID)
	if err != nil {
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}
