//go:build testing && unit

package configuration_test

import (
	"log/slog"
	"math"
	"os"
	"testing"

	"microsoft-apps-exporter/internal/configuration"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// setupTestEnv initializes environment variables for testing.
func setupTestEnv() {
	vars := map[string]string{
		"LOG_LEVEL":                 "LOG_LEVEL",
		"GRAPH_CLIENT_ID":           "GRAPH_CLIENT_ID",
		"GRAPH_TENANT_ID":           "GRAPH_TENANT_ID",
		"GRAPH_CLIENT_SECRET":       "GRAPH_CLIENT_SECRET",
		"GRAPH_APP_SCOPES":          "GRAPH_APP_SCOPES",
		"DB_HOST":                   "DB_HOST",
		"DB_PORT":                   "DB_PORT",
		"DB_USER":                   "DB_USER",
		"DB_PASSWORD":               "DB_PASSWORD",
		"DB_NAME":                   "DB_NAME",
		"WEBHOOK_LISTEN_IP":         "WEBHOOK_LISTEN_IP",
		"WEBHOOK_LISTEN_PORT":       "WEBHOOK_LISTEN_PORT",
		"WEBHOOK_EXTERNAL_BASE_URL": "WEBHOOK_EXTERNAL_BASE_URL",
	}
	for key, value := range vars {
		os.Setenv(key, value)
	}
}

// setupTestResourcesYaml configures Viper to load test YAML resources.
func setupTestResourcesYaml() {
	configuration.ResetConfig()

	viper.Reset()
	viper.SetConfigName("resources")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
}

// setupProdResourcesYaml configures Viper to load production YAML resources.
func setupProdResourcesYaml() {
	configuration.ResetConfig()

	viper.Reset()
	viper.AddConfigPath("../../..") // Root dir
}

// TestUnmarshalConfig verifies correct unmarshalling of YAML resources.
func TestUnmarshalConfig(t *testing.T) {
	setupTestResourcesYaml()

	assert.NoError(t, viper.ReadInConfig(), "resources.yaml should be parsed without syntax errors")

	var config configuration.Configuration
	assert.NoError(t, viper.Unmarshal(&config), "Resources should unmarshal into config correctly")

	validateSharepointResource(t, config)
}

// TestGetConfig verifies correct environment variables loading.
func TestGetConfig(t *testing.T) {
	setupTestEnv()
	setupProdResourcesYaml()

	slog.SetLogLoggerLevel(math.MaxInt) // Disable logging

	config := configuration.GetConfig()

	validateResourcesYaml(t, config)
	validateEnv(t, config)
}

// validateResourcesYaml checks if resources configuration are properly structured.
func validateResourcesYaml(t *testing.T, config configuration.Configuration) {
	validateSharepointResource(t, config)
}

// validateSharepointResource checks if SharePoint resource is properly structured.
func validateSharepointResource(t *testing.T, config configuration.Configuration) {
	if config.Sharepoint == nil {
		t.Fatalf("Sharepoint resource is expected to be specified")
	}
	assert.NotEmpty(t, config.Sharepoint.Lists, "Sharepoint.Lists should not be empty")

	for _, list := range config.Sharepoint.Lists {
		assert.NotEmpty(t, list.SiteID, "site_id is required")
		assert.NotEmpty(t, list.ListID, "list_id is required")
		assert.NotEmpty(t, list.DbTableName, "table_name is required")
		assert.NotEmpty(t, list.ColumnsMap, "columns_map must have at least one key-value pair")

		for key, value := range list.ColumnsMap {
			assert.NotEmpty(t, key, "columns_map key cannot be empty")
			assert.NotEmpty(t, value, "columns_map value cannot be empty")
		}
	}
}

// validateEnv ensures all required configuration fields are set.
func validateEnv(t *testing.T, config configuration.Configuration) {
	fields := []struct {
		name  string
		value string
	}{
		{"LOG_LEVEL", config.LOG_LEVEL},
		{"GRAPH_CLIENT_ID", config.GRAPH_CLIENT_ID},
		{"GRAPH_TENANT_ID", config.GRAPH_TENANT_ID},
		{"GRAPH_CLIENT_SECRET", config.GRAPH_CLIENT_SECRET},
		{"GRAPH_APP_SCOPES", config.GRAPH_APP_SCOPES},
		{"DB_HOST", config.DB_HOST},
		{"DB_PORT", config.DB_PORT},
		{"DB_USER", config.DB_USER},
		{"DB_PASSWORD", config.DB_PASSWORD},
		{"DB_NAME", config.DB_NAME},
		{"WEBHOOK_LISTEN_IP", config.WEBHOOK_LISTEN_IP},
		{"WEBHOOK_LISTEN_PORT", config.WEBHOOK_LISTEN_PORT},
		{"WEBHOOK_EXTERNAL_BASE_URL", config.WEBHOOK_EXTERNAL_BASE_URL},
	}

	for _, field := range fields {
		assert.NotEmpty(t, field.value, "%s should not be empty", field.name)
	}

	// Check the constructed DSN string
	expectedDSN := "postgres://DB_USER:DB_PASSWORD@DB_HOST:DB_PORT/DB_NAME?sslmode=disable"
	assert.Equal(t, expectedDSN, config.DB_DSN)
}

// TestConfigSingleton ensures GetConfig() returns the same instance every time.
func TestConfigSingleton(t *testing.T) {
	setupTestEnv()

	config1 := configuration.GetConfig()
	config2 := configuration.GetConfig()

	assert.Equal(t, &config1, &config2, "GetConfig should return a singleton instance")
}
