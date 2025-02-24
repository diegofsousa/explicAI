package configuration

import (
	"context"
	"net/http"

	"github.com/diegofsousa/explicAI/internal/application/service"
	"github.com/diegofsousa/explicAI/internal/infrastructure/api"
	"github.com/diegofsousa/explicAI/internal/infrastructure/db"
	"github.com/diegofsousa/explicAI/internal/infrastructure/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Application struct {
	server  *echo.Echo
	config  *viper.Viper
	clients *Clients
}

func NewApplication(config *viper.Viper, clients *Clients) *Application {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	logger := log.StartLog()
	initMiddlewares(server, logger)

	return &Application{
		server:  server,
		config:  config,
		clients: clients,
	}
}

func (a *Application) Start() {
	a.registerControllers()

	ctx := context.Background()

	host := a.config.GetString("server.host")

	log.LogInfo(ctx, a.config.GetString("app.name")+" is starting on "+host+"...")
	log.LogError(ctx, "server fatal error", a.server.Start(host))
}

func (a *Application) registerControllers() {
	summary := service.NewSummary(
		a.clients.AudioTranscript,
		a.clients.Summarize,
		db.NewSummary(a.config.GetString("database.url")),
	)

	api.NewExplicaServer(summary).Register(a.server)
}

func initMiddlewares(server *echo.Echo, logger *zap.Logger) {
	server.Use(middleware.Recover())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	server.Use(loggerMiddleware(logger))
}

func loggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := uuid.NewString()

			ctxLogger := logger.With(
				zap.String("requestId", requestID),
				zap.String("remoteIp", c.RealIP()),
				zap.String("path", c.Path()),
			)

			ctx := context.WithValue(c.Request().Context(), "logger", ctxLogger)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
