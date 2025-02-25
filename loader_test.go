package easycfg

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Test configuration struct
type TestConfig struct {
	Server struct {
		Host string `yaml:"host" mapstructure:"host"`
		Port int    `yaml:"port" mapstructure:"port"`
	} `yaml:"server" mapstructure:"server"`
	Database struct {
		URL      string `yaml:"url" mapstructure:"url"`
		Username string `yaml:"username" mapstructure:"username"`
		Password string `yaml:"password" mapstructure:"password"`
	} `yaml:"database" mapstructure:"database"`
	Logging struct {
		Level  string `yaml:"level" mapstructure:"level"`
		Format string `yaml:"format" mapstructure:"format"`
	} `yaml:"logging" mapstructure:"logging"`
}

func TestLoadConfig(t *testing.T) {
	// Create test YAML file
	yamlContent := `
server:
  host: localhost
  port: 8080
database:
  url: mysql://localhost:3306/testdb
  username: testuser
  password: testpass
logging:
  level: debug
  format: json
`
	tempDir := t.TempDir()
	yamlPath := filepath.Join(tempDir, "test_config.yml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Create configuration struct
	cfg := &TestConfig{}

	// Load configuration
	err := LoadConfig(yamlPath, cfg)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Check if configuration is loaded correctly
	if cfg.Server.Host != "localhost" {
		t.Errorf("cfg.Server.Host = %q, expected \"localhost\"", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("cfg.Server.Port = %d, expected 8080", cfg.Server.Port)
	}
	if cfg.Database.URL != "mysql://localhost:3306/testdb" {
		t.Errorf("cfg.Database.URL = %q, expected \"mysql://localhost:3306/testdb\"", cfg.Database.URL)
	}
	if cfg.Database.Username != "testuser" {
		t.Errorf("cfg.Database.Username = %q, expected \"testuser\"", cfg.Database.Username)
	}
	if cfg.Database.Password != "testpass" {
		t.Errorf("cfg.Database.Password = %q, expected \"testpass\"", cfg.Database.Password)
	}
	if cfg.Logging.Level != "debug" {
		t.Errorf("cfg.Logging.Level = %q, expected \"debug\"", cfg.Logging.Level)
	}
	if cfg.Logging.Format != "json" {
		t.Errorf("cfg.Logging.Format = %q, expected \"json\"", cfg.Logging.Format)
	}
}

func TestWatchConfig(t *testing.T) {
	// Create test YAML file
	yamlContent := `
server:
  host: localhost
  port: 8080
database:
  url: mysql://localhost:3306/testdb
  username: testuser
  password: testpass
logging:
  level: debug
  format: json
`
	tempDir := t.TempDir()
	yamlPath := filepath.Join(tempDir, "test_config.yml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Create configuration struct
	cfg := &TestConfig{}

	// Set change flag
	var changed bool
	changeDetected := make(chan bool)

	// Monitor configuration changes
	err := WatchConfig(yamlPath, cfg, func() {
		changed = true
		changeDetected <- true
	})
	if err != nil {
		t.Fatalf("WatchConfig failed: %v", err)
	}

	// Modify configuration file
	updatedYamlContent := `
server:
  host: 127.0.0.1
  port: 9090
database:
  url: mysql://localhost:3306/testdb
  username: testuser
  password: testpass
logging:
  level: info
  format: text
`
	// Wait for a while to ensure the file watcher has started
	time.Sleep(100 * time.Millisecond)

	// Write updated configuration
	if err := os.WriteFile(yamlPath, []byte(updatedYamlContent), 0644); err != nil {
		t.Fatalf("Failed to update test YAML file: %v", err)
	}

	// Wait for configuration change to be detected, or timeout
	select {
	case <-changeDetected:
		// Check if configuration has been updated
		if cfg.Server.Host != "127.0.0.1" {
			t.Errorf("cfg.Server.Host = %q, expected \"127.0.0.1\"", cfg.Server.Host)
		}
		if cfg.Server.Port != 9090 {
			t.Errorf("cfg.Server.Port = %d, expected 9090", cfg.Server.Port)
		}
		if cfg.Logging.Level != "info" {
			t.Errorf("cfg.Logging.Level = %q, expected \"info\"", cfg.Logging.Level)
		}
		if cfg.Logging.Format != "text" {
			t.Errorf("cfg.Logging.Format = %q, expected \"text\"", cfg.Logging.Format)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout: configuration change not detected")
	}

	if !changed {
		t.Error("Configuration change callback was not called")
	}
}
