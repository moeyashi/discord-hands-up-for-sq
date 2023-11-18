package handler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func ListSQ(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	guildID := i.GuildID
	guild, err := repository.GetGuild(ctx, guildID)
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
	if len(guild.SQList) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "SQが登録されていません。",
			},
		})
		return
	}

	embedFields := []*discordgo.MessageEmbedField{}
	components := []discordgo.MessageComponent{}
	for _, sq := range guild.SQList {
		embedFields = append(embedFields, &discordgo.MessageEmbedField{
			Name:  sq.Title,
			Value: "なし",
		})
		if len(components) < 5 {
			components = append(components, discordgo.Button{
				CustomID: "button_" + sq.Title,
				Label:    sq.Title,
				Style:    discordgo.DangerButton,
			})
		}
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{{Fields: embedFields}},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: components,
				},
			},
		},
	}); err != nil {
		fmt.Println(err)
	}
}
