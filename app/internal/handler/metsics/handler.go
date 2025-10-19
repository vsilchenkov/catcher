package handler

import (
	"catcher/app/internal/config"
	"catcher/app/internal/models"
	"catcher/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	config *config.Config
	logger logging.Logger
	tmpl   models.Template
}

func New(appCtx models.AppContext) *Handler {
	return &Handler{
		config: appCtx.Config,
		logger: appCtx.Logger,
		tmpl:   appCtx.Tmpl}
}

func (h *Handler) Init() *gin.Engine {

	debug := h.config.UseDebug()
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	if h.config.Interactive {
		router.Use(gin.Logger())
	}

	router.GET("/", h.root)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router

}
