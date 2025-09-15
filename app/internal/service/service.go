package service

import (
	"catcher/app/internal/lib/caching"
	"catcher/pkg/logging"
	"catcher/app/internal/models"
	"context"
)

type ServiceService struct {
	ctx    context.Context
	cacher caching.Cacher
	logger logging.Logger
}

func NewServiceService(appCtx models.AppContext) *ServiceService {
	return &ServiceService{
		ctx:    appCtx.Ctx,
		cacher: appCtx.Cacher,
		logger: appCtx.Logger}
}

func (s ServiceService) ClearCache() error {

	const op = "service.ClearCache"

	err := s.cacher.Clear(s.ctx)
	if err != nil {
		return err
	}

	s.logger.Debug("Кэш очищен", s.logger.Op(op))
	return nil
}
