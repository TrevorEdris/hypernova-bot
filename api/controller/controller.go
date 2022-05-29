package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/TrevorEdris/hypernova-bot/api/context"
	"github.com/TrevorEdris/hypernova-bot/api/services"
)

type (
	// Controller stores the Container, providing the route handlers
	// with access to all necessary dependencies.
	Controller struct {
		// Container provides route handlers with access to all services.
		Container *services.Container
	}
)

// NewController creates a new Controller.
func NewController(c *services.Container) Controller {
	return Controller{
		Container: c,
	}
}

func (c *Controller) Redirect(ctx echo.Context, route string, routeParams ...interface{}) error {
	url := ctx.Echo().Reverse(route, routeParams)
	return ctx.Redirect(http.StatusFound, url)
}

func (c *Controller) RenderErrorResponse(ctx echo.Context, status int, err error, msg string) error {
	if context.IsCanceledError(err) {
		return nil
	}
	ctx.Logger().Errorf("%s: %v", msg, err)
	return ctx.JSON(status, errResponse{msg})
}

func (c *Controller) RenderJSONResponse(ctx echo.Context, resp JSONResponse) error {
	ctx.Response().Status = resp.StatusCode
	for k, v := range resp.Headers {
		ctx.Response().Header().Set(k, v)
	}

	return ctx.JSON(resp.StatusCode, resp.Body)
}
