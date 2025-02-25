package easycfg

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkYamlToStruct(b *testing.B) {
	// Create test YAML file
	yamlContent := `
app:
  name: "benchmark-app"
  version: "1.0.0"
server:
  host: "localhost"
  port: 8080
  timeout: 30
  ssl:
    enabled: true
    cert: "/path/to/cert"
    key: "/path/to/key"
database:
  master:
    url: "postgres://localhost:5432/db"
    max_connections: 10
    timeout: 5
  slave:
    url: "postgres://localhost:5433/db"
    max_connections: 5
    timeout: 3
logging:
  level: "info"
  format: "json"
  output: "stdout"
  file:
    enabled: true
    path: "/var/log/app.log"
    max_size: 100
    max_backups: 5
    max_age: 30
cache:
  redis:
    host: "localhost"
    port: 6379
    password: ""
    db: 0
  memory:
    enabled: true
    size: 1000
    ttl: 300
`
	tempDir := b.TempDir()
	yamlPath := filepath.Join(tempDir, "benchmark_config.yml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		b.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Set output directory
	outputDir := filepath.Join(tempDir, "generated")

	// Reset timer
	b.ResetTimer()

	// Run benchmark test
	for i := 0; i < b.N; i++ {
		err := YamlToStruct(yamlPath, outputDir, "config")
		if err != nil {
			b.Fatalf("YamlToStruct failed: %v", err)
		}
	}
}

func BenchmarkLoadConfig(b *testing.B) {
	// Create test YAML file
	yamlContent := `
app:
  name: "benchmark-app"
  version: "1.0.0"
server:
  host: "localhost"
  port: 8080
  timeout: 30
database:
  url: "postgres://localhost:5432/db"
  max_connections: 10
  timeout: 5
logging:
  level: "info"
  format: "json"
`
	tempDir := b.TempDir()
	yamlPath := filepath.Join(tempDir, "benchmark_config.yml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		b.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Reset timer
	b.ResetTimer()

	// Run benchmark test
	for i := 0; i < b.N; i++ {
		cfg := make(map[string]interface{})
		err := LoadConfig(yamlPath, &cfg)
		if err != nil {
			b.Fatalf("LoadConfig failed: %v", err)
		}
	}
}

// Benchmark the performance of toCamelCase function
func BenchmarkToCamelCase(b *testing.B) {
	testCases := []string{
		"simple_name",
		"complex_name_with_multiple_parts",
		"very-long-name-with-different-separators.and.more",
		"UPPERCASE_NAME_WITH_NUMBERS_123",
		"mixed_Case_name_With_Different_Styles",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			_ = toCamelCase(tc)
		}
	}
}
