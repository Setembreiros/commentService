package unit_test_create_comment

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"commentservice/internal/bus"
	"commentservice/internal/feature/create_comment"
	mock_create_comment "commentservice/internal/feature/create_comment/test/mock"
	model "commentservice/internal/model/domain"

	"github.com/go-playground/assert/v2"
)

var controllerLoggerOutput bytes.Buffer
var controllerService *mock_create_comment.MockService
var controllerBus *bus.EventBus
var controller *create_comment.CreateCommentController

func setUpHandler(t *testing.T) {
	setUp(t)
	controllerService = mock_create_comment.NewMockService(ctrl)
	controllerBus = &bus.EventBus{}
	controller = create_comment.NewCreateCommentController(controllerService)
}

func TestCreateComment(t *testing.T) {
	setUpHandler(t)
	comment := &model.Comment{
		Username: "usernameA",
		PostId:   "post1",
		Content:  "o meu comentario",
	}
	data, _ := serializeData(comment)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/comment", bytes.NewBuffer(data))
	controllerService.EXPECT().AddComment(comment).Return(nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`

	controller.CreateComment(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnCreateComment(t *testing.T) {
	setUpHandler(t)
	comment := &model.Comment{
		Username: "usernameA",
		PostId:   "post1",
		Content:  "o meu comentario",
	}
	data, _ := serializeData(comment)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/comment", bytes.NewBuffer(data))
	expectedError := errors.New("some error")
	controllerService.EXPECT().AddComment(comment).Return(expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.CreateComment(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
