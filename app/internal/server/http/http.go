package http

import (
	"catcher/app/internal/models"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
)

type Server struct {
	models.AppContext
	httpServer *http.Server
	stop       chan struct{}
	stopOnce   sync.Once
}

func New(appCtx models.AppContext) *Server {
	return &Server{AppContext: appCtx}
}

func (s *Server) Run(port string, handler http.Handler) error {

	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	s.stop = make(chan struct{})
	errCh := make(chan error)

	go func() {

		defer func() {
			if r := recover(); r != nil {				
				errCh <- fmt.Errorf("panic recovered: %v", r)
			}
		}()

		if err := s.httpServer.ListenAndServe(); err != nil {
			errCh <- errors.Wrap(err, "http.ListenAndServe")
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-s.stop:
		return nil
	}

}

func (s *Server) Shutdown(ctx context.Context) error {

	defer s.stopOnce.Do(func() {
		close(s.stop)
	})

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		err = errors.Wrap(err, "http.Shutdown")
	}

	return err
}
