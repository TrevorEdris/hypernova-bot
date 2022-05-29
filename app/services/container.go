package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"

	"github.com/TrevorEdris/hypernova-bot/app/config"
	"github.com/TrevorEdris/hypernova-bot/app/internal/repository"
)

type (
	// Container contains all services used by the application and provides an easy way to
	// handle dependency injection.
	Container struct {
		// Config stores the application configuration.
		Config *config.Config

		// Validator stores the validator.
		Validator *Validator

		// Web stores the API framework.
		Web *echo.Echo

		// ItemRepo provides access to the Item storage medium.
		ItemRepo repository.ItemRepo
	}
)

func NewContainer() *Container {
	c := new(Container)
	c.initConfig()
	c.initValidator()
	c.initWeb()
	c.initItemRepo()
	return c
}

func (c *Container) Shutdown() error {
	return nil
}

func (c *Container) initConfig() {
	cfg, err := config.New()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	c.Config = &cfg
}

func (c *Container) initValidator() {
	c.Validator = NewValidator()
}

func (c *Container) initWeb() {
	c.Web = echo.New()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("failed to create new zap logger: %w", err))
	}
	c.Web.Use(echozap.ZapLogger(zapLogger))

	// Configure the logger for the web framework
	switch c.Config.App.LogLevel {
	case config.LvlDbg:
		c.Web.Logger.SetLevel(log.DEBUG)
	case config.LvlInfo:
		c.Web.Logger.SetLevel(log.INFO)
	case config.LvlWarn:
		c.Web.Logger.SetLevel(log.WARN)
	case config.LvlErr:
		c.Web.Logger.SetLevel(log.ERROR)
	default:
		c.Web.Logger.SetLevel(log.DEBUG)
	}

	c.Web.Validator = c.Validator
}

func (c *Container) initItemRepo() {
	switch c.Config.App.Storage {
	case config.StorageLocal:
		c.Web.Logger.Info("Configured for local storage")
		c.ItemRepo = repository.NewItemRepoLocal()
	case config.StorageDynamoDB:
		c.Web.Logger.Info("Configured for DynamoDB storage")
		c.ItemRepo = repository.NewItemRepoDynamoDB(c.Config, dynamodb.NewFromConfig(c.Config.AWSCfg))
	default:
		c.Web.Logger.Warnf("Invalid app storage (%s); defaulting to local storage", c.Config.App.Storage)
		c.ItemRepo = repository.NewItemRepoLocal()
	}
}
