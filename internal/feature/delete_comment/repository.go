package delete_comment

import (
	database "commentservice/internal/db"
)

type DeleteCommentRepository struct {
	dataRepository *database.Database
}

func NewDeleteCommentRepository(dataRepository *database.Database) *DeleteCommentRepository {
	return &DeleteCommentRepository{
		dataRepository: dataRepository,
	}
}

func (r *DeleteCommentRepository) DeleteComment(commentId uint64) error {
	return r.dataRepository.Client.DeleteComment(commentId)
}
