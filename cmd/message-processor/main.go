// Package /////////////////////////////////////////////////////////////////////
package main

// Imports /////////////////////////////////////////////////////////////////////
import (
	"log"

	"golang-test-task/cmd/message-processor/app"
)

// Main ////////////////////////////////////////////////////////////////////////

// This function is the entry point to the Message Processor service.
// It takes care of:
// - loading the configuration
// - creating a new application instance
// - running the application
func main() {
	// Load configuration
	config, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create application
	app, err := app.New(config)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Run application
	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

////////////////////////////////////////////////////////////////////////////////
