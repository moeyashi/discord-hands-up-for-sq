package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func HandleSaveResult(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		fmt.Println(err)
		return
	}

	msg := i.ApplicationCommandData().Resolved.Messages[i.ApplicationCommandData().TargetID]
	sheat, err := convertMessageToSheat(msg)

	if err != nil {
		fmt.Println(err)
		if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: err.Error(),
		}); err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	guild, err := repository.GetGuild(ctx, i.GuildID)
	if err != nil {
		fmt.Println(err)
		if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "データベースからの取得に失敗しました。",
		}); err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	if guild.Spreadsheet == "" {
		if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "スプレッドシートが設定されていません。`/results url-set`でスプレッドシートを設定し、`discordbot@hands-up-for-sq.iam.gserviceaccount.com`に編集権限を付与してください。",
		}); err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	credential := option.WithCredentialsFile("./google-api-credential.json")
	srv, err := sheets.NewService(ctx, credential)
	if err != nil {
		fmt.Println(err)
		if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "スプレッドシートへの接続に失敗しました。",
		}); err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	values := [][]interface{}{}
	for _, result := range sheat.Results {
		values = append(values, []interface{}{
			toSpreadsheetTime(sheat.Timestamp),
			sheat.EnemyName,
			result.TrackName,
			result.Places[0],
			result.Places[1],
			result.Places[2],
			result.Places[3],
			result.Places[4],
			result.Places[5],
			toDifference(result.Places),
		})
	}

	_, err = srv.Spreadsheets.Values.Append(guild.Spreadsheet, "A1", &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do()
	if err != nil {
		fmt.Println(err)
		if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "スプレッドシートへの書き込みに失敗しました。",
		}); err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("%s vs %s 戦績を保存しました。", toDiscordTimestamp(sheat.Timestamp, "f"), sheat.EnemyName),
	}); err != nil {
		fmt.Println(err)
		return
	}
}

type sheat struct {
	Timestamp time.Time
	EnemyName string
	Results   [12]result
}

type result struct {
	TrackName string
	Places    [6]int8
}

func convertMessageToSheat(msg *discordgo.Message) (sheat, error) {
	if !strings.Contains(msg.Embeds[0].Title, "即時集計") {
		return sheat{}, fmt.Errorf("即時集計ではありません。")
	}

	if !strings.Contains(msg.Embeds[0].Title, "6v6") {
		return sheat{}, fmt.Errorf("6v6ではありません。")
	}

	if !strings.Contains(msg.Embeds[0].Description, "@0") {
		return sheat{}, fmt.Errorf("mogiが終了していません。")
	}

	results, err := extractResults(msg.Embeds[0].Fields)
	if err != nil {
		return sheat{}, err
	}
	return sheat{
		Timestamp: msg.Timestamp,
		EnemyName: extractEnemyName(msg.Embeds[0].Title),
		Results:   results,
	}, nil
}

// "即時集計 6v6\nAiZ - Bup" -> "Bup"
func extractEnemyName(title string) string {
	return strings.Split(title, " - ")[1]
}

func extractResults(fields []*discordgo.MessageEmbedField) ([12]result, error) {
	var results [12]result
	for i, field := range fields {
		places, err := extractPlaces(field.Value)
		if err != nil {
			return results, err
		}
		results[i] = result{
			TrackName: extractTrackName(field.Name),
			Places:    places,
		}
	}
	return results, nil
}

func extractTrackName(name string) string {
	splitted := strings.Split(name, " ")
	return splitted[len(splitted)-1]
}

func extractPlaces(value string) ([6]int8, error) {
	placesStr := strings.ReplaceAll(strings.Split(value, " | ")[1], "`", "")
	var places [6]int8
	for i, placeStr := range strings.Split(placesStr, ",") {
		place, err := strconv.Atoi(placeStr)
		if err != nil {
			return places, err
		}
		places[i] = int8(place)
	}
	return places, nil
}

func toSpreadsheetTime(t time.Time) float64 {
	// Google Spreadsheet's epoch is "1899-12-30T00:00:00Z"
	epoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	duration := t.Sub(epoch)
	return float64(duration) / float64(24*time.Hour)
}

// 1は15, 2は12, 3は10, 4は9, 5は8, 6は7, 7は6, 8は5, 9は4, 10は3, 11は2, 12は1点
// 全部で82点なので(合計 - (82 - 合計))が得失点
func toDifference(places [6]int8) int {
	score := 0
	for _, place := range places {
		switch place {
		case 1:
			score += 15
		case 2:
			score += 12
		case 3:
			score += 10
		case 4:
			score += 9
		case 5:
			score += 8
		case 6:
			score += 7
		case 7:
			score += 6
		case 8:
			score += 5
		case 9:
			score += 4
		case 10:
			score += 3
		case 11:
			score += 2
		case 12:
			score += 1
		}
	}
	return score - (82 - score)
}

func toDiscordTimestamp(t time.Time, format string) string {
	return fmt.Sprintf("<t:%d:%s>", t.Unix(), format)
}
