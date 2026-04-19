package repository

import "github.com/Sush1sui/sushi-vc-bot-go/internal/models"

type GuildSettingsInterface interface {
	GetByGuildID(guildID string) (*models.GuildSettingsModel, error)
	UpsertByGuildID(guildID string, voiceExceptionRoleIDs, protectedRoleIDs, ignoredChannelIDs []string) (*models.GuildSettingsModel, error)
	DeleteByGuildID(guildID string) (int, error)
}

var GuildSettingsService GuildSettingsInterface
