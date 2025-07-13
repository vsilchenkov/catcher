package service

import (
	"catcher/app/internal/config"
	"catcher/app/internal/lib/caching"
	"catcher/app/internal/lib/logging"
	"catcher/app/internal/models"
	"catcher/app/internal/service/redirect"
	"catcher/app/internal/service/replicate"
	"catcher/app/internal/service/sentry/sending"
	"context"

	"github.com/google/uuid"
)

type ProjectyService struct {
	ctx    context.Context
	config *config.Config
	cacher caching.Cacher
	logger logging.Logger
}

func NewProjectyService(appCtx models.AppContext) *ProjectyService {
	return &ProjectyService{
		ctx:    appCtx.Ctx,
		config: appCtx.Config,
		cacher: appCtx.Cacher,
		logger: appCtx.Logger}
}

func (p ProjectyService) SendEvent(projectId string, input models.Event) (*models.SendEventResult, error) {

	// op := "ProjectyService.SendEvent"
	prj, err := p.config.ProjectById(projectId)
	if err != nil {
		return nil, ErrBadProject
	}

	appCtx := models.NewAppContext(p.ctx, p.config, p.cacher, p.logger)

	repl := replicate.New(appCtx)
	var svc replicate.ConvertEventer = repl
	event, err := svc.ConvertEvent(prj, input)
	if err != nil {
		return nil, ErrBadConvert
	}

	id := string(event.EventID)
	if  id == "" {
		id = uuid.New().String()
	}
	svcEvent := sending.NewEvent(prj, id, event, appCtx)
	eventID, err := redirect.Send(svcEvent, p.config.Registry.Timeout)

	result := models.SendEventResult{
		ID:      id,
		EventID: eventID,
	}

	return &result, err

}
