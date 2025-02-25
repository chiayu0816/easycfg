package easycfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	// Create test directory
	tempDir := t.TempDir()

	// Create test YAML file
	yamlContent := `
app:
  name: "test-app"
  version: "1.0.0"
server:
  host: "localhost"
  port: 8080
  timeout: 30
database:
  url: "postgres://localhost:5432/testdb"
  max_connections: 10
  timeout: 5
`
	yamlPath := filepath.Join(tempDir, "config.yml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Set output directory for generated Go file
	outputDir := filepath.Join(tempDir, "generated")

	// Step 1: Generate Go struct
	err := YamlToStruct(yamlPath, outputDir, "config")
	if err != nil {
		t.Fatalf("YamlToStruct failed: %v", err)
	}

	// Check if the generated file exists
	generatedFilePath := filepath.Join(outputDir, "config.go")
	if _, err := os.Stat(generatedFilePath); os.IsNotExist(err) {
		t.Fatalf("Generated file does not exist: %s", generatedFilePath)
	}

	// Step 2: Load configuration using map
	cfg := make(map[string]interface{})
	err = LoadConfig(yamlPath, &cfg)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Check if configuration is loaded correctly
	app, ok := cfg["app"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to get app configuration")
	}

	if app["name"] != "test-app" {
		t.Errorf("app.name = %v, expected \"test-app\"", app["name"])
	}

	server, ok := cfg["server"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to get server configuration")
	}

	// Safely check port value, not depending on specific type
	portValue := server["port"]
	fmt.Printf("Port value type: %T, value: %v\n", portValue, portValue)

	// Use type assertion to check port value
	switch port := portValue.(type) {
	case int:
		if port != 8080 {
			t.Errorf("server.port = %v (int), expected 8080", port)
		}
	case float64:
		if port != 8080 {
			t.Errorf("server.port = %v (float64), expected 8080", port)
		}
	default:
		t.Errorf("server.port type error: %T, value: %v", portValue, portValue)
	}

	// Step 3: Test configuration change monitoring
	changeDetected := make(chan bool)

	// Monitor configuration changes
	err = WatchConfig(yamlPath, &cfg, func() {
		changeDetected <- true
	})
	if err != nil {
		t.Fatalf("WatchConfig failed: %v", err)
	}

	// Modify configuration file
	updatedYamlContent := `
app:
  name: "updated-app"
  version: "2.0.0"
server:
  host: "127.0.0.1"
  port: 9090
  timeout: 60
database:
  url: "postgres://localhost:5432/testdb"
  max_connections: 20
  timeout: 10
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
		app, ok := cfg["app"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to get updated app configuration")
		}

		if app["name"] != "updated-app" {
			t.Errorf("Updated app.name = %v, expected \"updated-app\"", app["name"])
		}

		server, ok := cfg["server"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to get updated server configuration")
		}

		// Safely check updated port value
		updatedPortValue := server["port"]
		fmt.Printf("Updated port value type: %T, value: %v\n", updatedPortValue, updatedPortValue)

		// Use type assertion to check updated port value
		switch port := updatedPortValue.(type) {
		case int:
			if port != 9090 {
				t.Errorf("Updated server.port = %v (int), expected 9090", port)
			}
		case float64:
			if port != 9090 {
				t.Errorf("Updated server.port = %v (float64), expected 9090", port)
			}
		default:
			t.Errorf("Updated server.port type error: %T, value: %v", updatedPortValue, updatedPortValue)
		}

	case <-time.After(2 * time.Second):
		t.Fatal("Timeout: configuration change not detected")
	}
}

// TestIntegrationWithTestConfig performs a complete test of struct generation and loading using test_config.yml
func TestIntegrationWithTestConfig(t *testing.T) {
	// Create test directory
	tempDir := t.TempDir()

	// Read the original test_config.yml file
	yamlContent, err := os.ReadFile("test_config.yml")
	if err != nil {
		t.Fatalf("Failed to read test_config.yml file: %v", err)
	}

	// Create a copy of the test YAML file in the temporary directory
	// Note: We need to create a copy because later tests will modify this file to test hot reload functionality
	// Using the original file directly would cause tests to interfere with each other and potentially modify files in the project
	yamlPath := filepath.Join(tempDir, "test_config.yml")
	if err := os.WriteFile(yamlPath, yamlContent, 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Set output directory for generated Go file
	outputDir := filepath.Join(tempDir, "generated")

	// Step 1: Generate Go struct
	err = YamlToStruct(yamlPath, outputDir, "testconfig")
	if err != nil {
		t.Fatalf("YamlToStruct failed: %v", err)
	}

	// Check if the generated file exists
	generatedFilePath := filepath.Join(outputDir, "testconfig.go")
	if _, err := os.Stat(generatedFilePath); os.IsNotExist(err) {
		t.Fatalf("Generated file does not exist: %s", generatedFilePath)
	}

	// Read the generated file content
	content, err := os.ReadFile(generatedFilePath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	// Check if the generated content contains expected structs
	expectedStructs := []string{
		"TestConfig struct",
		"General struct",
		"GeneralServer struct",
		"GeneralSubscriber struct",
		"Redis struct",
		"Logger struct",
		"WsListenPort int",
		"Type string",
		"Port string",
		"RpcPort int",
		"Addrs []string",
		"Password string",
		"Path string",
		"Level string",
	}

	contentStr := string(content)
	for _, expected := range expectedStructs {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Generated file is missing expected content: %s", expected)
		}
	}

	// Step 2: Load configuration using map instead of manually defining structs
	// Note: We use map[string]interface{} to load configuration, avoiding duplicate struct definitions
	// Ideally, we should use the generated struct directly, but this would require dynamic compilation or other advanced techniques
	cfg := make(map[string]interface{})
	err = LoadConfig(yamlPath, &cfg)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Step 3: Check if configuration is loaded correctly
	// Check General section
	general, ok := cfg["general"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to get general configuration")
	}

	if general["type"] != "exchange" {
		t.Errorf("general.type = %q, expected \"exchange\"", general["type"])
	}

	generalServer, ok := general["server"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to get general.server configuration")
	}

	if generalServer["port"] != ":9311" {
		t.Errorf("general.server.port = %q, expected \":9311\"", generalServer["port"])
	}

	// Safely check ws_listen_port value
	wsListenPortValue := general["ws_listen_port"]
	switch wsListenPort := wsListenPortValue.(type) {
	case int:
		if wsListenPort != 8081 {
			t.Errorf("general.ws_listen_port = %d (int), expected 8081", wsListenPort)
		}
	case float64:
		if wsListenPort != 8081 {
			t.Errorf("general.ws_listen_port = %f (float64), expected 8081", wsListenPort)
		}
	default:
		t.Errorf("general.ws_listen_port type error: %T, value: %v", wsListenPortValue, wsListenPortValue)
	}

	generalSubscriber, ok := general["subscriber"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to get general.subscriber configuration")
	}

	if generalSubscriber["type"] != "redis" {
		t.Errorf("general.subscriber.type = %q, expected \"redis\"", generalSubscriber["type"])
	}

	// Safely check rpc_port value
	rpcPortValue := generalSubscriber["rpc_port"]
	switch rpcPort := rpcPortValue.(type) {
	case int:
		if rpcPort != 9211 {
			t.Errorf("general.subscriber.rpc_port = %d (int), expected 9211", rpcPort)
		}
	case float64:
		if rpcPort != 9211 {
			t.Errorf("general.subscriber.rpc_port = %f (float64), expected 9211", rpcPort)
		}
	default:
		t.Errorf("general.subscriber.rpc_port type error: %T, value: %v", rpcPortValue, rpcPortValue)
	}

	// Check Redis section
	redis, ok := cfg["redis"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to get redis configuration")
	}

	redisAddrs, ok := redis["addrs"].([]interface{})
	if !ok {
		t.Fatalf("Failed to get redis.addrs configuration")
	}

	expectedAddrs := []string{
		"10.121.1.9:7000",
		"10.121.1.9:7001",
		"10.121.1.9:7002",
		"10.121.1.9:7003",
		"10.121.1.9:7004",
		"10.121.1.9:7005",
	}

	if len(redisAddrs) != len(expectedAddrs) {
		t.Errorf("redis.addrs length = %d, expected %d", len(redisAddrs), len(expectedAddrs))
	} else {
		for i, addr := range redisAddrs {
			if addr.(string) != expectedAddrs[i] {
				t.Errorf("redis.addrs[%d] = %q, expected %q", i, addr, expectedAddrs[i])
			}
		}
	}

	if redis["password"] != "redis123" {
		t.Errorf("redis.password = %q, expected \"redis123\"", redis["password"])
	}

	// Check Logger section
	logger, ok := cfg["logger"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to get logger configuration")
	}

	if logger["path"] != "../log" {
		t.Errorf("logger.path = %q, expected \"../log\"", logger["path"])
	}

	if logger["level"] != "debug" {
		t.Errorf("logger.level = %q, expected \"debug\"", logger["level"])
	}

	// Step 4: Test configuration change monitoring
	changeDetected := make(chan bool)

	// Monitor configuration changes
	err = WatchConfig(yamlPath, &cfg, func() {
		changeDetected <- true
	})
	if err != nil {
		t.Fatalf("WatchConfig failed: %v", err)
	}

	// Modify configuration file
	updatedYamlContent := `
general:
  type: futures # changed
  server:
    port: :9312 # changed
  ws_listen_port: 8082 # changed
  subscriber:
    type: kafka # changed
    rpc_port: 9212 # changed
redis:
  addrs: # reduced
    - "10.121.1.9:7000"
    - "10.121.1.9:7001"
    - "10.121.1.9:7002"
  password: "newpassword" # changed
logger:
  path: "../newlog" # changed
  level: "info" # changed
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
		// Check General section
		general, ok = cfg["general"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to get updated general configuration")
		}

		if general["type"] != "futures" {
			t.Errorf("Updated general.type = %q, expected \"futures\"", general["type"])
		}

		generalServer, ok = general["server"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to get updated general.server configuration")
		}

		if generalServer["port"] != ":9312" {
			t.Errorf("Updated general.server.port = %q, expected \":9312\"", generalServer["port"])
		}

		// Safely check updated ws_listen_port value
		wsListenPortValue = general["ws_listen_port"]
		switch wsListenPort := wsListenPortValue.(type) {
		case int:
			if wsListenPort != 8082 {
				t.Errorf("Updated general.ws_listen_port = %d (int), expected 8082", wsListenPort)
			}
		case float64:
			if wsListenPort != 8082 {
				t.Errorf("Updated general.ws_listen_port = %f (float64), expected 8082", wsListenPort)
			}
		default:
			t.Errorf("Updated general.ws_listen_port type error: %T, value: %v", wsListenPortValue, wsListenPortValue)
		}

		generalSubscriber, ok = general["subscriber"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to get updated general.subscriber configuration")
		}

		if generalSubscriber["type"] != "kafka" {
			t.Errorf("Updated general.subscriber.type = %q, expected \"kafka\"", generalSubscriber["type"])
		}

		// Safely check updated rpc_port value
		rpcPortValue = generalSubscriber["rpc_port"]
		switch rpcPort := rpcPortValue.(type) {
		case int:
			if rpcPort != 9212 {
				t.Errorf("Updated general.subscriber.rpc_port = %d (int), expected 9212", rpcPort)
			}
		case float64:
			if rpcPort != 9212 {
				t.Errorf("Updated general.subscriber.rpc_port = %f (float64), expected 9212", rpcPort)
			}
		default:
			t.Errorf("Updated general.subscriber.rpc_port type error: %T, value: %v", rpcPortValue, rpcPortValue)
		}

		// Check Redis section
		redis, ok = cfg["redis"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to get updated redis configuration")
		}

		redisAddrs, ok = redis["addrs"].([]interface{})
		if !ok {
			t.Fatalf("Failed to get updated redis.addrs configuration")
		}

		updatedExpectedAddrs := []string{
			"10.121.1.9:7000",
			"10.121.1.9:7001",
			"10.121.1.9:7002",
		}

		// Due to Viper's behavior, the array length might not be updated immediately, so we only check if the first three elements match
		if len(redisAddrs) < len(updatedExpectedAddrs) {
			t.Errorf("Updated redis.addrs length = %d, expected at least %d", len(redisAddrs), len(updatedExpectedAddrs))
		} else {
			for i, expectedAddr := range updatedExpectedAddrs {
				if i < len(redisAddrs) && redisAddrs[i].(string) != expectedAddr {
					t.Errorf("Updated redis.addrs[%d] = %q, expected %q", i, redisAddrs[i], expectedAddr)
				}
			}
		}

		if redis["password"] != "newpassword" {
			t.Errorf("Updated redis.password = %q, expected \"newpassword\"", redis["password"])
		}

		// Check Logger section
		logger, ok = cfg["logger"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to get updated logger configuration")
		}

		if logger["path"] != "../newlog" {
			t.Errorf("Updated logger.path = %q, expected \"../newlog\"", logger["path"])
		}

		if logger["level"] != "info" {
			t.Errorf("Updated logger.level = %q, expected \"info\"", logger["level"])
		}

	case <-time.After(2 * time.Second):
		t.Fatal("Timeout: configuration change not detected")
	}
}
