package commands

import (
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/common"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/config"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/handler"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func InitializeJTC(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}

	// check if there is an existing category
	existingCategory, err := repository.CategoryJTCService.GetAllJTCs()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create category.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		fmt.Println("Error fetching existing JTC categories:", err)
		return
	}
	if len(existingCategory) > 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "JTC is already initialized.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	category, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name: "VC",
		Type: discordgo.ChannelTypeGuildCategory,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{ // everyone role
				ID: i.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel | discordgo.PermissionVoiceConnect | discordgo.PermissionSendMessages,
			},
			{
				ID: "1299577480868528330", // music bots
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionVoiceConnect | discordgo.PermissionVoiceSpeak,
			},
			{
				ID: config.GlobalConfig.FinestRoleId,
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionCreatePublicThreads | discordgo.PermissionSendMessages | discordgo.PermissionCreatePrivateThreads | discordgo.PermissionSendMessagesInThreads | discordgo.PermissionAddReactions | discordgo.PermissionManageThreads | discordgo.PermissionReadMessageHistory | discordgo.PermissionVoiceSpeak | discordgo.PermissionVoiceStreamVideo | discordgo.PermissionUseEmbeddedActivities, 
			},
			{
				ID: "1303998295911436309", // lvl50
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
			},
			{
				ID: "1303998297538560060", // lvl60
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
			},
			{
				ID: "1303998299031736393", // lvl70
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
			},
			{
				ID: "1303998300671709186", // lvl80
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
			},
			{
				ID: "1303998302785900544", // lvl90
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
			},
			{
				ID: "1303998304710819940", // lvl100
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
			},
			{
				ID: "1303916681692839956", // pioneers
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
			},
			{
				ID: "1303924607555997776", // supporter
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles,
			},
			{
				ID: "1292420325002448930", // booster
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles | discordgo.PermissionUseExternalEmojis | discordgo.PermissionUseExternalStickers,
			},
			{
				ID: "1310186525606154340", // staff
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionCreateInstantInvite | discordgo.PermissionEmbedLinks | discordgo.PermissionAttachFiles | discordgo.PermissionUseExternalEmojis | discordgo.PermissionUseExternalStickers,
			},
		},
	})
	if err != nil {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create category.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}

	jtcChannel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name: "Join to Create",
		Type: discordgo.ChannelTypeGuildVoice,
		ParentID: category.ID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{ // everyone role
				ID: i.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel | discordgo.PermissionVoiceConnect | discordgo.PermissionSendMessages,
			},
			{
				ID: config.GlobalConfig.FinestRoleId,
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionVoiceConnect,
				Deny: discordgo.PermissionSendMessages,
			},
		},
	})
	if err != nil {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create JTC channel.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}

	interfaceChannel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name: "vc-interface",
		Type: discordgo.ChannelTypeGuildText,
		ParentID: category.ID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID: i.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
			},
			{
				ID: config.GlobalConfig.FinestRoleId,
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel,
				Deny: discordgo.PermissionSendMessages | discordgo.PermissionCreateEvents | discordgo.PermissionCreatePublicThreads | discordgo.PermissionCreatePrivateThreads | discordgo.PermissionAddReactions,
			},
		},
	})
	if err != nil {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create interface channel.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
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
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to initialize JTC.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("JTC initialized successfully in **%s** category, interface: **%s**.", category.Name, interfaceChannel.Mention()),
		},
	})
	if err != nil {
		fmt.Println("Error responding to interaction:", err)
		return
	}
}