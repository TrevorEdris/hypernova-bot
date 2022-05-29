package middleware

import (
	"mime"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	// ContentType is a type alias for a string.
	ContentType string
)

const (
	// JSON represents the application/json content type.
	JSON ContentType = "application/json"
)

// EncforceContentType is a middleware function that will reject requests which
// do not include the specified content type header.
func EnforceContentType(ctype ContentType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contentType := c.Request().Header.Get("Content-Type")
			t, _, err := mime.ParseMediaType(contentType)
			if err != nil || t != string(ctype) {
				return echo.NewHTTPError(http.StatusUnsupportedMediaType)
			}
			return next(c)
		}
	}
}
