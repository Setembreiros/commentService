package unit_test_delete_comment

import (
	"errors"
	"testing"

	database "commentservice/internal/db"
	mock_database "commentservice/internal/db/test/mock"
	"commentservice/internal/feature/delete_comment"

	"github.com/stretchr/testify/assert"
)

var dbClient *mock_database.MockDatabaseClient
var deleteCommentRepository *delete_comment.DeleteCommentRepository

func repositorySetUp(t *testing.T) {
	setUp(t)
	dbClient = mock_database.NewMockDatabaseClient(ctrl)
	deleteCommentRepository = delete_comment.NewDeleteCommentRepository(database.NewDatabase(dbClient))
}

func TestDeleteCommentInRepository_WhenItReturnsSuccess(t *testing.T) {
	repositorySetUp(t)
	expectedCommentId := uint64(5)
	dbClient.EXPECT().DeleteComment(expectedCommentId).Return(nil)

	err := deleteCommentRepository.DeleteComment(expectedCommentId)

	assert.Nil(t, err)
}

func TestErrorOnDeleteCommentInRepository_WhenDeleteRelationshipFails(t *testing.T) {
	repositorySetUp(t)
	expectedCommentId := uint64(5)
	dbClient.EXPECT().DeleteComment(expectedCommentId).Return(errors.New("some error"))

	err := deleteCommentRepository.DeleteComment(expectedCommentId)

	assert.NotNil(t, err)
}
