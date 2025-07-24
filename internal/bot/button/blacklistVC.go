package button

import (
	"fmt"
	"strings"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func BlacklistMenu(i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil {
		return
	}

	minValue := 1
	selectMenu := discordgo.SelectMenu{
		MenuType:    discordgo.UserSelectMenu,
		CustomID:   "blacklist_menu",
		Placeholder: "Select users to blacklist",
		MinValues:   &minValue, // wtf discordgo why is this a pointer?
		MaxValues:   5,
	}

	row := discordgo.ActionsRow{Components: []discordgo.MessageComponent{selectMenu}}

	err := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "**Select members to blacklist from this channel**",
			Components: []discordgo.MessageComponent{row},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create blacklist menu.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}
}

func HandleBlacklistSelection(i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil {
		return
	}
	if i.MessageComponentData().CustomID != "blacklist_menu" {
		return
	}

	res, err := repository.CustomVcService.GetByOwnerOrChannelId(i.Member.User.ID, "")
	if err != nil || res == nil {
		e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No custom VC found that you own.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}

	customVc, err := bot.Session.Channel(res.ChannelID)
	if err != nil {
		e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve custom VC channel.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}

	selectedUserIds := i.MessageComponentData().Values
	if len(selectedUserIds) == 0 {
		err := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No users selected for blacklisting.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}

	var failedBlacklistUserIds []string
	for _, userId := range selectedUserIds {
		go func(userId string) {
			err := bot.Session.ChannelPermissionSet(
				customVc.ID,
				userId,
				discordgo.PermissionOverwriteTypeMember,
				0,
				discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory | discordgo.PermissionVoiceConnect)
			if err != nil {
				failedBlacklistUserIds = append(failedBlacklistUserIds, fmt.Sprintf("<@%s>", userId))
			}
		}(userId)
	}

	if len(failedBlacklistUserIds) > 0 {
		err := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Failed to blacklist the following users: %s", strings.Join(failedBlacklistUserIds, ", ")),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	} else {
		err := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Users successfully blacklisted from the channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
	}
}