package config

type SentryMetrics struct {
	Port     string       `yaml:"Port"`
	Interval int          `yaml:"Interval"`
	Sentry   SentryConfig `yaml:"Sentry"`
	Metrics  []Metrics      `yaml:"Metrics"`
}

type Metrics struct {
	Use      bool             `yaml:"Use"`
	Name     string           `yaml:"Name"`
	Projects []ProjectMetrics `yaml:"Projects"`
}

type SentryConfig struct {
	Url     string `yaml:"Url"`
	Org     string `yaml:"Org"`
	Token   string `yaml:"Token"`
	Timeout int    `yaml:"Timeout"`
}

type ProjectMetrics struct {
	Id   string `yaml:"Id"`
	Name string `yaml:"Name"`
}
