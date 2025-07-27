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

	category, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name: "VC",
		Type: discordgo.ChannelTypeGuildCategory,
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
			Content: fmt.Sprintf("JTC initialized successfully in **%s** category.", category.Name),
	}})
	if err != nil {
		fmt.Println("Error responding to interaction:", err)
		return
	}
}