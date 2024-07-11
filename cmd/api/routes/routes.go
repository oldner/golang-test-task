package routes

import (
	"golang-test-task/cmd/api/handlers"
	"golang-test-task/cmd/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, messageHandler *handlers.MessageHandler) {
	router.Use(middleware.Logging())

	v1 := router.Group("/v1")
	{
		v1.POST("/message", messageHandler.PostMessage)
	}
}
