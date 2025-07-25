package button

import (
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/config"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func UnlockVC(i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" { return }

	res, err := repository.CustomVcService.GetByOwnerOrChannelId(i.Member.User.ID, "")
	if err != nil || res == nil {
		err := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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

	customVC, err := bot.Session.Channel(res.ChannelID)
	if err != nil || customVC == nil {
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

	err = bot.Session.ChannelPermissionSet(
		customVC.ID,
		i.GuildID,
		discordgo.PermissionOverwriteTypeRole,
		discordgo.PermissionVoiceConnect | discordgo.PermissionReadMessageHistory | discordgo.PermissionSendMessages,
		0,
	)
	if err != nil {
		err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to unlock the voice channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}

	err = bot.Session.ChannelPermissionSet(
		customVC.ID,
		config.GlobalConfig.FinestRoleId,
		discordgo.PermissionOverwriteTypeRole,
		discordgo.PermissionVoiceConnect | discordgo.PermissionReadMessageHistory | discordgo.PermissionSendMessages,
		0,
	)
	if err != nil {
		err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to unlock the voice channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}

	editedChannel, err := bot.Session.ChannelEdit(
		customVC.ID,
		&discordgo.ChannelEdit{
			Name: fmt.Sprintf("%s's VC", i.Member.User.Username),
		},
	)
	if err != nil || editedChannel == nil {
		err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to rename the voice channel due to hitting Discord API's rate limit, but the channel is successfully unlocked\nIf you want to rename the voice channel via interface, please do so again in 15 minutes or just manually rename the channel yourself.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
	}

	err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Voice channel unlocked successfully.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Println("Error responding to interaction:", err)
	}
}