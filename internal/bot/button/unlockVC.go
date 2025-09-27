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

    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Flags: discordgo.MessageFlagsEphemeral,
        },
    })

    res, err := repository.CustomVcService.GetByOwnerOrChannelId(i.Member.User.ID, "")
    if err != nil || res == nil {
        msg := "You don't own a custom voice channel."
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
            Content: &msg,
        })
        return
    }

    customVC, err := s.Channel(res.ChannelID)
    if err != nil || customVC == nil {
        msg := "Custom VC not found."
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
            Content: &msg,
        })
        return
    }

    // Unlock permissions
    if err := setUnlockPerms(s, customVC.ID, i.GuildID, config.GlobalConfig.FinestRoleId); err != nil {
        msg := "Failed to unlock the voice channel."
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
            Content: &msg,
        })
        return
    }

    // Cooldown check
    RenameCooldownMu.Lock()
    now := time.Now()
    timestamps := RenameCooldown[customVC.ID]
    var recent []time.Time
    for _, t := range timestamps {
        if now.Sub(t) < RenameCooldownDuration { recent = append(recent, t) }
    }
    if len(recent) >= 2 {
        nextAvailable := RenameCooldownDuration - now.Sub(recent[0])
        RenameCooldownMu.Unlock()
        msg := fmt.Sprintf("Successfully unlocked VC! Please wait %s before renaming the voice channel again. This is due to Discord API's rate limit. You can rename the channel manually if needed.", nextAvailable.Truncate(time.Second))
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
            Content: &msg,
        })
        return
    }
    recent = append(recent, now)
    RenameCooldown[customVC.ID] = recent
    RenameCooldownMu.Unlock()

    // Rename channel
    _, err = s.ChannelEdit(customVC.ID, &discordgo.ChannelEdit{
        Name: fmt.Sprintf("%s's VC", i.Member.User.Username),
    })
    if err != nil {
        msg := "Failed to rename the voice channel due to hitting Discord API's rate limit."
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
            Content: &msg,
        })
        return
    }

    msg := "Successfully unlocked VC!"
    s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
        Content: &msg,
    })
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