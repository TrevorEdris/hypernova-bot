package main

import (
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/TrevorEdris/hypernova-bot/bot/controller"
	"github.com/TrevorEdris/hypernova-bot/bot/services"
)

func main() {
	container, err := services.NewContainer()
	if err != nil {
		panic(err)
	}

	controller := controller.NewController(container)

	// After main() finishes execution, run this function
	defer func() {
		err := container.Shutdown()
		if err != nil {
			container.Logger.Sugar().Fatalf("Failed to cleanly shutdown: %v", err)
		}
	}()

	// Run the bot
	errChan := make(chan error)
	go controller.Run(errChan)

	// Wait for the interrupt signal in order to gracefully shutdown the server
	// with a configurable timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	for {
		select {
		case err := <-errChan:
			container.Logger.Error("Encountered unexpected error from c.Run", zap.Error(err))
		case <-quit:
			container.Logger.Info("Received exit signal; Attempting to cleanly shut down")
		}
	}
}
