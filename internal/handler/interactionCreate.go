package handler

import (
	"fmt"
	"sync"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot/button"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

var interfaceData = make(map[string]interface{})
var interfaceDataLock sync.RWMutex

func LoadInterfaceData() error {
	categories, err := repository.CategoryJTCService.GetAllJTCs()
	if err != nil {
		return err
	}

	interfaceDataLock.Lock()
	defer interfaceDataLock.Unlock()

	for _, category := range categories {
		interfaceData[category.InterfaceMessageID] = struct{}{}
	}

	return nil
}

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("Interaction received:", i.Type)

	member := i.Member
	if member == nil || i.GuildID == "" { return }

	fmt.Println("Member ID:", member.User.ID, "Guild ID:", i.GuildID)
	// Example: Check if user is in a voice channel
	voiceChannelID := ""
	guild, _ := s.State.Guild(i.GuildID)
	if guild != nil {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == member.User.ID {
				voiceChannelID = vs.ChannelID
				break
			}
		}
	}
	if voiceChannelID == "" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "**You need to be in a voice channel to use this button.**",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	fmt.Println("Voice Channel ID:", voiceChannelID)
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		customID := i.MessageComponentData().CustomID
		fmt.Println("CustomID:", customID)
		switch customID {
			case "lock_vc":
				button.LockVC(s, i)
			case "unlock_vc":
				button.UnlockVC(s, i)
			case "hide":
				button.HideUnhideVC(s, i, "hide")
			case "unhide":
				button.HideUnhideVC(s, i, "unhide")
			case "limit":
				button.LimitVC(s, i)
			case "invite":
				button.InviteUserMenu(s, i)
			case "vc_invite_menu":
				button.HandleInviteMenu(s, i)
			case "blacklist":
				button.BlacklistMenu(s, i)
			case "blacklist_menu":
				button.HandleBlacklistSelection(s, i)
			case "permit":
				button.PermitVC(s, i)
			case "permit_menu":
				button.HandleSelectedPermittedUsers(s, i)
			case "rename":
				button.RenameVC(s, i)
			case "claim_vc":
				button.ClaimVC(s, i)
			case "transfer_owner":
				button.TransferOwnership(s, i)
			default:
				e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "**Unknown button interaction.**",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if e != nil {
					fmt.Println("Failed to respond to interaction:", e)
				}
		}
	}
}