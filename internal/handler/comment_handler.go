package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Devashish08/post-comments-service/internal/model"
	"github.com/Devashish08/post-comments-service/internal/store"
	"github.com/go-chi/chi"
	"github.com/yuin/goldmark"
)

type CommentHandler struct {
	Store    store.Store
	Markdown goldmark.Markdown
}

func NewCommentHandler(s store.Store) *CommentHandler {
	md := goldmark.New(
		goldmark.WithRendererOptions(),
	)

	return &CommentHandler{Store: s, Markdown: md}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postId")

	var requestBody struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.Content == "" {
		http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		return
	}

	var htmlBuffer bytes.Buffer
	if err := h.Markdown.Convert([]byte(requestBody.Content), &htmlBuffer); err != nil {
		http.Error(w, "Failed to parse markdown", http.StatusInternalServerError)
		return
	}

	comment := model.Comment{
		PostID:      postID,
		ContentRaw:  requestBody.Content,
		ContentHTML: htmlBuffer.String(),
	}

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
