package events

import (
	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot/button"
	"github.com/bwmarrin/discordgo"
)

func OnTransferOwner(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}
	if i.Type != discordgo.InteractionType(discordgo.UserSelectMenu) {
		return
	}
	data := i.MessageComponentData()
	if data.ComponentType != discordgo.UserSelectMenuComponent {
		return
	}
	if data.CustomID != "transfer_owner_menu" {
		return
	}

	button.HandleTransferOwnership(s, i)
}