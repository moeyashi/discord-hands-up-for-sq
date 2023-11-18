package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	_repo "github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleClick(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository _repo.Repository) {
	// CustomIDが正しいかチェック
	messageComponentData := i.MessageComponentData()
	if !strings.HasPrefix(messageComponentData.CustomID, "button_") {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "不正なボタンです。",
			},
		})
		return
	}

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

	sqTitle := strings.Split(messageComponentData.CustomID, "button_")[1]

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
	isExist := false
	for index, member := range members {
		if member.UserID == i.Member.User.ID {
			isExist = true
			members = append(members[:index], members[index+1:]...)
			responseMessage = fmt.Sprintf("%s を %s から外しました。", i.Member.User.Username, sqTitle)
			break
		}
	}
	if !isExist {
		members = append(members, _repo.Member{
			UserID:     i.Member.User.ID,
			UserName:   i.Member.User.Username,
			MemberType: _repo.MemberTypesParticipant,
		})
		responseMessage = fmt.Sprintf("%s を %s に追加しました。", i.Member.User.Username, sqTitle)
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
	res.Data.Content = responseMessage
	if err := s.InteractionRespond(i.Interaction, res); err != nil {
		fmt.Println(err)
	}

	// 古いメッセージを削除
	reference := i.Message.Reference()
	if err := s.ChannelMessageDelete(reference.ChannelID, reference.MessageID); err != nil {
		fmt.Println(err)
		return
	}
}
