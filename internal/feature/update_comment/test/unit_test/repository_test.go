package unit_test_update_comment

import (
	"errors"
	"testing"

	database "commentservice/internal/db"
	mock_database "commentservice/internal/db/test/mock"
	"commentservice/internal/feature/update_comment"
	model "commentservice/internal/model/domain"

	"github.com/stretchr/testify/assert"
)

var dbClient *mock_database.MockDatabaseClient
var updateCommentRepository *update_comment.UpdateCommentRepository

func repositorySetUp(t *testing.T) {
	setUp(t)
	dbClient = mock_database.NewMockDatabaseClient(ctrl)
	updateCommentRepository = update_comment.NewUpdateCommentRepository(database.NewDatabase(dbClient))
}

func TestUpdateCommentInRepository_WhenItReturnsSuccess(t *testing.T) {
	repositorySetUp(t)
	comment := &model.Comment{
		Id:       uint64(1234),
		Username: "usernameA",
		PostId:   "post1",
		Content:  "o meu comentario",
	}
	dbClient.EXPECT().UpdateComment(comment).Return(nil)

	err := updateCommentRepository.UpdateComment(comment)

	assert.Nil(t, err)
}

func TestErrorOnUpdateCommentInRepository_WhenUpdateCommentFails(t *testing.T) {
	repositorySetUp(t)
	comment := &model.Comment{
		Id:       uint64(1234),
		Username: "usernameA",
		PostId:   "post1",
		Content:  "o meu comentario",
	}
	dbClient.EXPECT().UpdateComment(comment).Return(errors.New("some error"))

	err := updateCommentRepository.UpdateComment(comment)

	assert.NotNil(t, err)
}
