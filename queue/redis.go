package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/user/chatbot-whatsapp/config"
	"github.com/user/chatbot-whatsapp/models"
)

const (
	QueueOutbound = "queue:messages:outbound"
	QueueRetry    = "queue:messages:retry"
	MaxRetries    = 3
)

type QueueMessage struct {
	CustomerID  uint                `json:"customer_id"`
	Phone       string              `json:"phone"`
	MessageType models.MessageType  `json:"message_type"`
	Body        string              `json:"body,omitempty"`
	MediaURL    string              `json:"media_url,omitempty"`
	Template    *QueueTemplate      `json:"template,omitempty"`
	RetryCount  int                 `json:"retry_count"`
	CreatedAt   time.Time           `json:"created_at"`
}

type QueueTemplate struct {
	Name     string   `json:"name"`
	Language string   `json:"language"`
	Params   []string `json:"params"`
}

type Client struct {
	rdb *redis.Client
	log *zap.Logger
}

func NewClient(cfg *config.RedisConfig, log *zap.Logger) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	return &Client{rdb: rdb, log: log}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) Publish(ctx context.Context, msg *QueueMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal queue message: %w", err)
	}

	if err := c.rdb.LPush(ctx, QueueOutbound, data).Err(); err != nil {
		return fmt.Errorf("push to queue: %w", err)
	}

	c.log.Debug("message published to queue",
		zap.String("phone", msg.Phone),
		zap.String("type", string(msg.MessageType)),
	)

	return nil
}

func (c *Client) Consume(ctx context.Context) (*QueueMessage, error) {
	result, err := c.rdb.BRPop(ctx, 30*time.Second, QueueOutbound).Result()
	if err != nil {
		return nil, err
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("invalid queue result")
	}

	var msg QueueMessage
	if err := json.Unmarshal([]byte(result[1]), &msg); err != nil {
		return nil, fmt.Errorf("unmarshal queue message: %w", err)
	}

	return &msg, nil
}

func (c *Client) PublishRetry(ctx context.Context, msg *QueueMessage) error {
	msg.RetryCount++
	if msg.RetryCount > MaxRetries {
		c.log.Warn("message exceeded max retries, discarding",
			zap.String("phone", msg.Phone),
			zap.Int("retries", msg.RetryCount),
		)
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal retry message: %w", err)
	}

	if err := c.rdb.LPush(ctx, QueueRetry, data).Err(); err != nil {
		return fmt.Errorf("push to retry queue: %w", err)
	}

	return nil
}

func (c *Client) QueueLength(ctx context.Context) (int64, error) {
	return c.rdb.LLen(ctx, QueueOutbound).Result()
}
