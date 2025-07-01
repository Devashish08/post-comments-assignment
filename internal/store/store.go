package store

import "github.com/Devashish08/post-comments-service/internal/model"

type Store interface {
	CreatePost(post *model.Post) error
	GetPost(id string) (*model.Post, error)
	GetAllPosts() ([]*model.Post, error)
	CreateComment(comment *model.Comment) error
	GetCommentsByPostID(postID string) ([]*model.Comment, error)
}
