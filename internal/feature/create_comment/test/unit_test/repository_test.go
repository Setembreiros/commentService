package unit_test_create_comment

import (
	"errors"
	"testing"

	database "commentservice/internal/db"
	mock_database "commentservice/internal/db/test/mock"
	"commentservice/internal/feature/create_comment"
	model "commentservice/internal/model/domain"

	"github.com/stretchr/testify/assert"
)

var dbClient *mock_database.MockDatabaseClient
var createCommentRepository *create_comment.CreateCommentRepository

func repositorySetUp(t *testing.T) {
	setUp(t)
	dbClient = mock_database.NewMockDatabaseClient(ctrl)
	createCommentRepository = create_comment.NewCreateCommentRepository(database.NewDatabase(dbClient))
}

func TestAddCommentInRepository_WhenItReturnsSuccess(t *testing.T) {
	repositorySetUp(t)
	comment := &model.Comment{
		Username: "usernameA",
		PostId:   "post1",
		Content:  "o meu comentario",
	}
	expectedCommentId := uint64(5)
	dbClient.EXPECT().CreateComment(comment).Return(expectedCommentId, nil)

	commentId, err := createCommentRepository.AddComment(comment)

	assert.Nil(t, err)
	assert.Equal(t, expectedCommentId, commentId)
}

func TestErrorOnAddCommentInRepository_WhenCreateRelationshipFails(t *testing.T) {
	repositorySetUp(t)
	comment := &model.Comment{
		Username: "usernameA",
		PostId:   "post1",
		Content:  "o meu comentario",
	}
	expectedCommentId := uint64(0)
	dbClient.EXPECT().CreateComment(comment).Return(expectedCommentId, errors.New("some error"))

	commentId, err := createCommentRepository.AddComment(comment)

	assert.NotNil(t, err)
	assert.Equal(t, expectedCommentId, commentId)
}
