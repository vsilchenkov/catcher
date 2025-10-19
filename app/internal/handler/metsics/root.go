package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) root(c *gin.Context) {

	err := h.tmpl.ExecuteTemplate(c.Writer, "index.html", nil) 
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template")
		return
	}
}
