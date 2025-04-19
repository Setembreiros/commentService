package delete_comment

import (
	"commentservice/internal/api"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=test/mock/controller.go

type DeleteCommentController struct {
	service Service
}

type Service interface {
	RemoveComment(commentId uint64) error
}

func NewDeleteCommentController(service Service) *DeleteCommentController {
	return &DeleteCommentController{
		service: service,
	}
}

func (controller *DeleteCommentController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.DELETE("/comment", controller.DeleteComment)
}

func (controller *DeleteCommentController) DeleteComment(c *gin.Context) {
	log.Info().Msg("Handling Request DELETE DeleteComment")

	commentId := c.Param("commentId")
	if commentId == "" {
		api.SendBadRequest(c, "Missing commentId parameter")
		return
	}

	id, err := strconv.ParseUint(commentId, 10, 64)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("CommentId %s couldn't be parsed", commentId)
		api.SendBadRequest(c, "CommentId couldn't be parsed. CommentId hould be a positive number")
		return
	}

	err = controller.service.RemoveComment(id)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOK(c)
}
