package logging

import "log/slog"

// TestLogger реализует интерфейс logging.Logger для тестов
type TestLogger struct {
	debugMessages []string
	warnMessages  []string
	infoMessages  []string
	errorMessages []string
}

func (l *TestLogger) Debug(msg string, attrs ...slog.Attr) {
	l.debugMessages = append(l.debugMessages, msg)
}

func (l *TestLogger) Info(msg string, attrs ...slog.Attr) {
	l.infoMessages = append(l.infoMessages, msg)
}

func (l *TestLogger) Warn(msg string, attrs ...slog.Attr) {
	l.warnMessages = append(l.warnMessages, msg)
}

func (l *TestLogger) Error(msg string, attrs ...slog.Attr) {
	l.errorMessages = append(l.errorMessages, msg)
}

func (l *TestLogger) Err(err error) slog.Attr {
	return slog.Any("error", err)
}

func (l *TestLogger) Op(value string) slog.Attr {
	return slog.Attr{
		Key:   "op",
		Value: slog.StringValue(value),
	}
}

func (l *TestLogger) Str(key, value string) slog.Attr {
	return slog.Attr{
		Key:   key,
		Value: slog.StringValue(value),
	}
}