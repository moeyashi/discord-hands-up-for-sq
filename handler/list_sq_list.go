package handler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func ListSQ(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	guildID := i.GuildID
	res, err := createSQListInteractionResponse(ctx, guildID, repository)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprint(err),
			},
		})
		return
	}

	if err := s.InteractionRespond(i.Interaction, res); err != nil {
		fmt.Println(err)
	}
}
