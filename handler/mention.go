package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	_repo "github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleMention(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository _repo.Repository) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		fmt.Println(err)
	}

	guild, err := repository.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprint(err),
		})
		return
	}

	// 直近のSQイベントを取得
	sqTitle := ""
	for _, sq := range guild.SQList {
		if sq.Timestamp.After(time.Now()) {
			sqTitle = sq.Title
			break
		}
	}

	// SQ Member の取得
	members, err := repository.GetSQMembers(ctx, guild, sqTitle)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprint(err),
		})
		return
	}

	// SQ Memberにメンション
	message := ""
	for _, member := range members {
		if member.UserID == i.Member.User.ID {
			continue
		}
		message += fmt.Sprintf("<@%s> ", member.UserID)
	}
	if message == "" {
		message = "メンションする対象がいません。"
	}
	if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: message,
	}); err != nil {
		fmt.Println(err)
	}
}
