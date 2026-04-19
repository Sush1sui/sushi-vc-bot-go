package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (c *MongoClient) GetByGuildID(guildID string) (*models.GuildSettingsModel, error) {
	if guildID == "" {
		return nil, fmt.Errorf("guildID must be provided")
	}

	var settings models.GuildSettingsModel
	err := c.Client.FindOne(context.Background(), bson.M{"guild_id": guildID}).Decode(&settings)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find guild settings: %w", err)
	}

	return &settings, nil
}

func (c *MongoClient) UpsertByGuildID(guildID string, voiceExceptionRoleIDs, protectedRoleIDs, ignoredChannelIDs []string) (*models.GuildSettingsModel, error) {
	if guildID == "" {
		return nil, fmt.Errorf("guildID must be provided")
	}

	settings := &models.GuildSettingsModel{
		GuildID:               guildID,
		VoiceExceptionRoleIDs: voiceExceptionRoleIDs,
		ProtectedRoleIDs:      protectedRoleIDs,
		IgnoredChannelIDs:     ignoredChannelIDs,
	}

	update := bson.M{
		"$set": bson.M{
			"guild_id":                 guildID,
			"voice_exception_role_ids": settings.VoiceExceptionRoleIDs,
			"protected_role_ids":       settings.ProtectedRoleIDs,
			"ignored_channel_ids":      settings.IgnoredChannelIDs,
		},
	}

	opts := options.UpdateOne().SetUpsert(true)
	if _, err := c.Client.UpdateOne(context.Background(), bson.M{"guild_id": guildID}, update, opts); err != nil {
		return nil, fmt.Errorf("failed to upsert guild settings: %w", err)
	}

	return c.GetByGuildID(guildID)
}

func (c *MongoClient) DeleteByGuildID(guildID string) (int, error) {
	if guildID == "" {
		return 0, fmt.Errorf("guildID must be provided")
	}

	result, err := c.Client.DeleteOne(context.Background(), bson.M{"guild_id": guildID})
	if err != nil {
		return 0, fmt.Errorf("failed to delete guild settings: %w", err)
	}

	return int(result.DeletedCount), nil
}
