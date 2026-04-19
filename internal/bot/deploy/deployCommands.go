package deploy

import (
	"fmt"
	"log"

	"github.com/Sush1sui/sushi-vc-bot-go/internal/bot/commands"
	"github.com/bwmarrin/discordgo"
)

var SlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:                     "initialize-jtc",
		Description:              "Initialize Join to Create setup",
		Type:                     discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(),
	},
	{
		Name:                     "delete-jtc-setup",
		Description:              "Delete Join to Create setup",
		Type:                     discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(),
	},
	{
		Name:                     "guild-settings",
		Description:              "Manage guild settings for JTC role and channel exceptions",
		Type:                     discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "create",
				Description: "Create guild settings document",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "target",
						Description: "Which list to set",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "voice_exception_roles", Value: "voice_exception_roles"},
							{Name: "protected_roles", Value: "protected_roles"},
							{Name: "ignored_channels", Value: "ignored_channels"},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "ids",
						Description: "Comma-separated IDs",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "update",
				Description: "Replace one guild settings list",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "target",
						Description: "Which list to set",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "voice_exception_roles", Value: "voice_exception_roles"},
							{Name: "protected_roles", Value: "protected_roles"},
							{Name: "ignored_channels", Value: "ignored_channels"},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "ids",
						Description: "Comma-separated IDs",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "delete",
				Description: "Delete IDs from one list or clear list when ids omitted",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "target",
						Description: "Which list to edit",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "voice_exception_roles", Value: "voice_exception_roles"},
							{Name: "protected_roles", Value: "protected_roles"},
							{Name: "ignored_channels", Value: "ignored_channels"},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "ids",
						Description: "Comma-separated IDs to remove",
						Required:    false,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "view",
				Description: "View effective guild settings",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "reset",
				Description: "Delete guild settings document and fallback to defaults",
			},
		},
	},
}

var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"initialize-jtc":   commands.InitializeJTC,
	"delete-jtc-setup": commands.DeleteInitializedJTC,
	"guild-settings":   commands.ManageGuildSettings,
}

func DeployCommands(s *discordgo.Session) {
	globalCmds, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		for _, cmd := range globalCmds {
			err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
			if err != nil {
				log.Printf("Failed to delete command %s: %v", cmd.Name, err)
			} else {
				log.Printf("Deleted command %s", cmd.Name)
			}
		}
	}

	guilds := s.State.Guilds
	for _, guild := range guilds {
		_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guild.ID, SlashCommands)
		if err != nil {
			log.Fatalf("Failed to deploy commands to guild %s: %v", guild.ID, err)
		}
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return // Only handle slash commands here!
		}
		if handler, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		} else {
			fmt.Printf("Unknown command: %s\n", i.ApplicationCommandData().Name)
			fmt.Printf("Available commands: %v\n", CommandHandlers)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Unknown command: %s\n", i.ApplicationCommandData().Name),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	})

	log.Println("Commands deployed successfully.")
}
