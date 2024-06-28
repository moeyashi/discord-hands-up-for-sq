package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/constant"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/response"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
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

	memberType := constant.MogiListSelectCustomIDFromString(i.MessageComponentData().CustomID).ToMemberTypes()
	guild, err := repo.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	mogiTitle := i.MessageComponentData().Values[0]

	// Mogi Member の取得
	members, err := repo.GetMogiMembers(ctx, guild, mogiTitle)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	existsSameIndex := repository.IndexOfSameRegistered(members, i.Member.User.ID, memberType)
	if existsSameIndex >= 0 {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(errors.New("既に参加しています")))
		return
	}

	// Mogi Memberに追加
	userName := getDisplayUsername(i.Member)
	existsIndex := indexOfSameMember(members, i.Member.User.ID)
	if existsIndex >= 0 {
		// 既に参加している場合は一旦削除
		members = append(members[:existsIndex], members[existsIndex+1:]...)
	}
	members = append(members, repository.Member{
		UserID:     i.Member.User.ID,
		UserName:   userName,
		MemberType: memberType,
	})
	if err := repo.PutMogiMembers(ctx, guild, mogiTitle, members); err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	// メッセージの作成
	res, err := response.MakeMogiListInteractionResponse(guild.MogiList)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}
	responseMessage := fmt.Sprintf("%s を %s に追加しました。", userName, mogiTitle)
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
