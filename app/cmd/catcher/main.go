package main

import (
	"catcher/app/build"
	"catcher/app/internal/config"
	"catcher/app/internal/lib/caching"
	"catcher/app/internal/lib/caching/memory"
	"catcher/app/internal/lib/logging"
	"catcher/app/internal/models"
	"catcher/app/internal/server"
	"catcher/app/internal/server/http"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/kardianos/service"
	"github.com/cockroachdb/errors"
)

// @title Catcher
// @version 1.0
// @description Catcher API Service

// @host localhost:8000
// @BasePath /api
func main() {

	ctx := context.Background()
	name := build.ProjectName

	svcConfig := &service.Config{
		Name:        name + "Service",
		DisplayName: name + " API Service",
		Description: "API-приложение " + name,
	}

	flags := config.ParseFlags()
	c, err := config.LoadSettigs(flags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if c.Sentry.Use {
		err = sentry.Init(logging.SentryClientOptions(c))
		if err != nil {
			log.Fatal(err)
		}
		defer sentry.Flush(2 * time.Second)
	}

	logger := logging.Initlogger(c)
	defer func() {
		if r := recover(); r != nil {
			var err error
			switch e := r.(type) {
			case error:
				err = e
			case string:
				err = errors.New(e)
			default:
				err = fmt.Errorf("panic: %v", e)
			}

			err = errors.WithStackDepth(err,2)
			logger.Error("Panic recovered",
				logger.Err(err),
			)
			os.Exit(1)
		}
	}()

	cacher := memory.New()
	svcCacher := caching.New(cacher)

	appCtx := models.NewAppContext(ctx, c, svcCacher, logger)
	srv := http.New(appCtx)
	prg := server.NewProgram(srv, appCtx)

	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Error("Error on service start",
			logger.Err(err))
	}

	err = s.Run()
	if err != nil {
		logger.Error("Ошибка запуска сервера",
			logger.Err(err))
	}

}
