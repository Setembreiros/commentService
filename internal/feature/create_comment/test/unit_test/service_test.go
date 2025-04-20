package unit_test_create_comment

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"commentservice/internal/bus"
	mock_bus "commentservice/internal/bus/test/mock"
	"commentservice/internal/feature/create_comment"
	mock_create_comment "commentservice/internal/feature/create_comment/test/mock"
	model "commentservice/internal/model/domain"
	"commentservice/internal/model/event"

	"github.com/stretchr/testify/assert"
)

var timeService *mock_create_comment.MockTimeService
var serviceRepository *mock_create_comment.MockRepository
var serviceExternalBus *mock_bus.MockExternalBus
var serviceBus *bus.EventBus
var createCommentService *create_comment.CreateCommentService

func setUpService(t *testing.T) {
	setUp(t)
	timeService = mock_create_comment.NewMockTimeService(ctrl)
	serviceRepository = mock_create_comment.NewMockRepository(ctrl)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus = bus.NewEventBus(serviceExternalBus)
	createCommentService = create_comment.NewCreateCommentService(timeService, serviceRepository, serviceBus)
}

func TestCreateCommentWithService_WhenItReturnsSuccess(t *testing.T) {
	setUpService(t)
	comment := &model.Comment{
		Username: "usernameA",
		PostId:   "post1",
		Content:  "o meu comentario",
	}
	timeNowString := time.Now().UTC().Format(model.TimeLayout)
	timeNow, _ := time.Parse(model.TimeLayout, timeNowString)
	comment.CreatedAt = timeNow
	expectedCommmentId := uint64(1000)
	expectedCommentWasCreatedEvent := &event.CommentWasCreatedEvent{
		CommentId: expectedCommmentId,
		Username:  comment.Username,
		PostId:    comment.PostId,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format(model.TimeLayout),
	}
	expectedEvent, _ := createEvent(event.CommentWasCreatedEventName, expectedCommentWasCreatedEvent)
	timeService.EXPECT().GetTimeNowUtc().Return(timeNow)
	serviceRepository.EXPECT().CreateComment(comment).Return(expectedCommmentId, nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(nil)

	err := createCommentService.CreateComment(comment)

	assert.Nil(t, err)
	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Comment was created, username: %s -> postId: %s", comment.Username, comment.PostId))
}

func TestErrorOnCreateCommentWithService_WhenAddingToRepositoryFails(t *testing.T) {
	setUpService(t)
	comment := &model.Comment{
		Username: "usernameA",
		PostId:   "post1",
		Content:  "o meu comentario",
	}
	timeNowString := time.Now().UTC().Format(model.TimeLayout)
	timeNow, _ := time.Parse(model.TimeLayout, timeNowString)
	comment.CreatedAt = timeNow
	timeService.EXPECT().GetTimeNowUtc().Return(timeNow)
	serviceRepository.EXPECT().CreateComment(comment).Return(uint64(0), errors.New("some error"))

	err := createCommentService.CreateComment(comment)

	assert.NotNil(t, err)
	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Error creating comment, username: %s -> postId: %s", comment.Username, comment.PostId))
}

func TestErrorOnCreateCommentWithService_WhenPublishingEventFails(t *testing.T) {
	setUpService(t)
	comment := &model.Comment{
		Username: "usernameA",
		PostId:   "post1",
		Content:  "o meu comentario",
	}
	timeNowString := time.Now().UTC().Format(model.TimeLayout)
	timeNow, _ := time.Parse(model.TimeLayout, timeNowString)
	comment.CreatedAt = timeNow
	expectedCommmentId := uint64(5)
	expectedCommentWasCreatedEvent := &event.CommentWasCreatedEvent{
		CommentId: expectedCommmentId,
		Username:  comment.Username,
		PostId:    comment.PostId,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format(model.TimeLayout),
	}
	expectedEvent, _ := createEvent(event.CommentWasCreatedEventName, expectedCommentWasCreatedEvent)
	timeService.EXPECT().GetTimeNowUtc().Return(timeNow)
	serviceRepository.EXPECT().CreateComment(comment).Return(expectedCommmentId, nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(errors.New("some error"))

	err := createCommentService.CreateComment(comment)

	assert.NotNil(t, err)
	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Publishing %s failed, username: %s -> postId: %s", event.CommentWasCreatedEventName, expectedCommentWasCreatedEvent.Username, expectedCommentWasCreatedEvent.PostId))
}
