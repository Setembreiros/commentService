package event

type CommentWasCreatedEvent struct {
	CommentId uint64 `json:"commentId"`
	Username  string `json:"username"`
	PostId    string `json:"postId"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

var CommentWasCreatedEventName = "CommentWasCreatedEvent"
