package events

import (
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func OnJoinLocked(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.Member == nil || vs.GuildID == "" {
		return
	}

	customVC, err := repository.CustomVcService.GetByOwnerOrChannelId("", vs.ChannelID)
	if err != nil || customVC == nil {
		return
	}

	// check if vc has @everyone deny connect overwrite
	channel, err := s.State.Channel(vs.ChannelID)
	if err != nil || channel == nil {
		return
	}
	overwrites := channel.PermissionOverwrites
	var everyoneOverwrite *discordgo.PermissionOverwrite
	for _, ow := range overwrites {
		if ow.ID == channel.GuildID && ow.Type == discordgo.PermissionOverwriteTypeRole {
			everyoneOverwrite = ow
			break
		}
	}
	// if @everyone is denied connect
	if everyoneOverwrite != nil && (everyoneOverwrite.Deny&discordgo.PermissionVoiceConnect) != 0 {
		s.ChannelPermissionSet(
            customVC.ChannelID,
            vs.Member.User.ID,
            discordgo.PermissionOverwriteTypeMember,
            discordgo.PermissionViewChannel | discordgo.PermissionVoiceConnect,
            0,
        )
	}
}