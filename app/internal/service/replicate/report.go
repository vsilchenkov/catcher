package replicate

import (
	"catcher/app/internal/lib/iofiles"
	"catcher/app/internal/lib/logging"
	"catcher/app/internal/models"
	"catcher/app/internal/service/redirect"
	"catcher/app/internal/service/project/userinfo"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	reportJson = "report.json"
	reportZip  = "report.zip"
)

var errBadProject = errors.New("no project setting")

type Service struct {
	models.AppContext
}

func New(appCtx models.AppContext) Service {
	return Service{
		AppContext: appCtx}
}

func (s Service) ConvertReport(r redirect.Report) (*models.RepportData, error) {

	const op = "redirect.convert"

	dir := filepath.Join(s.Config.WorkingDir, models.DirWeb, models.DirTemp)

	fileName := fmt.Sprintf("%s_%s", r.ID, reportZip)
	src := filepath.Join(dir, fileName)

	if err := iofiles.SaveDataToFile(r.Data, src); err != nil {
		s.Logger.Error("Не удалось выгрузить файл",
			s.Logger.Op(op),
			s.Logger.Str("name", src),
			s.Logger.Err(err))
		return nil, err
	}
	s.Logger.Debug("Сохранен файл",
		s.Logger.Str("name", src))

	dest := filepath.Join(dir, r.ID)
	if err := iofiles.Unzip(src, dest); err != nil {
		s.Logger.Error("Не удалось распаковать файл",
			s.Logger.Op(op),
			s.Logger.Str("name", src),
			s.Logger.Err(err))
		return nil, err
	}
	s.Logger.Debug("Распакован файл",
		s.Logger.Str("name", src),
		s.Logger.Str("dest", dest))

	srcjson := filepath.Join(dest, reportJson)
	repportData, err := fileToData(srcjson, s.Logger)
	if err != nil {
		return nil, err
	}
	s.Logger.Debug("Прочитан файл",
		s.Logger.Str("name", srcjson))

	files, err := listFiles(dest, s.Logger)
	if err != nil {
		return nil, err
	}

	prj, err := repportData.ProjectByConfig(s.Config)
	if err != nil {
		s.Logger.Error("Нет настройки для проекта",
			s.Logger.Op(op),
			s.Logger.Err(err))
		return nil, errBadProject
	}

	rd := &models.RepportData{
		ID:          r.ID,
		Prj:         prj,
		Data:        repportData,
		Files:       files,
		Src:         src,
		SrcDirFiles: dest,
	}

	defer s.CleanReport(rd)

	svc := prj.Service
	if svc.Use || rd.Data.SessionInfo.UserInfo.Empty() {
		creds := userinfo.NewCredintials(svc.Credintials.UserName, svc.Credintials.Password)
		svcUserInfo := userinfo.NewService(svc.Url, svc.IimeOut, creds, s.Logger)
		s.AddUserInfo(svcUserInfo, prj, rd)
	}

	return rd, nil

}

func (s Service) CleanReport(rd *models.RepportData) error {

	const op = "redirect.clean"

	if !s.Config.DeleteTempFiles {
		return nil
	}

	s.Logger.Debug("Очистка",
		s.Logger.Str("src", rd.Src))

	if err := os.Remove(rd.Src); err != nil {
		s.Logger.Error("Ошибка удаления файла",
			s.Logger.Op(op),
			s.Logger.Str("Src", rd.Src),
			s.Logger.Err(err))
		return err
	}

	if err := os.RemoveAll(rd.SrcDirFiles); err != nil {
		s.Logger.Error("Ошибка удаления каталога",
			s.Logger.Op(op),
			s.Logger.Str("SrcDirFiles", rd.SrcDirFiles),
			s.Logger.Err(err))
		return err
	}

	return nil

}

func fileToData(src string, logger logging.Logger) (*models.Repport, error) {

	const op = "redirect.fileToData"

	byteValue, err := iofiles.FileToByte(src)
	if err != nil {
		logger.Error("Не удалось прочитать файл",
			logger.Op(op),
			logger.Str("name", src),
			logger.Err(err))
		return nil, err
	}

	var repport models.Repport
	err = json.Unmarshal(byteValue, &repport)
	if err != nil {
		logger.Error("Неверный формат файла",
			logger.Op(op),
			logger.Str("name", src),
			logger.Err(err))
		return nil, err
	}

	return &repport, nil

}

func listFiles(dest string, logger logging.Logger) ([]models.FileData, error) {

	const op = "redirect.listFiles"

	filesDir, err := os.ReadDir(dest)
	if err != nil {
		logger.Error("ReadDir",
			logger.Op(op),
			logger.Str("dest", dest),
			logger.Err(err))
		return nil, err
	}

	var files []models.FileData

	for _, f := range filesDir {

		if f.IsDir() {
			continue
		}

		fileName := f.Name()
		if fileName == reportJson {
			continue
		}

		src := filepath.Join(dest, fileName)
		byteValue, err := iofiles.FileToByte(src)

		if err != nil {
			logger.Error("Не удалось прочитать файл",
				logger.Op(op),
				logger.Str("name", src),
				logger.Err(err))
			return nil, err
		}

		files = append(files, models.FileData{
			Name: fileName,
			Data: byteValue,
		})

		logger.Debug("Получен файл",
			logger.Str("name", src))

	}

	return files, nil

}
