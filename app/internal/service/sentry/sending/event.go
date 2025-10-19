package sending

import (
	"catcher/app/internal/config"
	"catcher/app/internal/models"
	sentryhub "catcher/app/internal/sentry/hub"
	"catcher/app/internal/service/sentry/eventing"
	"context"
	"log/slog"

	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
)

type Event struct {
	models.AppContext
	Prj   config.Project
	Event *sentry.Event
	ID    string
}

func NewEvent(prj config.Project, Id string, e *sentry.Event, appCtx models.AppContext) Event {
	return Event{
		Prj:        prj,
		Event:      e,
		ID:         Id,
		AppContext: appCtx,
	}
}

type eventer interface {
	IsEventSent(ctx context.Context) (*models.EventID, bool)
	IncrErrors(opCtx string) error
}

func (e Event) Send() (*models.EventID, error) {

	const op = "sending.event.send"

	const opKey models.CtxKey = "op"
	ctx := context.WithValue(e.Ctx, opKey, op)

	// cache
	svcEventing := eventing.New(e.Event, e.Prj, e.AppContext)
	var eventer eventer = svcEventing

	x, found := eventer.IsEventSent(ctx)
	if found {
		return x, nil
	}

	// IncrErrors
	if err := eventer.IncrErrors(models.OpEventer); err != nil {
		e.Logger.Error("Ошибка увеличение счетчика отправок",
			e.Logger.Op(op),
			e.Logger.Err(err))
	}

	var err error

	hub, err := sentryhub.Get(e.Prj, e.AppContext)
	if err != nil {
		e.Logger.Error("Ошибка получение sentry hub",
			e.Logger.Op(op),
			e.Logger.Err(err))
		return nil, err
	}

	hub.Scope().ClearBreadcrumbs()

	eventID := hub.CaptureEvent(e.Event)
	if eventID == nil {
		e.Logger.Error("Ошибка отправки соообщения в Sentry",
			e.Logger.Op(op),
			e.Logger.Str("ID", e.ID))
		return nil, errors.Errorf("sending error %s", e.ID)
	}

	e.Logger.Info("Отправлено в Sentry сообщение",
		e.Logger.Str("ID", e.ID),
		slog.Any("eventID", eventID),
		e.Logger.Op(op))

	res := models.EventID(*eventID)
	svcEventing.AddCacheSending(ctx, &res)

	return &res, nil

}
