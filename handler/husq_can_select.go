package handler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	_repo "github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleSelectCan(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository _repo.Repository) {
	const memberType = _repo.MemberTypesParticipant
	guild, err := repository.GetGuild(ctx, i.GuildID)
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

	sqTitle := i.MessageComponentData().Values[0]

	// SQ Member の取得
	members, err := repository.GetSQMembers(ctx, guild, sqTitle)
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

	// SQ Memberに追加 or 削除
	responseMessage := ""
	existsIndex := indexOfSameRegistered(members, i.Member.User.ID, memberType)
	if existsIndex >= 0 {
		members = append(members[:existsIndex], members[existsIndex+1:]...)
		responseMessage = fmt.Sprintf("%s を %s から外しました。", i.Member.Nick, sqTitle)
	} else {
		members = append(members, _repo.Member{
			UserID:     i.Member.User.ID,
			UserName:   i.Member.Nick,
			MemberType: memberType,
		})
		responseMessage = fmt.Sprintf("%s を %s に追加しました。", i.Member.Nick, sqTitle)
	}
	if err := repository.PutSQMembers(ctx, guild, sqTitle, members); err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprint(err),
			},
		})
		return
	}

	// 選択肢を使用不可にする
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: responseMessage,
		},
	})

	// メッセージの作成
	res, err := createSQListInteractionResponse(ctx, guild.SQList, repository)
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
