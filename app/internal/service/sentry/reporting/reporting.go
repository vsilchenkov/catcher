package reporting

import (
	"catcher/app/internal/config"
	"catcher/app/internal/models"
	"catcher/app/internal/sentryhub"
	"catcher/app/internal/service/sentry/eventing"
	"context"
	"log/slog"

	"github.com/cockroachdb/errors"
)

type Report struct {
	models.AppContext
	ID    string
	Prj   config.Project
	Data  models.Repport
	Files []models.FileData
}

func NewReport(id string, prj config.Project, data models.Repport, files []models.FileData, appCtx models.AppContext) Report {
	return Report{
		ID:         id,
		Prj:        prj,
		Data:       data,
		Files:      files,
		AppContext: appCtx,
	}
}

type eventer interface {
	IsEventSent(ctx context.Context) (*models.EventID, bool)
	EventNeedSend() bool
}

func (r Report) Send() (*models.EventID, error) {

	const op = "reporting.send"

	const opKey models.CtxKey = "op"
	ctx := context.WithValue(r.Ctx, opKey, op)

	data := r.Data

	var prj config.Project
	var err error

	additionalInfo := data.AdditionalInfo
	if additionalInfo != "" {
		prj, err = r.Config.ProjectById(additionalInfo)
	} else {
		prj, err = r.Config.ProjectByName(data.ConfigInfo.Name)
	}

	if err != nil {
		r.Logger.Error("Ошибка получение настроек проекта",
			r.Logger.Op(op),
			r.Logger.Err(err))
		return nil, err
	}

	event, err := r.event()
	if err != nil {
		r.Logger.Error("Ошибка сборки event",
			r.Logger.Op(op),
			r.Logger.Str("ID", r.ID),
			r.Logger.Err(err))
		return nil, err
	}

	// cache
	svcEventing := eventing.New(event, prj, r.AppContext)
	var eventer eventer = svcEventing
	x, found := eventer.IsEventSent(ctx)
	if found {
		return x, nil
	}

	// nonexept
	if !eventer.EventNeedSend() {
		r.Logger.Debug("Сообщение пропущено, не требует отправки",
			r.Logger.Str("ID", r.ID))
			res := models.EventID(r.ID) // возвращаем входящий ID
		return &res, nil
	}

	hub, err := sentryhub.Get(prj, r.AppContext)
	if err != nil {
		r.Logger.Error("Ошибка получение настроек проекта",
			r.Logger.Op(op),
			r.Logger.Err(err))
		return nil, err
	}

	hub.Scope().ClearBreadcrumbs()

	eventID := hub.CaptureEvent(event)
	if eventID == nil {
		r.Logger.Error("Ошибка отправки соообщения в Sentry",
			r.Logger.Op(op),
			r.Logger.Str("ID", r.ID))
		return nil, errors.Errorf("sending error %s", r.ID)
	}

	r.Logger.Debug("Отправлено в Sentry сообщение",
		r.Logger.Str("ID", r.ID),
		slog.Any("eventID", eventID),
		r.Logger.Op(op))

	res := models.EventID(*eventID)
	svcEventing.AddCacheSending(ctx, &res)

	return &res, nil

}
