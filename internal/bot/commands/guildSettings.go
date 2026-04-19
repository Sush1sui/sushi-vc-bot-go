package commands

import (
	"fmt"
	"strings"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/common"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

const (
	settingsTargetVoiceExceptions = "voice_exception_roles"
	settingsTargetProtectedRoles  = "protected_roles"
	settingsTargetIgnoredChannels = "ignored_channels"
)

func ManageGuildSettings(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}

	if repository.GuildSettingsService == nil {
		respondEphemeral(s, i, "Guild settings service is not initialized.")
		return
	}

	data := i.ApplicationCommandData()
	if len(data.Options) == 0 {
		respondEphemeral(s, i, "Usage: /guild-settings <create|update|delete|view|reset>")
		return
	}

	sub := data.Options[0]
	switch sub.Name {
	case "view":
		handleViewGuildSettings(s, i)
	case "create":
		handleCreateGuildSettings(s, i, sub)
	case "update":
		handleUpdateGuildSettings(s, i, sub)
	case "delete":
		handleDeleteGuildSettingsValue(s, i, sub)
	case "reset":
		handleResetGuildSettings(s, i)
	default:
		respondEphemeral(s, i, "Unknown subcommand.")
	}
}

func handleViewGuildSettings(s *discordgo.Session, i *discordgo.InteractionCreate) {
	settings := common.GetEffectiveGuildSettings(i.GuildID)
	content := fmt.Sprintf(
		"Guild Settings\nvoice_exception_roles: %s\nprotected_roles: %s\nignored_channels: %s",
		formatList(settings.VoiceExceptionRoleIDs),
		formatList(settings.ProtectedRoleIDs),
		formatList(settings.IgnoredChannelIDs),
	)
	respondEphemeral(s, i, content)
}

func handleCreateGuildSettings(s *discordgo.Session, i *discordgo.InteractionCreate, sub *discordgo.ApplicationCommandInteractionDataOption) {
	existing, err := repository.GuildSettingsService.GetByGuildID(i.GuildID)
	if err != nil {
		respondEphemeral(s, i, "Failed to check existing guild settings.")
		return
	}
	if existing != nil {
		respondEphemeral(s, i, "Guild settings already exist. Use /guild-settings update instead.")
		return
	}

	target := getSubOptionString(sub, "target")
	ids := parseCSVUnique(getSubOptionString(sub, "ids"))
	if len(ids) == 0 {
		respondEphemeral(s, i, "Please provide at least one ID.")
		return
	}

	voice, protected, ignored := defaults()
	voice, protected, ignored = applySet(target, ids, voice, protected, ignored)
	voice, protected = normalizeRoleLists(voice, protected)

	_, err = repository.GuildSettingsService.UpsertByGuildID(i.GuildID, voice, protected, ignored)
	if err != nil {
		respondEphemeral(s, i, "Failed to create guild settings.")
		return
	}

	respondEphemeral(s, i, "Guild settings created.")
}

func handleUpdateGuildSettings(s *discordgo.Session, i *discordgo.InteractionCreate, sub *discordgo.ApplicationCommandInteractionDataOption) {
	existing, err := repository.GuildSettingsService.GetByGuildID(i.GuildID)
	if err != nil {
		respondEphemeral(s, i, "Failed to retrieve guild settings.")
		return
	}
	if existing == nil {
		respondEphemeral(s, i, "No guild settings found. Use /guild-settings create first.")
		return
	}

	target := getSubOptionString(sub, "target")
	ids := parseCSVUnique(getSubOptionString(sub, "ids"))
	if len(ids) == 0 {
		respondEphemeral(s, i, "Please provide at least one ID.")
		return
	}

	voice := uniqueNonEmpty(existing.VoiceExceptionRoleIDs)
	protected := uniqueNonEmpty(existing.ProtectedRoleIDs)
	ignored := uniqueNonEmpty(existing.IgnoredChannelIDs)

	voice, protected, ignored = applySet(target, ids, voice, protected, ignored)
	voice, protected = normalizeRoleLists(voice, protected)

	_, err = repository.GuildSettingsService.UpsertByGuildID(i.GuildID, voice, protected, ignored)
	if err != nil {
		respondEphemeral(s, i, "Failed to update guild settings.")
		return
	}

	respondEphemeral(s, i, "Guild settings updated.")
}

func handleDeleteGuildSettingsValue(s *discordgo.Session, i *discordgo.InteractionCreate, sub *discordgo.ApplicationCommandInteractionDataOption) {
	existing, err := repository.GuildSettingsService.GetByGuildID(i.GuildID)
	if err != nil {
		respondEphemeral(s, i, "Failed to retrieve guild settings.")
		return
	}
	if existing == nil {
		respondEphemeral(s, i, "No guild settings found.")
		return
	}

	target := getSubOptionString(sub, "target")
	rawIDs := getSubOptionString(sub, "ids")

	voice := uniqueNonEmpty(existing.VoiceExceptionRoleIDs)
	protected := uniqueNonEmpty(existing.ProtectedRoleIDs)
	ignored := uniqueNonEmpty(existing.IgnoredChannelIDs)

	if strings.TrimSpace(rawIDs) == "" {
		voice, protected, ignored = applySet(target, []string{}, voice, protected, ignored)
	} else {
		ids := parseCSVUnique(rawIDs)
		if len(ids) == 0 {
			respondEphemeral(s, i, "Please provide valid IDs to delete.")
			return
		}
		voice, protected, ignored = applyDelete(target, ids, voice, protected, ignored)
	}

	voice, protected = normalizeRoleLists(voice, protected)

	_, err = repository.GuildSettingsService.UpsertByGuildID(i.GuildID, voice, protected, ignored)
	if err != nil {
		respondEphemeral(s, i, "Failed to delete values from guild settings.")
		return
	}

	respondEphemeral(s, i, "Guild settings updated after delete operation.")
}

func handleResetGuildSettings(s *discordgo.Session, i *discordgo.InteractionCreate) {
	deleted, err := repository.GuildSettingsService.DeleteByGuildID(i.GuildID)
	if err != nil {
		respondEphemeral(s, i, "Failed to reset guild settings.")
		return
	}
	if deleted == 0 {
		respondEphemeral(s, i, "No guild settings document found to reset.")
		return
	}

	respondEphemeral(s, i, "Guild settings reset. Bot will use fallback defaults.")
}

func applySet(target string, ids, voice, protected, ignored []string) ([]string, []string, []string) {
	switch target {
	case settingsTargetVoiceExceptions:
		voice = uniqueNonEmpty(ids)
	case settingsTargetProtectedRoles:
		protected = uniqueNonEmpty(ids)
	case settingsTargetIgnoredChannels:
		ignored = uniqueNonEmpty(ids)
	}
	return voice, protected, ignored
}

func applyDelete(target string, ids, voice, protected, ignored []string) ([]string, []string, []string) {
	switch target {
	case settingsTargetVoiceExceptions:
		voice = removeAll(voice, ids)
	case settingsTargetProtectedRoles:
		protected = removeAll(protected, ids)
	case settingsTargetIgnoredChannels:
		ignored = removeAll(ignored, ids)
	}
	return voice, protected, ignored
}

func defaults() ([]string, []string, []string) {
	return []string{}, []string{}, []string{}
}

func normalizeRoleLists(voice, protected []string) ([]string, []string) {
	voice = uniqueNonEmpty(voice)
	protected = uniqueNonEmpty(append(protected, voice...))
	return voice, protected
}

func getSubOptionString(sub *discordgo.ApplicationCommandInteractionDataOption, name string) string {
	for _, option := range sub.Options {
		if option.Name == name {
			return option.StringValue()
		}
	}
	return ""
}

func parseCSVUnique(value string) []string {
	parts := strings.Split(value, ",")
	ids := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			ids = append(ids, trimmed)
		}
	}
	return uniqueNonEmpty(ids)
}

func uniqueNonEmpty(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}

func removeAll(values, remove []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	blocked := make(map[string]struct{}, len(remove))
	for _, value := range remove {
		blocked[value] = struct{}{}
	}

	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := blocked[value]; ok {
			continue
		}
		result = append(result, value)
	}
	return result
}

func formatList(values []string) string {
	if len(values) == 0 {
		return "(empty)"
	}
	return strings.Join(values, ", ")
}

func respondEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
