package main

import (
	"context"
	"log"
	"net"
	user "small-mic/user/domain/protobuf"
	repocache "small-mic/user/repoCache"
	"small-mic/user/repositories"
	userservice "small-mic/user/userService"

	"github.com/go-redis/redis/v7"
	"google.golang.org/grpc"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// momgodb connection
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	//repo constuctor
	userRepoObj := repositories.NewUserRepoConstructor(client)
	// redis repo constuctor
	redisUserRepoObj := repocache.NewRedisClient(redisClient)
	// svc constructor
	/*
	   As userRepoStruct is a class/struct, which Implements all methods of the UserRepoContracts, So here we can pass the object of that class which implements all methods of a interface contracts ex UserRepoContracts.
	*/
	userSvcObj := userservice.NewServiceConstructor(userRepoObj, redisUserRepoObj)

	// GIN SERVER
	// eng := gin.New()
	// eng.Use(gin.Recovery())                // middleware. it is use for handling Panic
	// userservice.NewRoutes(eng, userSvcObj) // NewRoutes is a Initializer of Transport layer
	// eng.Run(":3000")

	// GRPC SERVER
	listener, err := net.Listen("tcp", "3000")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	user.RegisterUserServer(grpcServer, userSvcObj)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
