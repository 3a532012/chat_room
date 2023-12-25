package main

import (
	"context"
	_friendHandlerHttpDelivery "friend/delivery/http"
	_friendRepository "friend/repository/mongodb"
	_friendUsecase "friend/usecase/friend"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Println("DB init ...")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("Failed to connect to MongoDB:", err)
	}

	friendRepository := _friendRepository.NewMongodbFriendRepository(client)
	friendUsecase := _friendUsecase.NewFriendUsecase(friendRepository)

	r := gin.Default()
	_friendHandlerHttpDelivery.NewFriendHandler(r, friendUsecase)

	log.Fatal(r.Run())
}
