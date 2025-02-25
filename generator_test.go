package easycfg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestYamlToStruct(t *testing.T) {
	// Create test YAML file
	yamlContent := `
test:
  string_value: "hello"
  int_value: 123
  bool_value: true
  nested:
    value: "nested value"
  array:
    - "item1"
    - "item2"
`
	tempDir := t.TempDir()
	yamlPath := filepath.Join(tempDir, "test_config.yml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	// Set output directory
	outputDir := filepath.Join(tempDir, "generated")

	// Execute test
	err := YamlToStruct(yamlPath, outputDir, "testconfig")
	if err != nil {
		t.Fatalf("YamlToStruct failed: %v", err)
	}

	// Check if generated file exists
	generatedFilePath := filepath.Join(outputDir, "testconfig.go")
	if _, err := os.Stat(generatedFilePath); os.IsNotExist(err) {
		t.Fatalf("Generated file does not exist: %s", generatedFilePath)
	}

	// Read generated file content
	content, err := os.ReadFile(generatedFilePath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	// Check if generated content contains expected structs
	expectedStructs := []string{
		"TestConfig struct",
		"Test struct",
		"TestNested struct",
		"StringValue string",
		"IntValue int",
		"BoolValue bool",
		"Nested TestNested",
		"Array []string",
	}

	contentStr := string(content)
	for _, expected := range expectedStructs {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Generated file is missing expected content: %s", expected)
		}
	}
}

func TestToCamelCase(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"test", "Test"},
		{"test_case", "TestCase"},
		{"test-case", "TestCase"},
		{"test.case", "TestCase"},
		{"TEST_CASE", "TESTCASE"},
		{"123test", "123test"},
		{"test123", "Test123"},
		{"", ""},
	}

	for _, tc := range testCases {
		result := toCamelCase(tc.input)
		if result != tc.expected {
			t.Errorf("toCamelCase(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

// Helper function: check if string contains substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
