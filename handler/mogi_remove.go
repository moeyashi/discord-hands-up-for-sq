package handler

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleMogiRemove(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	guild, err := repository.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, makeErrorResponse(err))
		return
	}

	mogiList, err := repository.GetMogiList(ctx, guild)
	if err != nil {
		s.InteractionRespond(i.Interaction, makeErrorResponse(err))
		return
	}
	options := []discordgo.SelectMenuOption{}
	for _, mogi := range mogiList {
		options = append(options, discordgo.SelectMenuOption{
			Label: mogi.Title(),
			Value: mogi.Title(),
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
							CustomID: "mogi_remove_select",
							Options:  options,
						},
					},
				},
			},
		},
	})

	if err != nil {
		s.InteractionRespond(i.Interaction, makeErrorResponse(err))
		return
	}
}
