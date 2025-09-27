package button

import (
	"fmt"
	"time"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func RenameVC(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" { return }

	modal := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "rename_vc_modal",
			Title:    "Rename Voice Channel",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.TextInput{
							CustomID:    "vc_new_name",
							Label:       "Enter a new name for the voice channel",
							Style:       discordgo.TextInputShort,
							Required:    true,
							Placeholder: "e.g., Gaming Lounge",
							MinLength: 1,
							MaxLength: 32,
						},
					},
				},
			},
		},
	}

	err := s.InteractionRespond(i.Interaction, modal)
	if err != nil {
		e := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to open rename VC modal.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Error responding to interaction:", e)
		}
		return
	}
}

func HandleRenameVC(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.GuildID == "" || i.Member == nil { return }

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

	var newName string
	for _, row := range i.ModalSubmitData().Components {
		for _, comp := range row.(*discordgo.ActionsRow).Components {
			if input, ok := comp.(*discordgo.TextInput); ok && input.CustomID == "vc_new_name" {
				newName = input.Value
				break
			}
		}
	}
	if newName == "" {
		msg := "Invalid channel name."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return
	}
	if newName == customVc.Name {
		msg := "The new name is the same as the current name."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
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
		msg := fmt.Sprintf("Please wait %s before renaming the voice channel again. This is due to Discord API's rate limit. You can rename the channel manually if needed.", nextAvailable.Truncate(time.Second))
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return
	}
	recent = append(recent, now)
	RenameCooldown[customVc.ID] = recent
	RenameCooldownMu.Unlock()


	_, err = s.ChannelEdit(customVc.ID, &discordgo.ChannelEdit{
		Name: newName,
	})
	if err != nil {
		msg := "Failed to rename the voice channel due to hitting Discord API's rate limit."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return
	}

	msg := fmt.Sprintf("Successfully renamed VC to: %s", newName)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}