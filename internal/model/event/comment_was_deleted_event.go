package event

type CommentWasDeletedEvent struct {
	PostId    string `json:"postId"`
	CommentId uint64 `json:"commentId"`
}

var CommentWasDeletedEventName = "CommentWasDeletedEvent"
