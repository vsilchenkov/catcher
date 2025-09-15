package logging

import (
	"context"
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
)

func SentryHandler() slog.Handler {

	return sentryslog.Option{
		Level:     slog.LevelWarn,
		AddSource: true,
	}.NewSentryHandler(context.Background())

}

func SentryClientOptions(s *SentryConfig) sentry.ClientOptions {

	sentrySyncTransport := sentry.NewHTTPSyncTransport()
	sentrySyncTransport.Timeout = time.Second * 3

	return sentry.ClientOptions{
		Transport:        sentrySyncTransport,
		Dsn:              s.Dsn,
		Release:          s.ProjectName + "@" + s.Version,
		Environment:      s.Environment,
		Debug:            s.UseDebug(),
		AttachStacktrace: s.AttachStacktrace,
		TracesSampleRate: s.TracesSampleRate,
		EnableTracing:    s.EnableTracing,
	}
}
