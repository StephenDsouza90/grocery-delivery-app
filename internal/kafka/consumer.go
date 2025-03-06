package kafka

import (
	"context"
	"log"
	"os"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
}

type OrderCreatedMessage struct {
	OrderID    int     `json:"order_id"`
	CustomerID int     `json:"customer_id"`
	Amount     float64 `json:"amount"`
}

// Initialize Kafka consumer
func InitializeKafkaConsumer(brokers string, groupID string) *Consumer {
	kafkaConsumer, err := NewConsumer(
		[]string{os.Getenv(brokers)},
		groupID,
	)
	if err != nil {
		log.Fatal("Failed to create Kafka consumer:", err)
	}
	return kafkaConsumer
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, groupID string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0 // Match Kafka version
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumerGroup: consumerGroup,
	}, nil
}

// Consume starts consuming messages from Kafka
func (consumer *Consumer) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) {
	for {
		err := consumer.consumerGroup.Consume(ctx, topics, handler)
		if err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}
}
