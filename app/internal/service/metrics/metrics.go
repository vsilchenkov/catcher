package metrics

import (
	"catcher/app/internal/config"
	"catcher/app/internal/models"
	sentryclient "catcher/app/internal/sentry/client"
	"catcher/app/internal/service/metrics/session"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/jinzhu/copier"
)

const MetricUnknown = "неизвестная метрика"

type Metric interface {
	Run() error
}

// Карта фабрик (регистрируем статично)
var metricFactories = map[string]func(config config.Metrics, client *sentryclient.Client, appCtx models.AppContext) Metric{
	"Sessions": func(config config.Metrics, client *sentryclient.Client, appCtx models.AppContext) Metric {
		return session.New(config, client, appCtx)
	},
}

type Config struct {
	Interval int
}
type Metrics struct {
	*Config
	models.AppContext
	SentryClient *sentryclient.Client
}

func New(config *Config, appCtx models.AppContext) Metrics {

	configSentry := sentryclient.Config{}
	copier.Copy(&configSentry, appCtx.Config.SentryMetrics.Sentry)
	SentryClient := sentryclient.New(configSentry)

	return Metrics{
		Config:       config,
		AppContext:   appCtx,
		SentryClient: SentryClient,
	}
}

func (m Metrics) Start() error {

	const op = "service.metrics.Start"

	errCh := make(chan error, 1)

	for {
		go func() {
			if err := m.fetch(); err != nil {
				errCh <- err
			}
		}()

		select {
		case err := <-errCh:
			m.Logger.Error("Ошибка запуска метрик",
				m.Logger.Err(err),
				m.Logger.Op(op))
		default:
		}

		time.Sleep(time.Duration(m.Interval) * time.Second)
	}

}

func (m Metrics) fetch() error {

	configMetrics := m.AppContext.Config.SentryMetrics

	type res struct {
		err error
		op  string
	}

	resCh := make(chan res)
	var wg sync.WaitGroup

	for _, config := range configMetrics.Metrics {
		if !config.Use {
			continue
		}

		name := config.Name
		op := "metrics.fetch." + name

		factory, ok := metricFactories[name]
		if !ok {
			m.handleError(errors.Errorf(MetricUnknown+" %s", name), op)
			continue
		}

		instance := factory(config, m.SentryClient, m.AppContext) // Вызов New через фабрику

		wg.Go(func() {
			if err := instance.Run(); err != nil {
				resCh <- res{err: err, op: op}
			}
		})
	}

	go func() {
		for res := range resCh {
			m.handleError(res.err, res.op)
		}
	}()

	wg.Wait()
	close(resCh)

	return nil
}

func (m Metrics) handleError(err error, op string) {

	if err != nil {
		m.Logger.Info("Ошибка получения метрики",
			m.Logger.Err(err),
			m.Logger.Op(op))
	}
}
