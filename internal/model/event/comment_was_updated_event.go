package event

type CommentWasUpdatedEvent struct {
	CommentId uint64 `json:"commentId"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updatedAt"`
}

var CommentWasUpdatedEventName = "CommentWasDeletedEvent"
