package integration_test_get_create_comment

import (
	"bytes"
	"commentservice/internal/bus"
	mock_bus "commentservice/internal/bus/test/mock"
	database "commentservice/internal/db"
	"commentservice/internal/feature/create_comment"
	mock_create_comment "commentservice/internal/feature/create_comment/test/mock"
	model "commentservice/internal/model/domain"
	"commentservice/internal/model/event"
	integration_test_arrange "commentservice/test/integration_test_common/arrange"
	integration_test_assert "commentservice/test/integration_test_common/assert"
	integration_test_builder "commentservice/test/integration_test_common/builder"
	"commentservice/test/test_common"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

var timeService *mock_create_comment.MockTimeService
var serviceExternalBus *mock_bus.MockExternalBus
var db *database.Database
var controller *create_comment.CreateCommentController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUp(t *testing.T) {
	// Mocks
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
	ctrl := gomock.NewController(t)
	timeService = mock_create_comment.NewMockTimeService(ctrl)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus := bus.NewEventBus(serviceExternalBus)

	// Real infrastructure and services
	db = integration_test_arrange.CreateTestDatabase(ginContext)
	repository := create_comment.NewCreateCommentRepository(db)
	service := create_comment.NewCreateCommentService(timeService, repository, serviceBus)
	controller = create_comment.NewCreateCommentController(service)
}

func tearDown() {
	db.Client.Clean()
}

func TestCreateComment_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	comment := &model.Comment{
		Username: "usernameA",
		PostId:   "post1",
		Content:  "o meu comentario",
	}
	data, _ := test_common.SerializeData(comment)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/comment", bytes.NewBuffer(data))
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`
	timeNowString := time.Now().UTC().Format(model.TimeLayout)
	timeNow, _ := time.Parse(model.TimeLayout, timeNowString)
	timeService.EXPECT().GetTimeNowUtc().Return(timeNow)
	commentId := integration_test_arrange.GetNextCommentId()
	expectedCommentWasCreatedEvent := &event.CommentWasCreatedEvent{
		CommentId: commentId,
		Username:  comment.Username,
		PostId:    comment.PostId,
		Content:   comment.Content,
		CreatedAt: timeNow.Format(model.TimeLayout),
	}
	expectedEvent := integration_test_builder.NewEventBuilder(t).WithName(event.CommentWasCreatedEventName).WithData(expectedCommentWasCreatedEvent).Build()
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(nil)

	controller.CreateComment(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
	integration_test_assert.AssertCommentExists(t, db, commentId, comment)
}
