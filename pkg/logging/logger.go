package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"log/slog"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
)

// Logger — интерфейс логгера, используемый в проекте
type Logger interface {
	Debug(msg string, attrs ...slog.Attr)
	Info(msg string, attrs ...slog.Attr)
	Warn(msg string, attrs ...slog.Attr)
	Error(msg string, attrs ...slog.Attr)

	// Вспомогательные методы для создания атрибутов
	Err(err error) slog.Attr
	Op(value string) slog.Attr
	Str(key, value string) slog.Attr
}

var logger *slog.Logger

type defaultlogger struct {
	*slog.Logger
}

var levelMap = map[int]slog.Level{
	5: slog.LevelDebug,
	4: slog.LevelInfo,
	3: slog.LevelWarn,
	2: slog.LevelError,
}

func GetLogger() Logger {
	return NewLogger(logger)
}

func Initlogger(c *Config, sConfig *SentryConfig) Logger {

	var handler slog.Handler

	level, ok := levelMap[c.Level]
	if !ok {
		level = slog.LevelInfo
	}

	if !c.OutputInFile {
		// Вывод в stdout с текстовым форматом
		handler = tint.NewHandler(colorable.NewColorable(os.Stderr), &tint.Options{
			Level:      level,
			TimeFormat: "15:04:05",
			AddSource:  false,
			NoColor:    false,
		})
	} else {
		// Вывод в файл с JSON форматом
		fileName := "app.log"
		file, err := GetOutputLogFile(c.WorkingDir, c.Dir, fileName)
		if err != nil {
			log.Printf("Не удалось открыть файл логов %q, используется стандартный stderr\n%v", fileName, err)
			handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
				Level:     level,
				AddSource: false,
			})
		} else {
			handler = slog.NewJSONHandler(file, &slog.HandlerOptions{
				Level:     level,
				AddSource: false,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						// Преобразуем время в строку с нужным форматом
						t := a.Value.Time()
						formatted := t.Format("2006-01-02 15:04:05")
						return slog.String(a.Key, formatted)
					}
					return a
				},
			})
		}
	}

	if sConfig.Use {
		multiHandler := NewMultiHandler(handler, SentryHandler())
		logger = slog.New(multiHandler)
	} else {
		logger = slog.New(handler)
	}

	return NewLogger(logger)
}

// NewLogger создаёт Logger из slog.Logger
func NewLogger(slogger *slog.Logger) Logger {
	return &defaultlogger{slogger}
}

func (l *defaultlogger) Debug(msg string, attrs ...slog.Attr) {
	l.Logger.Debug(msg, convertAttrsToAny(attrs)...)
}

func (l *defaultlogger) Info(msg string, attrs ...slog.Attr) {
	l.Logger.Info(msg, convertAttrsToAny(attrs)...)
}

func (l *defaultlogger) Warn(msg string, attrs ...slog.Attr) {
	l.Logger.Warn(msg, convertAttrsToAny(attrs)...)
}

func (l *defaultlogger) Error(msg string, attrs ...slog.Attr) {
	l.Logger.Error(msg, convertAttrsToAny(attrs)...)
}

func (l *defaultlogger) Err(err error) slog.Attr {
	return slog.Any("error", err)
}

func (l *defaultlogger) Op(value string) slog.Attr {
	return slog.Attr{
		Key:   "op",
		Value: slog.StringValue(value),
	}
}

func (l *defaultlogger) Str(key, value string) slog.Attr {
	return slog.Attr{
		Key:   key,
		Value: slog.StringValue(value),
	}
}

func GetOutputLogFile(workingDir, logDir, fileName string) (*os.File, error) {

	fullLogDir := filepath.Join(workingDir, logDir)

	err := os.MkdirAll(fullLogDir, 0755)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", fullLogDir, fileName)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	return file, err
}

func convertAttrsToAny(attrs []slog.Attr) []any {
	res := make([]any, len(attrs))
	for i, v := range attrs {
		res[i] = v
	}
	return res
}
