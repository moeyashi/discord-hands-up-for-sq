package handler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleMogiRemoveSelect(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repo repository.Repository) {
	// 選択肢を使用不可にする
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "処理中です...",
		},
	})

	guild, err := repo.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, makeErrorFollowupResponse(err))
		return
	}

	mogiTitle := i.MessageComponentData().Values[0]
	err = repo.DeleteMogi(ctx, guild, mogiTitle)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, makeErrorFollowupResponse(err))
		return
	}

	// メッセージの作成
	res, err := createMogiListInteractionResponse(guild.MogiList)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, makeErrorFollowupResponse(err))
		return
	}
	responseMessage := fmt.Sprintf("%s を内戦リストから削除しました。\n%s", mogiTitle, res.Data.Content)
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content:    responseMessage,
		Embeds:     res.Data.Embeds,
		Components: res.Data.Components,
	})
	if err != nil {
		fmt.Println(err)
	}

	// 最後のメッセージを削除
	if err := deleteOldMessages(s, i.ChannelID); err != nil {
		fmt.Println(err)
	}
}
