package messagequeue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQConfig holds the RabbitMQ configuration
type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	VHost    string
}

// RabbitMQClient represents a RabbitMQ client
type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQClient creates a new RabbitMQ client
func NewRabbitMQClient(config RabbitMQConfig) (*RabbitMQClient, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.User, config.Password, config.Host, config.Port, config.VHost)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	log.Println("Successfully connected to RabbitMQ")
	return &RabbitMQClient{conn: conn, channel: channel}, nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQClient) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// DeclareQueue declares a queue
func (r *RabbitMQClient) DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name,       // name
		durable,    // durable
		autoDelete, // delete when unused
		exclusive,  // exclusive
		noWait,     // no-wait
		nil,        // arguments
	)
}

// DeclareExchange declares an exchange
func (r *RabbitMQClient) DeclareExchange(name, kind string, durable, autoDelete, internal, noWait bool) error {
	return r.channel.ExchangeDeclare(
		name,       // name
		kind,       // type
		durable,    // durable
		autoDelete, // auto-deleted
		internal,   // internal
		noWait,     // no-wait
		nil,        // arguments
	)
}

// BindQueue binds a queue to an exchange
func (r *RabbitMQClient) BindQueue(queueName, key, exchangeName string, noWait bool) error {
	return r.channel.QueueBind(
		queueName,    // queue name
		key,          // routing key
		exchangeName, // exchange
		noWait,       // no-wait
		nil,          // arguments
	)
}

// Publish publishes a message to an exchange
func (r *RabbitMQClient) Publish(exchange, routingKey string, mandatory, immediate bool, message interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return r.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		mandatory,  // mandatory
		immediate,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
}

// Consume consumes messages from a queue
func (r *RabbitMQClient) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queue,     // queue
		consumer,  // consumer
		autoAck,   // auto-ack
		exclusive, // exclusive
		noLocal,   // no-local
		noWait,    // no-wait
		nil,       // args
	)
}

// ProcessMessages processes messages from a queue with a handler function
func (r *RabbitMQClient) ProcessMessages(ctx context.Context, queueName, consumerName string, handler func([]byte) error) error {
	msgs, err := r.Consume(queueName, consumerName, false, false, false, false)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("channel closed")
			}

			err := handler(msg.Body)
			if err != nil {
				// Nack the message to requeue it
				msg.Nack(false, true)
				log.Printf("Error processing message: %v", err)
			} else {
				// Ack the message
				msg.Ack(false)
			}
		}
	}
}
