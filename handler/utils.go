package handler

import (
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

func sqListInFuture(sqInfo string, now time.Time) []repository.SQ {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Println(err)
		return []repository.SQ{}
	}
	nowUnix := now.Unix()
	re := regexp.MustCompile("`#(\\d+)` \\*\\*(\\dv\\d):\\*\\* <t:(\\d+):F>")
	results := re.FindAllStringSubmatch(sqInfo, -1)
	sqList := []repository.SQ{}
	for _, submatches := range results {
		timestamp, err := strconv.ParseInt(submatches[3], 10, 64)
		if err != nil {
			log.Println(err)
			return []repository.SQ{}
		}
		if nowUnix <= timestamp {
			hourContent := time.Unix(timestamp, 0).In(jst).Format("2日15:04")
			mogiFormat := submatches[2]
			sqList = append(sqList, repository.SQ{ID: submatches[1], Title: fmt.Sprintf("%s %s", hourContent, mogiFormat), Format: mogiFormat, Timestamp: time.Unix(timestamp, 0)})
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

func createMogiListInteractionResponse(mogiList []repository.Mogi) (*discordgo.InteractionResponse, error) {
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

func indexOfSameMember(members []repository.Member, userID string) int {
	for index, member := range members {
		if member.UserID == userID {
			return index
		}
	}
	return -1
}

func deleteOldMessages(s *discordgo.Session, channelID string) error {
	messages, err := s.ChannelMessages(channelID, 10, "", "", "")
	if err != nil {
		return err
	}
	lastMessage := ""
	for _, message := range messages {
		if message.Author.ID == s.State.User.ID && message.Flags&discordgo.MessageFlagsEphemeral == 0 {
			if lastMessage == "" {
				lastMessage = message.ID
			} else {
				if err := s.ChannelMessageDelete(channelID, message.ID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getDisplayUsername(member *discordgo.Member) string {
	userName := member.Nick
	if userName == "" {
		userName = member.User.Username
	}
	return userName
}

func mogiRoleName(mogi *repository.Mogi) string {
	return fmt.Sprintf("内戦 %s", mogi.Title())
}

func findMogiRole(s *discordgo.Session, guildID string, mogiRoleName string) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if role.Name == mogiRoleName {
			return role, nil
		}
	}
	return nil, nil
}
