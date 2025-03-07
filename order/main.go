package main

import (
	"log"

	"github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	"github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
	"github.com/StephenDsouza90/grocery-delivery-app/order/handler"
	"github.com/gin-gonic/gin"
)

const (
	Port = ":8080"
)

func main() {
	db := repository.ConnectToDatabase()

	repository.AutoMigrate(db, &repository.Order{})
	repository.AutoMigrate(db, &repository.Item{})

	producer := kafka.InitializeKafkaProducer(kafka.Brokers, kafka.OrderCreatedTopic)
	consumer := kafka.InitializeKafkaConsumer(kafka.Brokers, kafka.PaymentGroupID)

	repo := repository.NewDBRepository(db)
	handler := handler.NewHandler(repo, producer, consumer)

	// Start the server
	router := gin.Default()
	router.POST("/orders", handler.CreateOrder)

	if err := router.Run(Port); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}

}
