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

    customVc, err := s.Channel(res.ChannelID)
    if err != nil || customVc == nil {
        msg := "VC not found."
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
            Content: &msg,
        })
        return
    }

    // Lock permissions
    if err := setLockPerms(s, customVc.ID, i.GuildID, config.GlobalConfig.FinestRoleId); err != nil {
        msg := "Failed to set lock permissions."
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
            Content: &msg,
        })
        fmt.Println("Error setting lock permissions:", err)
        return
    }

    // add permission to the users in the vc after locking
    usersInVC := make(map[string]struct{})
    guild, err := s.State.Guild(i.GuildID)
    if err == nil && guild != nil {
        for _, vs := range guild.VoiceStates {
            if vs.ChannelID == customVc.ID {
                usersInVC[vs.UserID] = struct{}{}
            }
        }
    }
    for userID := range usersInVC {
        s.ChannelPermissionSet(
            customVc.ID,
            userID,
            discordgo.PermissionOverwriteTypeMember,
            discordgo.PermissionViewChannel | discordgo.PermissionVoiceConnect,
            0,
        )
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
        msg := fmt.Sprintf("Successfully locked VC! Please wait %s before renaming the voice channel again. This is due to Discord API's rate limit. You can rename the channel manually if needed.", nextAvailable.Truncate(time.Second))
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
            Content: &msg,
        })
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
        msg := "Failed to rename the voice channel due to hitting Discord API's rate limit."
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
            Content: &msg,
        })
        return
    }

    msg := "Successfully locked and renamed VC!"
    s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
        Content: &msg,
    })
}

func setLockPerms(s *discordgo.Session, channelID, guildID, roleID string) error {
    if err := s.ChannelPermissionSet(channelID, guildID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionVoiceConnect); err != nil {
        return err
    }
    if err := s.ChannelPermissionSet(
        channelID, 
        roleID, 
        discordgo.PermissionOverwriteTypeRole, 
        discordgo.PermissionViewChannel | 
        discordgo.PermissionCreateInstantInvite | 
        discordgo.PermissionVoiceSpeak | 
        discordgo.PermissionVoiceStreamVideo |
        discordgo.PermissionSendMessages | 
        discordgo.PermissionAddReactions | 
        discordgo.PermissionReadMessageHistory | 
        discordgo.PermissionUseApplicationCommands, 
        discordgo.PermissionManageEvents | 
        discordgo.PermissionCreateEvents | 
        discordgo.PermissionSendVoiceMessages,
    ); err != nil {
        return err
    }
    return nil
}