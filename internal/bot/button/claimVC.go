package button

import (
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func ClaimVC(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}

	member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil || member == nil {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You must be a member of the server to claim a VC.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Failed to respond to interaction:", e)
		}
		return
	}

	// Check if the user is in a voice channel
	var voiceChannelID string
	var voiceStates []*discordgo.VoiceState
	for _, vs := range s.State.Guilds {
		if vs.ID == i.GuildID {
			voiceStates = vs.VoiceStates
			break
		}
	}
	for _, vs := range voiceStates {
		if vs.UserID == i.Member.User.ID {
			voiceChannelID = vs.ChannelID
			break
		}
	}

	if voiceChannelID == "" {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You must be in a voice channel to claim it.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Failed to respond to interaction:", err)
		}
		return
	}

	res, err := repository.CustomVcService.GetByOwnerOrChannelId("", voiceChannelID)
	if err != nil || res == nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve custom VC data.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Failed to respond to interaction:", err)
		}
		return
	}

	if res.OwnerID == i.Member.User.ID {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You already own this VC.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Failed to respond to interaction:", err)
		}
		return
	}

	customVc, err := s.Channel(res.ChannelID)
	if err != nil || customVc == nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Custom VC not found.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Failed to respond to interaction:", err)
		}
		return
	}

	// check if owner is in the custom vc
	for _, vs := range voiceStates {
		if vs.UserID == customVc.OwnerID && vs.ChannelID == customVc.ID {
			e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "**Owner is still in the VC, you can't claim it.**",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if e != nil {
				fmt.Println("Failed to respond to interaction:", e)
			}
			return
		}
	}

	// Set permissions for the new owner in the custom VC
	err = s.ChannelPermissionSet(
		customVc.ID,
		i.Member.User.ID,
		discordgo.PermissionOverwriteTypeMember,
		discordgo.PermissionViewChannel | discordgo.PermissionManageChannels | discordgo.PermissionVoiceMoveMembers | discordgo.PermissionSendMessages | discordgo.PermissionAddReactions | discordgo.PermissionAttachFiles | discordgo.PermissionReadMessageHistory | discordgo.PermissionVoiceConnect,
		0,
	)
	if err != nil {
		fmt.Println("Failed to set permissions for the custom VC:", err)
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to set permissions for the custom VC.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Failed to respond to interaction:", e)
		}
		return
	}

	// Set permissions for the old owner of the custom VC
	err = s.ChannelPermissionSet(
		customVc.ID,
		res.OwnerID,
		discordgo.PermissionOverwriteTypeMember,
		discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory | discordgo.PermissionVoiceConnect,
		discordgo.PermissionManageChannels | discordgo.PermissionVoiceMoveMembers | discordgo.PermissionAddReactions | discordgo.PermissionAttachFiles,
	)
	if err != nil {
		fmt.Println("Failed to set permissions for the old owner of the custom VC:", err)
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to set permissions for the custom VC owner.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Failed to respond to interaction:", e)
		}
		return
	}

	count, err := repository.CustomVcService.ChangeOwnerByChannelId(customVc.ID, i.Member.User.ID)
	if err != nil || count == 0 {
		fmt.Println("Failed to change owner of the custom VC:", err)
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to claim the custom VC.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Failed to respond to interaction:", e)
		}
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You have successfully claimed the VC **%s**.", customVc.Name),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Println("Failed to respond to interaction:", err)
		return
	}
}