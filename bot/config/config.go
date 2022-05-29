package config

import (
	"os"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type (
	// Config is the aggregation of all necessary configurations.
	Config struct {
		App
		Metrics
		Discord
	}

	// App defines the configs needed for the application itself.
	App struct {
		Name        string `env:"APP_NAME,default=hypernova-bot"`
		Environment string `env:"APP_ENVIRONMENT,default=local"`
	}

	// Metrics defines the configs needed for collecting metrics.
	Metrics struct {
		Enabled bool   `env:"METRICS_ENABLED,default=false"`
		Addr    string `env:"METRICS_ADDRESS"`
		BufLen  int    `env:"METRICS_BUFFER,default=5"`
	}

	// Discord defines the configs needed for interacting with the discord API.
	Discord struct {
		Token string `env:"DISCORD_BOT_TOKEN,required"`
	}
)

func New() (*Config, error) {
	var cfg Config
	err := godotenv.Load()

	// If a .env file exists but was unable to be loaded
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	err = envdecode.StrictDecode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
