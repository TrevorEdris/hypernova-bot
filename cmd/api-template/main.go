package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/TrevorEdris/api-template/app/routes"
	"github.com/TrevorEdris/api-template/app/services"
)

func main() {
	c := services.NewContainer()
	defer func() {
		err := c.Shutdown()
		if err != nil {
			c.Web.Logger.Fatalf("Error during shutdown: %v", err)
		}
	}()

	// Build the router, responsible for creating the various HTTP handlers
	routes.BuildRouter(c)

	// Run the HTTP server in a goroutine
	go func() {
		srv := http.Server{
			Addr:         fmt.Sprintf("%s:%d", c.Config.HTTP.Hostname, c.Config.HTTP.Port),
			Handler:      c.Web,
			ReadTimeout:  c.Config.HTTP.ReadTimeout,
			WriteTimeout: c.Config.HTTP.WriteTimeout,
			IdleTimeout:  c.Config.HTTP.IdleTimeout,
		}

		if c.Config.HTTP.TLS.Enabled {
			certs, err := tls.LoadX509KeyPair(c.Config.HTTP.TLS.Certificate, c.Config.HTTP.TLS.Key)
			if err != nil {
				c.Web.Logger.Fatalf("failed to load TLS certificate: %v", err)
			}

			srv.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{certs},
			}
		}

		c.Web.Logger.Info("Starting HTTP server")
		err := c.Web.StartServer(&srv)
		if err != nil {
			c.Web.Logger.Fatalf("failed to start server; shutting down: %v", err)
		}
	}()

	// Wait for the interrupt signal in order to gracefully shutdown the server
	// with a configurable timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), c.Config.App.Timeout)
	defer cancel()
	err := c.Web.Shutdown(ctx)
	if err != nil {
		c.Web.Logger.Fatal(err)
	}
}
