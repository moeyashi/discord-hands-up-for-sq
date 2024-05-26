package response

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func MakeErrorInteractionResponse(err error) *discordgo.InteractionResponse {
	fmt.Println(err)
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprint(err),
		},
	}
}

func MakeErrorWebhookParams(err error) *discordgo.WebhookParams {
	fmt.Println(err)
	return &discordgo.WebhookParams{
		Flags:   discordgo.MessageFlagsEphemeral,
		Content: fmt.Sprint(err),
	}
}
