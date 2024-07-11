package routes

import (
	"golang-test-task/cmd/reporting-api/handlers"
	"golang-test-task/cmd/reporting-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func SetupRoutes(router *gin.Engine, redisClient *redis.Client) {
	router.Use(middleware.Logging())

	v1 := router.Group("/v1")
	{
		v1.GET("/message/list", handlers.GetMessageList(redisClient))
	}
}
