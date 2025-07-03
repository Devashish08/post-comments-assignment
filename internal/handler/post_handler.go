package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Devashish08/post-comments-service/internal/model"
	"github.com/Devashish08/post-comments-service/internal/store"
	"github.com/go-chi/chi"
)

// PostHandler handles HTTP requests related to post operations.
// This handler provides endpoints for creating, retrieving, and listing posts.
type PostHandler struct {
	Store store.Store
}

// PostResponse represents the response format for individual post requests.
// It includes the post data along with all associated comments.
type PostResponse struct {
	model.Post
	Comments []*model.Comment `json:"comments"`
}

// NewPostHandler creates a new PostHandler instance with the provided store.
func NewPostHandler(s store.Store) *PostHandler {
	return &PostHandler{
		Store: s,
	}
}

// CreatePost handles POST /posts requests to create a new post.
// Expects JSON body with post content and returns the created post with generated ID.
//
// Request: POST /posts
// Body: {"content": "Post content here"}
// Response: 201 Created with post JSON, or 400/500 on error
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

// GetAllPosts handles GET /posts requests to retrieve all posts.
// Returns an array of all posts in the system.
//
// Request: GET /posts
// Response: 200 OK with posts array JSON, or 500 on error
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

// GetPostByID handles GET /posts/{postId} requests to retrieve a specific post.
// Returns the post along with all its comments.
//
// Request: GET /posts/{postId}
// Response: 200 OK with post+comments JSON, 404 if not found, or 500 on error
func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postId")

	post, err := h.Store.GetPost(postID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	comments, err := h.Store.GetCommentsByPostID(postID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve comments for post")
		return
	}

	response := PostResponse{
		Post:     *post,
		Comments: comments,
	}

	RespondWithJSON(w, http.StatusOK, response)
}
