package create_comment

import (
	database "commentservice/internal/db"
	model "commentservice/internal/model/domain"
)

type CreateCommentRepository struct {
	dataRepository *database.Database
}

func NewCreateCommentRepository(dataRepository *database.Database) *CreateCommentRepository {
	return &CreateCommentRepository{
		dataRepository: dataRepository,
	}
}

func (r *CreateCommentRepository) AddComment(data *model.Comment) (uint64, error) {
	return r.dataRepository.Client.CreateComment(data)
}
