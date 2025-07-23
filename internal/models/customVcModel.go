package models

import "go.mongodb.org/mongo-driver/v2/bson"

type CustomVcModel struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	ChannelID string        `bson:"channel_id"`
	OwnerID   string        `bson:"owner_id"`
}