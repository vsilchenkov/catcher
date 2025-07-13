package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"catcher/app/build"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/cockroachdb/errors"
	"gopkg.in/yaml.v3"
)

const (
	logLevel_debug int = 5
	sentry_timeout int = 5
)

type Config struct {
	build.Option
	Server struct {
		Port string `yaml:"Port" binding:"required"`
	} `yaml:"Server" binding:"required"`

	Registry struct {
		UserMessage string `yaml:"UserMessage" binding:"required"`
		DumpType    int    `yaml:"DumpType" binding:"required"`
		Timeout     int    `yaml:"Timeout" binding:"required"`
	} `yaml:"Registry"`

	Projects []Project `yaml:"Projects" validate:"required,dive"`

	Log struct {
		Debug        bool   `yaml:"Debug" binding:"required"`
		Level        int    `yaml:"Level" binding:"required"`
		Dir          string `yaml:"Dir" binding:"required"`
		OutputInFile bool   `yaml:"OutputInFile" binding:"required"`
	} `yaml:"Log" binding:"required"`

	DeleteTempFiles bool `yaml:"DeleteTempFiles" binding:"required"`

	Sentry struct {
		Use              bool    `yaml:"Use"`
		Dsn              string  `yaml:"Dsn"`
		AttachStacktrace bool    `yaml:"AttachStacktrace"`
		TracesSampleRate float64 `yaml:"TracesSampleRate"`
		EnableTracing    bool    `yaml:"EnableTracing"`
	} `yaml:"Sentry"`
}

type flags struct {
	configPath string
}

func newConfig(b build.Option) *Config {
	return &Config{
		Option: b,
	}
}

func ParseFlags() flags {

	var debug bool
	var configPath string
	flag.BoolVar(&debug, "debug", false, "Use debug")
	flag.StringVar(&configPath, "config", "config/config.yml", "Путь к файлу настроек")
	flag.Parse()

	flags := flags{
		configPath: configPath,
	}

	return flags

}

func LoadSettigs(flags flags) (*Config, error) {

	const op = "config.LoadSettigs"

	b := build.NewOption()

	path := flags.configPath
	fullPath := filepath.Join(b.WorkingDir, path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, errors.WithMessagef(err, "%s - файл настроек %q не найден", op, fullPath)
	}
	file, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "%s - ошибка чтения файла %q", op, fullPath)
	}

	c := newConfig(*b)
	if err := yaml.Unmarshal(file, &c); err != nil {
		return nil, errors.WithMessagef(err, "%s - ошибка десириализации настроек", op)
	}

	validate := validator.New()
	err = validate.Struct(c)
	if err != nil {
		return nil, errors.WithMessagef(err, "%s - ошибка валидации настроек", op)
	}

	if err := defaults.Set(c); err != nil {
		return nil, errors.WithMessagef(err, "%s - ошибка установки настроек по умолчанию", op)
	}

	if c.Registry.Timeout == 0 {
		c.Registry.Timeout = sentry_timeout
	}

	loadEnv(c)

	fmt.Printf("Settings loaded: %s\n", fullPath)

	return c, nil
}

func loadEnv(c *Config) {

	godotenv.Load()
	Dsn := os.Getenv("SENTRY_DSN")
	if Dsn != "" {
		c.Sentry.Dsn = Dsn
	}
}

func (c Config) UseDebug() bool {
	return c.Log.Debug
}

func (c Config) ServerPort() string {
	return c.Server.Port
}

func (c Config) RegistryUserMessage() string {
	return c.Registry.UserMessage
}

func (c Config) RegistryDumpType() int {
	return c.Registry.DumpType
}

func (c Config) ProjectByName(name string) (Project, error) {

	for _, v := range c.Projects {
		if v.Name == name {
			return v, nil
		}
	}

	return Project{}, errors.New(fmt.Sprintf("Настройки проекта не найдены по имени: %s", name))
}

func (c Config) ProjectById(id string) (Project, error) {

	for _, v := range c.Projects {
		if v.Id == id {
			return v, nil
		}
	}
	return Project{}, errors.New(fmt.Sprintf("Настройки проекта не найдены по ID: %s", id))
}
