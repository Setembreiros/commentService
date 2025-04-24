package unit_test_delete_comment

import (
	"errors"
	"fmt"
	"testing"

	"commentservice/internal/bus"
	mock_bus "commentservice/internal/bus/test/mock"
	"commentservice/internal/feature/delete_comment"
	mock_delete_comment "commentservice/internal/feature/delete_comment/test/mock"
	"commentservice/internal/model/event"

	"github.com/stretchr/testify/assert"
)

var serviceRepository *mock_delete_comment.MockRepository
var serviceExternalBus *mock_bus.MockExternalBus
var serviceBus *bus.EventBus
var deleteCommentService *delete_comment.DeleteCommentService

func setUpService(t *testing.T) {
	setUp(t)
	serviceRepository = mock_delete_comment.NewMockRepository(ctrl)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus = bus.NewEventBus(serviceExternalBus)
	deleteCommentService = delete_comment.NewDeleteCommentService(serviceRepository, serviceBus)
}

func TestDeleteCommentWithService_WhenItReturnsSuccess(t *testing.T) {
	setUpService(t)
	expectedCommmentId := uint64(1000)
	expectedCommentWasDeletedEvent := &event.CommentWasDeletedEvent{
		CommentId: expectedCommmentId,
	}
	expectedEvent, _ := createEvent(event.CommentWasDeletedEventName, expectedCommentWasDeletedEvent)
	serviceRepository.EXPECT().DeleteComment(expectedCommmentId).Return(nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(nil)

	err := deleteCommentService.DeleteComment(expectedCommmentId)

	assert.Nil(t, err)
	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Comment was deleted, commentId: %d", expectedCommmentId))
}

func TestErrorOnDeleteCommentWithService_WhenAddingToRepositoryFails(t *testing.T) {
	setUpService(t)
	expectedCommmentId := uint64(1000)
	serviceRepository.EXPECT().DeleteComment(expectedCommmentId).Return(errors.New("some error"))

	err := deleteCommentService.DeleteComment(expectedCommmentId)

	assert.NotNil(t, err)
	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Error deleting comment, commentId: %d", expectedCommmentId))
}

func TestErrorOnDeleteCommentWithService_WhenPublishingEventFails(t *testing.T) {
	setUpService(t)
	expectedCommmentId := uint64(1000)
	expectedCommentWasDeletedEvent := &event.CommentWasDeletedEvent{
		CommentId: expectedCommmentId,
	}
	expectedEvent, _ := createEvent(event.CommentWasDeletedEventName, expectedCommentWasDeletedEvent)
	serviceRepository.EXPECT().DeleteComment(expectedCommmentId).Return(nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(errors.New("some error"))

	err := deleteCommentService.DeleteComment(expectedCommmentId)

	assert.NotNil(t, err)
	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Publishing %s failed, commentId: %d", event.CommentWasDeletedEventName, expectedCommmentId))
}
