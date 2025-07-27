package common

import "github.com/bwmarrin/discordgo"

func InterfaceEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: 0xFFFFFF,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Finesse VC Interface",
			IconURL: "https://images-ext-1.discordapp.net/external/3QmLnkyUjiyS6EAm51WT-Yyqe7bcDoF9QRTpsfECbII/https/media.tenor.com/ZjZcvkBzoNMAAAAi/pepe-scucha.gif",
		},
		Title:       "ENJOY UNLIMITED VC INTERFACE ACCESS!",
		Description: "Hey, Thank you for supporting our server! Here at Finesse we make sure that our members make the most of it when they're on a VC with friends! With having the freedom to do what you please with your own voice channel!",
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Use the buttons below to manage your voice channel",
			IconURL: "https://images-ext-1.discordapp.net/external/w1oTKGUTTcVtkkPbAEF-0CkhMwuugjhfnzKoX5UCVBE/%3Fsize%3D96%26quality%3Dlossless/https/cdn.discordapp.com/emojis/1293411594621157458.gif",
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://media.tenor.com/iJklJd0dfrcAAAAi/cat-cats.gif",
		},
	}
}

func InterfaceButtonsRow1() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					CustomID: "lock_vc",
					Label:    "Lock VC",
					Style:    discordgo.DangerButton,
					Emoji:    &discordgo.ComponentEmoji{ID: "1293802497735135243", Name: "lock_vc"},
				},
				&discordgo.Button{
					CustomID: "unlock_vc",
					Label:    "Unlock VC",
					Style:    discordgo.SuccessButton,
					Emoji:    &discordgo.ComponentEmoji{ID: "1293802495407030272", Name: "unlock_vc"},
				},
				&discordgo.Button{
					CustomID: "hide",
					Label:    "Hide",
					Style:    discordgo.DangerButton,
					Emoji:    &discordgo.ComponentEmoji{ID: "1293803740113010738", Name: "hide"},
				},
				&discordgo.Button{
					CustomID: "unhide",
					Label:    "Unhide",
					Style:    discordgo.SuccessButton,
					Emoji:    &discordgo.ComponentEmoji{ID: "1293802561887010827", Name: "unhide"},
				},
				&discordgo.Button{
					CustomID: "limit",
					Label:    "Limit",
					Style:    discordgo.PrimaryButton,
					Emoji:    &discordgo.ComponentEmoji{ID: "1293802599614648462", Name: "limit"},
				},
			},
		},
	}
}

func InterfaceButtonsRow2() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					CustomID: "invite",
					Label:    "Invite",
					Style:    discordgo.PrimaryButton,
					Emoji:    &discordgo.ComponentEmoji{ID: "1293802491967836246", Name: "invite"},
				},
				&discordgo.Button{
					CustomID: "blacklist",
					Label:    "Blacklist",
					Style:    discordgo.DangerButton,
					Emoji:    &discordgo.ComponentEmoji{ID: "1293802490956873738", Name: "blacklist"},
				},
				&discordgo.Button{
					CustomID: "permit",
					Label:    "Permit",
					Style:    discordgo.SuccessButton,
					Emoji:    &discordgo.ComponentEmoji{ID: "1293802489711431740", Name: "permit"},
				},
				&discordgo.Button{
					CustomID: "rename",
					Label:    "Rename",
					Style:    discordgo.PrimaryButton,
					Emoji:    &discordgo.ComponentEmoji{ID: "1293802483046678529", Name: "rename"},
				},
				&discordgo.Button{
					CustomID: "claim_vc",
					Label:    "Claim VC",
					Style:    discordgo.PrimaryButton,
					Emoji:    &discordgo.ComponentEmoji{ID: "1293802473789718531", Name: "claim_vc"},
				},
			},
		},
	}
}

func InterfaceButtonsRow3() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					CustomID: "transfer_owner",
					Label:    "Transfer Owner",
					Style:    discordgo.PrimaryButton,
					Emoji:    &discordgo.ComponentEmoji{ID: "1293802472560660512", Name: "transfer_owner"},
				},
			},
		},
	}
}