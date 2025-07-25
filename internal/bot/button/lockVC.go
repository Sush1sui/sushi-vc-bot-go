package button

import (
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/config"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func LockVC(i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" { return }

	res, err := repository.CustomVcService.GetByOwnerOrChannelId(i.Member.User.ID, "")
	if err != nil || res == nil {
		err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You are not the owner of a custom voice channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}
	
	customVc, err := bot.Session.Channel(res.ChannelID)
	if err != nil || customVc == nil {
		err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Custom VC not found.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}

	// Edit the channel permissions to lock it
	err = bot.Session.ChannelPermissionSet(
		customVc.ID,
		i.GuildID,
		discordgo.PermissionOverwriteTypeRole,
		0,
		discordgo.PermissionVoiceConnect | discordgo.PermissionReadMessageHistory | discordgo.PermissionSendMessages,
	)
	if err != nil {
		err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to lock the voice channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}

	err = bot.Session.ChannelPermissionSet(
		customVc.ID,
		config.GlobalConfig.FinestRoleId,
		discordgo.PermissionOverwriteTypeRole,
		0,
		discordgo.PermissionVoiceConnect | discordgo.PermissionReadMessageHistory | discordgo.PermissionSendMessages,
	)
	if err != nil {
		err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to lock the voice channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}

	err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Voice channel locked successfully.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Println("Error responding to interaction:", err)
	}
}