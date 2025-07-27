package events

import (
	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot/button"
	"github.com/bwmarrin/discordgo"
)

func OnBlacklistEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
  if i.Member == nil || i.GuildID == "" {
    return
  }
  if i.Type != discordgo.InteractionMessageComponent {
    return
  }
  data := i.MessageComponentData()
  if data.ComponentType != discordgo.UserSelectMenuComponent {
    return
  }
  if data.CustomID != "blacklist_menu" {
    return
  }

  button.HandleBlacklistSelection(s, i)
}