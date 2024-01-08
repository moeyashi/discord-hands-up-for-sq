package handler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	_repo "github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleLoungeNameSelect(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository _repo.Repository) {
	// 選択肢を使用不可にする
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "処理中です...",
		},
	})

	guild, err := repository.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprint(err),
		})
		return
	}

	sqTitle := i.MessageComponentData().Values[0]

	// SQ Member の取得
	members, err := repository.GetSQMembers(ctx, guild, sqTitle)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprint(err),
		})
		return
	}

	loungeRepo, err := _repo.NewLoungeRepository()
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprint(err),
		})
		return
	}

	embedFields := []*discordgo.MessageEmbedField{}
	for _, member := range members {
		// Discordのユーザー名を取得
		nameForLounge, err := loungeRepo.GetLoungeName(ctx, member.UserID)
		if err != nil {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprint(err),
			})
			return
		}
		nameForGuild := member.UserName
		embedFields = append(embedFields, &discordgo.MessageEmbedField{
			Name:  nameForGuild,
			Value: nameForLounge.Name,
		})
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "Loungeサーバーでの名前",
		Embeds:  []*discordgo.MessageEmbed{{Fields: embedFields}},
	})
	if err != nil {
		fmt.Println(err)
	}
}
