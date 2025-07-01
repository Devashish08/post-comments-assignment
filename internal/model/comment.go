package model

import "time"

type Comment struct {
	ID          string    `json:"id"`
	PostID      string    `json:"post_id"`
	ContentRaw  string    `json:"content"`
	ContentHTML string    `json:"content_html"`
	CreatedAt   time.Time `json:"created_at"`
}
