package events

import (
	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot/button"
	"github.com/bwmarrin/discordgo"
)

func OnInviteVCEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}

	// Only handle user select menus
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}
	data := i.MessageComponentData()
	if data.ComponentType != discordgo.UserSelectMenuComponent {
		return
	}
	if data.CustomID != "invite_vc_menu" {
		return
	}

	button.HandleInviteMenu(s, i)
}