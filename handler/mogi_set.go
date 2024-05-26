package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/response"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleMogiSet(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repo repository.Repository) {
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

	mogi := repository.MakeMogi(
		time.Now(),
		i.ApplicationCommandData().Options[0].Options[0].IntValue(),
		i.ApplicationCommandData().Options[0].Options[1].IntValue(),
		i.ApplicationCommandData().Options[0].Options[2].IntValue(),
	)
	err = repo.AppendMogiList(ctx, guild, *mogi)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	// discordのロールを作成する
	// ロール名は「内戦:月/日」
	mentionable := true
	_, err = s.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{
		Name:        mogiRoleName(mogi),
		Mentionable: &mentionable,
	})
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	res, err := response.MakeMogiListInteractionResponse(guild.MogiList)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content:    res.Data.Content,
		Embeds:     res.Data.Embeds,
		Components: res.Data.Components,
	}); err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	if err := deleteOldMessages(s, i.ChannelID); err != nil {
		fmt.Println(err)
		return
	}
}
