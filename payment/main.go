package main

import (
	"context"
	"log"

	k "github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	r "github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
	h "github.com/StephenDsouza90/grocery-delivery-app/payment/handler"
)

func main() {
	db := r.ConnectToDatabase()

	r.AutoMigrate(db, &r.Payment{})

	producer := k.InitializeProducer(k.Brokers, k.PaymentStatusTopic)
	consumer := k.InitializeConsumer(k.Brokers, k.OrderGroupID)

	repo := r.NewDBRepository(db)
	handler := h.NewHandler(repo, producer, consumer)

	ctx := context.Background()
	go consumer.Consume(ctx, []string{k.OrderCreatedTopic}, handler)

	log.Println("Payment service started")

	select {}
}
