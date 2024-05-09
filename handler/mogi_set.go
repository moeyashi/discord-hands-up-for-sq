package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func HandleMogiSet(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repo repository.Repository) {
	guild, err := repo.GetGuild(ctx, i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, makeErrorResponse(err))
		return
	}

	month := i.ApplicationCommandData().Options[0].Options[0].IntValue()
	date := i.ApplicationCommandData().Options[0].Options[1].IntValue()
	hour := int64(0)
	if len(i.ApplicationCommandData().Options[0].Options) == 3 {
		hour = i.ApplicationCommandData().Options[0].Options[2].IntValue()
	}

	repo.AppendMogiList(ctx, guild, *repository.MakeMogi(time.Now(), month, date, hour))

	res, err := createMogiListInteractionResponse(guild.MogiList)
	if err != nil {
		s.InteractionRespond(i.Interaction, makeErrorResponse(err))
		return
	}

	if err := s.InteractionRespond(i.Interaction, res); err != nil {
		fmt.Println(err)
		return
	}

	if err := deleteOldMessages(s, i.ChannelID); err != nil {
		fmt.Println(err)
		return
	}
}
