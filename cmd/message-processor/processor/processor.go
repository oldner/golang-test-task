package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"golang-test-task/pkg/models"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

type Processor struct {
	rabbitmqURL  string
	rabbitmqConn *amqp.Connection
	redisClient  *redis.Client
}

func NewProcessor(rabbitmqURL string, redisClient *redis.Client) (*Processor, error) {
	return &Processor{
		rabbitmqURL: rabbitmqURL,
		redisClient: redisClient,
	}, nil
}

func (p *Processor) Start(ctx context.Context) error {
	var err error
	p.rabbitmqConn, err = amqp.Dial(p.rabbitmqURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := p.rabbitmqConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("messages", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %v", err)
	}

	for {
		select {
		case d := <-msgs:
			var msg models.Message
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				fmt.Printf("Error unmarshaling message: %v\n", err)
				continue
			}

			if err := p.saveMessage(ctx, msg); err != nil {
				fmt.Printf("Error saving message: %v\n", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (p *Processor) saveMessage(ctx context.Context, msg models.Message) error {
	key := fmt.Sprintf("messages:%s:%s", msg.Sender, msg.Receiver)
	value, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	err = p.redisClient.ZAdd(ctx, key, &redis.Z{
		Score:  float64(msg.Timestamp.Unix()),
		Member: value,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to save message to Redis: %v", err)
	}

	return nil
}

func (p *Processor) Shutdown() error {
	if p.rabbitmqConn != nil {
		if err := p.rabbitmqConn.Close(); err != nil {
			return fmt.Errorf("error closing RabbitMQ connection: %v", err)
		}
	}
	return nil
}
