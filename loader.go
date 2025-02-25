package easycfg

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// LoadConfig loads configuration from YAML file to the specified struct using Viper
func LoadConfig(configPath string, configStruct interface{}) error {
	// Get file name and extension
	ext := filepath.Ext(configPath)
	fileName := filepath.Base(configPath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, ext)
	dirPath := filepath.Dir(configPath)

	// Initialize Viper
	v := viper.New()
	v.SetConfigName(fileNameWithoutExt)
	v.SetConfigType(strings.TrimPrefix(ext, "."))
	v.AddConfigPath(dirPath)

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read configuration file: %v", err)
	}

	// Map configuration to struct
	if err := v.Unmarshal(configStruct); err != nil {
		return fmt.Errorf("failed to map configuration to struct: %v", err)
	}

	return nil
}

// WatchConfig monitors configuration file changes and automatically reloads
func WatchConfig(configPath string, configStruct interface{}, onChange func()) error {
	// Get file name and extension
	ext := filepath.Ext(configPath)
	fileName := filepath.Base(configPath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, ext)
	dirPath := filepath.Dir(configPath)

	// Initialize Viper
	v := viper.New()
	v.SetConfigName(fileNameWithoutExt)
	v.SetConfigType(strings.TrimPrefix(ext, "."))
	v.AddConfigPath(dirPath)

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read configuration file: %v", err)
	}

	// Map configuration to struct
	if err := v.Unmarshal(configStruct); err != nil {
		return fmt.Errorf("failed to map configuration to struct: %v", err)
	}

	// Monitor configuration file changes
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		// Reload configuration
		if err := v.Unmarshal(configStruct); err != nil {
			fmt.Printf("failed to reload configuration: %v\n", err)
			return
		}

		// Call callback function
		if onChange != nil {
			onChange()
		}

		fmt.Println("configuration reloaded")
	})

	return nil
}
