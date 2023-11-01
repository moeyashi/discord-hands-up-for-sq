package handler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func ListSQ(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	guildID := i.GuildID
	guild, err := repository.GetGuild(ctx, guildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprint(err),
			},
		})
		return
	}
	content := "SQが登録されていません。"
	if len(guild.SQList) > 0 {
		content = ""
		for _, v := range guild.SQList {
			content += v.Title + "\n"
		}
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Content: content,
		},
	})
}