package events

import (
	"fmt"

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

            // Check if the user already owns a custom VC (search by owner)
            userVC, _ := repository.CustomVcService.GetByOwnerOrChannelId(vs.UserID, "")
            if userVC != nil {
                // move the user to their existing VC
				err = s.GuildMemberMove(vs.GuildID, vs.UserID, &userVC.ChannelID)
				if err != nil {
					fmt.Println("Error moving member to existing custom VC:", err)
				}
                continue
            }

			// create the voice channel for the user
			vcName := fmt.Sprintf("%s's VC", member.User.Username)
			newChannel, err := s.GuildChannelCreateComplex(vs.GuildID, discordgo.GuildChannelCreateData{
				Name: 	 vcName,
				Type: discordgo.ChannelTypeGuildVoice,
				ParentID: category.CategoryID,
				// PermissionOverwrites omitted -> inherits category overwrites
			})
			if err != nil {
				fmt.Println("Error creating voice channel:", err)
				return
			}

			// Optionally give the owner a member-specific overwrite (this does not affect other role overwrites)
             _ = s.ChannelPermissionSet(
				newChannel.ID, 
				vs.UserID, 
				discordgo.PermissionOverwriteTypeMember,
                discordgo.PermissionVoiceConnect|
				discordgo.PermissionManageChannels|
				discordgo.PermissionVoiceMoveMembers|
				discordgo.PermissionVoiceSpeak|
				discordgo.PermissionReadMessageHistory, 0,
			)

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