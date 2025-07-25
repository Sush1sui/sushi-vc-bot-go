package button

import (
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot"
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