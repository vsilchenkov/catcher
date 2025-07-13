package iofiles

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
)

func SaveDataToFile(data []byte, src string) error {

	return os.WriteFile(src, data, 0644)

}

func FileToByte(src string) ([]byte, error) {

	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	byteValue, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

func Unzip(src string, dest string) error {

	rcloser, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer rcloser.Close()

	for _, f := range rcloser.File {

		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return errors.Wrap(err, "Не допустимый путь файла")
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
