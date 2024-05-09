package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleResultsUrl(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		fmt.Println(err)
	}

	guild, err := repository.GetGuild(ctx, i.GuildID)
	if err != nil {
		fmt.Println(err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "データベースからの取得に失敗しました。",
		})
		return
	}

	if guild.Spreadsheet == "" {
		if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "スプレッドシートが設定されていません。`/results url-set`でスプレッドシートを設定し、`discordbot@hands-up-for-sq.iam.gserviceaccount.com`に編集権限を付与してください。",
		}); err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: toURL(guild.Spreadsheet),
	}); err != nil {
		fmt.Println(err)
	}
}

func HandleResultsSetURL(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		fmt.Println(err)
	}

	url := i.ApplicationCommandData().Options[0].Options[0].StringValue()

	guild, err := repository.GetGuild(ctx, i.GuildID)
	if err != nil {
		fmt.Println(err)
		if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "データベースからの取得に失敗しました。",
		}); err != nil {
			fmt.Println(err)
		}
		return
	}

	err = repository.PutResultsSpreadsheet(ctx, guild, toID(url))
	if err != nil {
		fmt.Println(err)
		if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "データベースへの書き込みに失敗しました。",
		}); err != nil {
			fmt.Println(err)
		}
		return
	}

	if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "スプレッドシートのURLを設定しました。",
	}); err != nil {
		fmt.Println(err)
	}
}

func toURL(spreadsheetID string) string {
	return fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", spreadsheetID)
}

func toID(url string) string {
	if !strings.HasPrefix(url, "https://docs.google.com/spreadsheets/d/") {
		return ""
	}

	withoutQueryParameter := strings.Split(url, "?")[0]
	splitedBySlash := strings.Split(withoutQueryParameter, "/")
	if len(splitedBySlash) < 6 {
		return ""
	}

	return splitedBySlash[5]
}
