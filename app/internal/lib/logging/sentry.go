package logging

import (
	"catcher/app/internal/config"
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
)

func SentryHandler() slog.Handler {

	return sentryslog.Option{
		Level:     slog.LevelWarn,
		AddSource: true,
	}.NewSentryHandler()
}

func SentryClientOptions(c *config.Config) sentry.ClientOptions {

	var environment string
	if c.UseDebug() {
		environment = "Debug"
	} else {
		environment = "Production"
	}

	sentrySyncTransport := sentry.NewHTTPSyncTransport()
	sentrySyncTransport.Timeout = time.Second * 3

	return sentry.ClientOptions{
		Transport:        sentrySyncTransport,
		Dsn:              c.Sentry.Dsn,
		Release:          c.ProjectName + "@" + c.Version,
		Environment:      environment,
		Debug:            c.UseDebug(),
		AttachStacktrace: c.Sentry.AttachStacktrace,
		TracesSampleRate: c.Sentry.TracesSampleRate,
		EnableTracing:    c.Sentry.EnableTracing,
	}
}
