package handler

import (
	"fmt"
	"io"
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

func (h *Handler) test(c *gin.Context) {

	const op = "handler.service.test"

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, op+".ReadBody", err)
		return
	}

	fmt.Printf("%s\n", string(body))

	if h.config.Log.OutputInFile {
		h.logger.Info(string(body))
	}

	queryParams := c.Request.URL.Query()
	for key, values := range queryParams {
		for _, value := range values {
			h.logger.Info("Query parameters", h.logger.Str(key, value))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}
