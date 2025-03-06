package kafka

import (
	"log"
	"os"

	"github.com/IBM/sarama"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
	topic         string
}

// Initialize Kafka producer
func InitializeKafkaProducer(brokers string, topic string) *Producer {
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
func (producer *Producer) SendOrderMessage(orderData map[string]string) {
	msg := &sarama.ProducerMessage{
		Topic: producer.topic,
		Key:   sarama.StringEncoder(orderData["order_id"]),
		Value: sarama.StringEncoder(toJSONOrder(orderData)),
	}
	producer.asyncProducer.Input() <- msg
}

// SendPaymentMessage sends a payment status event to the Kafka topic
func (producer *Producer) SendPaymentMessage(paymentData map[string]string) {
	msg := &sarama.ProducerMessage{
		Topic: producer.topic,
		Key:   sarama.StringEncoder(paymentData["order_id"]),
		Value: sarama.StringEncoder(toJSONPayment(paymentData)),
	}
	producer.asyncProducer.Input() <- msg
}

// SendDeliveryMessage sends a delivery status event to the Kafka topic
func (producer *Producer) SendDeliveryMessage(deliveryData map[string]string) {
	msg := &sarama.ProducerMessage{
		Topic: producer.topic,
		Key:   sarama.StringEncoder(deliveryData["order_id"]),
		Value: sarama.StringEncoder(toJSONDelivery(deliveryData)),
	}
	producer.asyncProducer.Input() <- msg
}

// Helper function (implement proper JSON marshaling)
func toJSONOrder(o map[string]string) string {
	return `{"order_id": "` + o["order_id"] + `", "customer_id": "` + o["customer_id"] + `", "order_date": "` + o["order_date"] + `", "total_amount": "` + o["total_amount"] + `", "order_status": "` + o["order_status"] + `", "delivery_date": "` + o["delivery_date"] + `", "delivery_time": "` + o["delivery_time"] + `"}`
}

// Helper function (implement proper JSON marshaling)
func toJSONPayment(p map[string]string) string {
	return `{"order_id": "` + p["order_id"] + `", "payment_status": "` + p["payment_status"] + `", "payment_date": "` + p["payment_date"] + `", "payment_id": "` + p["payment_id"] + `"}`
}

// Helper function (implement proper JSON marshaling)
func toJSONDelivery(d map[string]string) string {
	return `{"order_id": "` + d["order_id"] + `", "delivery_status": "` + d["delivery_status"] + `", "delivery_date": "` + d["delivery_date"] + `", "delivery_time": "` + d["delivery_time"] + `"}`
}
