package middlewares

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func NewRelicMiddleware(app *newrelic.Application) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			txn := app.StartTransaction(c.Request().Method + " " + c.Path())
			defer txn.End()

			ctx := newrelic.NewContext(c.Request().Context(), txn)
			c.SetRequest(c.Request().WithContext(ctx))

			correlationId := c.Request().Header.Get("correlationId")
			if correlationId == "" {
				correlationId = uuid.New().String()
				c.Request().Header.Set("correlationId", correlationId)
			}

			txn.AddAttribute("correlationId", correlationId)
			txn.AddAttribute("method", c.Request().Method)
			txn.AddAttribute("path", c.Request().URL.Path)
			txn.AddAttribute("ip", c.Request().RemoteAddr)

			if userId := c.Request().Header.Get("userId"); userId != "" {
				txn.AddAttribute("userId", userId)
			}

			txn.SetWebResponse(c.Response().Writer)
			txn.SetWebRequestHTTP(c.Request())

			return next(c)
		}
	}
}
