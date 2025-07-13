package userinfo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"catcher/app/internal/lib/times"
	"catcher/app/internal/testutil/logging"

	"github.com/stretchr/testify/assert"
)

func TestService_Get_Success(t *testing.T) {
	// Подготовка тестового ответа
	resp := Response{
		Result: true,
	}
	resp.UserInfo.Id = "123"
	resp.UserInfo.City = "Москва"
	resp.UserInfo.Branch = "Центральный"
	resp.UserInfo.Position = "Инженер"
	resp.UserInfo.Started = times.TimeTZ{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}
	resp.UserInfo.SessionInfo.IP = "127.0.0.1"
	resp.UserInfo.SessionInfo.Device = "PC"
	resp.UserInfo.SessionInfo.Session = 42
	resp.UserInfo.SessionInfo.Connection = 1

	body, _ := json.Marshal(resp)

	// Имитация внешнего сервера
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/userInfo", r.URL.Path)
		assert.Equal(t, "testuser", r.URL.Query().Get("username"))
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer server.Close()

	creds := NewCredintials("user", "pass")
	svc := NewService(server.URL, 2, creds, &logging.TestLogger{})

	ctx := context.Background()
	info, err := svc.Get(ctx, "testuser")

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "123", info.Id)
	assert.Equal(t, "Москва", info.City)
	assert.Equal(t, "Центральный", info.Branch)
	assert.Equal(t, "Инженер", info.Position)
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), info.Started)
}

func TestService_Get_BadStatus(t *testing.T) {
	// Имитация ответа с ошибкой
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	creds := NewCredintials("user", "pass")
	svc := NewService(server.URL, 2, creds, &logging.TestLogger{})

	ctx := context.Background()
	info, err := svc.Get(ctx, "testuser")

	assert.Error(t, err)
	assert.Nil(t, info)
	assert.Equal(t, ErrBadRequest, err)
}

func TestService_Get_InvalidJSON(t *testing.T) {
	// Имитация некорректного JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	creds := NewCredintials("user", "pass")
	svc := NewService(server.URL, 2, creds, &logging.TestLogger{})

	ctx := context.Background()
	info, err := svc.Get(ctx, "testuser")

	assert.Error(t, err)
	assert.Nil(t, info)
}
