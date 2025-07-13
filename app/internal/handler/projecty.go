package handler

import (
	"catcher/app/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Send Event
// @Tags Event
// @Description Отправка Event в Sentry
// @ID sendEvent
// @Accept  json
// @Produce  json
// @Param input body models.Event true "Данные события"
// @Success 200 {object} models.SendEventResult
// @Failure default {object} errorResponse
// @Router /api/prj/:id/sendEvent [post]
func (h *Handler) sendEvent(c *gin.Context) {

	const op = "projecty.sendEvent"

	var input models.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, op+".ShouldBindJSON", err)
		return
	}

	projectId := c.Param("id")

	eventId, err := h.services.Projecty.SendEvent(projectId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, op, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"event_id": eventId})

}
