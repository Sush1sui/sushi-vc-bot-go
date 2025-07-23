package mongodb

import (
	"context"
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (c *MongoClient) GetAllJTCs() ([]*models.CategoryJTCModel, error) {
	var categories []*models.CategoryJTCModel

	cursor, err := c.Client.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find category JTCs: %w", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var category models.CategoryJTCModel
		if err := cursor.Decode(&category); err != nil {
			return nil, fmt.Errorf("failed to decode category JTC: %w", err)
		}
		categories = append(categories, &category)
	}
	return categories, nil
}

func (c *MongoClient) CreateCategoryJTC(interfaceId, interfaceMessageId, jtcChannelId, categoryId string) (*models.CategoryJTCModel, error) {
	category := &models.CategoryJTCModel{
		InterfaceID:        interfaceId,
		InterfaceMessageID: interfaceMessageId,
		JTCChannelID:       jtcChannelId,
		CategoryID:         categoryId,
	}

	if _, err := c.Client.InsertOne(context.Background(), category); err != nil {
		return nil, fmt.Errorf("failed to create category JTC: %w", err)
	}
	return category, nil
}

func (c *MongoClient) DeleteAll() (int, error) {
	result, err := c.Client.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to delete all category JTCs: %w", err)
	}
	return int(result.DeletedCount), nil
}