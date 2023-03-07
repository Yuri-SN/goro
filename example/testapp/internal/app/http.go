package app

import (
	"fmt"
	"github.com/hanagantig/gracy"
	"go.uber.org/zap"
	"testapp/internal/handler/http"
	"testapp/internal/handler/http/api/v1"
)

func (a *App) StartHTTPServer() error {
	go func() {
		a.startHTTPServer()
	}()

	err := gracy.Wait()
	if err != nil {
		a.logger.Error("failed to gracefully shutdown server", zap.Error(err))
		return err
	}
	a.logger.Info("server gracefully stopped")
	return nil
}

func (a *App) startHTTPServer() {
	handler := v1.NewHandler(a.c.GetUseCase(), a.logger)

	router := http.NewRouter()
	router.
		//WithMetrics().
		//WithHealthChecks(app.hc).
		WithSwagger().
		WithHandler(handler, a.logger).
		WithProfiler()

	srv := http.NewServer(a.cfg.HTTP)
	srv.RegisterRoutes(router)

	gracy.AddCallback(func() error {
		return srv.Stop()
	})

	a.logger.Info(fmt.Sprintf("starting HTTP server at %s:%s", a.cfg.HTTP.Host, a.cfg.HTTP.Port))
	err := srv.Start()
	if err != nil {
		a.logger.Fatal("Fail to start %s http server:", zap.String("app", a.cfg.App.Name), zap.Error(err))
	}
}