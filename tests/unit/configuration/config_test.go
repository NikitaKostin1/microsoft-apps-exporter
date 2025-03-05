//go:build testing && unit

package configuration_test

import (
	"os"
	"testing"

	"microsoft-apps-exporter/internal/configuration"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// setupTestEnv initializes environment variables for testing.
func setupTestEnv() {
	vars := map[string]string{
		"GRAPH_CLIENT_ID":     "GRAPH_CLIENT_ID",
		"GRAPH_TENANT_ID":     "GRAPH_TENANT_ID",
		"GRAPH_CLIENT_SECRET": "GRAPH_CLIENT_SECRET",
		"GRAPH_APP_SCOPES":    "GRAPH_APP_SCOPES",
		"DB_HOST":             "DB_HOST",
		"DB_PORT":             "DB_PORT",
		"DB_USER":             "DB_USER",
		"DB_PASSWORD":         "DB_PASSWORD",
		"DB_NAME":             "DB_NAME",
	}
	for key, value := range vars {
		os.Setenv(key, value)
	}
}

// setupTestYaml configures Viper to load test YAML configurations.
func setupTestYaml() {
	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
	viper.AddConfigPath("./tests/unit")
}

// TestUnmarshalConfig verifies correct unmarshalling of YAML configuration.
func TestUnmarshalConfig(t *testing.T) {
	setupTestYaml()

	assert.NoError(t, viper.ReadInConfig(), "YAML should be parsed without syntax errors")

	var config configuration.Configuration
	assert.NoError(t, viper.Unmarshal(&config), "Config should unmarshal correctly")

	validateSharepointLists(t, config)
}

// TestGetConfig verifies correct environment variable loading and configuration retrieval.
func TestGetConfig(t *testing.T) {
	setupTestEnv()

	config := configuration.GetConfig()

	validateSharepointLists(t, config)
	validateEnv(t, config)
}

// validateSharepointLists checks if SharePoint configuration is properly structured.
func validateSharepointLists(t *testing.T, config configuration.Configuration) {
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

// validateEnv ensures all required environment variables are set in the configuration.
func validateEnv(t *testing.T, config configuration.Configuration) {
	envVars := []struct {
		name  string
		value string
	}{
		{"GRAPH_CLIENT_ID", config.GRAPH_CLIENT_ID},
		{"GRAPH_TENANT_ID", config.GRAPH_TENANT_ID},
		{"GRAPH_CLIENT_SECRET", config.GRAPH_CLIENT_SECRET},
		{"GRAPH_APP_SCOPES", config.GRAPH_APP_SCOPES},
		{"DB_HOST", config.DB_HOST},
		{"DB_PORT", config.DB_PORT},
		{"DB_USER", config.DB_USER},
		{"DB_PASSWORD", config.DB_PASSWORD},
		{"DB_NAME", config.DB_NAME},
	}

	for _, env := range envVars {
		assert.NotEmpty(t, env.value, "%s should not be empty", env.name)
	}
}

// TestConfigSingleton ensures GetConfig() returns the same instance every time.
func TestConfigSingleton(t *testing.T) {
	setupTestEnv()

	config1 := configuration.GetConfig()
	config2 := configuration.GetConfig()

	assert.Equal(t, &config1, &config2, "GetConfig should return a singleton instance")
}
