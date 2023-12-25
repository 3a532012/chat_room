package main

import (
	"context"
	"log"
	_userHandlerHttpDelivery "user/delivery/http"
	_userRepository "user/repository/mongodb"
	_userUsecase "user/usecase/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Println("DB init ...")
	clientOptions := options.Client().ApplyURI("mongodb://test:test@localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("Failed to connect to MongoDB:", err)
	}

	userRepository := _userRepository.NewMongodbUserRepository(client)
	userUsecase := _userUsecase.NewUserUsecase(userRepository)

	r := gin.Default()
	_userHandlerHttpDelivery.NewUserHandler(r, userUsecase)

	log.Fatal(r.Run())
}
