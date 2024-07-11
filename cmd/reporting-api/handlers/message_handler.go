package handlers

import (
	"encoding/json"
	"net/http"

	"golang-test-task/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func GetMessageList(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sender := c.Query("sender")
		receiver := c.Query("receiver")

		if sender == "" || receiver == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Both sender and receiver are required"})
			return
		}

		key := "messages:" + sender + ":" + receiver
		messages, err := redisClient.ZRevRange(c, key, 0, -1).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
			return
		}

		var result []models.Message
		for _, msg := range messages {
			var message models.Message
			if err := json.Unmarshal([]byte(msg), &message); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal message"})
				return
			}
			result = append(result, message)
		}

		c.JSON(http.StatusOK, result)
	}
}
