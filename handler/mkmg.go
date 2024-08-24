package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/response"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleMKMGTeam(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repo repository.Repository) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		fmt.Println(err)
		return
	}

	options := i.ApplicationCommandData().Options[0].Options

	if len(options) == 0 {
		// チーム名を確認
		guild, err := repo.GetGuild(ctx, i.GuildID)
		if err != nil {
			s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
			return
		}

		if guild.Name == "" {
			s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(errors.New("チーム名が設定されていません。`/mkmg team`でチーム名を設定してください。")))
			return
		}

		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("現在のチーム名は `%s` です", guild.Name),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		// チーム名を設定
		name := options[0].StringValue()
		_, err := repo.PutGuildName(ctx, i.GuildID, name)
		if err != nil {
			s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
			return
		}

		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("`/mkmg`コマンドでのチーム名を `%s` に設定しました", name),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func HandleMKMGPost(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repo repository.Repository) {
	time := i.ApplicationCommandData().Options[0].Options[0].StringValue()
	avg := i.ApplicationCommandData().Options[0].Options[1].IntValue()

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		fmt.Println(err)
		return
	}

	guild, err := repo.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(err))
		return
	}

	if guild.Name == "" {
		s.FollowupMessageCreate(i.Interaction, true, response.MakeErrorWebhookParams(errors.New("チーム名が設定されていません。`/mkmg team`でチーム名を設定してください。")))
		return
	}

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("%s時交流戦お相手募集\nこちら%s、平均%d\n主催可能、ID開設\n#mkmg", time, guild.Name, avg),
	})
}
