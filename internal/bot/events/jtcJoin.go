package events

import (
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/config"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func OnJoinVCEvent(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.Member == nil || vs.GuildID == "" {
		return
	}

	categories, err := repository.CategoryJTCService.GetAllJTCs()
	if err != nil || len(categories) == 0 {
		fmt.Println("No JTC categories found or error retrieving them:", err)
		return
	}

	for _, category := range categories {
		// user joined or was moved to a JTC voice channel
		if(vs.BeforeUpdate == nil || vs.BeforeUpdate.ChannelID != vs.ChannelID) && vs.ChannelID == category.JTCChannelID {
			// fetch the member
			member, err := s.GuildMember(vs.GuildID, vs.UserID)
			if err != nil {
				fmt.Println("Error fetching member:", err)
				return
			}

			// create the voice channel for the user
			vcName := fmt.Sprintf("%s's VC", member.User.Username)
			newChannel, err := s.GuildChannelCreateComplex(vs.GuildID, discordgo.GuildChannelCreateData{
				Name: 	 vcName,
				Type: discordgo.ChannelTypeGuildVoice,
				ParentID: category.CategoryID,
				PermissionOverwrites: []*discordgo.PermissionOverwrite{
					{
						ID:    vs.GuildID, // @everyone
						Type:  discordgo.PermissionOverwriteTypeRole,
						Deny:  discordgo.PermissionVoiceConnect | discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
						Allow: 0,
					},
					{
						ID:    vs.UserID, // The user
						Type:  discordgo.PermissionOverwriteTypeMember,
						Allow: discordgo.PermissionVoiceConnect | discordgo.PermissionVoiceSpeak | discordgo.PermissionManageChannels | discordgo.PermissionReadMessageHistory | discordgo.PermissionVoiceMoveMembers,
						Deny:  0,
					},
					{
						ID:    config.GlobalConfig.FinestRoleId, // Finest Role
						Type:  discordgo.PermissionOverwriteTypeRole,
						Allow: discordgo.PermissionVoiceConnect | discordgo.PermissionVoiceSpeak | discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory,
						Deny:  0,
					},
				},
			})
			if err != nil {
				fmt.Println("Error creating voice channel:", err)
				return
			}

			// move the user to the new channel
			err = s.GuildMemberMove(vs.GuildID, vs.UserID, &newChannel.ID)
			if err != nil {
				fmt.Println("Error moving member to new channel:", err)
				return
			}

			_, err = repository.CustomVcService.CreateVc(vs.UserID, newChannel.ID)
			if err != nil {
				fmt.Println("Error saving custom VC:", err)
				return
			}
		}
	}
}