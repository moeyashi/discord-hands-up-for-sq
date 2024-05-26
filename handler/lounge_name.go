package handler

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/response"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleLoungeName(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	guild, err := repository.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, response.MakeErrorInteractionResponse(err))
		return
	}

	filteredSQList := filterSQListForDisplay(guild.SQList, time.Now())
	options := []discordgo.SelectMenuOption{}
	for _, sq := range filteredSQList {
		options = append(options, discordgo.SelectMenuOption{
			Label: sq.Title,
			Value: sq.Title,
		})
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.SelectMenu{
							CustomID: "lounge_name_select",
							Options:  options,
						},
					},
				},
			},
		},
	})

	if err != nil {
		s.InteractionRespond(i.Interaction, response.MakeErrorInteractionResponse(err))
	}
}
