package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victor-felix/chat-service/app/models"
)

func (h *handler) postUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.responseProblem(c, err)
		return
	}

	createdUser, err := h.userService.Create(c, &user)

	if err != nil {
		h.responseProblem(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}