package configuration

import (
	"fmt"
	"log/slog"
	"microsoft-apps-exporter/internal/models"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type Configuration struct {
	LOG_LEVEL string

	GRAPH_CLIENT_ID     string
	GRAPH_TENANT_ID     string
	GRAPH_CLIENT_SECRET string
	GRAPH_APP_SCOPES    string

	Sharepoint *models.SharepointResource `mapstructure:"sharepoint"`

	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_DSN      string

	WEBHOOK_LISTEN_IP         string
	WEBHOOK_LISTEN_PORT       string
	WEBHOOK_EXTERNAL_BASE_URL string
}

var (
	config Configuration
	once   sync.Once
)

// GetConfig initializes and returns the singleton app configuration.
func GetConfig() Configuration {
	once.Do(func() {
		loadResourcesYaml()
		loadEnvVars()
		buildPostgresDSN()
		unmarshalYaml()
	})
	return config
}

// loadResourcesYaml reads YAML configuration from resources.yaml.
func loadResourcesYaml() {
	viper.SetConfigName("resources")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		slog.Debug("Failed to read resources.yaml", "error", err, "operation", "config")
	} else {
		slog.Debug("resources.yaml loaded", "file", viper.ConfigFileUsed(), "operation", "config")
	}
}

// loadEnvVars loads remaining configuration from environment variables.
func loadEnvVars() {
	config.LOG_LEVEL = os.Getenv("LOG_LEVEL")

	config.GRAPH_CLIENT_ID = os.Getenv("GRAPH_CLIENT_ID")
	config.GRAPH_TENANT_ID = os.Getenv("GRAPH_TENANT_ID")
	config.GRAPH_CLIENT_SECRET = os.Getenv("GRAPH_CLIENT_SECRET")
	config.GRAPH_APP_SCOPES = os.Getenv("GRAPH_APP_SCOPES")

	config.DB_HOST = os.Getenv("DB_HOST")
	config.DB_PORT = os.Getenv("DB_PORT")
	config.DB_USER = os.Getenv("DB_USER")
	config.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	config.DB_NAME = os.Getenv("DB_NAME")

	config.WEBHOOK_LISTEN_IP = os.Getenv("WEBHOOK_LISTEN_IP")
	config.WEBHOOK_LISTEN_PORT = os.Getenv("WEBHOOK_LISTEN_PORT")
	config.WEBHOOK_EXTERNAL_BASE_URL = os.Getenv("WEBHOOK_EXTERNAL_BASE_URL")
}

// buildPostgresDSN constructs the connection string for PostgreSQL.
func buildPostgresDSN() {
	config.DB_DSN = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME,
	)
}

// unmarshalYaml binds values from resources.yaml into the config struct.
func unmarshalYaml() {
	if err := viper.Unmarshal(&config); err != nil {
		slog.Error("Failed to unmarshal resources.yaml", "error", err, "operation", "config")
	}
}
