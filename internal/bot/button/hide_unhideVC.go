package button

import (
	"fmt"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/config"
	"github.com/Sush1sui/sushi-vc-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func HideUnhideVC(i *discordgo.InteractionCreate, action string) {
	if i.Member == nil || i.GuildID == "" { return }
	if action != "hide" && action != "unhide" { return } 

	res, err := repository.CustomVcService.GetByOwnerOrChannelId(i.Member.User.ID, "")
	if err != nil || res == nil {
		e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not own a custom VC.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Failed to respond to interaction:", e)
		}
		return
	}

	// Check if the user is the owner of the custom VC
	if res.OwnerID != i.Member.User.ID {
		e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have permission to hide/unhide this VC.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Failed to respond to interaction:", e)
		}
		return
	}

	customVC, err := bot.Session.Channel(res.ChannelID)
	if err != nil || customVC == nil {
		e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Custom VC not found.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if e != nil {
			fmt.Println("Failed to respond to interaction:", e)
		}
		return
	}

	switch action {
		case "hide":
			e := bot.Session.ChannelPermissionSet(
				customVC.ID,
				config.GlobalConfig.FinestRoleId,
				discordgo.PermissionOverwriteTypeRole,
				0,
				discordgo.PermissionViewChannel,
			)
			if e != nil {
				e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to hide the VC.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if e != nil {
					fmt.Println("Failed to respond to interaction:", e)
				}
				return
			}

			e = bot.Session.ChannelPermissionSet(
				customVC.ID,
				i.GuildID,
				discordgo.PermissionOverwriteTypeRole,
				0,
				discordgo.PermissionViewChannel,
			)
			if e != nil {
				e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to hide the VC.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if e != nil {
					fmt.Println("Failed to respond to interaction:", e)
				}
				return
			}
		case "unhide":
			e := bot.Session.ChannelPermissionSet(
				customVC.ID,
				config.GlobalConfig.FinestRoleId,
				discordgo.PermissionOverwriteTypeRole,
				0,
				discordgo.PermissionViewChannel,
			)
			if e != nil {
				e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to hide the VC.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if e != nil {
					fmt.Println("Failed to respond to interaction:", e)
				}
				return
			}

			e = bot.Session.ChannelPermissionSet(
				customVC.ID,
				i.GuildID,
				discordgo.PermissionOverwriteTypeRole,
				0,
				discordgo.PermissionViewChannel,
			)
			if e != nil {
				e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to hide the VC.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if e != nil {
					fmt.Println("Failed to respond to interaction:", e)
				}
				return
			}
		default:
			e := bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Invalid action specified.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if e != nil {
				fmt.Println("Failed to respond to interaction:", e)
			}
			return
	}
	
	err = bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Custom VC has been successfully %s.", action),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Println("Failed to respond to interaction:", err)
		return
	}
}