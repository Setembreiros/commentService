package model

import "time"

type Comment struct {
	CommentId uint64    `json:"commentId"`
	Username  string    `json:"username"`
	PostId    string    `json:"postId"`
	Content   string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}
