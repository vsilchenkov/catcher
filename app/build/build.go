package build

import (
	"os"
	"path/filepath"

	"github.com/kardianos/service"
)

const ProjectName = "Cather"

var Version = "1.5.4"

var Time string
var User string

type Option struct {
	Version     string
	ProjectName string
	WorkingDir  string
	Interactive bool
}

func NewOption() *Option {
	return &Option{
		ProjectName: ProjectName,
		Version:     Version,
		Interactive: service.Interactive(),
		WorkingDir:  WorkingDir(),
	}
}

// Получить абсолютный путь к папке, где лежит .exe или запущена программа
func WorkingDir() string {

	if service.Interactive() {

		cwd, err := os.Getwd()
		if err != nil {
			return "."
		}

		return cwd

	} else {

		exePath, err := os.Executable()
		if err != nil {
			return "."
		}
		return filepath.Dir(exePath)
	}
}
