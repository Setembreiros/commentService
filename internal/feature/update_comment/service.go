package update_comment

import (
	"commentservice/internal/bus"
	model "commentservice/internal/model/domain"
	"commentservice/internal/model/event"
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	UpdateComment(data *model.Comment) error
}

type TimeService interface {
	GetTimeNowUtc() time.Time
}

type UpdateCommentService struct {
	timeService TimeService
	repository  Repository
	bus         *bus.EventBus
}

func NewUpdateCommentService(timeService TimeService, repository Repository, bus *bus.EventBus) *UpdateCommentService {
	return &UpdateCommentService{
		timeService: timeService,
		repository:  repository,
		bus:         bus,
	}
}

func (s *UpdateCommentService) UpdateComment(comment *model.Comment) error {
	comment.UpdatedAt = s.timeService.GetTimeNowUtc()
	err := s.repository.UpdateComment(comment)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error updating comment, commentId: %d", comment.Id)
		return err
	}

	err = s.publishCommentWasUpdatedEvent(comment)
	if err != nil {
		return err
	}

	log.Info().Msgf("Comment was updated, commentId: %d", comment.Id)

	return nil
}

func (s *UpdateCommentService) publishCommentWasUpdatedEvent(data *model.Comment) error {
	commentWasUpdatedEvent := &event.CommentWasUpdatedEvent{
		CommentId: data.Id,
		Content:   data.Content,
		UpdatedAt: data.UpdatedAt.Format(model.TimeLayout),
	}

	err := s.bus.Publish(event.CommentWasUpdatedEventName, commentWasUpdatedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Publishing %s failed, commentId: %d", event.CommentWasUpdatedEventName, commentWasUpdatedEvent.CommentId)
		return err
	}

	return nil
}
