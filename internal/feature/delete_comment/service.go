package delete_comment

import (
	"commentservice/internal/bus"
	"commentservice/internal/model/event"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	DeleteComment(commentId uint64) error
}

type DeleteCommentService struct {
	repository Repository
	bus        *bus.EventBus
}

func NewDeleteCommentService(repository Repository, bus *bus.EventBus) *DeleteCommentService {
	return &DeleteCommentService{
		repository: repository,
		bus:        bus,
	}
}

func (s *DeleteCommentService) DeleteComment(postId string, commentId uint64) error {
	err := s.repository.DeleteComment(commentId)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error deleting comment, commentId: %d", commentId)
		return err
	}

	err = s.publishCommentWasDeletedEvent(postId, commentId)
	if err != nil {
		return err
	}

	log.Info().Msgf("Comment was deleted, commentId: %d", commentId)

	return nil
}

func (s *DeleteCommentService) publishCommentWasDeletedEvent(postId string, commentId uint64) error {
	commentWasDeletedEvent := &event.CommentWasDeletedEvent{
		PostId:    postId,
		CommentId: commentId,
	}

	err := s.bus.Publish(event.CommentWasDeletedEventName, commentWasDeletedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Publishing %s failed, commentId: %d", event.CommentWasDeletedEventName, commentId)
		return err
	}

	return nil
}
