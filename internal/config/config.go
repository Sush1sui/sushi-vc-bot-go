package config

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Config struct {
	Port string
	BotToken string
	AppID string
	ServerUrl string
	MongoDBName string
	CategoryJTCCollectionName string
	CustomVcCollectionName string
	FinestRoleId string
}

var GlobalConfig *Config

func New() (error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("BOT_TOKEN is not set in the environment variables")
	}

	appID := os.Getenv("APP_ID")
	if appID == "" {
		return fmt.Errorf("APP_ID is not set in the environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7694"
	}

	mongoDBName := os.Getenv("MONGODB_NAME")
	if mongoDBName == "" {
		return fmt.Errorf("MONGODB_NAME is not set in the environment variables")
	}

	categoryJTCCollectionName := os.Getenv("CATEGORY_JTC_COLLECTION_NAME")
	if categoryJTCCollectionName == "" {
		return fmt.Errorf("CATEGORY_JTC_COLLECTION_NAME is not set in the environment variables")
	}

	customVcCollectionName := os.Getenv("CUSTOM_VC_COLLECTION_NAME")
	if customVcCollectionName == "" {
		return fmt.Errorf("CUSTOM_VC_COLLECTION_NAME is not set in the environment variables")
	}

	finestRoleId := os.Getenv("FINEST_ROLE_ID")
	if finestRoleId == "" {
		return fmt.Errorf("FINEST_ROLE_ID is not set in the environment variables")
	}

	GlobalConfig = &Config{
		Port:     port,
		BotToken: botToken,
		AppID:   appID,
		ServerUrl: os.Getenv("SERVER_URL"),
		MongoDBName: mongoDBName,
		CategoryJTCCollectionName: categoryJTCCollectionName,
		CustomVcCollectionName: customVcCollectionName,
		FinestRoleId: finestRoleId,
	}
	return nil
}

func MongoConnection() *mongo.Client {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
  client, err := mongo.Connect(opts)
  if err != nil {
    panic(err)
  }

  // Send a ping to confirm a successful connection
  if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
    panic(err)
  }
  fmt.Println("DB Connected!")

	return client
}