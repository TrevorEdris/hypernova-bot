package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"github.com/TrevorEdris/hypernova-bot/app/controller"
	"github.com/TrevorEdris/hypernova-bot/app/middleware"
	"github.com/TrevorEdris/hypernova-bot/app/services"
)

// BuildRouter builds the HTTP router for all endpoint handlers.
func BuildRouter(c *services.Container) {
	g := c.Web.Group("")

	// Force HTTPS if enabled
	if c.Config.HTTP.TLS.Enabled {
		g.Use(echomw.HTTPSRedirect())
	}

	// Use various other middleware functions here
	g.Use(
		echomw.RemoveTrailingSlashWithConfig(echomw.TrailingSlashConfig{
			RedirectCode: http.StatusMovedPermanently,
		}),
		echomw.Recover(),
		echomw.Secure(),
		echomw.RequestID(),
		echomw.Gzip(),
		middleware.LogRequestID(),
		echomw.TimeoutWithConfig(echomw.TimeoutConfig{
			Timeout: c.Config.App.Timeout,
		}),
	)

	ctr := controller.NewController(c)

	defaultRoutes(c, g, ctr)

	// Create a group of routes where the json content type is enforced
	itemGroup := g.Group("/item", middleware.EnforceContentType(middleware.JSON))
	itemRoutes(c, itemGroup, ctr)
}

func defaultRoutes(c *services.Container, g *echo.Group, ctr controller.Controller) {
	hello := Hello{Controller: ctr}
	g.GET("/", hello.Get).Name = "helloworld"
}

func itemRoutes(c *services.Container, g *echo.Group, ctr controller.Controller) {
	item := Item{Controller: ctr}
	g.GET("/:id", item.Get).Name = "itemget"
	g.POST("", item.Post).Name = "itempost"
	g.PUT("/:id", item.Put).Name = "itemput"
	g.DELETE("/:id", item.Delete).Name = "itemdelete"
}
