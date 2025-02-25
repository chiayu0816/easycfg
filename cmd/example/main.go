package main

import (
	"fmt"
	"log"

	"github.com/chiayu0816/easycfg"
)

func main() {
	// Check if configuration file has been generated
	fmt.Println("EasyCfg Example Program")
	fmt.Println("Attempting to load configuration file...")

	// Create configuration struct instance
	// Note: You should import the generated configuration package and use the correct struct name
	// For example: cfg := &generated.WsClientConfig{}
	// Since we cannot directly import the generated package, we use a map instead
	cfg := make(map[string]interface{})

	// Load configuration
	if err := easycfg.LoadConfig("test_config.yml", &cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}

	fmt.Println("Configuration successfully loaded!")
	fmt.Println("In actual use, you can access the configuration like this:")
	fmt.Println("cfg.General.Type")
	fmt.Println("cfg.General.Server.Port")
	fmt.Println("cfg.General.WsListenPort")
	fmt.Println("cfg.Redis.Addrs")
	fmt.Println("cfg.Logger.Level")

	fmt.Println("\nYou can also monitor configuration file changes:")
	fmt.Println("easycfg.WatchConfig(\"test_config.yml\", cfg, func() {")
	fmt.Println("    // Logic to handle configuration changes")
	fmt.Println("})")
}
