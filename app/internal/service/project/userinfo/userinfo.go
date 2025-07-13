package userinfo

import (
	"catcher/app/internal/lib/logging"
	"catcher/app/internal/lib/times"
	"catcher/app/internal/models"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/cockroachdb/errors"
)

const method = "userInfo"

var ErrBadRequest = errors.New("bad request")

type Response struct {
	Result   bool `json:"result" binding:"required"`
	UserInfo struct {
		Id          string       `json:"id"`
		City        string       `json:"city"`
		Branch      string       `json:"branch"`
		Position    string       `json:"position"`
		Started     times.TimeTZ `json:"started"`
		SessionInfo struct {
			IP         string `json:"IP"`
			Device     string `json:"device"`
			Session    int    `json:"session"`
			Connection int    `json:"connection"`
		} `json:"session_info"`
	} `json:"userinfo"`
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

func (s Service) Get(ctx context.Context, input string) (*models.UserInfo, error) {

	const op = "userinfo.Get"

	ctxTimeOut, cancel := context.WithTimeout(ctx, time.Duration(s.timeOut)*time.Second)
	defer cancel()

	path, _ := url.JoinPath(s.url, method)
	req, err := http.NewRequestWithContext(ctxTimeOut, http.MethodGet, path, nil)
	if err != nil {
		s.logger.Error("Обшибка создяния контекста userinfo",
			s.logger.Err(err),
			s.logger.Str("name", input),
			s.logger.Op(op))
		return nil, errors.New("userinfo.http.NewRequestWithContext")
	}

	req.SetBasicAuth(s.credintials.userName, s.credintials.password)
	
	q := req.URL.Query()
	q.Add("username", input)
	encoded := q.Encode()
	encoded = strings.ReplaceAll(encoded, "+", "%20") // Возвращаем назад пробелы
	req.URL.RawQuery = encoded
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("Обшибка получения userinfo",
			s.logger.Err(err),
			s.logger.Str("name", input),
			s.logger.Op(op))
		return nil, errors.New("userinfo.http.DefaultClient.Do")
	}

	defer func() {
		errBodyClose := resp.Body.Close()
		if errBodyClose != nil {
			s.logger.Error("ошибка закрытия response body",
				s.logger.Str("name", input),
				s.logger.Str("name", input),
				s.logger.Err(errBodyClose))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Обшибка чтения userinfo",
			s.logger.Err(err),
			s.logger.Str("name", input),
			s.logger.Op(op))
		return nil, errors.New("userinfo.io.ReadAll")
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Ошибка получения userinfo",
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
		s.logger.Error("Ошибка парсинга userinfo",
			s.logger.Err(err),
			s.logger.Str("name", input),
			s.logger.Op(op))
		return nil, errors.New("userinfo.json.Unmarshal")
	}

	var resInfo models.UserInfo
	if res.Result {
		copier.Copy(&resInfo, &res.UserInfo)
		resInfo.Started = res.UserInfo.Started.Time
	} else {
		s.logger.Warn("userinfo - нет данных пользователя",
			s.logger.Str("name", input),
			s.logger.Op(op))
	}

	return &resInfo, nil
}
