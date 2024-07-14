package handler

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/response"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
	"github.com/moeyashi/discord-hands-up-for-sq/usecase"
)

func SetSQ(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	guild, err := repository.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, response.MakeErrorInteractionResponse(err))
		return
	}

	now := time.Now()
	msgContent := i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content
	sqList := usecase.NewSQList(
		guild.SQList,
		sqListInFuture(msgContent, now),
		now,
	)
	if err := repository.PutSQList(ctx, guild, sqList); err != nil {
		s.InteractionRespond(i.Interaction, response.MakeErrorInteractionResponse(err))
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "SQリストを更新しました",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
