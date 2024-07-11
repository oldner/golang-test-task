// Package /////////////////////////////////////////////////////////////////////
package app

// Imports /////////////////////////////////////////////////////////////////////
import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang-test-task/cmd/message-processor/processor"
	"golang-test-task/cmd/message-processor/storage"
)

// Types ///////////////////////////////////////////////////////////////////////

// App represents the main application structure for the Message Processor service.
type App struct {
	config    *Config
	processor *processor.Processor
}

// Functions ///////////////////////////////////////////////////////////////////

// New creates a new instance of the App with the given configuration.
// It sets up the Redis client and message processor.
func New(config *Config) (*App, error) {
	redisClient, err := storage.NewRedisClient(config.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %v", err)
	}

	proc, err := processor.NewProcessor(config.RabbitMQURL, redisClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create processor: %v", err)
	}

	return &App{
		config:    config,
		processor: proc,
	}, nil
}

// Run starts the message processor and sets up graceful shutdown.
func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the processor in a goroutine
	go func() {
		if err := a.processor.Start(ctx); err != nil {
			fmt.Printf("Processor error: %v\n", err)
			cancel()
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the processor
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down message processor...")
	cancel()

	if err := a.processor.Shutdown(); err != nil {
		return fmt.Errorf("error during shutdown: %v", err)
	}

	fmt.Println("Message processor exited")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
