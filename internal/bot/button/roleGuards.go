package button

import (
	"github.com/Sush1sui/sushi-vc-bot-go/internal/common"
	"github.com/bwmarrin/discordgo"
)

func memberHasProtectedRole(member *discordgo.Member, protectedRoleIDs []string) bool {
	if member == nil {
		return false
	}

	for _, roleID := range member.Roles {
		if common.ContainsString(protectedRoleIDs, roleID) {
			return true
		}
	}

	return false
}
