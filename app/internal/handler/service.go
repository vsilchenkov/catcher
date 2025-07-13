package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) swagger(c *gin.Context) {

	c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")

}

// @Summary Clear cache
// @Tags Cache
// @Description Очищает все данные, сохранённые в кэше
// @ID clearCache
// @Produce  json
// @Success 200 {object} map[string]string
// @Failure default {object} errorResponse
// @Router /api/service/clearCache [get]
func (h *Handler) clearCache(c *gin.Context) {

	const op = "handler.service.clearCache"

	if err := h.services.Service.ClearCache(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, op, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}
