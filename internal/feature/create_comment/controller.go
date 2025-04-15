package create_comment

import (
	"commentservice/internal/api"
	model "commentservice/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=test/mock/controller.go

type CreateCommentController struct {
	service Service
}

type Service interface {
	AddComment(comment *model.Comment) error
}

func NewCreateCommentController(service Service) *CreateCommentController {
	return &CreateCommentController{
		service: service,
	}
}

func (controller *CreateCommentController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/comment", controller.CreateComment)
}

func (controller *CreateCommentController) CreateComment(c *gin.Context) {
	log.Info().Msg("Handling Request POST CreateComment")
	var comment model.Comment

	if err := c.BindJSON(&comment); err != nil {
		log.Error().Stack().Err(err).Msg("Invalid Data")
		api.SendBadRequest(c, "Invalid Json Request")
		return
	}

	err := controller.service.AddComment(&comment)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOK(c)
}
