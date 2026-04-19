package commands

import (
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/common"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/handler"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func buildCategoryOverwrites(guildID string, voiceExceptionRoleIDs []string) []*discordgo.PermissionOverwrite {
	overwrites := []*discordgo.PermissionOverwrite{
		{
			ID:   guildID, // @everyone
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: discordgo.PermissionSendMessages,
		},
	}

	for _, roleID := range voiceExceptionRoleIDs {
		overwrites = append(overwrites, &discordgo.PermissionOverwrite{
			ID:   roleID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Allow: discordgo.PermissionCreateInstantInvite |
				discordgo.PermissionCreatePublicThreads |
				discordgo.PermissionSendMessages |
				discordgo.PermissionCreatePrivateThreads |
				discordgo.PermissionSendMessagesInThreads |
				discordgo.PermissionAddReactions |
				discordgo.PermissionManageThreads |
				discordgo.PermissionReadMessageHistory |
				discordgo.PermissionVoiceSpeak |
				discordgo.PermissionVoiceStreamVideo |
				discordgo.PermissionUseEmbeddedActivities |
				discordgo.PermissionViewChannel |
				discordgo.PermissionVoiceConnect,
		})
	}

	return overwrites
}

func buildJoinToCreateOverwrites(guildID string, voiceExceptionRoleIDs []string) []*discordgo.PermissionOverwrite {
	overwrites := []*discordgo.PermissionOverwrite{
		{
			ID:   guildID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: discordgo.PermissionViewChannel | discordgo.PermissionVoiceConnect | discordgo.PermissionSendMessages,
		},
	}

	for _, roleID := range voiceExceptionRoleIDs {
		overwrites = append(overwrites, &discordgo.PermissionOverwrite{
			ID:    roleID,
			Type:  discordgo.PermissionOverwriteTypeRole,
			Allow: discordgo.PermissionViewChannel | discordgo.PermissionVoiceConnect,
			Deny:  discordgo.PermissionSendMessages,
		})
	}

	return overwrites
}

func buildInterfaceOverwrites(guildID string, voiceExceptionRoleIDs []string) []*discordgo.PermissionOverwrite {
	overwrites := []*discordgo.PermissionOverwrite{
		{
			ID:   guildID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
		},
	}

	for _, roleID := range voiceExceptionRoleIDs {
		overwrites = append(overwrites, &discordgo.PermissionOverwrite{
			ID:    roleID,
			Type:  discordgo.PermissionOverwriteTypeRole,
			Allow: discordgo.PermissionViewChannel,
			Deny: discordgo.PermissionSendMessages |
				discordgo.PermissionCreateEvents |
				discordgo.PermissionCreatePublicThreads |
				discordgo.PermissionCreatePrivateThreads |
				discordgo.PermissionAddReactions,
		})
	}

	return overwrites
}

func InitializeJTC(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}

	// acknowledge immediately so we have time to do work
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	// check if there is an existing category
	existingCategory, err := repository.CategoryJTCService.GetAllJTCs()
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func(s string) *string { return &s }("Failed to check existing JTC categories."),
		})
		return
	}
	if len(existingCategory) > 0 {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func(s string) *string { return &s }("A JTC category already exists. Please delete the existing one before creating a new one."),
		})
		return
	}

	settings := common.GetEffectiveGuildSettings(i.GuildID)

	category, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name:                 "VC",
		Type:                 discordgo.ChannelTypeGuildCategory,
		PermissionOverwrites: buildCategoryOverwrites(i.GuildID, settings.VoiceExceptionRoleIDs),
	})
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func(s string) *string { return &s }("Failed to create category channel."),
		})
		return
	}

	jtcChannel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name:                 "Join to Create",
		Type:                 discordgo.ChannelTypeGuildVoice,
		ParentID:             category.ID,
		PermissionOverwrites: buildJoinToCreateOverwrites(i.GuildID, settings.VoiceExceptionRoleIDs),
	})
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func(s string) *string { return &s }("Failed to create JTC voice channel."),
		})
		return
	}

	interfaceChannel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name:                 "vc-interface",
		Type:                 discordgo.ChannelTypeGuildText,
		ParentID:             category.ID,
		PermissionOverwrites: buildInterfaceOverwrites(i.GuildID, settings.VoiceExceptionRoleIDs),
	})
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func(s string) *string { return &s }("Failed to create interface text channel."),
		})
		return
	}

	interfaceMessage, err := s.ChannelMessageSendComplex(interfaceChannel.ID, &discordgo.MessageSend{
		Embed: common.InterfaceEmbed(),
		Components: append(
			common.InterfaceButtonsRow1(),
			append(
				common.InterfaceButtonsRow2(),
				common.InterfaceButtonsRow3()...,
			)...,
		),
	})
	if err != nil {
		fmt.Println("Error sending interface message:", err)
		return
	}

	// setup handlers
	s.AddHandler(handler.InteractionHandler)

	res, err := repository.CategoryJTCService.CreateCategoryJTC(interfaceChannel.ID, interfaceMessage.ID, jtcChannel.ID, category.ID)
	if err != nil || res == nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: func(s string) *string { return &s }("Failed to save JTC configuration to the database."),
		})
		return
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: func(s string) *string { return &s }("JTC has been successfully initialized!"),
	})
}
