package response

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/constant"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func MakeSQListInteractionResponse(sqList []repository.SQ, now time.Time) (*discordgo.InteractionResponse, error) {
	embedFields := []*discordgo.MessageEmbedField{}
	components := []discordgo.MessageComponent{}

	filteredSQList := filterSQListForDisplay(sqList, now)

	for _, sq := range filteredSQList {
		embedFieldsValue := "なし"
		userNames := []string{}
		for _, member := range sq.Members {
			userName := member.UserName
			if member.MemberType == repository.MemberTypesTemporary {
				userName = userName + "(仮)"
			} else if member.MemberType == repository.MemberTypesSub {
				userName = userName + "(sub)"
			}
			userNames = append(userNames, userName)
		}
		if len(userNames) > 0 {
			embedFieldsValue = strings.Join(userNames, ",")
		}

		embedFields = append(embedFields, &discordgo.MessageEmbedField{
			Name:  makeSQListEmbedFieldName(sq),
			Value: embedFieldsValue,
		})
		components = append(components, discordgo.Button{
			CustomID: "button_" + sq.Title,
			Label:    sq.Title,
			Style:    discordgo.SecondaryButton,
		})
	}

	// componentsが5つまでしか入らないため、5つごとにRowを分ける
	actionsRows := []discordgo.ActionsRow{}
	for index, component := range components {
		if index%5 == 0 {
			actionsRows = append(actionsRows, discordgo.ActionsRow{})
		}
		actionsRows[len(actionsRows)-1].Components = append(actionsRows[len(actionsRows)-1].Components, component)
	}

	rows := []discordgo.MessageComponent{}
	for _, actionsRow := range actionsRows {
		rows = append(rows, actionsRow)
	}

	ret := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "SQリスト",
		},
	}
	if len(embedFields) != 0 {
		ret.Data.Embeds = []*discordgo.MessageEmbed{{Fields: embedFields}}
	} else {
		ret.Data.Content = "SQリストはありません。\nsq-infoのメッセージから`husq set`コマンドを実行してSQリストを設定してください。"
	}
	if len(rows) != 0 {
		ret.Data.Components = rows
	}
	return ret, nil
}

func MakeSQListSelect(userID string, sqList []repository.SQ, customID constant.SQListSelectCustomID, now time.Time) *discordgo.SelectMenu {
	filteredSQList := filterSQListForDisplay(sqList, now)
	memberType := customID.ToMemberTypes()
	options := []discordgo.SelectMenuOption{}
	for _, sq := range filteredSQList {
		if repository.IndexOfSameRegistered(sq.Members, userID, memberType) >= 0 {
			continue
		}
		options = append(options, discordgo.SelectMenuOption{
			Label: sq.Title,
			Value: sq.Title,
		})
	}
	return &discordgo.SelectMenu{
		CustomID: string(customID),
		Options:  options,
	}
}

func filterSQListForDisplay(sqList []repository.SQ, now time.Time) []repository.SQ {
	nowUnix := now.Unix()
	filteredSQList := []repository.SQ{}
	for _, sq := range sqList {
		sqUnix := sq.Timestamp.Unix()
		// now ~ 3日後までのSQを表示
		if nowUnix > sqUnix || nowUnix+60*60*24*3 < sqUnix {
			continue
		}
		filteredSQList = append(filteredSQList, sq)
	}
	return filteredSQList
}

func makeSQListEmbedFieldName(
	sq repository.SQ,
) string {
	canMembersCount := 0
	tempMembersCount := 0
	subMembersCount := 0
	for _, member := range sq.Members {
		switch member.MemberType {
		case repository.MemberTypesParticipant:
			canMembersCount++
		case repository.MemberTypesTemporary:
			tempMembersCount++
		case repository.MemberTypesSub:
			subMembersCount++
		}
	}
	members := []string{}
	if canMembersCount > 0 {
		members = append(members, fmt.Sprintf("can %d", canMembersCount))
	}
	if tempMembersCount > 0 {
		members = append(members, fmt.Sprintf("temp %d", tempMembersCount))
	}
	if subMembersCount > 0 {
		members = append(members, fmt.Sprintf("sub %d", subMembersCount))
	}
	if len(members) == 0 {
		return sq.Title
	}
	return fmt.Sprintf("%s (%s)", sq.Title, strings.Join(members, ", "))
}
