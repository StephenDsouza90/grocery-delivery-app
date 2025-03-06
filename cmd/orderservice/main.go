package main

import (
	"log"

	"github.com/StephenDsouza90/grocery-delivery-app/cmd/utils"
	"github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	"github.com/StephenDsouza90/grocery-delivery-app/internal/order"
	"github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
	"github.com/gin-gonic/gin"
)

const (
	Port = ":8080"
)

func main() {
	utils.LoadEnv()

	db := repository.ConnectToDatabase()

	repository.AutoMigrate(db, &repository.Order{})
	repository.AutoMigrate(db, &repository.Item{})

	producer := kafka.InitializeKafkaProducer(kafka.Brokers, kafka.OrderCreatedTopic)
	consumer := kafka.InitializeKafkaConsumer(kafka.Brokers, kafka.PaymentGroupID)

	repo := repository.NewDBRepository(db)
	handler := order.NewHandler(repo, producer, consumer)

	// Start the server
	router := gin.Default()
	router.POST("/orders", handler.CreateOrder)

	if err := router.Run(Port); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}

}
