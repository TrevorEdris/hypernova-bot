package services

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/TrevorEdris/hypernova-bot/bot/config"
)

type (
	// Container contains all services used by the bot and provides an easy way to
	// handle dependency injection.
	Container struct {
		// Config stores the application configuration.
		Config *config.Config

		// Logger provies access to logging functions.
		Logger *zap.Logger

		// DiscordSession provides access to the underlying discordgo library.
		DiscordSession *discordgo.Session
	}
)

func NewContainer() (*Container, error) {
	c := new(Container)
	err := c.initConfig()
	if err != nil {
		return nil, err
	}

	err = c.initLogger()
	if err != nil {
		return nil, err
	}

	err = c.initDiscord()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Container) initConfig() error {
	conf, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to init config: %w", err)
	}

	c.Config = conf
	return nil
}

func (c *Container) initLogger() error {
	l, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	c.Logger = l
	return nil
}

func (c *Container) initDiscord() error {
	dg, err := discordgo.New("Bot " + c.Config.Discord.Token)
	if err != nil {
		return fmt.Errorf("failed to init discord: %w", err)
	}

	c.DiscordSession = dg
	return nil
}

func (c *Container) Shutdown() error {
	err := c.DiscordSession.Close()
	if err != nil {
		return fmt.Errorf("failed to cleanly close discord session: %w", err)
	}

	err = c.Logger.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync logger: %w", err)
	}

	return nil
}
