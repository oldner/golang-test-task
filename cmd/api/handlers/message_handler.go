// Package /////////////////////////////////////////////////////////////////////
package handlers

// Imports /////////////////////////////////////////////////////////////////////
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang-test-task/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/streadway/amqp"
)

// Types ///////////////////////////////////////////////////////////////////////

// PostMessageRequest represents the structure of the incoming message request.
type PostMessageRequest struct {
	Sender   string `json:"sender" validate:"required"`
	Receiver string `json:"receiver" validate:"required"`
	Message  string `json:"message" validate:"required"`
}

// MessageHandler handles the RabbitMQ connection and message publishing.
type MessageHandler struct {
	rabbitMQConn *amqp.Connection
	rabbitMQChan *amqp.Channel
}

// Variables ///////////////////////////////////////////////////////////////////

var validate *validator.Validate

// Init ////////////////////////////////////////////////////////////////////////

func init() {
	validate = validator.New()
}

// Functions ///////////////////////////////////////////////////////////////////

// NewMessageHandler creates a new MessageHandler with the given RabbitMQ URL.
func NewMessageHandler(rabbitMQURL string) (*MessageHandler, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	_, err = ch.QueueDeclare(
		"messages",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	return &MessageHandler{
		rabbitMQConn: conn,
		rabbitMQChan: ch,
	}, nil
}

// Close closes the RabbitMQ channel and connection.
func (h *MessageHandler) Close() {
	if h.rabbitMQChan != nil {
		h.rabbitMQChan.Close()
	}
	if h.rabbitMQConn != nil {
		h.rabbitMQConn.Close()
	}
}

// PostMessage handles the posting of a new message to RabbitMQ.
func (h *MessageHandler) PostMessage(c *gin.Context) {
	var req PostMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := models.Message{
		Sender:    req.Sender,
		Receiver:  req.Receiver,
		Content:   req.Message,
		Timestamp: time.Now(),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal message"})
		return
	}

	err = h.rabbitMQChan.Publish(
		"",
		"messages",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgBytes,
		})
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Message sent successfully"})
}

////////////////////////////////////////////////////////////////////////////////
