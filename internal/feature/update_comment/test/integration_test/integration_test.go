package integration_test_get_update_comment

import (
	"bytes"
	"commentservice/internal/bus"
	mock_bus "commentservice/internal/bus/test/mock"
	database "commentservice/internal/db"
	"commentservice/internal/feature/update_comment"
	mock_update_comment "commentservice/internal/feature/update_comment/test/mock"
	model "commentservice/internal/model/domain"
	"commentservice/internal/model/event"
	integration_test_arrange "commentservice/test/integration_test_common/arrange"
	integration_test_assert "commentservice/test/integration_test_common/assert"
	integration_test_builder "commentservice/test/integration_test_common/builder"
	"commentservice/test/test_common"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

var timeService *mock_update_comment.MockTimeService
var serviceExternalBus *mock_bus.MockExternalBus
var db *database.Database
var controller *update_comment.UpdateCommentController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUp(t *testing.T) {
	// Mocks
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
	ctrl := gomock.NewController(t)
	timeService = mock_update_comment.NewMockTimeService(ctrl)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus := bus.NewEventBus(serviceExternalBus)

	// Real infrastructure and services
	db = integration_test_arrange.CreateTestDatabase(ginContext)
	repository := update_comment.NewUpdateCommentRepository(db)
	service := update_comment.NewUpdateCommentService(timeService, repository, serviceBus)
	controller = update_comment.NewUpdateCommentController(service)
}

func tearDown() {
	db.Client.Clean()
}

func TestUpdateComment_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	expectedCommentId := populateDb(t)
	comment := &model.Comment{
		Content: "o meu comentario actualizado",
	}
	data, _ := test_common.SerializeData(comment)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/comment", bytes.NewBuffer(data))
	ginContext.Params = []gin.Param{{Key: "commentId", Value: strconv.FormatUint(expectedCommentId, 10)}}
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`
	timeNowString := time.Now().UTC().Format(model.TimeLayout)
	timeNow, _ := time.Parse(model.TimeLayout, timeNowString)
	timeService.EXPECT().GetTimeNowUtc().Return(timeNow)
	commentId := integration_test_arrange.GetNextCommentId()
	expectedCommentWasUpdatedEvent := &event.CommentWasUpdatedEvent{
		CommentId: commentId,
		Username:  comment.Username,
		PostId:    comment.PostId,
		Content:   comment.Content,
		UpdatedAt: timeNow.Format(model.TimeLayout),
	}
	expectedEvent := integration_test_builder.NewEventBuilder(t).WithName(event.CommentWasUpdatedEventName).WithData(expectedCommentWasUpdatedEvent).Build()
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(nil)

	controller.UpdateComment(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
	integration_test_assert.AssertCommentExists(t, db, commentId, comment)
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
