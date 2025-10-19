package models

import (
	"catcher/app/internal/config"
	"catcher/app/internal/lib/caching"
	"catcher/pkg/logging"
	"context"
	"io"
)

const (
	DirWeb  = "web"
	DirTemp = "temp"
)

type CtxKey string

type Template interface {
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

type AppContext struct {
	Ctx    context.Context
	Config *config.Config
	Cacher caching.Cacher
	Logger logging.Logger
	Tmpl   Template
}

func NewAppContext(ctx context.Context, c *config.Config, cacher caching.Cacher, loger logging.Logger, tmpl Template) AppContext {
	return AppContext{
		Ctx:    ctx,
		Config: c,
		Cacher: cacher,
		Logger: loger,
		Tmpl:   tmpl}
}
