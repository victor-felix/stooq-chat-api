package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victor-felix/chat-service/app/models"
)

func (h *handler) postRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		h.responseProblem(c, err)
		return
	}

	createdRoom, err := h.roomService.Create(c, &room)

	if err != nil {
		h.responseProblem(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdRoom)
}

func (h *handler) getRooms(c *gin.Context) {
	rooms, err := h.roomService.FindAll(c)

	if err != nil {
		h.responseProblem(c, err)
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *handler) getRoomMessages(c *gin.Context) {
	roomID := c.Param("id")

	messages, err := h.messageService.FindByRoomId(c, roomID)

	if err != nil {
		h.responseProblem(c, err)
		return
	}

	c.JSON(http.StatusOK, messages)
}
