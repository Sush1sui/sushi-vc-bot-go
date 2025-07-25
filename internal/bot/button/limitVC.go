package button

import (
	"fmt"
	"strconv"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func LimitVC(i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" { return }

	modal := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "limit_vc_modal",
			Title:    "Set Voice Channel Limit",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.TextInput{
							CustomID:    "vc_limit",
							Label:       "Enter a user limit (1-99)",
							Style:       discordgo.TextInputShort,
							Required:    true,
							Placeholder: "e.g., 10",
						},
					},
				},
			},
		},
	}

	err := bot.Session.InteractionRespond(i.Interaction, modal)
	if err != nil {
		e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to open limit VC modal.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}	
}

func HandleLimitVC(i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil { return }

	var limit int
	for _, row := range i.ModalSubmitData().Components {
		for _, comp := range row.(*discordgo.ActionsRow).Components {
			if input, ok := comp.(*discordgo.TextInput); ok && input.CustomID == "vc_limit" {
				limit, e := strconv.Atoi(input.Value)
				if e != nil || limit < 1 || limit > 99 {
					err := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Invalid limit. Please enter a number between 1 and 99.",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
					if err != nil {
						fmt.Println("Error responding to interaction:", err)
					}
					return
				}
				break
			}
		}
	}

	res, err := repository.CustomVcService.GetByOwnerOrChannelId(i.Member.User.ID, "")
	if err != nil || res == nil {
		err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not own a custom voice channel.",
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
				Content: "Failed to retrieve custom VC channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}

	_, err = bot.Session.ChannelEdit(
		res.ChannelID,
		&discordgo.ChannelEdit{
			UserLimit: limit,
		},
	)
	if err != nil {
		err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to update voice channel limit.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
	}
}