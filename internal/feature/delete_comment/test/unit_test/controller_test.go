package unit_test_delete_comment

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"commentservice/internal/bus"
	"commentservice/internal/feature/delete_comment"
	mock_delete_comment "commentservice/internal/feature/delete_comment/test/mock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var controllerLoggerOutput bytes.Buffer
var controllerService *mock_delete_comment.MockService
var controllerBus *bus.EventBus
var controller *delete_comment.DeleteCommentController

func setUpController(t *testing.T) {
	setUp(t)
	controllerService = mock_delete_comment.NewMockService(ctrl)
	controllerBus = &bus.EventBus{}
	controller = delete_comment.NewDeleteCommentController(controllerService)
}

func TestDeleteComment(t *testing.T) {
	setUpController(t)
	ginContext.Request = httptest.NewRequest(http.MethodDelete, "/comment", nil)
	expectedPostId := "post1"
	expectedCommentId := uint64(1234)
	commentId := strconv.FormatUint(expectedCommentId, 10)
	ginContext.Params = []gin.Param{
		{Key: "postId", Value: expectedPostId},
		{Key: "commentId", Value: commentId},
	}
	controllerService.EXPECT().DeleteComment(expectedPostId, expectedCommentId).Return(nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`

	controller.DeleteComment(ginContext)

	assert.Equal(t, 200, apiResponse.Code)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequesErrorOnDeleteComment_WhenMissingPostId(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("DELETE", "/comment", nil)
	expectedBodyResponse := `{
		"error": true,
		"message": "Missing postId parameter",
		"content": null
	}`

	controller.DeleteComment(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequesErrorOnDeleteComment_WhenMissingCommentId(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("DELETE", "/comment", nil)
	expectedPostId := "post1"
	ginContext.Params = []gin.Param{
		{Key: "postId", Value: expectedPostId},
	}
	expectedBodyResponse := `{
		"error": true,
		"message": "Missing commentId parameter",
		"content": null
	}`

	controller.DeleteComment(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequesErrorOnDeleteComment_WhenCommentIdIsNotUint64(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("DELETE", "/comment", nil)
	expectedPostId := "post1"
	expectedCommentId := "no uint64"
	ginContext.Params = []gin.Param{
		{Key: "postId", Value: expectedPostId},
		{Key: "commentId", Value: expectedCommentId},
	}
	expectedBodyResponse := `{
		"error": true,
		"message": "CommentId couldn't be parsed. CommentId should be a positive number",
		"content": null
	}`

	controller.DeleteComment(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("CommentId %v couldn't be parsed", expectedCommentId))
}

func TestInternalServerErrorOnDeleteComment(t *testing.T) {
	setUpController(t)
	ginContext.Request = httptest.NewRequest(http.MethodDelete, "/comment", nil)
	expectedPostId := "post1"
	expectedCommentId := uint64(1234)
	commentId := strconv.FormatUint(expectedCommentId, 10)
	ginContext.Params = []gin.Param{
		{Key: "postId", Value: expectedPostId},
		{Key: "commentId", Value: commentId},
	}
	expectedError := errors.New("some error")
	controllerService.EXPECT().DeleteComment(expectedPostId, expectedCommentId).Return(expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.DeleteComment(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
