package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	JSONResponse struct {
		StatusCode int
		Headers    map[string]string
		RequestID  string
		Path       string
		URL        string
		Context    echo.Context
		ToURL      func(name string, params ...interface{}) string
		Body       interface{}
	}

	errResponse struct {
		Message string `json:"message"`
	}
)

func NewJSONResponse(ctx echo.Context) JSONResponse {
	return JSONResponse{
		Context:    ctx,
		ToURL:      ctx.Echo().Reverse,
		Path:       ctx.Request().URL.Path,
		URL:        ctx.Request().URL.String(),
		StatusCode: http.StatusOK,
		Headers:    make(map[string]string),
		RequestID:  ctx.Response().Header().Get(echo.HeaderXRequestID),
	}
}
