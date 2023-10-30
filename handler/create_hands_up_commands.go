package handler

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func CreateHandsUpCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msgContent := i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content
	commands := createHandsUpCommandsForBetweenNowAndTomorrow(msgContent, time.Now())
	content := "今日明日開催のSQイベントが見つかりませんでした。"
	if len(commands) > 0 {
		content = strings.Join(commands, "\n")
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}
