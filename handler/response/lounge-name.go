package response

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

type MakeLoungeNameResponseWebhookParamsParameterItem struct {
	Member     repository.Member
	LoungeName *repository.GetLoungeNameResponse
}

func MakeLoungeNameResponseWebhookParams(members []MakeLoungeNameResponseWebhookParamsParameterItem) *discordgo.WebhookParams {
	sumOfMMR := 0
	embedFields := []*discordgo.MessageEmbedField{}
	for _, member := range members {
		nameForGuild := member.Member.UserName
		embedFields = append(embedFields, &discordgo.MessageEmbedField{
			Name:  nameForGuild,
			Value: fmt.Sprintf("%s (%d)", member.LoungeName.Name, member.LoungeName.MMR),
		})
		sumOfMMR += member.LoungeName.MMR
	}

	title := ""
	if len(members) > 0 {
		title = "平均MMR: " + fmt.Sprintf("%.1f", float64(sumOfMMR)/float64(len(members)))
	}

	return &discordgo.WebhookParams{
		Content: "Loungeサーバーでの名前",
		Embeds: []*discordgo.MessageEmbed{{
			Title:  title,
			Fields: embedFields,
		}},
	}
}
