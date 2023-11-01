package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func SetSQ(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	guildID := i.GuildID
	msgContent := i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID].Content
	sqList := sqListInFuture(msgContent, time.Now())
	if err:= repository.PutSQList(ctx, guildID, sqList); err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprint(err),
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "SQリストを更新しました",
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}
