package event

type CommentWasDeletedEvent struct {
	CommentId uint64 `json:"commentId"`
}

var CommentWasDeletedEventName = "CommentWasDeletedEvent"
