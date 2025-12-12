package routes

import (
	"notifier/clients"
	"notifier/configs"
	"notifier/handlers"
	"notifier/logger"

	"github.com/labstack/echo/v4"
)

func SetupMailRoutes(app *echo.Echo, logger *logger.Logger) {
	mailClient := clients.NewMailClient(configs.InitBrevo(configs.GetConfig().BrevoAPIKey), configs.GetConfig().FromEmail)

	mailHandler := handlers.NewMailHandler(mailClient, logger)
	mailGroup := app.Group("/mail")

	mailGroup.POST("/send", mailHandler.SendEmail)
}
