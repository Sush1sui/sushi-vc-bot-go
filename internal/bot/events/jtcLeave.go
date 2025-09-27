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
	if err != nil {
		fmt.Println("No custom VCs found or error retrieving them:", err)
		return
	}

	for _, customVC := range customVCs {
		if customVC.ChannelID == vs.BeforeUpdate.ChannelID {
			channel, err := s.State.Channel(vs.BeforeUpdate.ChannelID)
			if err != nil || channel == nil {
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

			if memberCount == 1 {
				// get the only user left in the vc
				var lastUserID string
				if guild != nil {
					for _, voiceState := range guild.VoiceStates {
						if voiceState.ChannelID == channel.ID {
							lastUserID = voiceState.UserID
							break
						}
					}
				}
				
				// only one user left and is not the owner, automatically transfer ownership
				// Set the new owner
				s.ChannelPermissionSet(
					customVC.ChannelID,
					lastUserID,
					discordgo.PermissionOverwriteTypeMember,
					discordgo.PermissionViewChannel | discordgo.PermissionManageChannels | discordgo.PermissionVoiceMoveMembers | discordgo.PermissionSendMessages | discordgo.PermissionAddReactions | discordgo.PermissionAttachFiles | discordgo.PermissionReadMessageHistory | discordgo.PermissionVoiceConnect,
					0,
				)

				// Delete all permissions for the old owner
				s.ChannelPermissionDelete(customVC.ChannelID, customVC.OwnerID)

				// Change the owner in the database
				repository.CustomVcService.ChangeOwnerByChannelId(customVC.ChannelID, lastUserID)

				// get the user name
				lastUser, _ := s.GuildMember(vs.GuildID, lastUserID)

				// Rename the channel
				s.ChannelEditComplex(customVC.ChannelID, &discordgo.ChannelEdit{
					Name: fmt.Sprintf("%s's VC", lastUser.User.Username),
				})
				return
			}

			// check if the user is the owner of the custom VC
			if customVC.OwnerID != vs.UserID {
				// find member overwrite (if any)
				var memberOw *discordgo.PermissionOverwrite
				for _, ow := range channel.PermissionOverwrites {
					if ow.Type == discordgo.PermissionOverwriteTypeMember && ow.ID == vs.UserID {
						memberOw = ow
						break
					}
				}

				// consider the user "permitted" if they have discordgo.PermissionSendVoiceMessages
				permitted := false
				if memberOw != nil {
					const checkPerms = discordgo.PermissionSendVoiceMessages
					if (memberOw.Allow & checkPerms) != 0 {
						permitted = true
					}
				}

				if !permitted {
					// remove any member-specific overwrite to reset permissions
					if err := s.ChannelPermissionDelete(channel.ID, vs.UserID); err != nil {
						fmt.Println("Failed to reset permissions for user on channel:", err)
					} else {
						fmt.Printf("Reset permissions for user %s on channel %s\n", vs.UserID, channel.ID)
					}
				}
			}
		}
	}
}