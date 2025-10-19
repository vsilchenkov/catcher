package session

import (
	"catcher/app/internal/config"
	"catcher/app/internal/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const query = "sessions/?field=sum(session)" +
	"&groupBy=project&groupBy=environment&groupBy=release&groupBy=session.status" +
	"&interval=1d" +
	"&statsPeriod=24h"

type Metric struct {
	Client
	models.AppContext
	config config.Metrics
}

type response struct {
	Groups []struct {
		By struct {
			Project       int    `json:"project"`
			Environment   string `json:"environment"`
			Release       string `json:"release"`
			SessionStatus string `json:"session.status"`
		} `json:"by"`
		Totals map[string]float64 `json:"totals"`
	} `json:"groups"`
}

type Client interface {
	Request(endpoint string) ([]byte, error)
}

var (
	// Метрика для сессий по проектам и статусам
	sentrySessionTotalName = "sentry_session_total"
	// sentrySessionTotal     = promauto.NewSummaryVec(prometheus.SummaryOpts{
	// 	Name: sentrySessionTotalName,
	// 	Help: "Total number of sessions from Sentry per project and status",
	// 	Objectives: map[float64]float64{
	// 		0.5:  0.05,
	// 		0.9:  0.01,
	// 		0.95: 0.01,
	// 		0.99: 0.005,
	// 	},
	// }, []string{"project", "environment", "release", "session_status"})

	sentrySessionTotal     = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: sentrySessionTotalName,
		Help: "Total number of sessions from Sentry per project and status",		
	}, []string{"project", "environment", "release", "session_status"})

	// Общая метрика для суммы по всем проектам и статусам
	sentrySessionAllTotalName = "sentry_session_all_total"
	sentrySessionAllTotal     = promauto.NewGauge(prometheus.GaugeOpts{
		Name: sentrySessionAllTotalName,
		Help: "Total number of sessions from all Sentry projects and statuses",
	})
)

func New(config config.Metrics, client Client, appCtx models.AppContext) Metric {
	return Metric{
		config:     config,
		Client:     client,
		AppContext: appCtx,
	}
}

func (m Metric) Run() error {

	const op = "metric.session.run"

	m.Logger.Debug("Запускаю метрику", m.Logger.Str("name", m.config.Name))

	urlParts := []string{}
	for _, prj := range m.config.Projects {
		urlParts = append(urlParts, fmt.Sprintf("project=%s", prj.Id))
	}
	prjQuery := strings.Join(urlParts, "&")

	endPoint := query + "&" + prjQuery
	body, err := m.Client.Request(endPoint)
	if err != nil {
		return errors.WithMessagef(err, "Ошибка получения метрики %s", op)
	}

	var data response
	if err := json.Unmarshal(body, &data); err != nil {
		return errors.WithMessagef(err, "Ошибка unmarshal метрики %s", op)
	}

	var total float64
	var mu sync.Mutex

	for _, group := range data.Groups {

		if val, ok := group.Totals["sum(session)"]; ok {

			projectID := strconv.Itoa(group.By.Project)
			environment := group.By.Environment
			release := group.By.Release
			status := group.By.SessionStatus
			// sentrySessionTotal.WithLabelValues(projectID, environment, release, status).Observe(val)
			sentrySessionTotal.WithLabelValues(projectID, environment, release, status).Set(val)
			mu.Lock()
			total += val
			mu.Unlock()
			m.Logger.Debug("Обновление метрики",
				m.Logger.Str("name", sentrySessionTotalName),
				m.Logger.Str("project", projectID),
				m.Logger.Str("environment", environment),
				m.Logger.Str("release", release),
				m.Logger.Str("status", status),
				m.Logger.Any("val", val),
				m.Logger.Op(op))
		}
	}

	sentrySessionAllTotal.Set(total)
	m.Logger.Debug("Обновление метрики",
		m.Logger.Str("name", sentrySessionAllTotalName),
		m.Logger.Any("val", total),
		m.Logger.Op(op))

	return nil
}
