package handler

import (
	"catcher/app/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Start Session
// @Tags Session
// @Description Отправка начало сессии в Sentry
// @ID startSession
// @Accept  json
// @Produce  json
// @Param input body models.Session true "Данные сессии"
// @Success 200 {object} map[string]string
// @Failure default {object} errorResponse
// @Router /api/prj/:id/session/start [post]
func (h *Handler) startSession(c *gin.Context) {

	const op = "session.start"
	h.startEndSession(true, op, c)

}

// @Summary End Session
// @Tags Session
// @Description Отправка окончания сессии в Sentry
// @ID endSession
// @Accept  json
// @Produce  json
// @Param input body models.Session true "Данные сессии"
// @Success 200 {object} map[string]string
// @Failure default {object} errorResponse
// @Router /api/prj/:id/session/end [post]
func (h *Handler) endSession(c *gin.Context) {

	const op = "session.end"
	h.startEndSession(false, op, c)
}

func (h *Handler) startEndSession(start bool, op string, c *gin.Context) {

	var input models.Session
	var err error
	if err = c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, op+".ShouldBindJSON", err)
		return
	}

	projectId := c.Param("id")

	if start {
		err = h.services.Session.Start(projectId, input)
	} else {
		err = h.services.Session.End(projectId, input)
	}

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, op, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})

}
