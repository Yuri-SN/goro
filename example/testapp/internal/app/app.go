// Code generated by goro; DO NOT EDIT.

package app

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"sync"

	"go.uber.org/zap/zapcore"
	"testapp/internal/config"
)

// This file was generated by the goro tool.
// Editing this file might prove futile when you re-run the goro generate command

type Logger interface {
	Debug(string, ...zapcore.Field)
	Info(string, ...zapcore.Field)
	Error(string, ...zapcore.Field)
	Fatal(string, ...zapcore.Field)
}

type App struct {
	cfg config.Config

	c     *Container
	cOnce *sync.Once

	//hc     health.Checker
	//hcOnce *sync.Once

	mysql    *sql.DB
	mysqlx   *sqlx.DB
	postgres *sql.DB

	logger Logger
}

var a *App

func NewApp(configPath string) (*App, error) {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return nil, err
	}

	app := &App{
		cOnce: &sync.Once{},
		//hcOnce: &sync.Once{},
		cfg: cfg,
	}

	//goro:init logger
	app.initLogger()

	//goro:init healthChecker
	//app.initHealthChecker()

	mysqlConnect, err := app.newMySQLConnect(cfg.MainDB)
	if err != nil {
		return nil, err
	}
	app.mysql = mysqlConnect
	mysqlxConn, err := app.newMySQLxConnect(cfg.MainDB)
	if err != nil {
		return nil, err
	}
	app.mysqlx = mysqlxConn
	postgresConn, err := app.newPostgresConnect(cfg.MainDB)
	if err != nil {
		return nil, err
	}
	app.postgres = postgresConn

	//goro:init dependencies
	app.c = NewContainer(app.mysql, app.mysqlx, app.postgres)

	return app, nil
}

func SetGlobalApp(app *App) {
	a = app
}

func GetGlobalApp() (*App, error) {
	if a == nil {
		return nil, errors.New("global app is not initialized")
	}

	return a, nil
}
