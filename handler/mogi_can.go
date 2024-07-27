package handler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/domain/discord"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/constant"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/response"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
	"github.com/moeyashi/discord-hands-up-for-sq/usecase"
)

func HandleMogiCan(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	handleMogi(ctx, s, i, repository, constant.MogiListSelectCustomIDCan)
}

func HandleMogiTemp(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	handleMogi(ctx, s, i, repository, constant.MogiListSelectCustomIDTemp)
}

func HandleMogiSub(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	handleMogi(ctx, s, i, repository, constant.MogiListSelectCustomIDSub)
}

func handleMogi(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository, handsUpType constant.MogiListSelectCustomID) {
	guild, err := repository.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, response.MakeErrorInteractionResponse(err))
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						response.MakeMogiListSelect(i.Member.User.ID, guild.MogiList, handsUpType),
					},
				},
			},
		},
	})

	if err != nil {
		s.InteractionRespond(i.Interaction, response.MakeErrorInteractionResponse(err))
	}
}

func HandleMogiSelect(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repo repository.Repository) {
	// 選択肢を使用不可にする
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "処理中です...",
		},
	})

	// 入力データの取得
	mogiTitle := i.MessageComponentData().Values[0]
	memberType := constant.MogiListSelectCustomIDFromString(i.MessageComponentData().CustomID).ToMemberTypes()

	guild, err := repo.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	// Mogi Memberに追加
	if err := usecase.AppendMogiMember(ctx, repo, repository.NewDiscordRepository(s), guild, mogiTitle, i.Member, memberType); err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	// メッセージの作成
	res, err := response.MakeMogiListInteractionResponse(guild.MogiList)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}
	responseMessage := fmt.Sprintf("%s を %s に追加しました。", discord.GetDisplayUsername(i.Member), mogiTitle)
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
