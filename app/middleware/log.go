package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const (
	logHeaderFormat = `{"time":"${time.rfc3339_nano}","id":"%s","level":"${level}","prefix":"${prefix}","file":"$[short_file}","line":${line}"}`
)

// LogRequestID includes the request ID in all logs for the given request.
// This requires that the middleware that includes the request ID be executed
// before this middleware func.
func LogRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rID := c.Response().Header().Get(echo.HeaderXRequestID)
			c.Logger().SetHeader(fmt.Sprintf(logHeaderFormat, rID))
			return next(c)
		}
	}
}
