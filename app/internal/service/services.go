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
}

type Service interface {
	ClearCache() error
}

type Services struct {
	Registry
	Projecty
	Service
}

func NewService(appCtx models.AppContext) *Services {
	return &Services{
		Registry: NewRegistryService(appCtx),
		Projecty: NewProjectyService(appCtx),
		Service: NewServiceService(appCtx),
	}
}


