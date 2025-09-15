package logging

type Config struct {
	BuildConfig
	Debug        bool   // Включить режим Отладка
	Level        int    // Уровень логирования от 2 до 5, где 2 - ошибка, 3 - предупреждение, 4 - информация, 5 - дебаг
	OutputInFile bool   // Каталог логов
	Dir          string // Включить логирование в каталог LogDir
}

type SentryConfig struct {
	BuildConfig
	Use              bool // Включить sentry для отлова внутренних ошибок
	Dsn              string
	Environment      string
	AttachStacktrace bool
	TracesSampleRate float64
	EnableTracing    bool
	Debug            bool // Включить режим Отладка
}

type BuildConfig struct {
	Version     string
	ProjectName string
	WorkingDir  string
}

func (c Config) UseDebug() bool {
	return c.Debug
}

func (s SentryConfig) UseDebug() bool {
	return s.Debug
}
