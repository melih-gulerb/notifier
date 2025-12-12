package handlers

import (
	"notifier/clients"
	"notifier/logger"
	"notifier/models"
	"notifier/templates"

	"github.com/labstack/echo/v4"
)

type MailHandler struct {
	mailClient *clients.MailClient
	logger     *logger.Logger
}

func NewMailHandler(mailClient *clients.MailClient, logger *logger.Logger) *MailHandler {
	return &MailHandler{
		mailClient: mailClient,
		logger:     logger,
	}
}

func (h *MailHandler) SendEmail(c echo.Context) error {
	var req models.SendEmailRequest

	correlationId := c.Request().Header.Get("correlationId")

	h.logger.Info(c.Request().Context()).WithCorrelationId(correlationId).Log("Send mail request received")

	if err := c.Bind(&req); err != nil {
		h.logger.Error(c.Request().Context()).WithError(err).WithCorrelationId(correlationId).Log("Failed to bind request")
		return c.JSON(400, models.Response{
			Success: false,
			Message: models.ServerError,
			Data:    nil,
		})
	}

	subject, htmlContent, err := templates.RenderTemplate(req.TemplateCode, req.Data)
	if err != nil {
		h.logger.Error(c.Request().Context()).WithError(err).WithCorrelationId(correlationId).Log("Failed to render template")
		return c.JSON(400, models.Response{
			Success: false,
			Message: models.ServerError,
			Data:    nil,
		})
	}

	h.mailClient.SendEmail(req.To, subject, htmlContent)

	h.logger.Info(c.Request().Context()).WithCorrelationId(correlationId).WithLogData("to", req.To).WithLogData("subject", subject).Log("Mail successfully sent")
	return c.JSON(200, models.Response{
		Success: true,
		Message: models.Success,
		Data:    nil,
	})
}
