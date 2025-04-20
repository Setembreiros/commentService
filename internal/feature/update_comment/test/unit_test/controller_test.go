package unit_test_update_comment

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"commentservice/internal/bus"
	"commentservice/internal/feature/update_comment"
	mock_update_comment "commentservice/internal/feature/update_comment/test/mock"
	model "commentservice/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var controllerLoggerOutput bytes.Buffer
var controllerService *mock_update_comment.MockService
var controllerBus *bus.EventBus
var controller *update_comment.UpdateCommentController

func setUpController(t *testing.T) {
	setUp(t)
	controllerService = mock_update_comment.NewMockService(ctrl)
	controllerBus = &bus.EventBus{}
	controller = update_comment.NewUpdateCommentController(controllerService)
}

func TestUpdateComment(t *testing.T) {
	setUpController(t)
	expectedCommentId := "1234"
	comment := &model.Comment{
		Content: "o meu comentario",
	}
	data, _ := serializeData(comment)
	comment.Id, _ = strconv.ParseUint(expectedCommentId, 10, 64)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/comment", bytes.NewBuffer(data))
	ginContext.Params = []gin.Param{{Key: "commentId", Value: expectedCommentId}}
	controllerService.EXPECT().UpdateComment(comment).Return(nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`

	controller.UpdateComment(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequesErrorOnUpdateComment_WhenMissingCommentID(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("PUT", "/comment", nil)
	expectedBodyResponse := `{
		"error": true,
		"message": "Missing commentId parameter",
		"content": null
	}`

	controller.UpdateComment(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequesErrorOnUpdateComment_WhenCommentIdIsNotUint64(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("PUT", "/comment", nil)
	expectedCommentId := "no uint64"
	ginContext.Params = []gin.Param{{Key: "commentId", Value: expectedCommentId}}
	expectedBodyResponse := `{
		"error": true,
		"message": "CommentId couldn't be parsed. CommentId hould be a positive number",
		"content": null
	}`

	controller.UpdateComment(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("CommentId %v couldn't be parsed", expectedCommentId))
}

func TestInternalServerErrorOnUpdateComment(t *testing.T) {
	setUpController(t)
	expectedCommentId := "1234"
	comment := &model.Comment{
		Username: "usernameA",
		PostId:   "post1",
		Content:  "o meu comentario",
	}
	data, _ := serializeData(comment)
	comment.Id, _ = strconv.ParseUint(expectedCommentId, 10, 64)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/comment", bytes.NewBuffer(data))
	ginContext.Params = []gin.Param{{Key: "commentId", Value: expectedCommentId}}
	expectedError := errors.New("some error")
	controllerService.EXPECT().UpdateComment(comment).Return(expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.UpdateComment(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
