// Package store provides data persistence abstractions for the post-comments service.
package store

import "github.com/Devashish08/post-comments-service/internal/model"

// Store defines the contract for data persistence operations.
// Implementations should be thread-safe and handle errors consistently.
type Store interface {
	// CreatePost persists a new post with auto-generated ID and timestamp.
	CreatePost(post *model.Post) error

	// GetPost retrieves a post by its ID.
	GetPost(id string) (*model.Post, error)

	// GetAllPosts retrieves all posts from storage.
	GetAllPosts() ([]*model.Post, error)

	// CreateComment persists a new comment with auto-generated ID and timestamp.
	// Validates that the associated post exists.
	CreateComment(comment *model.Comment) error

	// GetCommentsByPostID retrieves all comments for a specific post.
	GetCommentsByPostID(postID string) ([]*model.Comment, error)
}
