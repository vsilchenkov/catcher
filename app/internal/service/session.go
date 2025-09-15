package service

// https://develop.sentry.dev/sdk/telemetry/sessions

import (
	"catcher/app/internal/config"
	"catcher/app/internal/lib/caching"
	"catcher/app/internal/models"
	"catcher/app/internal/sentryhub"
	"catcher/pkg/logging"
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type SessionService struct {
	ctx    context.Context
	config *config.Config
	cacher caching.Cacher
	logger logging.Logger
}

func NewSessionService(appCtx models.AppContext) *SessionService {
	return &SessionService{
		ctx:    appCtx.Ctx,
		config: appCtx.Config,
		cacher: appCtx.Cacher,
		logger: appCtx.Logger}
}

func (s SessionService) Start(projectId string, input models.Session) error {

	const op = "service.session.start"

	prj, err := s.config.ProjectById(projectId)
	if err != nil {
		return errors.WithMessage(ErrBadProject, op)
	}

	if input.Started.IsZero() {
		input.Started = time.Now().UTC()
	}

	err = s.StartEnd(prj, input, true, op)
	if err == nil {
		key := input.Key(prj)
		s.logger.Debug("Запись сессии в кэш",
			s.logger.Str("key", key),
			s.logger.Op(op))

		s.cacher.Set(s.ctx, key, input, time.Duration(prj.Session.Duration)*time.Minute)
	}

	return err

}

func (s SessionService) End(projectId string, input models.Session) error {

	const op = "service.session.end"

	prj, err := s.config.ProjectById(projectId)
	if err != nil {
		return errors.WithMessage(ErrBadProject, op)
	}

	if input.Started.IsZero() {
		// Заполним Started из кэша
		cach, err := s.sessionCache(prj, input)
		if err != nil {
			s.logger.Warn("Cессия в кэше не найдена",
				s.logger.Op(op),
				s.logger.Str("did", input.Did),
				s.logger.Str("sid", input.Sid),
				s.logger.Err(err))
			return errors.WithMessage(ErrBadSesion, op)
		}
		input.Started = cach.Started
	}

	// Заполним Errors
	errCount, err := s.errorsCount(prj, input)
	if err != nil {
		s.logger.Error("Ошибка получения количества ошибок из кэша",
			s.logger.Op(op),
			s.logger.Str("did", input.Did),
			s.logger.Str("sid", input.Sid),
			s.logger.Err(err))
		return errors.WithMessage(ErrBadSesionErrors, op)

	}
	input.ErrorsReporter = errCount[models.OpReporter]
	input.ErrorsEventer = errCount[models.OpEventer]

	return s.StartEnd(prj, input, false, op)

}

func (s SessionService) StartEnd(prj config.Project, input models.Session, start bool, op string) error {

	var err error

	appCtx := models.NewAppContext(s.ctx, s.config, s.cacher, s.logger)

	hub, err := sentryhub.Get(prj, appCtx)
	if err != nil {
		s.logger.Error("Ошибка получение sentry hub",
			s.logger.Op(op),
			s.logger.Err(err))
		return ErrBadSentryHub
	}

	sentrySession := sentryhub.SentrySession{
		Sid:     uuidv5(input.Did + input.Sid),
		Did:     input.Did,
		Started: input.Started.Format(time.RFC3339),
		Errors:  input.ErrorsReporter + input.ErrorsEventer,
		Status:  sesionStatus(start, input),

		Attrs: sentryhub.SessionAttrs{
			Release:     input.Release,
			Environment: input.Environment,
		},
	}

	if start {
		sentrySession.Init = true
		err = hub.StartSession(sentrySession)
	} else {
		err = hub.EndSession(sentrySession)
	}

	if err != nil {
		s.logger.Error("Ошибка отправки sesion в Sentry",
			s.logger.Op(op),
			s.logger.Err(err))
		return ErrBadRequestSentry
	}

	return nil

}

func (s SessionService) sessionCache(prj config.Project, input models.Session) (*models.Session, error) {

	const op = "sessionCache"

	key := input.Key(prj)
	
	var res *models.Session
	found, err := s.cacher.Get(s.ctx, key, &res)
	if err != nil {
		return nil, errors.WithMessage(err, key)
	} else if !found {
		return nil, errors.WithMessage(ErrBadSesion, key)
	}

	s.logger.Debug("Получена сессия из кэша",
		s.logger.Str("key", key),
		s.logger.Op(op))

	return res, nil
}

func (s SessionService) errorsCount(prj config.Project, input models.Session) (map[string]int, error) {

	const op = "errors"

	result := map[string]int{
		models.OpEventer:  0,
		models.OpReporter: 0}

	for k := range result {

		key := fmt.Sprintf("%s:%s:%s:%s:%s:%s", prj.Id, k, op, input.Environment, input.Did, input.Sid)

		var res int64
		found, err := s.cacher.Get(s.ctx, key, &res)
		if err != nil {
			return result, errors.WithMessage(err, key)
		} else if found {
			result[k] = int(res)
		}

		s.logger.Debug("Получено количество ошибок из кэша",
			s.logger.Str("key", key),
			s.logger.Op(op),
			s.logger.Str("Count", fmt.Sprintf("%v", result[k])))

	}

	return result, nil
}

// ok: The session is currently in progress but healthy. This can be the terminal state of a session.
// exited: The session terminated normally.
// crashed: The session terminated in a crash.
// abnormal: The session encountered a non crash related abnormal exit.
func sesionStatus(start bool, input models.Session) string {

	switch {
	case start:
		return "ok"
	case input.ErrorsReporter > 0:
		return "crashed" // считаем, что если есть ошибки через сервис report, то это crash
	default:
		return "exited"
	}

}

// uuid5 — генерирует UUID версии 5 из строки
func uuidv5(seed string) string {

	// Берём стандартный namespace (можно любой свой)
	namespace := uuid.NameSpaceDNS // или NameSpaceURL, NameSpaceOID, NameSpaceX500

	// Генерируем детерминированный UUID v5
	u := uuid.NewSHA1(namespace, []byte(seed))

	return u.String()
}
