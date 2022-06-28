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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"time"
)

type config struct {
	Port       int    `envconfig:"HTTP_PORT" default:"8080"`
	JaegerHost string `envconfig:"JAEGER_HOST" default:"http://localhost"`
	JaegerPort int    `envconfig:"JAEGER_PORT" default:"14268"`
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
	logger := logrus.New()

	var conf config
	if err := envconfig.Process("", &conf); err != nil {
		return fmt.Errorf("error while processing env variables: %w", err)
	}

	tp, err := traceProvider(conf)
	if err != nil {
		return err
	}

	otel.SetTracerProvider(tp)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			logger.Errorf(fmt.Errorf("error while shutdowning OT: %w", err).Error())
		}
	}(ctx)

	swagger, err := api.GetSwagger()
	if err != nil {
		return fmt.Errorf("error while getting swagger documentation: %w", err)
	}
	a.e.Use(middleware.OapiRequestValidator(swagger))

	cStore, err := store.NewCatalogStore()
	if err != nil {
		return fmt.Errorf("error while create store: %w", err)
	}

	api.RegisterHandlers(a.e, handler.NewHandler(logger, cStore))

	return a.e.Start(fmt.Sprintf(":%d", conf.Port))
}

// Shutdown performs graceful shutdown
func (a *App) Shutdown() error {
	return a.e.Shutdown(context.Background())
}

func traceProvider(conf config) (*trace.TracerProvider, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(fmt.Sprintf("%s:%d/api/traces", conf.JaegerHost, conf.JaegerPort))))
	if err != nil {
		return nil, fmt.Errorf("error while creating OTel exporter: %w", err)
	}
	return trace.NewTracerProvider(
			trace.WithBatcher(exporter),
			trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("catalog")))),
		nil
}
