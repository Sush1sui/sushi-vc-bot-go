package deploy

import (
	"log"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot/events"
	"github.com/bwmarrin/discordgo"
)

var EventHandlers = []any{
	events.OnJoinVCEvent,
	events.OnLeaveVCEvent,
	events.OnLimitVC,
	events.OnRenameVC,
	events.OnJoinLocked,
}

func DeployEvents(s *discordgo.Session) {
	for _, handler := range EventHandlers {
		s.AddHandler(handler)
	}

	log.Println("Event handlers deployed successfully.")
}