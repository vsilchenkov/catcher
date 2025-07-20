package reporting

import (
	"catcher/app/internal/lib/gitbsl"
	"catcher/app/internal/sentryhub/normalize"
	"catcher/app/internal/service/sentry/stacking"
	"fmt"
	"strings"

	"github.com/getsentry/sentry-go"
)

type Gitter interface {
	GetFileContent(filePath string) (*string, error)
}

func (r Report) exception() []sentry.Exception {

	message := r.exceptionMessage()
	if len(message) == 0 {
		return nil
	}

	result := make([]sentry.Exception, 1)

	//	Разделим и получим Type и Value
	// "{ОбщийМодуль.ОбменСSentry.Модуль(1114)}: Ошибка при вызове метода контекста (Вставить)",
	part := normalize.SplitString(message[0], "{", "}")
	t := part[0]
	v := part[1]

	t = strings.TrimPrefix(t, ":")
	t = strings.TrimSpace(t)

	exeption := sentry.Exception{
		Type:  t,
		Value: v,
	}

	// Подставим Module без номера строки
	//ОбщийМодуль.ОбменСSentry.Модуль(1114)
	part = normalize.SplitString(v, "(", ")")
	m := part[0]
	if m != "" {
		exeption.Module = m
	}

	exeption.Stacktrace = r.stacktrace()

	result[0] = exeption
	return result
}

// frame
// "ВнешняяОбработка.Ошибка.Форма.Форма.Форма", 9, "СоСтекомНаСервере()"
func (r Report) stacktrace() *sentry.Stacktrace {

	stack := r.Data.ErrorInfo.ApplicationErrorInfo.Stack
	lenstack := len(stack)

	if lenstack == 0 {
		return nil
	}

	frames := make([]sentry.Frame, lenstack)

	for i, v := range stack {

		moduleAbs := fmt.Sprintf("%v", v[0])
		module := normalize.RemoveModuleNameSuffix(moduleAbs)

		function := fmt.Sprintf("%v", v[2])
		function = normalize.RemoveModuleNameSuffix(function)

		isExternal := gitbsl.IsExternalModule(module) || gitbsl.IsExpansion(module)

		inApp := !isExternal
		fileName := fileNameByModule(module, isExternal)
		contextLine := contextByFunction(function)

		var lineno int
		l, ok := v[1].(float64)
		if ok {
			lineno = int(l)
		}

		absPath := r.absPath(moduleAbs, isExternal)

		var stackStart bool
		if len(frames)-1 == i {
			stackStart = true
		}

		frames[i] = sentry.Frame{
			Function:    function,
			Module:      module,
			Lineno:      lineno,
			InApp:       inApp,
			Filename:    fileName,
			ContextLine: contextLine,
			StackStart:  stackStart,
			AbsPath:     absPath,
			Platform:    r.Prj.Sentry.Platform,
		}

	}

	stacktrace := &sentry.Stacktrace{
		Frames: frames,
	}

	svcStack := stacking.New(r.Prj, r.AppContext)
	svcStack.AddContextAround(stacktrace)

	return stacktrace

}

func (r Report) exceptionMessage() []string {

	var result []string

	errorInfo := r.Data.ErrorInfo.ApplicationErrorInfo.Errors

	// текст ошибки на втором уровне
	for _, v := range errorInfo {
		if len(v) > 0 {
			e := strings.TrimSpace(fmt.Sprintf("%v", v[0]))
			if e != "" {
				result = append(result, e)
			}
		}
	}

	return result

}

func (r Report) absPath(m string, isExternal bool) string {

	if isExternal {
		return ""
	}

	result, _ := gitbsl.NewPath(m, r.Prj.Git.SourceCodeRoot, r.Logger).AbsPath()
	return result

}

func fileNameByModule(m string, isExternal bool) string {

	if isExternal {
		part := strings.Split(m, ".")
		if len(part) > 1 {
			return fmt.Sprintf("%v.%v", part[0], part[1])
		} else {
			return m
		}

	} else {
		return m
	}
}

func contextByFunction(f string) string {
	return strings.TrimSpace(f)
	// part := strings.Split(f, ".")
	// return part[len(part)-1]
}
