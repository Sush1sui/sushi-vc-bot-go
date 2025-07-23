package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/common"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/config"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository/mongodb"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/server"
)

func main() {
	err := config.New()
	if err != nil {
		fmt.Println("Error initializing configuration:", err)
	}

	mongoClient := config.MongoConnection()
	defer mongoClient.Disconnect(context.Background())
	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		panic(fmt.Errorf("failed to connect to MongoDB: %w", err))
	}

	categoryJTCCollection := mongoClient.Database(config.GlobalConfig.MongoDBName).Collection(config.GlobalConfig.CategoryJTCCollectionName)
	customVcCollection := mongoClient.Database(config.GlobalConfig.MongoDBName).Collection(config.GlobalConfig.CustomVcCollectionName)

	repository.CategoryJTCService = &mongodb.MongoClient{
		Client: categoryJTCCollection,
	}
	repository.CustomVcService = &mongodb.MongoClient{
		Client: customVcCollection,
	}


	addr := fmt.Sprintf(":%s", config.GlobalConfig.Port)
	router := server.NewRouter()
	fmt.Printf("Server is listening on Port: %s\n", config.GlobalConfig.Port)

	go func() {
		if err := http.ListenAndServe(addr, router); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	go func() {
		common.PingServerLoop(config.GlobalConfig.ServerUrl)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Shutting down server gracefully...")
}