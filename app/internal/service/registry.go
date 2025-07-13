package service

import (
	"catcher/app/internal/config"
	"catcher/app/internal/lib/caching"
	"catcher/app/internal/lib/logging"
	"catcher/app/internal/models"
	"catcher/app/internal/service/redirect"
	"catcher/app/internal/service/replicate"
	"catcher/app/internal/service/sentry/reporting"
	"context"
	"errors"
)

var ErrBadRequest = errors.New("bad request")
var ErrBadProject = errors.New("no project setting")
var ErrBadConvert = errors.New("bad convert")

type ConvertReporter interface {
	ConvertReport(r redirect.Report) (*models.RepportData, error)
}

type RegistryService struct {
	ctx    context.Context
	config *config.Config
	cacher caching.Cacher
	logger logging.Logger
}

func NewRegistryService(appCtx models.AppContext) *RegistryService {
	return &RegistryService{
		ctx:    appCtx.Ctx,
		config: appCtx.Config,
		cacher: appCtx.Cacher,
		logger: appCtx.Logger}
}

func (r *RegistryService) GetInfo(input models.RegistryInput) models.RegistryInfo {

	var needSendReport bool

	if _, err := r.config.ProjectByName(input.ConfigName); err == nil {
		needSendReport = true
	} else {
		needSendReport = false
	}

	return models.RegistryInfo{
		NeedSendReport: needSendReport,
		UserMessage:    r.config.RegistryUserMessage(),
		DumpType:       r.config.RegistryDumpType(),
	}
}

func (r *RegistryService) PushReport(input models.RegistryPushReportInput) (*models.RegistryPushReportResult, error) {

	appCtx := models.NewAppContext(r.ctx, r.config, r.cacher, r.logger)
	
	report := redirect.NewReport(input.ID, input.Data, appCtx)

	repl := replicate.New(appCtx)
	rData, err := ConvertReport(repl, report)
	if err != nil {
		return nil, ErrBadConvert
	}

	svcReport := reporting.NewReport(input.ID, rData.Prj, *rData.Data, rData.Files, appCtx)
	eventID, err := redirect.Send(svcReport, r.config.Registry.Timeout)

	result := models.RegistryPushReportResult{
		ID:      input.ID,
		EventID: eventID,
	}

	return &result, err
}

func ConvertReport(cnv ConvertReporter, r redirect.Report) (*models.RepportData, error) {
	return cnv.ConvertReport(r)
}
