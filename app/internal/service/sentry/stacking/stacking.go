package stacking

import (
	"catcher/app/internal/config"
	"catcher/app/internal/git"
	"catcher/app/internal/models"

	"github.com/getsentry/sentry-go"
)

type Service struct {
	models.AppContext
	prj config.Project
}

func New(prj config.Project, appCtx models.AppContext) Service {
	return Service{
		prj:        prj,
		AppContext: appCtx,
	}
}

func (s Service) AddContextAround(stacktrace *sentry.Stacktrace) error {

	const op = "stack.AddContextAround"

	gitter, err := s.prj.GetGit()
	if err != nil {
		s.Logger.Error("Не удалось создать объект git",
			s.Logger.Err(err),
			s.Logger.Op(op))
		return err
	}

	if gitter == nil {
		return nil // выключено
	}

	svcGit := git.New(gitter, s.prj, s.AppContext)

	for i, frame := range stacktrace.Frames {

		ctxAround := svcGit.GetContextAround(frame.AbsPath, frame.Lineno)
		stacktrace.Frames[i].PreContext = ctxAround.Pre
		stacktrace.Frames[i].PostContext = ctxAround.Post

	}

	return nil

}
