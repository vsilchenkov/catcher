package main

import (
	"catcher/app/build"
	"catcher/app/internal/config"
	handler "catcher/app/internal/handler/metsics"
	"catcher/app/internal/models"
	"catcher/app/internal/server"
	"catcher/app/internal/server/http"
	"catcher/app/internal/service/metrics"
	"catcher/pkg/errors"
	"catcher/pkg/logging"
	"context"
	"embed"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jinzhu/copier"
	svc "github.com/kardianos/service"
)

const projectName = "Sentry exporter"

//go:embed versioninfo.json
var versionInfoData []byte

//go:embed templates/*
var content embed.FS

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
			fmt.Printf("init sentry error: %v\n", err)
			os.Exit(1)
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

	tmpl := template.Must(template.ParseFS(content, "templates/*.html"))
	appCtx := models.NewAppContext(ctx, c, nil, logger, tmpl)

	// Сборщик метрик
	configMetrics := &metrics.Config{
		Interval: c.SentryMetrics.Interval}
	collector := metrics.New(configMetrics, appCtx)
	go collector.Start()

	// Сервер
	srv := http.New(appCtx)

	handler := handler.New(appCtx).Init()
	prg := server.NewProgram(srv, handler, c.SentryMetrics.Port, appCtx)

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
