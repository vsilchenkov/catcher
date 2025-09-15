package service

import (
	"catcher/app/internal/models"
)

type Registry interface {
	GetInfo(input models.RegistryInput) models.RegistryInfo
	PushReport(input models.RegistryPushReportInput) (*models.RegistryPushReportResult, error)
}

type Projecty interface {
	SendEvent(projectId string, input models.Event) (*models.SendEventResult, error)
	ClearCache(projectId string) error
}

type Session interface {
	Start(projectId string, input models.Session) error
	End(projectId string, input models.Session) error
}

type Service interface {
	ClearCache() error
}

type Services struct {
	Registry
	Projecty
	Session
	Service
}

func New(appCtx models.AppContext) *Services {
	return &Services{
		Registry: NewRegistryService(appCtx),
		Projecty: NewProjectyService(appCtx),
		Session: NewSessionService(appCtx),
		Service:  NewServiceService(appCtx),
	}
}
