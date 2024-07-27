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
