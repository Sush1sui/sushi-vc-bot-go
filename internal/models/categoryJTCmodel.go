package models

import "go.mongodb.org/mongo-driver/v2/bson"

type CategoryJTCModel struct {
	ID bson.ObjectID `bson:"_id,omitempty"`
	CategoryID      string        `bson:"category_id"`
	JTCChannelID    string        `bson:"jtc_channel_id"`
	InterfaceID      string        `bson:"interface_id"`
	InterfaceMessageID string    `bson:"interface_message_id"`
	CustomVCsID       []string    `bson:"custom_vcs_id"`
}