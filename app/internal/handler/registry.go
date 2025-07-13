package handler

import (
	"io"
	"net/http"

	"catcher/app/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Get Info
// @Tags info
// @Description Проверка работы метода getInfo
// @ID getInfo
// @Produce  json
// @Success 200 {object} models.RegistryInfo
// @Failure default {object} errorResponse
// @Router /api/reg [get]
func (h *Handler) getInfo(c *gin.Context) {

	input := models.NewRegistryInput()

	info := h.services.Registry.GetInfo(input)
	c.JSON(http.StatusOK, info)

}

// @Summary Get Info Post
// @Tags Info
// @Description Получение информации для отчета об ошибки
// @ID getInfoPost
// @Accept  json
// @Produce  json
// @Param input body models.RegistryInput true "Значения для отчета об ошибке"
// @Success 200 {object} models.RegistryInfo
// @Failure default {object} errorResponse
// @Router /api/reg/getInfo [post]
func (h *Handler) getInfoPost(c *gin.Context) {

	const op = "registry.getInfoPost"

	var input models.RegistryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, op, err)
		return
	}

	info := h.services.Registry.GetInfo(input)
	c.JSON(http.StatusOK, info)

}

// @Summary Push Report
// @Tags Report
// @Description Отправка отчета об ошибки
// @ID pushReport
// @Accept multipart/form-data
// @Produce  json
// @Param file formData file true "Файл в архиве формата https://its.1c.ru/db/v8327doc#bookmark:dev:TI000002558"
// @Success 200 {object} models.RegistryPushReportResult
// @Failure default {object} errorResponse
// @Router /api/reg/pushReport [post]
func (h *Handler) pushReport(c *gin.Context) {

	const op = "registry.pushReport"

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, op+".io.ReadAll", err)
		return
	}

	id := uuid.New().String()

	input := models.RegistryPushReportInput{
		ID:   id,
		Data: data,
	}

	result, err := h.services.Registry.PushReport(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, op, err)
		return
	}

	c.JSON(http.StatusOK, result)

}
