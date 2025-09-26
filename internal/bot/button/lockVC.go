package button

import (
	"fmt"
	"sync"
	"time"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/config"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

var (
    RenameCooldown   = make(map[string][]time.Time)
    RenameCooldownMu sync.Mutex
)
const RenameCooldownDuration = 10 * time.Minute

func LockVC(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if i.Member == nil || i.GuildID == "" {
        return
    }

    res, err := repository.CustomVcService.GetByOwnerOrChannelId(i.Member.User.ID, "")
    if err != nil || res == nil {
        respond(s, i, "You are not an owner of a custom vc channel.")
        return
    }

    customVc, err := s.Channel(res.ChannelID)
    if err != nil || customVc == nil {
        respond(s, i, "Custom VC not found.")
        return
    }

    // Lock permissions
    if err := setLockPerms(s, customVc.ID, i.GuildID, config.GlobalConfig.FinestRoleId); err != nil {
        respond(s, i, "Failed to lock the voice channel.")
        return
    }

    // Cooldown check
    RenameCooldownMu.Lock()
    now := time.Now()
    timestamps := RenameCooldown[customVc.ID]
    var recent []time.Time
    for _, t := range timestamps {
        if now.Sub(t) < RenameCooldownDuration { recent = append(recent, t) }
    }
    if len(recent) >= 2 {
        nextAvailable := RenameCooldownDuration - now.Sub(recent[0])
        RenameCooldownMu.Unlock()
        respond(s, i, fmt.Sprintf("Successfully locked VC! Please wait %s before renaming the voice channel again. This is due to Discord API's rate limit. You can rename the channel manually if needed.", nextAvailable.Truncate(time.Second)))
        return
    }
    recent = append(recent, now)
    RenameCooldown[customVc.ID] = recent
    RenameCooldownMu.Unlock()

    // Rename channel
    _, err = s.ChannelEdit(customVc.ID, &discordgo.ChannelEdit{
        Name: fmt.Sprintf("🔒 | %s's VC", i.Member.User.Username),
    })
    if err != nil {
        respond(s, i, "Failed to rename the voice channel due to hitting Discord API's rate limit.")
        return
    }

    respond(s, i, "Voice channel locked successfully.")
}

func setLockPerms(s *discordgo.Session, channelID, guildID, roleID string) error {
    if err := s.ChannelPermissionSet(channelID, guildID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionVoiceConnect|discordgo.PermissionReadMessageHistory|discordgo.PermissionSendMessages); err != nil {
        return err
    }
    if err := s.ChannelPermissionSet(channelID, roleID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionVoiceConnect); err != nil {
        return err
    }
    return nil
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
    _ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: msg,
            Flags:   discordgo.MessageFlagsEphemeral,
        },
    })
}