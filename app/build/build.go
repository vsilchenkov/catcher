package build

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/kardianos/service"
)

var Time string
var User string

func NewOption(versionInfoData []byte) (*Option, error) {

	var vi VersionInfo
	if err := json.Unmarshal(versionInfoData, &vi); err != nil {
		return nil, errors.WithMessage(err, "ошибка парсинга файла versioninfo.json")
	}

	return &Option{
		VersionInfo: vi,
		Version:     vi.StringFileInfo.ProductVersion,
		Interactive: service.Interactive(),
		WorkingDir:  WorkingDir(),
	}, nil
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
