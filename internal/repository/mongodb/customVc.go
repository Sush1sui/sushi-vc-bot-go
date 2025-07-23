package mongodb

import (
	"context"
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (c *MongoClient) GetAllVcs() ([]*models.CustomVcModel, error) {
	var channels []*models.CustomVcModel

	cursor, err := c.Client.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find custom VCs: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var channel models.CustomVcModel
		if err := cursor.Decode(&channel); err != nil {
			return nil, fmt.Errorf("failed to decode custom VC: %w", err)
		}
		channels = append(channels, &channel)
	}
	return channels, nil
}

func (c * MongoClient) CreateVc(ownerId string, channelId string) (*models.CustomVcModel, error) {
	channel := &models.CustomVcModel{
		OwnerID:   ownerId,
		ChannelID: channelId,
	}

	if _, err := c.Client.InsertOne(context.Background(), channel); err != nil {
		return nil, fmt.Errorf("failed to create custom VC: %w", err)
	}
	return channel, nil
}

func (c *MongoClient) GetByOwnerOrChannelId(ownerId string, channelId string) (*models.CustomVcModel, error) {
	var channel models.CustomVcModel
	if ownerId == "" && channelId == "" {
		return nil, fmt.Errorf("either ownerId or channelId must be provided")
	}

	if ownerId != "" {
		if err := c.Client.FindOne(context.Background(), bson.M{"owner_id": ownerId}).Decode(&channel); err != nil {
			return nil, fmt.Errorf("failed to find custom VC by owner ID: %w", err)
		}
	} else if channelId != "" {
		if err := c.Client.FindOne(context.Background(), bson.M{"channel_id": channelId}).Decode(&channel); err != nil {
			return nil, fmt.Errorf("failed to find custom VC by channel ID: %w", err)
		}
	}
	return &channel, nil
}

func (c *MongoClient) DeleteByOwnerOrChannelId(ownerId string, channelId string) (int, error) {
	if ownerId == "" && channelId == "" {
		return 0, fmt.Errorf("either ownerId or channelId must be provided")
	}

	filter := bson.M{}
	if ownerId != "" {
		filter["owner_id"] = ownerId
	} else if channelId != "" {
		filter["channel_id"] = channelId
	}

	result, err := c.Client.DeleteOne(context.Background(), filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete custom VC: %w", err)
	}
	if result.DeletedCount == 0 {
		return 0, fmt.Errorf("no custom VC found to delete")
	}
	return int(result.DeletedCount), nil
}

func (c *MongoClient) ChangeOwnerByChannelId(channelId string, newOwnerId string) (int, error) {
	if channelId == "" || newOwnerId == "" {
		return 0, fmt.Errorf("channelId and newOwnerId must be provided")
	}

	update := bson.M{"$set": bson.M{"owner_id": newOwnerId}}
	result, err := c.Client.UpdateOne(context.Background(), bson.M{"channel_id": channelId}, update)
	if err != nil {
		return 0, fmt.Errorf("failed to change owner of custom VC: %w", err)
	}
	if result.MatchedCount == 0 {
		return 0, fmt.Errorf("no custom VC found with channel ID: %s", channelId)
	}
	return int(result.ModifiedCount), nil
}