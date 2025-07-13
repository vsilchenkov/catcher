package sentryhub

import (
	"catcher/app/internal/config"
	"catcher/app/internal/models"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
)

type Hub struct {
	*sentry.Hub
	appCtx models.AppContext
}

type project struct {
	name string
	id   string
}

type hubsmap map[project]Hub

var hubs hubsmap
var hubsMu sync.RWMutex

func init() {
	hubs = make(hubsmap)
}

func New(prj config.Project, appCtx models.AppContext) (*Hub, error) {

	c, err := sentry.NewClient(clientOptions(prj, appCtx))
	if err != nil {
		return nil, err
	}

	h := sentry.NewHub(c, sentry.NewScope())
	hub := Hub{h, appCtx}

	hubsMu.Lock()
	hubs[toProject(prj)] = hub
	hubsMu.Unlock()

	return &hub, nil

}

func Get(prj config.Project, appCtx models.AppContext) (*Hub, error) {

	hubsMu.RLock()
	hub, ok := hubs[toProject(prj)]
	hubsMu.RUnlock()

	if ok {
		return &hub, nil
	}

	return New(prj, appCtx)
}

func toProject(prj config.Project) project {
	return project{
		name: prj.Name,
		id:   prj.Name,
	}
}

func clientOptions(prj config.Project, appCtx models.AppContext) sentry.ClientOptions {

	sentrySyncTransport := sentry.NewHTTPSyncTransport()
	sentrySyncTransport.Timeout = time.Second * 3

	return sentry.ClientOptions{
		Transport:        sentrySyncTransport,
		Dsn:              prj.Sentry.Dsn,
		Release:          prj.Release(),
		Environment:      prj.Sentry.Environment,
		Debug:            appCtx.Config.UseDebug(),
		AttachStacktrace: false,
		EnableTracing:    false,
		Integrations: func([]sentry.Integration) []sentry.Integration {
			return make([]sentry.Integration, 0)
		},
	}
}
