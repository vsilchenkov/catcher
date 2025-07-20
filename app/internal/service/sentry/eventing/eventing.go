package eventing

import (
	"catcher/app/internal/config"
	"catcher/app/internal/models"
	"catcher/app/internal/service/project/nonexept"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

const (
	opSending       = "hub.CaptureEvent"
	expirationCache = 5
)

type Service struct {
	event *sentry.Event
	models.AppContext
	prj config.Project
	key string
}

func New(e *sentry.Event, prj config.Project, appCtx models.AppContext) *Service {
	s := &Service{
		event:      e,
		prj:        prj,
		AppContext: appCtx,
	}
	s.generateKey()
	return s
}

func (s *Service) generateKey() bool {

	exception := s.event.Exception

	if len(exception) == 0 {
		return false
	}

	stacktrace := exception[0].Stacktrace
	if stacktrace == nil {
		return false
	}

	frames := stacktrace.Frames
	if len(frames) == 0 {
		return false
	}

	var key string
	for _, f := range frames {
		if f.StackStart {
			
			var m string
			switch {
			case f.AbsPath != "":
				m = f.AbsPath
			case f.Module != "":
				m = f.Module
			default:
				m = f.Filename
			}
			
			key = fmt.Sprintf("%s:%d", m, f.Lineno)
		}
	}

	if key == "" {
		return false
	}

	s.key = fmt.Sprintf("%s:%s:%s", s.prj.Name, s.event.User.Username, key)
	return true

}

type nonExepter interface {
	Get(ctx context.Context) ([]string, error)
}

func (s *Service) EventNeedSend() bool {

	prj := s.prj

	if !prj.Service.Exeptions.Use {
		return true
	}

	var exepts []string

	// Проверим на исключения
	key := fmt.Sprintf("%s:exeptions", prj.Id)
	x, found := s.Cacher.Get(s.Ctx, key)
	if found {
		exepts = x.([]string)
	} else {

		svc := prj.Service
		creds := nonexept.NewCredintials(svc.Credintials.UserName, svc.Credintials.Password)
		svcNonExept := nonexept.NewService(svc.Url, svc.IimeOut, creds, s.Logger)

		var nonExepter nonExepter = svcNonExept
		e, err := nonExepter.Get(s.Ctx)
		if err != nil {
			return true
		}

		s.Cacher.Set(s.Ctx, key, e, time.Duration(s.prj.Service.Exeptions.Cache.Expiration)*time.Minute)
		exepts = e

	}

	t, v := s.valueExeption()

	if t == "" && v == "" {
		return true
	}

	s.Logger.Debug("Поиск значения в списке пропускаемых исключений",
		s.Logger.Str("type", t),
		s.Logger.Str("value", v))

	for _, e := range exepts {

		if t != "" && strings.Contains(t, e) {
			s.Logger.Debug("Найдено значение (type) в списке пропускаемых исключений",
				s.Logger.Str("type", t),
				s.Logger.Str("exept", e))
			return false
		}

		if v != "" && strings.Contains(v, e) {
			s.Logger.Debug("Найдено значение (value) в списке пропускаемых исключений",
				s.Logger.Str("value", v),
				s.Logger.Str("exept", e))
			return false
		}

	}

	return true
}

func (s *Service) IsEventSent(ctx context.Context) (*models.EventID, bool) {

	const op = "eventing.IsEventSent"
	key := s.keySending()

	if !s.prj.Sentry.SendingCache.Use || key == "" {
		return nil, false
	}

	opCtx := opKey(ctx, op)

	var eventID *models.EventID
	x, found := s.Cacher.Get(s.Ctx, key)
	if found {
		eventID = x.(*models.EventID)
		s.Logger.Debug("Используем кэш Event. Cooбщение в Sentry уже было отправлено",
			s.Logger.Op(opCtx),
			s.Logger.Str("key", key),
			s.Logger.Str("eventID", eventID.String()))

	}

	return eventID, found
}

func (s *Service) AddCacheSending(ctx context.Context, eventID *models.EventID) {

	const op = "eventing.AddCacheSending"
	key := s.keySending()

	if !s.prj.Sentry.SendingCache.Use || key == "" {
		return
	}

	opCtx := opKey(ctx, op)
	s.Cacher.Set(s.Ctx, key, eventID, time.Duration(s.prj.Sentry.SendingCache.Expiration)*time.Minute)

	s.Logger.Debug("Добавлен кэш сообщения",
		s.Logger.Op(opCtx),
		s.Logger.Str("key", key),
		s.Logger.Str("eventID", eventID.String()))

}

func (s *Service) keySending() string {
	return fmt.Sprintf("%s:%s", opSending, s.key)
}

func (s *Service) valueExeption() (string, string) {

	for _, e := range s.event.Exception {
		return e.Type, e.Value
	}

	return "", ""
}

func opKey(ctx context.Context, op string) string {

	const opKey models.CtxKey = "op"

	var res string
	opCtx, ok := ctx.Value(opKey).(string)
	if ok {
		res = fmt.Sprintf("%s:%s", opCtx, op)
	} else {
		res = op
	}

	return res
}
