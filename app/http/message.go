package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/victor-felix/chat-service/app/errors"
	"github.com/victor-felix/chat-service/app/models"
)

func (h *handler) postMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		h.responseProblem(c, err)
		return
	}

	_, err := h.roomService.FindByID(c, message.RoomID)

	if err != nil {
		h.responseProblem(c, errors.NewErrNotFound(models.ErrorRoomNotFound).WithMessage("room not found"))
		return
	}

	jwtClaims, ok := c.Get(JwtClaimsAttribute)

	if !ok {
		h.responseProblem(c, errors.NewErrUnauthorized(models.ErrorAuthTokenMissing).WithMessage("auth token is missing"))
		c.Abort()
		return
	}

	message.UserID = jwtClaims.(jwt.MapClaims)["user_id"].(string)

	createdMessage, err := h.messageService.Create(c, &message)

	if err != nil {
		h.responseProblem(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdMessage)
}