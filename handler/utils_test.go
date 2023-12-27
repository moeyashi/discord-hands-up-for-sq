package handler

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

// 2023-10-24 ~ 2023-10-29 のSQイベント
// @everyone last SQ events of S9:
// #1760 2v2: 2023年10月24日 火曜日 19:00
// #1761 3v3: 2023年10月25日 水曜日 03:00
// #1762 2v2: 2023年10月25日 水曜日 10:00
// #1763 2v2: 2023年10月25日 水曜日 23:00
// #1764 2v2: 2023年10月26日 木曜日 06:00
// #1765 6v6: 2023年10月26日 木曜日 11:00
// #1766 3v3: 2023年10月26日 木曜日 21:00
// #1767 2v2: 2023年10月27日 金曜日 03:00
// #1768 3v3: 2023年10月27日 金曜日 12:00
// #1769 4v4: 2023年10月27日 金曜日 23:00
// #1770 3v3: 2023年10月28日 土曜日 06:00
// #1771 6v6: 2023年10月28日 土曜日 12:00
// #1772 2v2: 2023年10月28日 土曜日 19:00
// #1773 2v2: 2023年10月29日 日曜日 03:00
// #1774 2v2: 2023年10月29日 日曜日 09:00
var sampleSQInfo = "@everyone last SQ events of S9:\n`#1760` **2v2:** <t:1698141600:F>\n`#1761` **3v3:** <t:1698170400:F>\n`#1762` **2v2:** <t:1698195600:F>\n`#1763` **2v2:** <t:1698242400:F>\n`#1764` **2v2:** <t:1698267600:F>\n`#1765` **6v6:** <t:1698285600:F>\n`#1766` **3v3:** <t:1698321600:F>\n`#1767` **2v2:** <t:1698343200:F>\n`#1768` **3v3:** <t:1698375600:F>\n`#1769` **4v4:** <t:1698415200:F>\n`#1770` **3v3:** <t:1698440400:F>\n`#1771` **6v6:** <t:1698462000:F>\n`#1772` **2v2:** <t:1698487200:F>\n`#1773` **2v2:** <t:1698516000:F>\n`#1774` **2v2:** <t:1698537600:F>"

func Test_createSetCommandsInFuture_今日明日のイベントが取得できる(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	result := createSetCommandsInFuture(sampleSQInfo, time.Date(2023, 10, 27, 0, 0, 0, 0, jst))
	expected := []string{
		"/hands-up set hour:27日03:00 2v2 number:12",
		"/hands-up set hour:27日12:00 3v3 number:12",
		"/hands-up set hour:27日23:00 4v4 number:12",
		"/hands-up set hour:28日06:00 3v3 number:12",
		"/hands-up set hour:28日12:00 6v6 number:12",
		"/hands-up set hour:28日19:00 2v2 number:12",
		"/hands-up set hour:29日03:00 2v2 number:12",
		"/hands-up set hour:29日09:00 2v2 number:12",
	}
	if len(result) != len(expected) {
		t.Fatalf("len(result) = %d, want %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("result[%d] = %s, want %s", i, v, expected[i])
		}
	}
}

func Test_createSetCommandsInFuture_未来のイベントがない場合_空のsliceを返却する(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	tests := []struct {
		name string
		now  time.Time
	}{
		{
			name: "nowが2023-10-29 09:00:01の場合",
			now:  time.Date(2023, 10, 29, 9, 0, 1, 0, jst),
		},
		{
			name: "nowが2023-10-30の場合",
			now:  time.Date(2023, 10, 30, 0, 0, 0, 0, jst),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createSetCommandsInFuture(sampleSQInfo, tt.now)
			expected := []string{}
			if len(result) != 0 {
				t.Errorf("len(result) = %d, want %d", len(result), len(expected))
			}
		})
	}
}

func Test_createOutCommandsForAll(t *testing.T) {
	sampleHandsUpNow := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label: "27日03:00",
				},
				&discordgo.Button{
					Label: "27日12:00",
				},
			},
		},
	}
	actual := createOutCommandsForAll(sampleHandsUpNow)
	expected := []string{
		"/hands-up out hour:27日03:00",
		"/hands-up out hour:27日12:00",
	}
	if len(actual) != len(expected) {
		t.Fatalf("len(actual) = %d, want %d", len(actual), len(expected))
	}
	for i, v := range actual {
		if v != expected[i] {
			t.Errorf("actual[%d] = %s, want %s", i, v, expected[i])
		}
	}
}

func Test_sqListInFuture_現在日時以降のものがすべて取得できる(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	actual := sqListInFuture(sampleSQInfo, time.Date(2023, 10, 26, 0, 0, 0, 0, jst))
	expected := []repository.SQ{
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
	if len(actual) != len(expected) {
		t.Fatalf("len(actual) = %d, want %d", len(actual), len(expected))
	}
	for i, v := range actual {
		if !assertSQ(v, expected[i]) {
			t.Errorf("actual[%d] = %s, want %s", i, v.Title, expected[i].Title)
		}
	}
}

func assertSQ(actual repository.SQ, expected repository.SQ) bool {
	return actual.ID == expected.ID &&
		actual.Title == expected.Title &&
		actual.Format == expected.Format &&
		actual.Timestamp.Equal(expected.Timestamp)
}

func Test_makeSQListEmbedFieldName(t *testing.T) {
	tests := []struct {
		name     string
		sq       repository.SQ
		expected string
	}{
		{name: "memberがいなければtitleのみ", sq: repository.SQ{Title: "26日06:00 2v2"}, expected: "26日06:00 2v2"},
		{name: "memberがいればtitleとmember数", sq: repository.SQ{Title: "26日06:00 2v2", Members: []repository.Member{
			{MemberType: repository.MemberTypesParticipant},
			{MemberType: repository.MemberTypesTemporary},
			{MemberType: repository.MemberTypesTemporary},
			{MemberType: repository.MemberTypesSub},
			{MemberType: repository.MemberTypesSub},
			{MemberType: repository.MemberTypesSub},
		}}, expected: "26日06:00 2v2 (can 1, temp 2, sub 3)"},
		{name: "memberがいればtitleとmember数 canのみ", sq: repository.SQ{Title: "26日06:00 2v2", Members: []repository.Member{
			{MemberType: repository.MemberTypesParticipant},
		}}, expected: "26日06:00 2v2 (can 1)"},
		{name: "memberがいればtitleとmember数 tempのみ", sq: repository.SQ{Title: "26日06:00 2v2", Members: []repository.Member{
			{MemberType: repository.MemberTypesTemporary},
		}}, expected: "26日06:00 2v2 (temp 1)"},
		{name: "memberがいればtitleとmember数 subのみ", sq: repository.SQ{Title: "26日06:00 2v2", Members: []repository.Member{
			{MemberType: repository.MemberTypesSub},
		}}, expected: "26日06:00 2v2 (sub 1)"},
		{name: "memberがいればtitleとmember数 canとtemp", sq: repository.SQ{Title: "26日06:00 2v2", Members: []repository.Member{
			{MemberType: repository.MemberTypesParticipant},
			{MemberType: repository.MemberTypesTemporary},
			{MemberType: repository.MemberTypesTemporary},
		}}, expected: "26日06:00 2v2 (can 1, temp 2)"},
		{name: "memberがいればtitleとmember数 canとsub", sq: repository.SQ{Title: "26日06:00 2v2", Members: []repository.Member{
			{MemberType: repository.MemberTypesParticipant},
			{MemberType: repository.MemberTypesSub},
			{MemberType: repository.MemberTypesSub},
		}}, expected: "26日06:00 2v2 (can 1, sub 2)"},
		{name: "memberがいればtitleとmember数 tempとsub", sq: repository.SQ{Title: "26日06:00 2v2", Members: []repository.Member{
			{MemberType: repository.MemberTypesTemporary},
			{MemberType: repository.MemberTypesSub},
			{MemberType: repository.MemberTypesSub},
		}}, expected: "26日06:00 2v2 (temp 1, sub 2)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := makeSQListEmbedFieldName(tt.sq)
			if actual != tt.expected {
				t.Errorf("actual = %s, want %s", actual, tt.expected)
			}
		})
	}
}

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
	actual, _ := createSQListInteractionResponse(sqList, time.Date(2023, 10, 26, 0, 0, 0, 0, jst))
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
