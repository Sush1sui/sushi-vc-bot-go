package button

import (
	"fmt"
	"time"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/config"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func UnlockVC(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if i.Member == nil || i.GuildID == "" {
        return
    }

    res, err := repository.CustomVcService.GetByOwnerOrChannelId(i.Member.User.ID, "")
    if err != nil || res == nil {
        respond(s, i, "You are not the owner of a custom voice channel.")
        return
    }

    customVC, err := s.Channel(res.ChannelID)
    if err != nil || customVC == nil {
        respond(s, i, "Custom VC not found.")
        return
    }

    // Unlock permissions
    if err := setUnlockPerms(s, customVC.ID, i.GuildID, config.GlobalConfig.FinestRoleId); err != nil {
        respond(s, i, "Failed to unlock the voice channel.")
        return
    }

    // Cooldown check
    RenameCooldownMu.Lock()
    lastRename, exists := RenameCooldown[customVC.ID]
    now := time.Now()
    if exists && now.Sub(lastRename) < RenameCooldownDuration {
        remaining := RenameCooldownDuration - now.Sub(lastRename)
        RenameCooldownMu.Unlock()
        respond(s, i, fmt.Sprintf("Successfully unlocked VC! Please wait %s before renaming the voice channel again. This is due to Discord API's rate limit. You can rename the channel manually if needed.", remaining.Truncate(time.Second)))
        return
    }
    RenameCooldown[customVC.ID] = now
    RenameCooldownMu.Unlock()

    // Rename channel
    _, err = s.ChannelEdit(customVC.ID, &discordgo.ChannelEdit{
        Name: fmt.Sprintf("%s's VC", i.Member.User.Username),
    })
    if err != nil {
        respond(s, i, "Failed to rename the voice channel due to hitting Discord API's rate limit.")
        return
    }

    respond(s, i, "Successfully unlocked VC!")
}

func setUnlockPerms(s *discordgo.Session, channelID, guildID, roleID string) error {
    if err := s.ChannelPermissionSet(channelID, guildID, discordgo.PermissionOverwriteTypeRole, discordgo.PermissionVoiceConnect|discordgo.PermissionReadMessageHistory|discordgo.PermissionSendMessages, 0); err != nil {
        return err
    }
    if err := s.ChannelPermissionSet(channelID, roleID, discordgo.PermissionOverwriteTypeRole, discordgo.PermissionVoiceConnect|discordgo.PermissionReadMessageHistory|discordgo.PermissionSendMessages, 0); err != nil {
        return err
    }
    return nil
}