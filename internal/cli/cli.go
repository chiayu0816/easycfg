package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/chiayu0816/easycfg"
)

// Run executes the CLI command
func Run() {
	// Define command line parameters
	yamlPath := flag.String("yaml", "", "Path to YAML configuration file")
	outputDir := flag.String("output", "generated", "Output directory for generated Go files")
	packageName := flag.String("package", "config", "Package name for generated Go files")
	watch := flag.Bool("watch", false, "Whether to watch for configuration file changes")
	flag.Parse()

	// Check required parameters
	if *yamlPath == "" {
		fmt.Println("Error: YAML configuration file path must be specified")
		flag.Usage()
		os.Exit(1)
	}

	// Ensure YAML file exists
	if _, err := os.Stat(*yamlPath); os.IsNotExist(err) {
		fmt.Printf("Error: YAML file does not exist: %s\n", *yamlPath)
		os.Exit(1)
	}

	// Generate Go struct file
	if err := easycfg.YamlToStruct(*yamlPath, *outputDir, *packageName); err != nil {
		fmt.Printf("Error: Failed to generate Go struct: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully generated Go struct file")

	// If watch mode is enabled
	if *watch {
		fmt.Println("Watching for configuration file changes...")

		// Create a dummy config map to use with WatchConfig
		dummyConfig := make(map[string]interface{})

		// Watch for YAML file changes using the WatchConfig function
		if err := easycfg.WatchConfig(*yamlPath, &dummyConfig, func() {
			// Regenerate Go struct when changes are detected
			if err := easycfg.YamlToStruct(*yamlPath, *outputDir, *packageName); err != nil {
				fmt.Printf("Error: Failed to regenerate Go struct: %v\n", err)
			} else {
				fmt.Println("Configuration changes detected, Go struct file has been regenerated")
			}
		}); err != nil {
			fmt.Printf("Error: Failed to set up configuration file watcher: %v\n", err)
			os.Exit(1)
		}

		// Block main thread
		select {}
	}
}
