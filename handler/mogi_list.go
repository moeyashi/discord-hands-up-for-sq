package handler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleMogiList(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	guild, err := repository.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, makeErrorResponse(err))
		return
	}

	res, err := createMogiListInteractionResponse(guild.MogiList)
	if err != nil {
		s.InteractionRespond(i.Interaction, makeErrorResponse(err))
		return
	}

	if err := s.InteractionRespond(i.Interaction, res); err != nil {
		fmt.Println(err)
		return
	}

	if err := deleteOldMessages(s, i.ChannelID); err != nil {
		fmt.Println(err)
		return
	}
}
