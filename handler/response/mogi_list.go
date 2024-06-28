package response

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/constant"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func MakeMogiListInteractionResponse(mogiList []repository.Mogi) (*discordgo.InteractionResponse, error) {
	embedFields := []*discordgo.MessageEmbedField{}
	components := []discordgo.MessageComponent{}

	for _, mogi := range mogiList {
		embedFieldsValue := "なし"
		userNames := []string{}
		for _, member := range mogi.Members {
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
			Name:  makeMogiListEmbedFieldName(mogi),
			Value: embedFieldsValue,
		})
		components = append(components, discordgo.Button{
			CustomID: "button_mogi_" + mogi.Title(),
			Label:    mogi.Title(),
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
			Content: "内戦リスト",
		},
	}
	if len(embedFields) != 0 {
		ret.Data.Embeds = []*discordgo.MessageEmbed{{Fields: embedFields}}
	} else {
		ret.Data.Content = "内戦リストはありません。\n`mogi set`コマンドを実行して内戦リストを設定してください。"
	}
	if len(rows) != 0 {
		ret.Data.Components = rows
	}
	return ret, nil
}

func makeMogiListEmbedFieldName(
	mogi repository.Mogi,
) string {
	canMembersCount := 0
	tempMembersCount := 0
	subMembersCount := 0
	for _, member := range mogi.Members {
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
		return mogi.Title()
	}
	return fmt.Sprintf("%s (%s)", mogi.Title(), strings.Join(members, ", "))
}

func MakeMogiListSelect(userID string, mogiList []repository.Mogi, customID constant.MogiListSelectCustomID) *discordgo.SelectMenu {
	memberType := customID.ToMemberTypes()
	options := []discordgo.SelectMenuOption{}
	for _, mogi := range mogiList {
		if repository.IndexOfSameRegistered(mogi.Members, userID, memberType) >= 0 {
			continue
		}
		options = append(options, discordgo.SelectMenuOption{
			Label: mogi.Title(),
			Value: mogi.Title(),
		})
	}
	return &discordgo.SelectMenu{
		CustomID: string(customID),
		Options:  options,
	}
}
