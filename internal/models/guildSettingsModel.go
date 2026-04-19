package models

import "go.mongodb.org/mongo-driver/v2/bson"

type GuildSettingsModel struct {
	ID                    bson.ObjectID `bson:"_id,omitempty"`
	GuildID               string        `bson:"guild_id"`
	VoiceExceptionRoleIDs []string      `bson:"voice_exception_role_ids"`
	ProtectedRoleIDs      []string      `bson:"protected_role_ids"`
	IgnoredChannelIDs     []string      `bson:"ignored_channel_ids"`
}
