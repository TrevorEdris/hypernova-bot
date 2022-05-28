package config

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type (
	Environment string
	LogLevel    string
	Storage     string
)

const (
	EnvLocal   Environment = "local"
	EnvTest    Environment = "test"
	EnvDev     Environment = "dev"
	EnvStaging Environment = "staging"
	EnvLT      Environment = "lt"
	EnvQA      Environment = "qa"
	EnvProd    Environment = "prod"

	StorageLocal    Storage = "local"
	StorageDynamoDB Storage = "dynamodb"

	LvlDbg  LogLevel = "debug"
	LvlInfo LogLevel = "info"
	LvlWarn LogLevel = "warn"
	LvlErr  LogLevel = "error"
)

// SwitchEnvironment sets the environment variable used to dictate which environment the application is
// currently running in. This must be called prior to loading the configuration in order for it
// to take effect.
func SwitchEnvironment(env Environment) {
	if err := os.Setenv("APP_ENVIRONMENT", string(env)); err != nil {
		panic(err)
	}
}

type (
	// Config is the aggregation of all necessary configurations.
	Config struct {
		App
		HTTP
		Metrics
		AWS
		DynamoDB
	}

	// App defines the configs needed for the application itself.
	App struct {
		Name        string        `env:"APP_NAME,default=backfill"`
		Environment Environment   `env:"APP_ENVIRONMENT,default=local"`
		LogLevel    LogLevel      `env:"LOG_LEVEL,default=info"`
		Timeout     time.Duration `env:"APP_TIMEOUT,default=20s"`
		Storage     Storage       `env:"APP_STORAGE,default=local"`
	}

	// HTTP stores the configuration for the HTTP server.
	HTTP struct {
		Hostname     string        `env:"HTTP_HOSTNAME,default=0.0.0.0"`
		Port         uint16        `env:"HTTP_PORT,default=8000"`
		ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT,default=5s"`
		WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT,default=10s"`
		IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT,default=2m"`
		TLS          struct {
			Enabled     bool   `env:"HTTP_TLS_ENABLED,default=false"`
			Certificate string `env:"HTTP_TLS_CERTIFICATE"`
			Key         string `env:"HTTP_TLS_KEY"`
		}
	}

	// Metrics defines the configs needed for collecting metrics.
	Metrics struct {
		Enabled bool   `env:"METRICS_ENABLED,default=false"`
		Addr    string `env:"METRICS_ADDRESS"`
		BufLen  int    `env:"METRICS_BUFFER,default=5"`
	}

	// AWS defines the configs related to AWS services.
	AWS struct {
		AccessKeyID string `env:"AWS_ACCESS_KEY_ID"`
		Secret      string `env:"AWS_SECRET_ACCESS_KEY"`
		Region      string `env:"AWS_REGION"`
		Endpoint    string `env:"AWS_ENDPOINT"`
		AWSCfg      aws.Config
	}

	// DynamoDB defines the configs related specfically to the DynamoDB service.
	DynamoDB struct {
		ItemTable string `env:"DYNAMODB_ITEM_TABLE"`
	}
)

// New loads the configuration based on the environment variables.
func New() (Config, error) {
	var cfg Config
	err := godotenv.Load()

	// If a .env file exists but was unable to be loaded
	if err != nil && !os.IsNotExist(err) {
		return Config{}, err
	}

	err = envdecode.StrictDecode(&cfg)
	if err != nil {
		return Config{}, err
	}

	cfg.AWS.AWSCfg, err = loadAWSCfg(context.Background(), cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func loadAWSCfg(ctx context.Context, cfg Config) (aws.Config, error) {
	customResolver := aws.EndpointResolverWithOptions(
		aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if cfg.AWS.Endpoint != "" {
					return aws.Endpoint{
						URL:           cfg.AWS.Endpoint,
						SigningRegion: region,
						Source:        aws.EndpointSourceCustom,
					}, nil
				}
				// returning EndpointNotFoundError will allow the service to fallback to its default resolution
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			},
		),
	)
	return awsconfig.LoadDefaultConfig(ctx, awsconfig.WithEndpointResolverWithOptions(customResolver))
}
