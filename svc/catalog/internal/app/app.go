package app

import (
	"context"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/kelseyhightower/envconfig"
	"github.com/konrad945/eCommerce/svc/catalog/api"
	"github.com/konrad945/eCommerce/svc/catalog/internal/handler"
	"github.com/konrad945/eCommerce/svc/catalog/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type config struct {
	Port int `envconfig:"HTTP_PORT" default:"8080"`
}

type App struct {
	e *echo.Echo
}

// NewApp setups an App struct
func NewApp() *App {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return &App{
		e: e,
	}
}

// Run initialize handler and starts serving an app
func (a *App) Run() error {
	var conf config
	if err := envconfig.Process("", &conf); err != nil {
		return fmt.Errorf("error while processing env variables: %w", err)
	}

	swagger, err := api.GetSwagger()
	if err != nil {
		return fmt.Errorf("error while getting swagger documentation: %w", err)
	}
	a.e.Use(middleware.OapiRequestValidator(swagger))

	cStore, err := store.NewCatalogStore()
	if err != nil {
		return fmt.Errorf("error while create store: %w", err)
	}

	logger := logrus.New()

	api.RegisterHandlers(a.e, handler.NewHandler(logger, cStore))

	return a.e.Start(fmt.Sprintf(":%d", conf.Port))
}

// Shutdown performs graceful shutdown
func (a *App) Shutdown() error {
	return a.e.Shutdown(context.Background())
}
