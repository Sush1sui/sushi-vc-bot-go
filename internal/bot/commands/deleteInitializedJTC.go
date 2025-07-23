package commands

import (
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func DeleteInitializedJTC(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}

	count, err := repository.CategoryJTCService.DeleteAll()
	if err != nil {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to delete channels.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}
	if count == 0 {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No Initialized JTC setups to delete.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}

	e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Channels deleted successfully.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if e != nil {
		fmt.Println("Error responding to interaction:", e)
	}
}