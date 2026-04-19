package button

import (
	"fmt"
	"strings"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/common"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func PermitVC(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil {
		return
	}

	minValue := 1
	selectMenu := discordgo.SelectMenu{
		MenuType:    discordgo.UserSelectMenu,
		CustomID:    "permit_menu",
		Placeholder: "Select users to invite",
		MinValues:   &minValue,
		MaxValues:   5,
	}

	row := discordgo.ActionsRow{Components: []discordgo.MessageComponent{selectMenu}}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title:      "**Permit Users**",
			Components: []discordgo.MessageComponent{row},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create permit menu.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}
}

func HandleSelectedPermittedUsers(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil {
		return
	}
	if i.MessageComponentData().CustomID != "permit_menu" {
		return
	}

	selectedUserIds := i.MessageComponentData().Values
	if len(selectedUserIds) == 0 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No users selected.",
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
				Content: "Failed to retrieve VC data.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}

	customVc, err := s.Channel(res.ChannelID)
	if err != nil || customVc == nil {
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

	settings := common.GetEffectiveGuildSettings(i.GuildID)

	if common.ContainsString(settings.IgnoredChannelIDs, customVc.ID) {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This channel is configured as ignored and cannot be modified.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}

	var skippedProtectedUsers []string
	var permittedUsers []string

	for _, userId := range selectedUserIds {
		member, memberErr := s.GuildMember(i.GuildID, userId)
		if memberErr == nil && memberHasProtectedRole(member, settings.ProtectedRoleIDs) {
			skippedProtectedUsers = append(skippedProtectedUsers, fmt.Sprintf("<@%s>", userId))
			continue
		}

		err := s.ChannelPermissionSet(
			customVc.ID,
			userId,
			discordgo.PermissionOverwriteTypeMember,
			discordgo.PermissionViewChannel|
				discordgo.PermissionVoiceConnect|
				discordgo.PermissionReadMessageHistory|
				discordgo.PermissionSendMessages|
				discordgo.PermissionAddReactions|
				discordgo.PermissionUseApplicationCommands|
				discordgo.PermissionVoiceSpeak|
				discordgo.PermissionVoiceStreamVideo|
				discordgo.PermissionSendVoiceMessages,
			0,
		)
		if err != nil {
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Failed to permit user <@%s> to the voice channel.", userId),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				fmt.Println("Error responding to interaction:", err)
			}
			return
		}

		permittedUsers = append(permittedUsers, userId)
	}

	if len(permittedUsers) == 0 && len(skippedProtectedUsers) > 0 {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("No users were permitted. Skipped protected users: %s", strings.Join(skippedProtectedUsers, ", ")),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: buildPermitResponse(permittedUsers, skippedProtectedUsers, customVc.Name),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Println("Error responding to interaction:", err)
	}
}

func buildPermitResponse(permittedUsers, skippedProtectedUsers []string, channelName string) string {
	response := fmt.Sprintf("Successfully permitted users: %s to the voice channel %s.", strings.Join(permittedUsers, ", "), channelName)
	if len(skippedProtectedUsers) == 0 {
		return response
	}

	return response + " Skipped protected users: " + strings.Join(skippedProtectedUsers, ", ")
}
