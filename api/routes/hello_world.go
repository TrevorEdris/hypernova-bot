package routes

import (
	"net/http"

	"github.com/TrevorEdris/hypernova-bot/api/controller"
	"github.com/labstack/echo/v4"
)

type (
	Hello struct {
		controller.Controller
	}

	helloJSONResponse struct {
		Message string `json:"message"`
	}
)

func (c *Hello) Get(ctx echo.Context) error {
	resp := controller.NewJSONResponse(ctx)
	resp.StatusCode = http.StatusOK
	resp.Body = helloJSONResponse{"Hello world!"}
	return c.RenderJSONResponse(ctx, resp)
}
