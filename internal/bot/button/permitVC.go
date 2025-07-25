package button

import (
	"fmt"
	"strings"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func PermitVC(i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil { return }

	minValue := 1
	selectMenu := discordgo.SelectMenu{
		MenuType:    discordgo.UserSelectMenu,
		CustomID:   "permit_menu",
		Placeholder: "Select users to invite",
		MinValues:   &minValue,
		MaxValues:   5,
	}

	row := discordgo.ActionsRow{Components: []discordgo.MessageComponent{selectMenu}}

	err := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "Permit Users",
			Components: []discordgo.MessageComponent{row},
		},
	})
	if err != nil {
		e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create permit menu.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}
}

func HandleSelectedPermittedUsers(i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil { return }

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

	res, err := repository.CustomVcService.GetByOwnerOrChannelId(i.Member.User.ID, "")
	if err != nil || res == nil {
		err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve VC data.",
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
				Content: "Failed to retrieve custom VC channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}

	for _, userId := range selectedUserIds {
		err := bot.Session.ChannelPermissionSet(
			customVc.ID,
			userId,
			discordgo.PermissionOverwriteTypeMember,
			discordgo.PermissionViewChannel|discordgo.PermissionVoiceConnect|discordgo.PermissionReadMessageHistory,
			0,
		)
		if err != nil {
			err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Failed to permit user <@%s> to the voice channel.", userId),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				fmt.Println("Error responding to interaction:", err)
			}
			return
		}
	}

	err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Successfully permitted users: %s to the voice channel %s.",
				strings.Join(selectedUserIds, ", "),
				customVc.Name,
			),
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Println("Error responding to interaction:", err)
	}
}