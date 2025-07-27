package events

import (
	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot/button"
	"github.com/bwmarrin/discordgo"
)

func OnRenameVC(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" { return }
	if i.Type != discordgo.InteractionModalSubmit { return }
	if i.ModalSubmitData().CustomID != "rename_vc_modal" { return }

	button.HandleRenameVC(s, i)
}