package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/Devashish08/post-comments-service/internal/model"
	"github.com/google/uuid"
)

type InMemoryStore struct {
	mu       sync.RWMutex
	posts    map[string]*model.Post
	comments map[string][]*model.Comment
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		posts:    make(map[string]*model.Post),
		comments: make(map[string][]*model.Comment),
	}
}

func (s *InMemoryStore) CreatePost(post *model.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	post.ID = uuid.NewString()
	post.CreatedAt = time.Now()

	s.posts[post.ID] = post
	return nil
}

func (s *InMemoryStore) GetPost(id string) (*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, exists := s.posts[id]
	if !exists {
		return nil, fmt.Errorf("post not found")
	}

	return post, nil
}

func (s *InMemoryStore) GetAllPosts() ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	allPosts := make([]*model.Post, 0, len(s.posts))
	for _, post := range s.posts {
		allPosts = append(allPosts, post)
	}

	return allPosts, nil
}

func (s *InMemoryStore) CreateComment(comment *model.Comment) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	_, exists := s.posts[comment.PostID]
	if !exists {
		return fmt.Errorf("post with id %s not found", comment.PostID)
	}

	comment.ID = uuid.NewString()
	comment.CreatedAt = time.Now()

	s.comments[comment.PostID] = append(s.comments[comment.PostID], comment)
	return nil
}

func (s *InMemoryStore) GetCommentsByPostID(postID string) ([]*model.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	comments := s.comments[postID]
	if comments == nil {
		return []*model.Comment{}, nil
	}

	return comments, nil
}
