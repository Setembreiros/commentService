package unit_test_update_comment

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"commentservice/internal/bus"
	mock_bus "commentservice/internal/bus/test/mock"
	"commentservice/internal/feature/update_comment"
	mock_update_comment "commentservice/internal/feature/update_comment/test/mock"
	model "commentservice/internal/model/domain"
	"commentservice/internal/model/event"

	"github.com/stretchr/testify/assert"
)

var timeService *mock_update_comment.MockTimeService
var serviceRepository *mock_update_comment.MockRepository
var serviceExternalBus *mock_bus.MockExternalBus
var serviceBus *bus.EventBus
var updateCommentService *update_comment.UpdateCommentService

func setUpService(t *testing.T) {
	setUp(t)
	timeService = mock_update_comment.NewMockTimeService(ctrl)
	serviceRepository = mock_update_comment.NewMockRepository(ctrl)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus = bus.NewEventBus(serviceExternalBus)
	updateCommentService = update_comment.NewUpdateCommentService(timeService, serviceRepository, serviceBus)
}

func TestUpdateCommentWithService_WhenItReturnsSuccess(t *testing.T) {
	setUpService(t)
	comment := &model.Comment{
		Id:      uint64(1234),
		Content: "o meu comentario",
	}
	timeNowString := time.Now().UTC().Format(model.TimeLayout)
	timeNow, _ := time.Parse(model.TimeLayout, timeNowString)
	comment.UpdatedAt = timeNow
	expectedCommentWasUpdatedEvent := &event.CommentWasUpdatedEvent{
		CommentId: comment.Id,
		Content:   comment.Content,
		UpdatedAt: comment.UpdatedAt.Format(model.TimeLayout),
	}
	expectedEvent, _ := createEvent(event.CommentWasUpdatedEventName, expectedCommentWasUpdatedEvent)
	timeService.EXPECT().GetTimeNowUtc().Return(timeNow)
	serviceRepository.EXPECT().UpdateComment(comment).Return(nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(nil)

	err := updateCommentService.UpdateComment(comment)

	assert.Nil(t, err)
	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Comment was updated, commentId: %d", comment.Id))
}

func TestErrorOnUpdateCommentWithService_WhenUpdatingInRepositoryFails(t *testing.T) {
	setUpService(t)
	comment := &model.Comment{
		Id:      uint64(1234),
		Content: "o meu comentario",
	}
	timeNowString := time.Now().UTC().Format(model.TimeLayout)
	timeNow, _ := time.Parse(model.TimeLayout, timeNowString)
	comment.UpdatedAt = timeNow
	timeService.EXPECT().GetTimeNowUtc().Return(timeNow)
	serviceRepository.EXPECT().UpdateComment(comment).Return(errors.New("some error"))

	err := updateCommentService.UpdateComment(comment)

	assert.NotNil(t, err)
	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Error updating comment, commentId: %d", comment.Id))
}

func TestErrorOnUpdateCommentWithService_WhenPublishingEventFails(t *testing.T) {
	setUpService(t)
	comment := &model.Comment{
		Id:      uint64(1234),
		Content: "o meu comentario",
	}
	timeNowString := time.Now().UTC().Format(model.TimeLayout)
	timeNow, _ := time.Parse(model.TimeLayout, timeNowString)
	comment.UpdatedAt = timeNow
	expectedCommentWasUpdatedEvent := &event.CommentWasUpdatedEvent{
		CommentId: comment.Id,
		Content:   comment.Content,
		UpdatedAt: comment.UpdatedAt.Format(model.TimeLayout),
	}
	expectedEvent, _ := createEvent(event.CommentWasUpdatedEventName, expectedCommentWasUpdatedEvent)
	timeService.EXPECT().GetTimeNowUtc().Return(timeNow)
	serviceRepository.EXPECT().UpdateComment(comment).Return(nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(errors.New("some error"))

	err := updateCommentService.UpdateComment(comment)

	assert.NotNil(t, err)
	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Publishing %s failed, commentId: %d", event.CommentWasUpdatedEventName, expectedCommentWasUpdatedEvent.CommentId))
}
