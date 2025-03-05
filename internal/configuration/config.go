package configuration

import (
	"fmt"
	"log/slog"
	"microsoft-apps-exporter/internal/models"
	"os"
	"regexp"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Configuration holds all application settings
type Configuration struct {
	GRAPH_CLIENT_ID     string
	GRAPH_TENANT_ID     string
	GRAPH_CLIENT_SECRET string
	GRAPH_APP_SCOPES    string

	Sharepoint struct {
		Lists []models.ListConfig `mapstructure:"lists"`
	} `mapstructure:"sharepoint"`

	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_DSN      string

	WEBHOOK_BASE_URL    string
	WEBHOOK_SERVER_PORT string
}

var (
	config Configuration
	once   sync.Once
)

const ProjectRootDir = "microsoft-apps-exporter"

// loadEnvFile loads the .env file dynamically based on the current working directory.
func loadDotEnv() {
	projectPattern := regexp.MustCompile("^(.*" + ProjectRootDir + ")")
	currentDir, _ := os.Getwd()
	rootPath := projectPattern.Find([]byte(currentDir))

	if err := godotenv.Load(string(rootPath) + "/.env"); err != nil {
		slog.Debug("Failed to load .env file", "exception", err, "operation", "config")
	}
}

// loadConfigFile reads configuration from a YAML file using Viper.
func loadConfigYaml() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Error reading config file", "exception", err, "operation", "config")
	}
}

// GetConfig initializes and returns the application configuration as a singleton.
func GetConfig() Configuration {
	once.Do(func() {
		loadDotEnv()
		loadConfigYaml()

		config = Configuration{
			GRAPH_CLIENT_ID:     os.Getenv("GRAPH_CLIENT_ID"),
			GRAPH_TENANT_ID:     os.Getenv("GRAPH_TENANT_ID"),
			GRAPH_CLIENT_SECRET: os.Getenv("GRAPH_CLIENT_SECRET"),
			GRAPH_APP_SCOPES:    os.Getenv("GRAPH_APP_SCOPES"),

			Sharepoint: struct {
				Lists []models.ListConfig `mapstructure:"lists"`
			}{
				Lists: []models.ListConfig{},
			},

			DB_HOST:     os.Getenv("DB_HOST"),
			DB_PORT:     os.Getenv("DB_PORT"),
			DB_USER:     os.Getenv("DB_USER"),
			DB_PASSWORD: os.Getenv("DB_PASSWORD"),
			DB_NAME:     os.Getenv("DB_NAME"),
			DB_DSN:      "",

			WEBHOOK_BASE_URL:    os.Getenv("WEBHOOK_BASE_URL"),
			WEBHOOK_SERVER_PORT: os.Getenv("WEBHOOK_SERVER_PORT"),
		}

		// Construct the PostgreSQL DSN string
		config.DB_DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME)

		// Load yaml config to the structure
		if err := viper.Unmarshal(&config); err != nil {
			slog.Error("Error unmarshalling config", "error", err, "operation", "config")
		}
	})

	return config
}
