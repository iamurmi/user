package main

import (
	"context"
	"log"
	"net"

	user "github.com/iamurmi/user/domain/protobuf"
	repocache "github.com/iamurmi/user/repoCache"
	"github.com/iamurmi/user/repositories"
	userservice "github.com/iamurmi/user/userService"

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
