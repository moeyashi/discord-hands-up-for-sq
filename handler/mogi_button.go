package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/response"
	_repo "github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleMogiButtonClick(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository _repo.Repository) {

	guild, err := repository.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, response.MakeErrorInteractionResponse(err))
		return
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		fmt.Println(err)
	}

	mogiTitle := strings.Split(i.MessageComponentData().CustomID, "button_mogi_")[1]
	mogi, err := repository.GetMogi(ctx, guild, mogiTitle)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	// Mogi Member の取得
	members, err := repository.GetMogiMembers(ctx, guild, mogiTitle)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	// Mogi Memberに追加 or 削除
	roleName := mogiRoleName(mogi)
	role, err := findMogiRole(s, i.GuildID, roleName)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	userName := getDisplayUsername(i.Member)
	responseMessage := ""
	isExist := false
	for index, member := range members {
		if member.UserID == i.Member.User.ID {
			isExist = true
			members = append(members[:index], members[index+1:]...)
			responseMessage = fmt.Sprintf("%s を %s から外しました。", userName, mogiTitle)
			if role != nil {
				if err := s.GuildMemberRoleRemove(i.GuildID, member.UserID, role.ID); err != nil {
					s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
					return
				}
			}
			break
		}
	}
	if !isExist {
		members = append(members, _repo.Member{
			UserID:     i.Member.User.ID,
			UserName:   userName,
			MemberType: _repo.MemberTypesParticipant,
		})
		responseMessage = fmt.Sprintf("%s を %s に追加しました。", userName, mogiTitle)
		if role != nil {
			if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, role.ID); err != nil {
				s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
				return
			}
		}
	}
	if err := repository.PutMogiMembers(ctx, guild, mogiTitle, members); err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	// メッセージの作成
	res, err := createMogiListInteractionResponse(guild.MogiList)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}
	res.Data.Content = responseMessage
	if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content:    res.Data.Content,
		Embeds:     res.Data.Embeds,
		Components: res.Data.Components,
	}); err != nil {
		fmt.Println(err)
	}

	if err := deleteOldMessages(s, i.ChannelID); err != nil {
		fmt.Println(err)
	}
}
