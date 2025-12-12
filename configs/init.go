package configs

import (
	"log"
	"notifier/logger"
	"notifier/middlewares"
	"os"

	"github.com/getbrevo/brevo-go/lib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func InitEcho(logger *newrelic.Application) *echo.Echo {
	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middlewares.NewRelicMiddleware(logger))

	app.GET("/health", func(c echo.Context) error {
		return c.JSON(200, "its ok")
	})

	return app
}

func InitBrevo(apiKey string) *lib.APIClient {
	brevoCfg := lib.NewConfiguration()
	brevoCfg.AddDefaultHeader("api-key", apiKey)
	return lib.NewAPIClient(brevoCfg)
}

func InitLogging() (*logger.Logger, *newrelic.Application) {
	nrApp, err := logger.InitNewRelic(
		os.Getenv("NEW_RELIC_APP_NAME"),
		os.Getenv("NEW_RELIC_LICENSE_KEY"),
	)
	if err != nil {
		log.Fatal("New Relic couldn't got start:", err)
	}

	appLogger := logger.NewLogger(nrApp)
	appLogger.SetLevel("info")

	return appLogger, nrApp
}
