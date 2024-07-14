package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
	"github.com/moeyashi/discord-hands-up-for-sq/usecase"
)

func HandleLoungeSQInfo(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, repo repository.Repository) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Username != "MK8DX 150cc Lounge #sq-info" || m.Author.Discriminator != "0000" {
		return
	}

	guild, err := repo.GetGuild(ctx, m.GuildID)
	if err != nil {
		handleLoungeSQInfoError(s, m, err)
		return
	}

	now := time.Now()

	sqList := usecase.NewSQList(
		guild.SQList,
		sqListInFuture(m.Content, now),
		now,
	)

	if err := repo.PutSQList(ctx, guild, sqList); err != nil {
		handleLoungeSQInfoError(s, m, err)
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "SQリストを自動で更新しました")
	if err != nil {
		log.Printf("Cannot send message to channel %v: %v", m.ChannelID, err)
	}
}

func handleLoungeSQInfoError(s *discordgo.Session, m *discordgo.MessageCreate, err error) {
	if err == nil {
		return
	}
	_, _err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("sq-infoの自動処理に失敗しました。husq setコマンドを実行してください。: %v", err))
	if _err != nil {
		log.Printf("Cannot send message to channel %v: %v", m.ChannelID, _err)
	}
}
