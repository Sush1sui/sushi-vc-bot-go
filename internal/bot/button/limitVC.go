package button

import (
	"fmt"
	"strconv"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func LimitVC(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	err := s.InteractionRespond(i.Interaction, modal)
	if err != nil {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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

func HandleLimitVC(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil { return }

	var limit int
	found := false
	for _, row := range i.ModalSubmitData().Components {
		for _, comp := range row.(*discordgo.ActionsRow).Components {
			if input, ok := comp.(*discordgo.TextInput); ok && input.CustomID == "vc_limit" {
				val, e := strconv.Atoi(input.Value)
				if e != nil || val < 1 || val > 99 {
					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
				limit = val
				found = true
				break
			}
		}
	}
	if !found {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No limit input found.",
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
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
	
	customVC, err := s.Channel(res.ChannelID)
	if err != nil || customVC == nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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

	fmt.Println("Setting voice channel limit to:", limit)
	_, err = s.ChannelEdit(
		res.ChannelID,
		&discordgo.ChannelEdit{
			UserLimit: limit,
		},
	)
	if err != nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Voice channel user limit set to %d.", limit),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Println("Error responding to interaction:", err)
	}
}