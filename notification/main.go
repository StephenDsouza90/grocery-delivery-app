package main

import (
	"context"
	"log"

	k "github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	r "github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
	h "github.com/StephenDsouza90/grocery-delivery-app/notification/handler"
)

func main() {
	db := r.ConnectToNeo4jDatabase()

	consumer := k.InitializeConsumer(k.Brokers, k.NotificationGroupID)

	repo := r.NewNeo4jDBRepository(db)
	handler := h.NewHandler(repo, consumer)

	ctx := context.Background()
	go consumer.Consume(ctx, []string{k.OrderCreatedTopic, k.PaymentStatusTopic, k.DeliveryStatusTopic}, handler)

	log.Println("Notification service started")

	select {}
}
