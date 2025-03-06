package main

import (
	"context"
	"log"

	"github.com/StephenDsouza90/grocery-delivery-app/cmd/utils"
	"github.com/StephenDsouza90/grocery-delivery-app/internal/delivery"
	"github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	"github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
)

func main() {
	utils.LoadEnv()

	db := repository.ConnectToDatabase()

	repository.AutoMigrate(db, &repository.Delivery{})

	producer := kafka.InitializeKafkaProducer(kafka.Brokers, kafka.DeliveryStatusTopic)
	consumer := kafka.InitializeKafkaConsumer(kafka.Brokers, kafka.PaymentGroupID)

	repo := repository.NewDBRepository(db)
	handler := delivery.NewHandler(repo, producer, consumer)

	ctx := context.Background()
	go consumer.Consume(ctx, []string{kafka.PaymentStatusTopic, kafka.OrderCreatedTopic}, handler)

	log.Println("Delivery service started")

	select {}

	// TODO : API for delivery service to say if order is completed

}
