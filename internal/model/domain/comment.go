package model

import "time"

type Comment struct {
	Id        uint64    `json:"id"`
	Username  string    `json:"username"`
	PostId    string    `json:"postId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
