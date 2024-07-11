// Package /////////////////////////////////////////////////////////////////////
package app

// Imports /////////////////////////////////////////////////////////////////////
import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-test-task/cmd/api/handlers"
	"golang-test-task/cmd/api/routes"

	"github.com/gin-gonic/gin"
)

// Types ///////////////////////////////////////////////////////////////////////

// App represents the main application structure for the API service.
type App struct {
	router         *gin.Engine
	config         *Config
	messageHandler *handlers.MessageHandler
}

// Functions ///////////////////////////////////////////////////////////////////

// New creates a new instance of the App with the given configuration.
// It sets up the message handler and routes.
func New(config *Config) (*App, error) {
	messageHandler, err := handlers.NewMessageHandler(config.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create message handler: %v", err)
	}

	router := gin.Default()
	routes.SetupRoutes(router, messageHandler)

	return &App{
		router:         router,
		config:         config,
		messageHandler: messageHandler,
	}, nil
}

// Run starts the HTTP server and sets up graceful shutdown.
func (a *App) Run() error {
	srv := &http.Server{
		Addr:    ":" + a.config.Port,
		Handler: a.router,
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %v", err)
	}

	a.messageHandler.Close()

	fmt.Println("Server exiting")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
