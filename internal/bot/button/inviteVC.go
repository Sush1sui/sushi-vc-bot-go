package button

import (
	"fmt"
	"strings"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func InviteUserMenu(i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil { return }

	minValue := 1
	selectMenu := discordgo.SelectMenu{
		MenuType:    discordgo.UserSelectMenu,
		CustomID:   "vc_invite_menu",
		Placeholder: "Select users to invite",
		MinValues:   &minValue,
		MaxValues:   5,
	}

	row := discordgo.ActionsRow{Components: []discordgo.MessageComponent{selectMenu}}

	err := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Please select users to invite:",
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{row},
		},
	})
	if err != nil {
		err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create invite menu.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}
}

func HandleInviteMenu(i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil { return }
	if i.MessageComponentData().CustomID != "vc_invite_menu" { return }

	selectedUserIds := i.MessageComponentData().Values
	if len(selectedUserIds) == 0 {
		err := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No users selected.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}

	messageURL := fmt.Sprintf("https://discord.com/channels/%s/%s", i.GuildID, i.ChannelID)
	guild, err := bot.Session.Guild(i.GuildID)
	if err != nil || guild == nil {
		e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve guild information.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("**<@%s> has been invited to the voice channel in %s!**\n[Join Here](%s)", i.Member.User.ID, guild.Name, messageURL),
	}
	usersInvited := []string{}
	usersFailedToInvite := []string{}
	for _, userId := range selectedUserIds {
		dmChannel, err := bot.Session.UserChannelCreate(userId)
		if err != nil {
			usersFailedToInvite = append(usersFailedToInvite, "<@"+userId+">")
			continue
		}

		msg, err := bot.Session.ChannelMessageSendEmbed(dmChannel.ID, embed)
		if err != nil || msg == nil {
			usersFailedToInvite = append(usersFailedToInvite, "<@"+userId+">")
			continue
		} else {
			usersInvited = append(usersInvited, "<@"+userId+">")
		}
	}

	err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Invited: %s\nFailed to invite: %s",
				strings.Join(usersInvited, ", "),
				strings.Join(usersFailedToInvite, ", "),
			),
		},
	})
	if err != nil {
		fmt.Println("Error responding to interaction:", err)
	}
}