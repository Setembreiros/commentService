package integration_test_assert

import (
	database "commentservice/internal/db"
	model "commentservice/internal/model/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertCommentExists(t *testing.T, db *database.Database, expectedCommentId uint64, expectedComment *model.Comment) {
	comment, err := db.Client.GetCommentById(expectedCommentId)
	assert.Nil(t, err)
	assert.Equal(t, expectedCommentId, comment.Id)
	assert.Equal(t, expectedComment.PostId, comment.PostId)
	assert.Equal(t, expectedComment.Username, comment.Username)
	assert.Equal(t, expectedComment.Content, comment.Content)
}

func AssertCommentDoesNotExists(t *testing.T, db *database.Database, expectedCommentId uint64) {
	_, err := db.Client.GetCommentById(expectedCommentId)
	assert.NotNil(t, err)
}
