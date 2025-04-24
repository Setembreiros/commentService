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
	expectedCommentId := uint64(1234)
	commentId := strconv.FormatUint(expectedCommentId, 10)
	ginContext.Params = []gin.Param{{Key: "commentId", Value: commentId}}
	controllerService.EXPECT().DeleteComment(expectedCommentId).Return(nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`

	controller.DeleteComment(ginContext)

	assert.Equal(t, 200, apiResponse.Code)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequesErrortDeleteComment_WhenMissingCommentID(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("DELETE", "/comment", nil)
	expectedBodyResponse := `{
		"error": true,
		"message": "Missing commentId parameter",
		"content": null
	}`

	controller.DeleteComment(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequesErrortDeleteComment_WhenCommentIdIsNotUint64(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("DELETE", "/comment", nil)
	expectedCommentId := "no uint64"
	ginContext.Params = []gin.Param{{Key: "commentId", Value: expectedCommentId}}
	expectedBodyResponse := `{
		"error": true,
		"message": "CommentId couldn't be parsed. CommentId hould be a positive number",
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
	expectedCommentId := uint64(1234)
	ginContext.Params = []gin.Param{{Key: "commentId", Value: strconv.FormatUint(expectedCommentId, 10)}}
	expectedError := errors.New("some error")
	controllerService.EXPECT().DeleteComment(expectedCommentId).Return(expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.DeleteComment(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
