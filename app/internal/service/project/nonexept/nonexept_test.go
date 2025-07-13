package nonexept

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"catcher/app/internal/testutil/logging"

	"github.com/stretchr/testify/assert"
)

// Хелпер для создания сервиса с тестовым http.Client
func newTestService(serverURL string) Service {
	return NewService(
		serverURL,
		1, // timeout
		NewCredintials("user", "pass"),
		&logging.TestLogger{},
	)
}

func TestService_Get_Success(t *testing.T) {
	// Мокаем сервер
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/nonexception", r.URL.Path)
		user, pass, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "user", user)
		assert.Equal(t, "pass", pass)

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"result":true,"exeptions":["A","B"]}`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	svc := newTestService(server.URL)
	ctx := context.Background()

	exeptions, err := svc.Get(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []string{"A", "B"}, exeptions)
}

func TestService_Get_BadStatus(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "bad request")
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	svc := newTestService(server.URL)
	ctx := context.Background()

	exeptions, err := svc.Get(ctx)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Nil(t, exeptions)
}

func TestService_Get_InvalidJSON(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "not a json")
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	svc := newTestService(server.URL)
	ctx := context.Background()

	exeptions, err := svc.Get(ctx)
	assert.Error(t, err)
	assert.Nil(t, exeptions)
}

func TestService_Get_Timeout(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"result":true,"exeptions":[]}`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	svc := NewService(server.URL, 1, NewCredintials("user", "pass"), &logging.TestLogger{})
	ctx := context.Background()

	exeptions, err := svc.Get(ctx)
	assert.Error(t, err)
	assert.Nil(t, exeptions)
}
