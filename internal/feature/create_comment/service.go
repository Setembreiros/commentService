package create_comment

import (
	"commentservice/internal/bus"
	model "commentservice/internal/model/domain"
	"commentservice/internal/model/event"
	"time"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	CreateComment(data *model.Comment) (uint64, error)
}

type TimeService interface {
	GetTimeNowUtc() time.Time
}

type CreateCommentService struct {
	timeService TimeService
	repository  Repository
	bus         *bus.EventBus
}

func NewCreateCommentService(timeService TimeService, repository Repository, bus *bus.EventBus) *CreateCommentService {
	return &CreateCommentService{
		timeService: timeService,
		repository:  repository,
		bus:         bus,
	}
}

func (s *CreateCommentService) CreateComment(comment *model.Comment) error {
	comment.CreatedAt = s.timeService.GetTimeNowUtc()
	var err error
	comment.Id, err = s.repository.CreateComment(comment)

	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error adding comment, username: %s -> postId: %s", comment.Username, comment.PostId)
		return err
	}

	err = s.publishCommentWasCreatedEvent(comment)
	if err != nil {
		return err
	}

	log.Info().Msgf("Comment was created, username: %s -> postId: %s", comment.Username, comment.PostId)

	return nil
}

func (s *CreateCommentService) publishCommentWasCreatedEvent(data *model.Comment) error {
	commentWasCreatedEvent := &event.CommentWasCreatedEvent{
		CommentId: data.Id,
		Username:  data.Username,
		PostId:    data.PostId,
		Content:   data.Content,
		CreatedAt: data.CreatedAt.Format(model.TimeLayout),
	}

	err := s.bus.Publish(event.CommentWasCreatedEventName, commentWasCreatedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Publishing %s failed, username: %s -> postId: %s", event.CommentWasCreatedEventName, commentWasCreatedEvent.Username, commentWasCreatedEvent.PostId)
		return err
	}

	return nil
}
