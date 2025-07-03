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

// CommentHandler handles HTTP requests related to comment operations.
// This handler provides endpoints for creating comments and retrieving comments for posts.
// It includes markdown processing functionality to convert raw markdown to HTML.
type CommentHandler struct {
	Store    store.Store
	Markdown goldmark.Markdown
}

// NewCommentHandler creates a new CommentHandler instance with the provided store.
// Initializes the markdown processor with default settings for safe HTML generation.
func NewCommentHandler(s store.Store) *CommentHandler {
	md := goldmark.New(
		goldmark.WithRendererOptions(),
	)

	return &CommentHandler{Store: s, Markdown: md}
}

// CreateComment handles POST /posts/{postId}/comments requests to create a new comment.
// Processes markdown content and stores both raw and HTML versions of the comment.
//
// Request: POST /posts/{postId}/comments
// Body: {"content": "Comment content in **markdown**"}
// Response: 201 Created with comment JSON, 400/404/500 on error
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

// GetCommentsForPost handles GET /posts/{postId}/comments requests.
// Returns all comments associated with the specified post.
//
// Request: GET /posts/{postId}/comments
// Response: 200 OK with comments array JSON, 404 if post not found, or 500 on error
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
