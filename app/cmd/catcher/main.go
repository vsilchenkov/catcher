package main

import (
	"catcher/app/build"
	"catcher/app/internal/config"
	"catcher/app/internal/handler"
	"catcher/app/internal/lib/caching"
	"catcher/app/internal/lib/caching/memory"
	"catcher/app/internal/lib/caching/redis"
	"catcher/app/internal/models"
	"catcher/app/internal/server"
	"catcher/app/internal/server/http"
	"catcher/app/internal/service"
	"catcher/pkg/errors"
	"catcher/pkg/logging"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jinzhu/copier"
	svc "github.com/kardianos/service"
)

const projectName = "Cather"

//go:embed versioninfo.json
var versionInfoData []byte

// @title Catcher
// @version 1.0
// @description Catcher API Service

// @host localhost:8000
// @BasePath /api
func main() {

	ctx := context.Background()

	svcConfig := &svc.Config{
		Name:        projectName + "Service",
		DisplayName: projectName + " API Service",
		Description: "API-приложение " + projectName,
	}

	b, err := build.NewOption(versionInfoData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := config.New(*b)
	flags := config.ParseFlags()
	err = config.LoadSettigs(flags, c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sentryConfig := &logging.SentryConfig{}
	copier.Copy(sentryConfig, c.Option)
	copier.Copy(sentryConfig, c.Sentry)

	if sentryConfig.Use {
		err = sentry.Init(logging.SentryClientOptions(sentryConfig))
		if err != nil {
			log.Fatal(err)
		}
		defer sentry.Flush(2 * time.Second)
	}

	logConfig := &logging.Config{}
	copier.Copy(logConfig, c.Option)
	copier.Copy(logConfig, c.Log)
	logger := logging.Initlogger(logConfig, sentryConfig)
	defer func() {
		if err := errors.PanicRecovered(recover()); err != nil {
			logger.Error("Panic recovered",
				logger.Err(err),
			)
			os.Exit(1)
		}
	}()

	var svcCacher caching.Cacher
	useMemory := true
	if c.Redis.Use {
		cacher, err := redis.New(&redis.Option{
			Addr:     c.Redis.Addr,
			Username: c.Redis.Credintials.UserName,
			Password: c.Redis.Credintials.Password,
			DB:       c.Redis.DB,
		})
		if err == nil {
			svcCacher = caching.New(cacher)
			useMemory = false
		} else {
			logger.Error("Failed to initialize Redis cacher",
				logger.Err(err))
		}
	}

	if useMemory {
		cacher := memory.New()
		svcCacher = caching.New(cacher)
	}

	appCtx := models.NewAppContext(ctx, c, svcCacher, logger)
	srv := http.New(appCtx)

	service := service.New(appCtx)
	handler := handler.New(service, appCtx).Init()
	prg := server.NewProgram(srv, handler, appCtx)

	s, err := svc.New(prg, svcConfig)
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
