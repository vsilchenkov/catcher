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

// @Summary Сlear the project cache
// @Tags Cache
// @Description Очищает все данные проекта, сохранённые в кэше 
// @ID clearProjectCache
// @Produce  json
// @Success 200 {object} map[string]string
// @Failure default {object} errorResponse
// @Router /api/service/clearCache [get]
func (h *Handler) clearProjectCache(c *gin.Context) {

	const op = "handler.service.clearProjectCache"

	projectId := c.Param("id")

	if err := h.services.Projecty.ClearCache(projectId); err != nil {
		newErrorResponse(c, http.StatusBadRequest, op, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}
