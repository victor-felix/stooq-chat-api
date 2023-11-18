package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victor-felix/chat-service/app/models"
)

func (h *handler) postAuth(c *gin.Context) {
	var payload models.AuthPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.responseProblem(c, err)
		return
	}

	token, err := h.authService.Login(c, &payload)

	if err != nil {
		h.responseProblem(c, err)
		return
	}

	c.JSON(http.StatusOK, token)
}