package handler

import (
	"catcher/app/internal/config"
	"catcher/app/internal/lib/logging"
	"catcher/app/internal/models"
	"catcher/app/internal/service"
	_ "catcher/docs"
	"os"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Services
	config   *config.Config
	logger   logging.Logger
}

func NewHandler(service *service.Services, appCtx models.AppContext) *Handler {
	return &Handler{
		services: service,
		config:   appCtx.Config,
		logger:   appCtx.Logger}
}

func (h *Handler) InitHandler() *gin.Engine {

	debug := h.config.UseDebug()
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	if debug {
		router.Use(h.initCustumLogger())
	} else {
		if h.config.Interactive {
			router.Use(gin.Logger())
		}
	}

	if h.config.Sentry.Use {
		router.Use(h.initSentryGin())
	}

	router.GET("/", h.swagger)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	api.GET("/", h.swagger)
	{
		registry := api.Group("/reg")
		{
			registry.GET("/", h.getInfo)
			registry.POST("/getInfo", h.getInfoPost)
			registry.POST("/pushReport", h.pushReport)
		}
		prj := api.Group("/prj/:id")
		{
			prj.POST("/sendEvent", h.sendEvent)
		}

		service := api.Group("/service")
		{
			service.GET("/clearCache", h.clearCache)
		}

	}

	return router

}

func (h *Handler) initCustumLogger() gin.HandlerFunc {

	var out *os.File
	out = os.Stdout

	settngs := h.config
	if settngs.Log.OutputInFile {
		fileName := "api.log"
		file, err := logging.GetOutputLogFile(settngs.WorkingDir, settngs.Log.Dir, fileName)
		if err == nil {
			out = file
		} else {
			h.logger.Error("Не удалось открыть файл логов, используется стандартный stderr",
				h.logger.Str("name", fileName),
				h.logger.Err(err))
		}
	}

	custumlogger := gin.LoggerWithWriter(out)
	return custumlogger
}

func (h *Handler) initSentryGin() gin.HandlerFunc {

	if err := sentry.Init(logging.SentryClientOptions(h.config)); err != nil {
		h.logger.Error("Sentry initialization failed",
			h.logger.Err(err))
	}

	return sentrygin.New(sentrygin.Options{
		Repanic: true,
	})
}
