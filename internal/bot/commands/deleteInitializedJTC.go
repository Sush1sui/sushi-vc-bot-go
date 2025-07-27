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

	categories, err := repository.CategoryJTCService.GetAllJTCs()
	if err != nil {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve channels.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}
	if len(categories) == 0 {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No Initialized JTC setups to delete.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}

	// Delete all JTC categories and channels
	for _, category := range categories {
		_, err = s.ChannelDelete(category.JTCChannelID)
		if err != nil {
			e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to delete JTC channel.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if e != nil {
				fmt.Println("Error responding to interaction:", e)
			}
			return
		}
		_, err = s.ChannelDelete(category.InterfaceID)
		if err != nil {
			e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to delete JTC interface.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if e != nil {
				fmt.Println("Error responding to interaction:", e)
			}
			return
		}
		_, err = s.ChannelDelete(category.CategoryID)
		if err != nil {
			e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to delete JTC category.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if e != nil {
				fmt.Println("Error responding to interaction:", e)
			}
			return
		}
		fmt.Printf("Deleted category\n")
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