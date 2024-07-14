package handler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/response"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleLoungeNameSelect(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repo repository.Repository) {
	// 選択肢を使用不可にする
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "処理中です...",
		},
	})

	guild, err := repo.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	sqTitle := i.MessageComponentData().Values[0]

	// SQ Member の取得
	members, err := repo.GetSQMembers(ctx, guild, sqTitle)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	loungeRepo, err := repository.NewLoungeRepository()
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	items := []response.MakeLoungeNameResponseWebhookParamsParameterItem{}
	for _, member := range members {
		// Discordのユーザー名を取得
		nameForLounge, err := loungeRepo.GetLoungeName(ctx, member.UserID)
		if err != nil {
			s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
			return
		}
		items = append(items, response.MakeLoungeNameResponseWebhookParamsParameterItem{
			Member:     member,
			LoungeName: nameForLounge,
		})
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, response.MakeLoungeNameResponseWebhookParams(items))
	if err != nil {
		fmt.Println(err)
	}
}
