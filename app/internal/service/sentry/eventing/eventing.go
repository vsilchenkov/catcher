package eventing

import (
	"catcher/app/internal/config"
	"catcher/app/internal/models"
	"catcher/app/internal/service/project/nonexept"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
)

const (
	opSending       = "hub.CaptureEvent"	
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

	s.key = fmt.Sprintf("%s:%s", s.event.User.Username, key)
	return true

}

type nonExepter interface {
	Get(ctx context.Context) ([]string, error)
}

func (s *Service) EventNeedSend() bool {

	const op = "eventing.EventNeedSend"

	prj := s.prj

	if !prj.Service.Exeptions.Use {
		return true
	}

	var exepts []string

	// Проверим на исключения
	key := fmt.Sprintf("%s:exeptions", prj.Id)
	found, err := s.Cacher.Get(s.Ctx, key, &exepts)
	if !found {
		if err != nil {
			s.Logger.Error("Ошибка получения значения из кэша",
				s.Logger.Op(op),
				s.Logger.Str("key", key),
				s.Logger.Err(err))
		}

		svc := prj.Service
		creds := nonexept.NewCredintials(svc.Credintials.UserName, svc.Credintials.Password)
		svcNonExept := nonexept.NewService(svc.Url, svc.IimeOut, creds, s.Logger)

		var nonExepter nonExepter = svcNonExept
		e, err := nonExepter.Get(s.Ctx)
		if err != nil {
			return true
		}

		if len(e) > 0 {
			s.Cacher.Set(s.Ctx, key, e, time.Duration(s.prj.Service.Exeptions.Cache.Expiration)*time.Minute)
		} else {
			s.Logger.Warn("Попытка добавить в кэш пустой список исключений",
				s.Logger.Op(op),
				s.Logger.Str("key", key))
		}
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
	key := s.keySending(s.prj)

	if !s.prj.Sentry.SendingCache.Use || key == "" {
		return nil, false
	}

	opCtx := opKey(ctx, op)

	eventID := new(models.EventID)
	found, err := s.Cacher.Get(s.Ctx, key, eventID)
	if found {
		s.Logger.Debug("Используем кэш Event. Cooбщение в Sentry уже было отправлено",
			s.Logger.Op(opCtx),
			s.Logger.Str("key", key),
			s.Logger.Str("eventID", eventID.String()))
	} else {
		if err != nil {
			s.Logger.Error("Ошибка получения значения из кэша",
				s.Logger.Op(op),
				s.Logger.Str("key", key),
				s.Logger.Err(err))
		}
	}

	return eventID, found
}

func (s *Service) AddCacheSending(ctx context.Context, eventID *models.EventID) {

	const op = "eventing.AddCacheSending"
	key := s.keySending(s.prj)

	if !s.prj.Sentry.SendingCache.Use || key == "" {
		return
	}

	opCtx := opKey(ctx, op)
	if eventID != nil {
		s.Cacher.Set(s.Ctx, key, *eventID, time.Duration(s.prj.Sentry.SendingCache.Expiration)*time.Minute)
		s.Logger.Debug("Добавлен кэш сообщения",
			s.Logger.Op(opCtx),
			s.Logger.Str("key", key),
			s.Logger.Str("eventID", eventID.String()))
	} else {
		s.Logger.Warn("Попытка добавить в кэш nil eventID",
			s.Logger.Op(opCtx),
			s.Logger.Str("key", key))
	}

}

func (s *Service) IncrErrors(opCtx string) error {

	const op = "errors"
	key := fmt.Sprintf("%s:%s:%s:%s:%s:%s", s.prj.Id, opCtx, op, s.event.Environment, s.event.User.Username, s.sessionEvent())
	
	s.Logger.Debug("Запись количества ошибок в кэш",
			s.Logger.Str("key", key),
			s.Logger.Op(opCtx))

	_, err := s.Cacher.Incr(s.Ctx, key, time.Duration(s.prj.Session.Duration)*time.Minute)
	if err != nil {
		return errors.WithMessagef(err, op)
	}

	return nil
}

func (s *Service) keySending(prj config.Project) string {
	return fmt.Sprintf("%s:%s:%s", prj.Id, opSending, s.key)
}

func (s *Service) valueExeption() (string, string) {

	for _, e := range s.event.Exception {
		return e.Type, e.Value
	}

	return "", ""
}

func (s Service) sessionEvent() string {

	sessionCtx, ok := s.event.Contexts["Session Data"]
	if !ok {
		sessionCtx, ok = s.event.Contexts["Session_data"]
		if !ok {
			return ""
		}
	}

	value, found := sessionCtx["Session"]
	if !found {
		return ""
	}

	switch v := value.(type) {
	case float64:
		return fmt.Sprintf("%.0f", v) // без дробной части
	case int:
		return fmt.Sprintf("%v", v)
	case string:
		return v
	default:
		return ""
	}

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
