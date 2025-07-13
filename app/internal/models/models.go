package models

import (
	"catcher/app/internal/config"
	"catcher/app/internal/lib/caching"
	"catcher/app/internal/lib/logging"
	"context"
)

const (
	DirWeb  = "web"
	DirTemp = "temp"
)

type CtxKey string

type AppContext struct {
	Ctx    context.Context
	Config *config.Config
	Cacher caching.Cacher
	Logger logging.Logger
}

func NewAppContext(ctx context.Context, c *config.Config, cacher caching.Cacher, loger logging.Logger) AppContext {
	return AppContext{
		Ctx:    ctx,
		Config: c,
		Cacher: cacher,
		Logger: loger}
}
