package update_comment

import (
	database "commentservice/internal/db"
	model "commentservice/internal/model/domain"
)

type UpdateCommentRepository struct {
	dataRepository *database.Database
}

func NewUpdateCommentRepository(dataRepository *database.Database) *UpdateCommentRepository {
	return &UpdateCommentRepository{
		dataRepository: dataRepository,
	}
}

func (r *UpdateCommentRepository) UpdateComment(data *model.Comment) error {
	return r.dataRepository.Client.UpdateComment(data)
}
