package button

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

func InviteUserMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil { return }

	minValue := 1
	selectMenu := discordgo.SelectMenu{
		MenuType:    discordgo.UserSelectMenu,
		CustomID:   "vc_invite_menu",
		Placeholder: "Select users to invite",
		MinValues:   &minValue,
		MaxValues:   5,
	}

	row := discordgo.ActionsRow{Components: []discordgo.MessageComponent{selectMenu}}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Please select users to invite:",
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{row},
		},
	})
	if err != nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create invite menu.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			fmt.Println("Error responding to interaction:", err)
		}
		return
	}
}

func HandleInviteMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil { return }
	if i.MessageComponentData().CustomID != "vc_invite_menu" { return }

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

	messageURL := fmt.Sprintf("https://discord.com/channels/%s/%s", i.GuildID, i.ChannelID)
	guild, err := s.Guild(i.GuildID)
	if err != nil || guild == nil {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to retrieve guild information.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("**<@%s> has been invited to the voice channel in %s!**\n[Join Here](%s)", i.Member.User.ID, guild.Name, messageURL),
	}
	usersInvited := []string{}
	usersFailedToInvite := []string{}
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, userId := range selectedUserIds {
		wg.Add(1)
		go func(userId string) {
			defer wg.Done()
			dmChannel, err := s.UserChannelCreate(userId)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				usersFailedToInvite = append(usersFailedToInvite, "<@"+userId+">")
				return
			}

			msg, err := s.ChannelMessageSendEmbed(dmChannel.ID, embed)
			if err != nil || msg == nil {
				usersFailedToInvite = append(usersFailedToInvite, "<@"+userId+">")
				return
			} else {
				usersInvited = append(usersInvited, "<@"+userId+">")
			}
		}(userId)
	}

	wg.Wait()

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Invited: %s\nFailed to invite: %s",
				strings.Join(usersInvited, ", "),
				strings.Join(usersFailedToInvite, ", "),
			),
		},
	})
	if err != nil {
		fmt.Println("Error responding to interaction:", err)
	}
}