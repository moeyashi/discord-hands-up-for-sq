package response

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

func Test_createSQListInteractionResponse_現在時刻から3日後まで取得できる_日付をまたいだ直後の場合(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	sqList := []repository.SQ{
		{ID: "1760", Title: "24日19:00 2v2", Format: "2v2", Timestamp: time.Date(2023, 10, 24, 19, 0, 0, 0, jst)},
		{ID: "1761", Title: "25日03:00 3v3", Format: "3v3", Timestamp: time.Date(2023, 10, 25, 3, 0, 0, 0, jst)},
		{ID: "1762", Title: "25日10:00 2v2", Format: "2v2", Timestamp: time.Date(2023, 10, 25, 10, 0, 0, 0, jst)},
		{ID: "1763", Title: "25日23:00 2v2", Format: "2v2", Timestamp: time.Date(2023, 10, 25, 23, 0, 0, 0, jst)},
		{ID: "1764", Title: "26日06:00 2v2", Format: "2v2", Timestamp: time.Date(2023, 10, 26, 6, 0, 0, 0, jst)},
		{ID: "1765", Title: "26日11:00 6v6", Format: "6v6", Timestamp: time.Date(2023, 10, 26, 11, 0, 0, 0, jst)},
		{ID: "1766", Title: "26日21:00 3v3", Format: "3v3", Timestamp: time.Date(2023, 10, 26, 21, 0, 0, 0, jst)},
		{ID: "1767", Title: "27日03:00 2v2", Format: "2v2", Timestamp: time.Date(2023, 10, 27, 3, 0, 0, 0, jst)},
		{ID: "1768", Title: "27日12:00 3v3", Format: "3v3", Timestamp: time.Date(2023, 10, 27, 12, 0, 0, 0, jst)},
		{ID: "1769", Title: "27日23:00 4v4", Format: "4v4", Timestamp: time.Date(2023, 10, 27, 23, 0, 0, 0, jst)},
		{ID: "1770", Title: "28日06:00 3v3", Format: "3v3", Timestamp: time.Date(2023, 10, 28, 6, 0, 0, 0, jst)},
		{ID: "1771", Title: "28日12:00 6v6", Format: "6v6", Timestamp: time.Date(2023, 10, 28, 12, 0, 0, 0, jst)},
		{ID: "1772", Title: "28日19:00 2v2", Format: "2v2", Timestamp: time.Date(2023, 10, 28, 19, 0, 0, 0, jst)},
		{ID: "1773", Title: "29日03:00 2v2", Format: "2v2", Timestamp: time.Date(2023, 10, 29, 3, 0, 0, 0, jst)},
		{ID: "1774", Title: "29日09:00 2v2", Format: "2v2", Timestamp: time.Date(2023, 10, 29, 9, 0, 0, 0, jst)},
	}
	actual, _ := MakeSQListInteractionResponse(sqList, time.Date(2023, 10, 26, 0, 0, 0, 0, jst))
	expected := &discordgo.InteractionResponseData{
		Content: "SQリスト",
		Embeds: []*discordgo.MessageEmbed{
			{
				Fields: []*discordgo.MessageEmbedField{
					{Name: "26日06:00 2v2", Value: "なし"},
					{Name: "26日11:00 6v6", Value: "なし"},
					{Name: "26日21:00 3v3", Value: "なし"},
					{Name: "27日03:00 2v2", Value: "なし"},
					{Name: "27日12:00 3v3", Value: "なし"},
					{Name: "27日23:00 4v4", Value: "なし"},
					{Name: "28日06:00 3v3", Value: "なし"},
					{Name: "28日12:00 6v6", Value: "なし"},
					{Name: "28日19:00 2v2", Value: "なし"},
				},
			},
		},
	}
	if !assertInteractionResponseData(actual.Data, expected) {
		t.Errorf("actual = %v, want %v", actual, expected)
	}
}

func assertInteractionResponseData(actual *discordgo.InteractionResponseData, expected *discordgo.InteractionResponseData) bool {
	if actual.Content != expected.Content {
		return false
	}
	if len(actual.Embeds) != len(expected.Embeds) {
		return false
	}
	for i, v := range actual.Embeds {
		if !assertEmbed(v, expected.Embeds[i]) {
			return false
		}
	}
	return true
}

func assertEmbed(actual *discordgo.MessageEmbed, expected *discordgo.MessageEmbed) bool {
	if len(actual.Fields) != len(expected.Fields) {
		return false
	}
	for i, v := range actual.Fields {
		if !assertEmbedField(v, expected.Fields[i]) {
			return false
		}
	}
	return true
}

func assertEmbedField(actual *discordgo.MessageEmbedField, expected *discordgo.MessageEmbedField) bool {
	return actual.Name == expected.Name &&
		actual.Value == expected.Value
}
