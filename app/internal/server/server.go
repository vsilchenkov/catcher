package server

import (
	"catcher/app/internal/handler"
	"catcher/app/internal/lib/logging"
	"catcher/app/internal/models"
	"catcher/app/internal/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	svc "github.com/kardianos/service"
	"github.com/cockroachdb/errors"
)

var errServerClosed = errors.New("http.ListenAndServe: http: Server closed")

type Program struct {
	models.AppContext
	Srv Srv
}

type Srv interface {
	Run(port string, handler http.Handler) error
	Shutdown(ctx context.Context) error
}

func NewProgram(srv Srv, appCtx models.AppContext) *Program {
	return &Program{Srv: srv,
		AppContext: appCtx}
}

func (p *Program) Start(s svc.Service) error {
	
	// Запускаем в отдельной горутине, чтобы не блокировать Start
	i := "Starting a web-server on port"
	port := p.Config.Server.Port

	p.Logger.Info(i,
		p.Logger.Str("port", port),
		p.Logger.Str("version", p.Config.Version))

	if p.Config.Log.OutputInFile {
		fmt.Printf("%s: %s\n", i, port)
	}

	go p.Run()
	return nil
}

func (p *Program) Run() {

	appCtx := p.AppContext
	service := service.NewService(appCtx)
	handlers := handler.NewHandler(service, appCtx)

	logger := logging.GetLogger()
	port := appCtx.Config.ServerPort()

	err := os.MkdirAll(filepath.Join(appCtx.Config.WorkingDir, models.DirWeb, models.DirTemp), 0755)
	if err != nil {
		logger.Error("Ошибка создание временных каталогов",
			logger.Err(err))
		os.Exit(1)
	}

	if err := p.Srv.Run(port, handlers.InitHandler()); err != nil && err.Error() != errServerClosed.Error() {
		logger.Error("Error running server",
			logger.Err(err))
		os.Exit(1)
	}

}

func (p *Program) Stop(s svc.Service) error {

	appCtx := p.AppContext
	appCtx.Logger.Debug("Server Shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := p.Srv.Shutdown(ctx); err != nil {
		appCtx.Logger.Error("Error on server shutting down",
			appCtx.Logger.Err(err))
		return err
	}

	i := "Server is stopped"
	appCtx.Logger.Info(i)
	if p.Config.Log.OutputInFile {
		fmt.Printf("%s\n", i)
	}

	return nil

}
