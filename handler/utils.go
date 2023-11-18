package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func sqListInFuture(sqInfo string, now time.Time) []string {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Println(err)
		return []string{}
	}
	nowUnix := now.Unix()
	re := regexp.MustCompile("`#(\\d+)` \\*\\*(\\dv\\d):\\*\\* <t:(\\d+):F>")
	results := re.FindAllStringSubmatch(sqInfo, -1)
	sqList := []string{}
	for _, submatches := range results {
		timestamp, err := strconv.ParseInt(submatches[3], 10, 64)
		if err != nil {
			log.Println(err)
			return []string{}
		}
		if nowUnix <= timestamp {
			hourContent := time.Unix(timestamp, 0).In(jst).Format("2日15:04")
			mogiFormat := submatches[2]
			sqList = append(sqList, fmt.Sprintf("%s %s", hourContent, mogiFormat))
		}
	}
	return sqList
}

func createSetCommandsInFuture(sqInfo string, now time.Time) []string {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Println(err)
		return []string{}
	}
	nowUnix := now.Unix()
	re := regexp.MustCompile("`#(\\d+)` \\*\\*(\\dv\\d):\\*\\* <t:(\\d+):F>")
	results := re.FindAllStringSubmatch(sqInfo, -1)
	commands := []string{}
	for _, submatches := range results {
		timestamp, err := strconv.ParseInt(submatches[3], 10, 64)
		if err != nil {
			log.Println(err)
			return []string{}
		}
		if nowUnix <= timestamp {
			hourContent := time.Unix(timestamp, 0).In(jst).Format("2日15:04")
			mogiFormat := submatches[2]
			command := fmt.Sprintf("/hands-up set hour:%s %s number:12", hourContent, mogiFormat)
			commands = append(commands, command)
		}
	}
	return commands
}

func createOutCommandsForAll(handsUpNow []discordgo.MessageComponent) []string {
	commands := []string{}
	for _, actionsRowComponent := range handsUpNow {
		if actionsRowComponent.Type() != discordgo.ActionsRowComponent {
			continue
		}
		rowJson, err := actionsRowComponent.MarshalJSON()
		if err != nil {
			return []string{}
		}
		var row discordgo.ActionsRow
		if err := row.UnmarshalJSON(rowJson); err != nil {
			return []string{}
		}
		for _, buttonComponent := range row.Components {
			if buttonComponent.Type() != discordgo.ButtonComponent {
				continue
			}
			buttonJson, err := buttonComponent.MarshalJSON()
			if err != nil {
				return []string{}
			}
			var button discordgo.Button
			if err := json.Unmarshal(buttonJson, &button); err != nil {
				return []string{}
			}
			commands = append(commands, fmt.Sprintf("/hands-up out hour:%s", button.Label))
		}
	}
	return commands
}

func createSQListInteractionResponse(ctx context.Context, guildID string, repository repository.Repository) (*discordgo.InteractionResponse, error) {
	guild, err := repository.GetGuild(ctx, guildID)
	if err != nil {
		return nil, err
	}
	if len(guild.SQList) == 0 {
		return nil, fmt.Errorf("SQが登録されていません。")
	}

	embedFields := []*discordgo.MessageEmbedField{}
	components := []discordgo.MessageComponent{}
	for _, sq := range guild.SQList {
		members := sq.Members
		embedFieldsValue := "なし"
		userNames := []string{}
		for _, member := range members {
			userNames = append(userNames, member.UserName)
		}
		if len(userNames) > 0 {
			embedFieldsValue = strings.Join(userNames, ",")
		}

		embedFields = append(embedFields, &discordgo.MessageEmbedField{
			Name:  sq.Title,
			Value: embedFieldsValue,
		})
		if len(components) < 5 {
			components = append(components, discordgo.Button{
				CustomID: "button_" + sq.Title,
				Label:    sq.Title,
				Style:    discordgo.DangerButton,
			})
		}
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{{Fields: embedFields}},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: components,
				},
			},
		},
	}, nil
}
