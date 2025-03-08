package main

import (
	"log"

	k "github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	r "github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
	h "github.com/StephenDsouza90/grocery-delivery-app/order/handler"
	"github.com/gin-gonic/gin"
)

const (
	Port = ":8080"
)

func main() {
	db := r.ConnectToDatabase()

	r.AutoMigrate(db, &r.Order{})
	r.AutoMigrate(db, &r.Item{})

	producer := k.InitializeProducer(k.Brokers, k.OrderCreatedTopic)
	consumer := k.InitializeConsumer(k.Brokers, k.PaymentGroupID)

	repo := r.NewDBRepository(db)
	handler := h.NewHandler(repo, producer, consumer)

	// Start the server
	router := gin.Default()
	router.POST("/orders", handler.CreateOrder)

	if err := router.Run(Port); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
