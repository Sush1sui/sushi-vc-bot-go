package common

import (
	"github.com/Sush1sui/sushi-vc-bot-go/internal/models"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
)

func GetEffectiveGuildSettings(guildID string) *models.GuildSettingsModel {
	settings := &models.GuildSettingsModel{
		GuildID:               guildID,
		VoiceExceptionRoleIDs: []string{},
		ProtectedRoleIDs:      []string{},
		IgnoredChannelIDs:     []string{},
	}

	if repository.GuildSettingsService == nil {
		return settings
	}

	fromDB, err := repository.GuildSettingsService.GetByGuildID(guildID)
	if err != nil || fromDB == nil {
		return settings
	}

	if len(fromDB.VoiceExceptionRoleIDs) > 0 {
		settings.VoiceExceptionRoleIDs = uniqueNonEmpty(fromDB.VoiceExceptionRoleIDs)
	}

	protected := append([]string{}, fromDB.ProtectedRoleIDs...)
	protected = append(protected, settings.VoiceExceptionRoleIDs...)
	settings.ProtectedRoleIDs = uniqueNonEmpty(protected)
	settings.IgnoredChannelIDs = uniqueNonEmpty(fromDB.IgnoredChannelIDs)

	return settings
}

func ContainsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func uniqueNonEmpty(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}
