package integration_test_get_delete_comment

import (
	"commentservice/internal/bus"
	mock_bus "commentservice/internal/bus/test/mock"
	database "commentservice/internal/db"
	"commentservice/internal/feature/delete_comment"
	model "commentservice/internal/model/domain"
	"commentservice/internal/model/event"
	integration_test_arrange "commentservice/test/integration_test_common/arrange"
	integration_test_assert "commentservice/test/integration_test_common/assert"
	integration_test_builder "commentservice/test/integration_test_common/builder"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

var serviceExternalBus *mock_bus.MockExternalBus
var db *database.Database
var controller *delete_comment.DeleteCommentController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUp(t *testing.T) {
	// Mocks
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
	ctrl := gomock.NewController(t)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus := bus.NewEventBus(serviceExternalBus)

	// Real infrastructure and services
	db = integration_test_arrange.CreateTestDatabase(ginContext)
	repository := delete_comment.NewDeleteCommentRepository(db)
	service := delete_comment.NewDeleteCommentService(repository, serviceBus)
	controller = delete_comment.NewDeleteCommentController(service)
}

func tearDown() {
	db.Client.Clean()
}

func TestDeleteComment_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	ginContext.Request = httptest.NewRequest(http.MethodDelete, "/comment", nil)
	expectedCommentId := populateDb(t)
	ginContext.Params = []gin.Param{{Key: "commentId", Value: strconv.FormatUint(expectedCommentId, 10)}}
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`
	expectedCommentWasDeletedEvent := &event.CommentWasDeletedEvent{
		CommentId: expectedCommentId,
	}
	expectedEvent := integration_test_builder.NewEventBuilder(t).WithName(event.CommentWasDeletedEventName).WithData(expectedCommentWasDeletedEvent).Build()
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(nil)

	controller.DeleteComment(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
	integration_test_assert.AssertCommentDoesNotExists(t, db, expectedCommentId)
}

func populateDb(t *testing.T) uint64 {
	existingComment := &model.Comment{
		Username:  "usernameA",
		PostId:    "post1",
		Content:   "o meu comentario",
		CreatedAt: time.Now(),
	}

	return integration_test_arrange.AddComment(t, existingComment)
}
