package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/joho/godotenv"
	pb "gitlab.com/lemmyGo/lemmyGoUsers/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var db *mongo.Database

func main() {
	godotenv.Load(".env")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	dbUri := os.Getenv("DB_URI")
	fmt.Println(dbUri)
	opts := options.Client().ApplyURI(dbUri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	db = client.Database("Likky")

	fmt.Println("Connected to MongoDB")
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on 8081!")
	serverRegistrar := grpc.NewServer()
	service := &mUserServer{}
	pb.RegisterUsersServer(serverRegistrar, service)
	sErr := serverRegistrar.Serve(lis)
	if err != nil {
		panic(sErr)
	}
}
