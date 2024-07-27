package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/domain/discord"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/response"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
	"github.com/moeyashi/discord-hands-up-for-sq/usecase"
)

func HandleMogiButtonClick(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repo repository.Repository) {
	guild, err := repo.GetGuild(ctx, i.GuildID)
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
	mogi, err := repo.GetMogi(ctx, guild, mogiTitle)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}
	members := mogi.Members

	// Mogi Memberに追加 or 削除
	userName := discord.GetDisplayUsername(i.Member)
	responseMessage := ""
	isExist := false
	for index, member := range members {
		if member.UserID == i.Member.User.ID {
			// 削除
			isExist = true
			responseMessage = fmt.Sprintf("%s を %s から外しました。", userName, mogiTitle)
			roleName := mogi.RoleName()
			role, err := repository.NewDiscordRepository(s).FindRoleByName(i.GuildID, roleName)
			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
				return
			}
			if role != nil {
				if err := s.GuildMemberRoleRemove(i.GuildID, member.UserID, role.ID); err != nil {
					s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
					return
				}
			}
			members = append(members[:index], members[index+1:]...)
			if err := repo.PutMogiMembers(ctx, guild, mogiTitle, members); err != nil {
				s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
				return
			}
			break
		}
	}
	if !isExist {
		// 追加
		responseMessage = fmt.Sprintf("%s を %s に追加しました。", userName, mogiTitle)
		if err := usecase.AppendMogiMember(ctx, repo, repository.NewDiscordRepository(s), guild, mogiTitle, i.Member, repository.MemberTypesParticipant); err != nil {
			s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
			return
		}
	}

	// メッセージの作成
	res, err := response.MakeMogiListInteractionResponse(guild.MogiList)
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
