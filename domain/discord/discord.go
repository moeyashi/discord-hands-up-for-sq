package discord

import "github.com/bwmarrin/discordgo"

func GetDisplayUsername(member *discordgo.Member) string {
	userName := member.Nick
	if userName == "" {
		userName = member.User.Username
	}
	return userName
}
