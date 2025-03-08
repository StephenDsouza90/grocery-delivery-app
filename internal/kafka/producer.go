package kafka

import (
	"encoding/json"
	"log"
	"os"

	"github.com/IBM/sarama"
	r "github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
	u "github.com/StephenDsouza90/grocery-delivery-app/internal/utils"
	"github.com/fatih/structs"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
	topic         string
}

// Initialize Kafka producer
func InitializeProducer(brokers string, topic string) *Producer {
	kafkaProducer, err := NewProducer(
		[]string{os.Getenv(brokers)},
		topic,
	)
	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}
	return kafkaProducer
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		asyncProducer: producer,
		topic:         topic,
	}, nil
}

// SendOrderMessage sends an order created event to the Kafka topic
func (p *Producer) SendOrderMessage(order r.Order) {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(u.ConverterIntToStr(order.OrderID)),
		Value: sarama.StringEncoder(toJSONOrder(order)),
	}
	p.asyncProducer.Input() <- msg
}

// SendPaymentMessage sends a payment status event to the Kafka topic
func (p *Producer) SendPaymentMessage(payment r.Payment) {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(u.ConverterIntToStr(payment.OrderID)),
		Value: sarama.StringEncoder(toJSONPayment(payment)),
	}
	p.asyncProducer.Input() <- msg
}

// SendDeliveryMessage sends a delivery status event to the Kafka topic
func (p *Producer) SendDeliveryMessage(delivery r.Delivery) {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(u.ConverterIntToStr(delivery.OrderID)),
		Value: sarama.StringEncoder(toJSONDelivery(delivery)),
	}
	p.asyncProducer.Input() <- msg
}

// Helper function (implement proper JSON marshaling)
func toJSONOrder(o r.Order) string {
	orderMap := structs.Map(o)
	jsonData, err := json.Marshal(orderMap)
	if err != nil {
		log.Printf("Error marshaling order to JSON: %v", err)
		return ""
	}
	return string(jsonData)
}

// Helper function (implement proper JSON marshaling)
func toJSONPayment(p r.Payment) string {
	orderMap := structs.Map(p)
	jsonData, err := json.Marshal(orderMap)
	if err != nil {
		log.Printf("Error marshaling order to JSON: %v", err)
		return ""
	}
	return string(jsonData)
}

// Helper function (implement proper JSON marshaling)
func toJSONDelivery(d r.Delivery) string {
	orderMap := structs.Map(d)
	jsonData, err := json.Marshal(orderMap)
	if err != nil {
		log.Printf("Error marshaling order to JSON: %v", err)
		return ""
	}
	return string(jsonData)
}
