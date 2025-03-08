package main

import (
	"context"
	"log"

	h "github.com/StephenDsouza90/grocery-delivery-app/delivery/handler"
	k "github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	r "github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
	"github.com/gin-gonic/gin"
)

const (
	Port = ":8081"
)

func main() {
	db := r.ConnectToDatabase()

	r.AutoMigrate(db, &r.Delivery{})

	producer := k.InitializeProducer(k.Brokers, k.DeliveryStatusTopic)
	consumer := k.InitializeConsumer(k.Brokers, k.PaymentGroupID)

	repo := r.NewDBRepository(db)
	handler := h.NewHandler(repo, producer, consumer)

	ctx := context.Background()
	go consumer.Consume(ctx, []string{k.PaymentStatusTopic, k.OrderCreatedTopic}, handler)

	// Start the server
	router := gin.Default()
	router.POST("/update-delivery", handler.UpdateDelivery)

	if err := router.Run(Port); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}

}
