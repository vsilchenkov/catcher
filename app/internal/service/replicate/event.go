package replicate

import (
	"catcher/app/internal/config"
	"catcher/app/internal/lib/gitbsl"
	"catcher/app/internal/models"
	"catcher/app/internal/sentryhub/normalize"
	"catcher/app/internal/service/sentry/stacking"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jinzhu/copier"
)

type ConvertEventer interface {
	ConvertEvent(prj config.Project, input models.Event) (*sentry.Event, error)
}

func (s Service) ConvertEvent(prj config.Project, input models.Event) (*sentry.Event, error) {

	event := input.Event

	if s.Config.UseDebug() {
		event.Timestamp = time.Now()
	} else {
		event.Timestamp = time.Time(input.Timestamp)
	}

	if event.Platform == "" {
		event.Platform = prj.Sentry.Platform
	}

	event.Exception = s.ConvertExeptions(prj, input.Exception)
	event.Message = normalize.RemoveFromSecondBrace(event.Message, true)
	event.Message = normalize.RemoveFromSecondBrace(event.Message, false) // два раза

	return &event, nil

}

func (s Service) ConvertExeptions(prj config.Project, input models.Exception) []sentry.Exception {

	exeptions := make([]sentry.Exception, len(input.Values))
	for i, value := range input.Values {

		exeption := sentry.Exception{}
		exeption.Type = value.Type
		exeption.Value = value.Value

		stacktrace := sentry.Stacktrace{}

		var exeptModule string
		frames := make([]sentry.Frame, len(value.Stacktrace.Frames))
		for j, valFrame := range value.Stacktrace.Frames {

			frame := sentry.Frame{}
			copier.Copy(&frame, &valFrame)

			if frame.Platform == "" {
				frame.Platform = prj.Sentry.Platform
			}

			module := valFrame.Module
			if exeptModule == "" {
				exeptModule = module
			}

			isExternal := gitbsl.IsExternalModule(module) || gitbsl.IsExpansion(module)
			inApp := !isExternal
			frame.InApp = inApp
			s.absPath(&prj, &valFrame, &frame)

			frames[j] = frame

		}
		stacktrace.Frames = frames
		exeption.Module = exeptModule
		exeption.Stacktrace = &stacktrace

		svcStack := stacking.New(prj, s.AppContext)
		svcStack.AddContextAround(exeption.Stacktrace)

		exeptions[i] = exeption
	}

	return exeptions
}

func (s Service) absPath(prj *config.Project, v *models.Frame, frame *sentry.Frame) {

	if !frame.InApp || frame.AbsPath != "" {
		return
	}

	var m string
	if v.ModuleAbs != "" {
		m = v.ModuleAbs
	} else {
		m = v.Module
	}

	result, _ := gitbsl.NewPath(m, prj.Git.SourceCodeRoot, s.Logger).AbsPath()
	frame.AbsPath = result

}
