package events

import (
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func OnLeaveVCEvent(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.Member == nil || vs.GuildID == "" {
		return
	}

	// Only proceed if the user left a channel
	if vs.BeforeUpdate == nil || vs.BeforeUpdate.ChannelID == "" {
		return
	}

	// Check if the user left a JTC voice channel
	customVCs, err := repository.CustomVcService.GetAllVcs()
	if err != nil || len(customVCs) == 0 {
		fmt.Println("No custom VCs found or error retrieving them:", err)
		return
	}

	for _, customVC := range customVCs {
		if customVC.ChannelID == vs.BeforeUpdate.ChannelID {
			channel, err := s.State.Channel(vs.BeforeUpdate.ChannelID)
			if err != nil || channel == nil {
				channel, err = s.Channel(vs.BeforeUpdate.ChannelID)
				fmt.Println("Error retrieving channel:", err)
				return
			}

			memberCount := 0
			guild, err := s.State.Guild(channel.GuildID)
			if err == nil && guild != nil {
				for _, voiceState := range guild.VoiceStates {
					if voiceState.ChannelID == channel.ID {
						memberCount++
					}
				}
			}

			if memberCount == 0 {
				count, err := repository.CustomVcService.DeleteByOwnerOrChannelId("", channel.ID)
				if err != nil || count == 0 {
					fmt.Println("Error deleting custom VC interface:", err)
				}
				_, err = s.ChannelDelete(channel.ID)
				if err != nil {
					fmt.Println("Error deleting voice channel:", err)
				} else {
					fmt.Printf("Deleted empty custom VC: %s\n", channel.Name)
				}
			}
		}
	}
}