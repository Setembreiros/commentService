package update_comment

import (
	"commentservice/internal/api"
	model "commentservice/internal/model/domain"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=test/mock/controller.go

type UpdateCommentController struct {
	service Service
}

type Service interface {
	UpdateComment(comment *model.Comment) error
}

func NewUpdateCommentController(service Service) *UpdateCommentController {
	return &UpdateCommentController{
		service: service,
	}
}

func (controller *UpdateCommentController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.PUT("/comment/:commentId", controller.UpdateComment)
}

func (controller *UpdateCommentController) UpdateComment(c *gin.Context) {
	log.Info().Msg("Handling Request PUT UpdateComment")

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

	var comment model.Comment
	if err := c.BindJSON(&comment); err != nil {
		log.Error().Stack().Err(err).Msg("Invalid Data")
		api.SendBadRequest(c, "Invalid Json Request")
		return
	}
	comment.Id = id

	err = controller.service.UpdateComment(&comment)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOK(c)
}
