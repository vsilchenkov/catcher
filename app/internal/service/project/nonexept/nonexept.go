package nonexept

import (
	"catcher/pkg/logging"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/cockroachdb/errors"
)

const method = "nonexception"

var ErrBadRequest = errors.New("bad request")

type Response struct {
	Result    bool     `json:"result" binding:"required"`
	Exeptions []string `json:"exeptions"`
}

type Service struct {
	url         string
	timeOut     int
	credintials Credintials
	logger      logging.Logger
}

type Credintials struct {
	userName string
	password string
}

func NewService(url string, timeOut int, creds Credintials, logger logging.Logger) Service {
	return Service{
		url:         url,
		timeOut:     timeOut,
		credintials: creds,
		logger:      logger,
	}
}

func NewCredintials(userName, password string) Credintials {
	return Credintials{
		userName: userName,
		password: password,
	}
}

func (s Service) Get(ctx context.Context) ([]string, error) {

	const op = "nonexept.Get"

	ctxTimeOut, cancel := context.WithTimeout(ctx, time.Duration(s.timeOut)*time.Second)
	defer cancel()

	path, _ := url.JoinPath(s.url, method)
	req, err := http.NewRequestWithContext(ctxTimeOut, http.MethodGet, path, nil)
	if err != nil {
		s.logger.Error("Обшибка создяния контекста nonexept",
			s.logger.Err(err),
			s.logger.Op(op))
		return nil, errors.New("nonexept.http.NewRequestWithContext")
	}

	req.SetBasicAuth(s.credintials.userName, s.credintials.password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("Обшибка получения nonexept",
			s.logger.Err(err),

			s.logger.Op(op))
		return nil, errors.New("nonexept.http.DefaultClient.Do")
	}

	defer func() {
		errBodyClose := resp.Body.Close()
		if errBodyClose != nil {
			s.logger.Error("ошибка закрытия response body",
				s.logger.Err(errBodyClose))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Обшибка чтения nonexept",
			s.logger.Err(err),
			s.logger.Op(op))
		return nil, errors.New("nonexept.io.ReadAll")
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Ошибка получения nonexept",
			s.logger.Str("url", s.url),
			s.logger.Str("status", resp.Status),
			slog.Int("status_code", resp.StatusCode),
			s.logger.Str("body", string(body)),
			s.logger.Op(op))
		return nil, ErrBadRequest
	}

	// Парсим ответ
	var res Response
	if err = json.Unmarshal(body, &res); err != nil {
		s.logger.Error("Ошибка парсинга nonexept",
			s.logger.Err(err),
			s.logger.Op(op))
		return nil, errors.New("nonexept.json.Unmarshal")
	}

	s.logger.Debug("Получен список пропускаемых exeptions",
		s.logger.Str("url", s.url),
		s.logger.Op(op))

	return res.Exeptions, nil
}
