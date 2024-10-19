package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/domain/discord"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/response"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleMMR(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repo repository.Repository) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		fmt.Println(err)
	}

	role := i.ApplicationCommandData().Options[0].RoleValue(s, i.GuildID)
	if role == nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(fmt.Errorf("ロールが見つかりませんでした")))
		return
	}

	discordRepo := repository.NewDiscordRepository(s)
	members, err := discordRepo.GuildMemberByRole(i.GuildID, role.ID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	loungeRepo, err := repository.NewLoungeRepository()
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	mmrSum := 0
	mmrCount := 0
	res := &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Fields: []*discordgo.MessageEmbedField{},
			},
		},
	}
	for i, member := range members {
		// Discordのユーザー名を取得
		nameForLounge, err := loungeRepo.GetLoungeName(ctx, member.User.ID)
		if err != nil || nameForLounge == nil || nameForLounge.MMR == 0 {
			res.Embeds[(i)/25].Fields = append(res.Embeds[(i)/25].Fields, &discordgo.MessageEmbedField{
				Name:  discord.GetDisplayUsername(member),
				Value: "MMRが取得できませんでした",
			})
		} else {
			res.Embeds[(i)/25].Fields = append(res.Embeds[(i)/25].Fields, &discordgo.MessageEmbedField{
				Name:  discord.GetDisplayUsername(member),
				Value: strconv.Itoa(nameForLounge.MMR),
			})
			mmrSum += nameForLounge.MMR
			mmrCount++
		}

		if ((i + 1) % 25) == 0 {
			res.Embeds = append(res.Embeds, &discordgo.MessageEmbed{Fields: []*discordgo.MessageEmbedField{}})
		}
	}
	res.Embeds[0].Title = fmt.Sprintf("平均MMR: %d", mmrSum/mmrCount)

	_, err = s.FollowupMessageCreate(i.Interaction, true, res)
	if err != nil {
		fmt.Println(err)
	}
}
