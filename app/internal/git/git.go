package git

import (
	"catcher/app/internal/config"
	"catcher/app/internal/models"
	"fmt"
	"strings"
	"time"
)

type Gitter interface {
	GetFileContent(filePath string) (*string, error)
}

type ContextAround struct {
	Pre  []string
	Post []string
}

type Service struct {
	Gitter
	models.AppContext
	prj config.Project
}

func New(g Gitter, prj config.Project, appCtx models.AppContext) *Service {

	if g == nil || !prj.Sentry.ContextAround.Use {
		return nil
	}

	return &Service{
		Gitter:     g,
		AppContext: appCtx,
		prj:        prj,
	}

}

func (s Service) GetContextAround(filePath string, lineno int) ContextAround {

	const op = "exceptionbsl.prePostContext"

	prj := s.prj

	if filePath == "" {
		return ContextAround{}
	}

	useCache := prj.Sentry.ContextAround.Cache.Use
	key := fmt.Sprintf("%s:%s:%s:%d", op, prj.Name, filePath, lineno)

	if useCache {
		if x, found := s.Cacher.Get(s.Ctx, key); found {
			s.Logger.Debug("Используем кэш ContextAround",
				s.Logger.Op(op),
				s.Logger.Str("key", key))
			return x.(ContextAround)
		}
	}

	content, err := s.Gitter.GetFileContent(filePath)
	if err != nil {
		s.Logger.Warn("Не удалось получить контент из git",
			s.Logger.Err(err),
			s.Logger.Str("filePath", filePath),
			s.Logger.Op(op))
		return ContextAround{}
	}

	lines := strings.Split(*content, "\n")
	lenLines := len(lines)
	if lineno > lenLines {
		return ContextAround{}
	}

	quantity := prj.Sentry.ContextAround.Quantity
	indexLineno := lineno - 1

	var pre, post []string

	start := indexLineno - quantity
	start = max(start, 0)
	start = min(start, lenLines-1)
	qSart := indexLineno - start
	if qSart > 0 {
		pre = make([]string, qSart)
		count := 0
		for i := start; i <= indexLineno-1; i++ {
			pre[count] = lines[i]
			count++
		}
	}

	end := indexLineno + quantity
	end = min(end, lenLines-1)
	qEnd := end - indexLineno
	if qEnd > 0 {
		post = make([]string, qEnd)
		count := 0
		for i := indexLineno + 1; i <= end; i++ {
			post[count] = lines[i]
			count++
		}
	}

	res := ContextAround{
		Pre:  pre,
		Post: post,
	}

	s.Logger.Debug("Получен контент из git",
			s.Logger.Op(op),
			s.Logger.Str("filePath", filePath))
						
	if useCache {
		s.Cacher.Set(s.Ctx, key, res, time.Duration(prj.Sentry.ContextAround.Cache.Expiration)*time.Minute)
	}

	return res
}
