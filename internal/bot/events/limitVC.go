package events

import (
	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot/button"
	"github.com/bwmarrin/discordgo"
)

func OnLimitVC(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" { return }
	if i.Type != discordgo.InteractionModalSubmit { return }
	if i.ModalSubmitData().CustomID != "limit_vc_modal" { return }

	button.HandleLimitVC(s, i)
}