package main

import (
	"context"
	"log"

	"github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	"github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
	"github.com/StephenDsouza90/grocery-delivery-app/payment/handler"
)

func main() {
	db := repository.ConnectToDatabase()

	repository.AutoMigrate(db, &repository.Payment{})

	producer := kafka.InitializeKafkaProducer(kafka.Brokers, kafka.PaymentStatusTopic)
	consumer := kafka.InitializeKafkaConsumer(kafka.Brokers, kafka.OrderGroupID)

	repo := repository.NewDBRepository(db)
	handler := handler.NewHandler(repo, producer, consumer)

	ctx := context.Background()
	go consumer.Consume(ctx, []string{kafka.OrderCreatedTopic}, handler)

	log.Println("Payment service started")

	select {}

}
