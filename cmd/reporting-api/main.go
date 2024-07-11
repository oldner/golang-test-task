package main

import (
	"log"

	"golang-test-task/cmd/reporting-api/app"
)

func main() {
	config, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	app, err := app.New(config)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
