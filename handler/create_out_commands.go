package handler

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func CreateOutCommands(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	msgComponents := i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Components
	commands := createOutCommandsForAll(msgComponents)
	content := "過去開催の募集が見つかりませんでした。"
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
