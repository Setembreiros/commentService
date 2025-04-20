package event

type CommentWasUpdatedEvent struct {
	CommentId uint64 `json:"commentId"`
	Username  string `json:"username"`
	PostId    string `json:"postId"`
	Content   string `json:"content"`
	UpdatedAt string `json:"createdAt"`
}

var CommentWasUpdatedEventName = "CommentWasUpdatedEvent"
