package handler

import (
	"context"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func CreateSetCommands(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	msgContent := i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content
	commands := createSetCommandsInFuture(msgContent, time.Now())
	content := "現在以降開催のSQイベントが見つかりませんでした。"
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
