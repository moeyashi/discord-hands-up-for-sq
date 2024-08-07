package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/domain/discord"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/constant"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/response"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
	"github.com/moeyashi/discord-hands-up-for-sq/util"
)

func HandleSelect(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repo repository.Repository) {
	// 選択肢を使用不可にする
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "処理中です...",
		},
	})

	memberType := constant.SQListSelectCustomIDFromString(i.MessageComponentData().CustomID).ToMemberTypes()
	guild, err := repo.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprint(err),
		})
		return
	}

	sqTitle := i.MessageComponentData().Values[0]

	// SQ Member の取得
	members, err := repo.GetSQMembers(ctx, guild, sqTitle)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprint(err),
		})
		return
	}

	existsSameIndex := repository.IndexOfSameRegistered(members, i.Member.User.ID, memberType)
	if existsSameIndex >= 0 {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "既に参加しています。",
		})
		return
	}

	// SQ Memberに追加
	userName := discord.GetDisplayUsername(i.Member)
	existsIndex := util.IndexOfSameMember(members, i.Member.User.ID)
	if existsIndex >= 0 {
		// 既に参加している場合は一旦削除
		members = append(members[:existsIndex], members[existsIndex+1:]...)
	}
	members = append(members, repository.Member{
		UserID:     i.Member.User.ID,
		UserName:   userName,
		MemberType: memberType,
	})
	if err := repo.PutSQMembers(ctx, guild, sqTitle, members); err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprint(err),
		})
		return
	}

	// メッセージの作成
	res, err := response.MakeSQListInteractionResponse(guild.SQList, time.Now())
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprint(err),
		})
		return
	}
	responseMessage := fmt.Sprintf("%s を %s に追加しました。", userName, sqTitle)
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
