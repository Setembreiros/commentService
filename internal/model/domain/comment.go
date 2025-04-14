package model

import "time"

type Comment struct {
	Username  string    `json:"username"`
	PostId    string    `json:"postId"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}
