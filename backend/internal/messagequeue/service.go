package messagequeue

import (
	"context"
	"fmt"
	"log"
	"time"
)

// MessageType represents the type of message
type MessageType string

const (
	// Market data message types
	MarketDataQuote     MessageType = "market.quote"
	MarketDataDepth     MessageType = "market.depth"
	MarketDataOHLC      MessageType = "market.ohlc"
	MarketDataTrade     MessageType = "market.trade"
	
	// Order message types
	OrderNew            MessageType = "order.new"
	OrderUpdate         MessageType = "order.update"
	OrderCancel         MessageType = "order.cancel"
	OrderExecution      MessageType = "order.execution"
	
	// Portfolio message types
	PortfolioUpdate     MessageType = "portfolio.update"
	PortfolioPosition   MessageType = "portfolio.position"
	
	// Strategy message types
	StrategySignal      MessageType = "strategy.signal"
	StrategyExecution   MessageType = "strategy.execution"
	
	// System message types
	SystemAlert         MessageType = "system.alert"
	SystemNotification  MessageType = "system.notification"
)

// Message represents a message in the system
type Message struct {
	Type      MessageType     `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   interface{}     `json:"payload"`
}

// MessageBroker is an interface for message brokers
type MessageBroker interface {
	// Publish publishes a message to a topic
	Publish(ctx context.Context, topic string, message interface{}) error
	
	// Subscribe subscribes to a topic and returns a channel for messages
	Subscribe(ctx context.Context, topic string, handler func([]byte) error) error
	
	// Close closes the connection
	Close() error
}

// MessageService manages message brokers
type MessageService struct {
	redis    *RedisClient
	rabbitmq *RabbitMQClient
}

// NewMessageService creates a new message service
func NewMessageService(redisConfig RedisConfig, rabbitmqConfig RabbitMQConfig) (*MessageService, error) {
	redis, err := NewRedisClient(redisConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %w", err)
	}
	
	rabbitmq, err := NewRabbitMQClient(rabbitmqConfig)
	if err != nil {
		redis.Close()
		return nil, fmt.Errorf("failed to create RabbitMQ client: %w", err)
	}
	
	// Declare exchanges
	exchanges := []struct {
		name string
		kind string
	}{
		{"market.data", "topic"},
		{"order.events", "topic"},
		{"portfolio.events", "topic"},
		{"strategy.events", "topic"},
		{"system.events", "topic"},
	}
	
	for _, exchange := range exchanges {
		err := rabbitmq.DeclareExchange(exchange.name, exchange.kind, true, false, false, false)
		if err != nil {
			redis.Close()
			rabbitmq.Close()
			return nil, fmt.Errorf("failed to declare exchange %s: %w", exchange.name, err)
		}
	}
	
	log.Println("Message service initialized successfully")
	return &MessageService{
		redis:    redis,
		rabbitmq: rabbitmq,
	}, nil
}

// Close closes all connections
func (s *MessageService) Close() error {
	var redisErr, rabbitErr error
	
	if s.redis != nil {
		redisErr = s.redis.Close()
	}
	
	if s.rabbitmq != nil {
		rabbitErr = s.rabbitmq.Close()
	}
	
	if redisErr != nil {
		return redisErr
	}
	
	return rabbitErr
}

// PublishMarketData publishes market data
func (s *MessageService) PublishMarketData(ctx context.Context, msgType MessageType, data interface{}) error {
	message := Message{
		Type:      msgType,
		Timestamp: time.Now(),
		Payload:   data,
	}
	
	// For high-frequency market data, use Redis for performance
	err := s.redis.Publish(ctx, string(msgType), message)
	if err != nil {
		return err
	}
	
	// Also publish to RabbitMQ for durable processing
	topic := fmt.Sprintf("market.data.%s", msgType[7:]) // Remove "market." prefix
	return s.rabbitmq.Publish("market.data", topic, false, false, message)
}

// PublishOrderEvent publishes an order event
func (s *MessageService) PublishOrderEvent(ctx context.Context, msgType MessageType, data interface{}) error {
	message := Message{
		Type:      msgType,
		Timestamp: time.Now(),
		Payload:   data,
	}
	
	// Order events need durability, use RabbitMQ
	topic := fmt.Sprintf("order.events.%s", msgType[6:]) // Remove "order." prefix
	return s.rabbitmq.Publish("order.events", topic, true, false, message)
}

// PublishPortfolioEvent publishes a portfolio event
func (s *MessageService) PublishPortfolioEvent(ctx context.Context, msgType MessageType, data interface{}) error {
	message := Message{
		Type:      msgType,
		Timestamp: time.Now(),
		Payload:   data,
	}
	
	// Portfolio events need durability, use RabbitMQ
	topic := fmt.Sprintf("portfolio.events.%s", msgType[10:]) // Remove "portfolio." prefix
	return s.rabbitmq.Publish("portfolio.events", topic, true, false, message)
}

// PublishStrategyEvent publishes a strategy event
func (s *MessageService) PublishStrategyEvent(ctx context.Context, msgType MessageType, data interface{}) error {
	message := Message{
		Type:      msgType,
		Timestamp: time.Now(),
		Payload:   data,
	}
	
	// Strategy events need durability, use RabbitMQ
	topic := fmt.Sprintf("strategy.events.%s", msgType[9:]) // Remove "strategy." prefix
	return s.rabbitmq.Publish("strategy.events", topic, true, false, message)
}

// PublishSystemEvent publishes a system event
func (s *MessageService) PublishSystemEvent(ctx context.Context, msgType MessageType, data interface{}) error {
	message := Message{
		Type:      msgType,
		Timestamp: time.Now(),
		Payload:   data,
	}
	
	// System events need durability, use RabbitMQ
	topic := fmt.Sprintf("system.events.%s", msgType[7:]) // Remove "system." prefix
	return s.rabbitmq.Publish("system.events", topic, true, false, message)
}

// SubscribeMarketData subscribes to market data
func (s *MessageService) SubscribeMarketData(ctx context.Context, msgType MessageType, handler func([]byte) error) error {
	// For high-frequency market data, use Redis for performance
	redisChannel := string(msgType)
	messages, err := s.redis.Subscribe(ctx, redisChannel)
	if err != nil {
		return err
	}
	
	go func() {
		for msg := range messages {
			if err := handler([]byte(msg.Payload)); err != nil {
				log.Printf("Error handling Redis message: %v", err)
			}
		}
	}()
	
	return nil
}

// SubscribeOrderEvents subscribes to order events
func (s *MessageService) SubscribeOrderEvents(ctx context.Context, msgType MessageType, consumer string, handler func([]byte) error) error {
	// Order events need durability, use RabbitMQ
	queueName := fmt.Sprintf("order.%s.%s", msgType[6:], consumer) // Remove "order." prefix
	routingKey := fmt.Sprintf("order.events.%s", msgType[6:])      // Remove "order." prefix
	
	queue, err := s.rabbitmq.DeclareQueue(queueName, true, false, false, false)
	if err != nil {
		return err
	}
	
	err = s.rabbitmq.BindQueue(queue.Name, routingKey, "order.events", false)
	if err != nil {
		return err
	}
	
	go func() {
		if err := s.rabbitmq.ProcessMessages(ctx, queue.Name, consumer, handler); err != nil {
			log.Printf("Error processing RabbitMQ messages: %v", err)
		}
	}()
	
	return nil
}

// SubscribePortfolioEvents subscribes to portfolio events
func (s *MessageService) SubscribePortfolioEvents(ctx context.Context, msgType MessageType, consumer string, handler func([]byte) error) error {
	// Portfolio events need durability, use RabbitMQ
	queueName := fmt.Sprintf("portfolio.%s.%s", msgType[10:], consumer) // Remove "portfolio." prefix
	routingKey := fmt.Sprintf("portfolio.events.%s", msgType[10:])      // Remove "portfolio." prefix
	
	queue, err := s.rabbitmq.DeclareQueue(queueName, true, false, false, false)
	if err != nil {
		return err
	}
	
	err = s.rabbitmq.BindQueue(queue.Name, routingKey, "portfolio.events", false)
	if err != nil {
		return err
	}
	
	go func() {
		if err := s.rabbitmq.ProcessMessages(ctx, queue.Name, consumer, handler); err != nil {
			log.Printf("Error processing RabbitMQ messages: %v", err)
		}
	}()
	
	return nil
}

// SubscribeStrategyEvents subscribes to strategy events
func (s *MessageService) SubscribeStrategyEvents(ctx context.Context, msgType MessageType, consumer string, handler func([]byte) error) error {
	// Strategy events need durability, use RabbitMQ
	queueName := fmt.Sprintf("strategy.%s.%s", msgType[9:], consumer) // Remove "strategy." prefix
	routingKey := fmt.Sprintf("strategy.events.%s", msgType[9:])      // Remove "strategy." prefix
	
	queue, err := s.rabbitmq.DeclareQueue(queueName, true, false, false, false)
	if err != nil {
		return err
	}
	
	err = s.rabbitmq.BindQueue(queue.Name, routingKey, "strategy.events", false)
	if err != nil {
		return err
	}
	
	go func() {
		if err := s.rabbitmq.ProcessMessages(ctx, queue.Name, consumer, handler); err != nil {
			log.Printf("Error processing RabbitMQ messages: %v", err)
		}
	}()
	
	return nil
}

// SubscribeSystemEvents subscribes to system events
func (s *MessageService) SubscribeSystemEvents(ctx context.Context, msgType MessageType, consumer string, handler func([]byte) error) error {
	// System events need durability, use RabbitMQ
	queueName := fmt.Sprintf("system.%s.%s", msgType[7:], consumer) // Remove "system." prefix
	routingKey := fmt.Sprintf("system.events.%s", msgType[7:])      // Remove "system." prefix
	
	queue, err := s.rabbitmq.DeclareQueue(queueName, true, false, false, false)
	if err != nil {
		return err
	}
	
	err = s.rabbitmq.BindQueue(queue.Name, routingKey, "system.events", false)
	if err != nil {
		return err
	}
	
	go func() {
		if err := s.rabbitmq.ProcessMessages(ctx, queue.Name, consumer, handler); err != nil {
			log.Printf("Error processing RabbitMQ messages: %v", err)
		}
	}()
	
	return nil
}
