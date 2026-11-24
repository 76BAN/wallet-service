package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	LoggingExchange   = "logging_exchange"
	LoggingRoutingKey = "logging.created"
)

type LoggingMessage struct {
	ServiceName  string         `json:"service_name"`
	UserID       *int64         `json:"user_id"`
	Level        string         `json:"level"`
	EventName    string         `json:"event_name"`
	Status       string         `json:"status"`
	Message      string         `json:"message"`
	ResourceType *string        `json:"resource_type"`
	ResourceID   *string        `json:"resource_id"`
	RequestID    *string        `json:"request_id"`
	Metadata     map[string]any `json:"metadata"`
	OccurredAt   time.Time      `json:"occurred_at"`
}

type LoggingPublisher struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewLoggingPublisher(amqpURL string) (*LoggingPublisher, error) {

	// Membuat koneksi RabbitMQ.
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal terhubung ke RabbitMQ logging: %w", err,  
		)
	}

	// Membuat channel khusus logging.
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf(
			"gagal membuat channel logging: %w", err,
		)
	}

	// Membuat exchange logging.
	err = channel.ExchangeDeclare(
		LoggingExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf(
			"gagal membuat logging exchange: %w", err,
		)
	}

	return &LoggingPublisher{
		Connection: conn,
		Channel:    channel,
	}, nil
}

func (p *LoggingPublisher) PublishLogging(
	ctx context.Context,
	logData LoggingMessage,
) error {

	if p.Channel == nil || p.Channel.IsClosed() {
		return fmt.Errorf("RabbitMQ logging channel tidak tersedia")
	}

	// Mengisi waktu kejadian jika belum tersedia.
	if logData.OccurredAt.IsZero() {
		logData.OccurredAt = time.Now()
	}

	// Mengubah data logging menjadi JSON.
	body, err := json.Marshal(logData)
	if err != nil {
		return fmt.Errorf(
			"gagal mengubah logging menjadi JSON: %w", err,
		)
	}

	// Mengirim pesan ke exchange.
	err = p.Channel.PublishWithContext(
		ctx,
		LoggingExchange,
		LoggingRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)

	if err != nil {
		return fmt.Errorf(
			"gagal mengirim logging: %w", err,
		)
	}

	return nil
}

func (p *LoggingPublisher) Close() {
	if p.Channel != nil {
		p.Channel.Close()
	}

	if p.Connection != nil {
		p.Connection.Close()
	}
}
