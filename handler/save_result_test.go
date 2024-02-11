package handler

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Test_convertMessageToSheat(t *testing.T) {
	actual, err := convertMessageToSheat(makeSampleMessage())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if actual.Timestamp != time.Date(2024, time.February, 9, 15, 47, 8, 614, time.UTC) {
		t.Errorf("actual: %s, expected: %s", actual.Timestamp, time.Date(2024, time.February, 9, 15, 47, 8, 614, time.UTC))
	}
	if actual.EnemyName != "Bup" {
		t.Errorf("actual: %s, expected: %s", actual.EnemyName, "Bup")
	}
	if len(actual.Results) != 12 {
		t.Errorf("actual: %d, expected: %d", len(actual.Results), 12)
	}
	assertResult(t, actual.Results[0], result{TrackName: "アムステルダム", Places: [6]int8{4, 5, 7, 9, 10, 11}})
	assertResult(t, actual.Results[1], result{TrackName: "ワリスタ", Places: [6]int8{2, 4, 5, 6, 8, 11}})
	assertResult(t, actual.Results[2], result{TrackName: "ワリスノ", Places: [6]int8{4, 5, 8, 9, 10, 11}})
	assertResult(t, actual.Results[3], result{TrackName: "ロサンゼルス", Places: [6]int8{3, 4, 5, 8, 9, 12}})
	assertResult(t, actual.Results[4], result{TrackName: "ジャングル", Places: [6]int8{4, 6, 7, 9, 10, 11}})
	assertResult(t, actual.Results[5], result{TrackName: "スノボ", Places: [6]int8{3, 6, 7, 9, 11, 12}})
	assertResult(t, actual.Results[6], result{TrackName: "64虹", Places: [6]int8{1, 3, 5, 6, 8, 12}})
	assertResult(t, actual.Results[7], result{TrackName: "GBAマリサ", Places: [6]int8{4, 6, 7, 9, 11, 12}})
	assertResult(t, actual.Results[8], result{TrackName: "ハイラル", Places: [6]int8{4, 5, 8, 10, 11, 12}})
	assertResult(t, actual.Results[9], result{TrackName: "7虹", Places: [6]int8{3, 7, 8, 9, 11, 12}})
	assertResult(t, actual.Results[10], result{TrackName: "ロンドン", Places: [6]int8{1, 2, 3, 7, 8, 9}})
	assertResult(t, actual.Results[11], result{TrackName: "しんでん", Places: [6]int8{1, 2, 6, 8, 11, 12}})
}

func Test_convertMessageToSheat_error_即時集計でない場合(t *testing.T) {
	msg := makeSampleMessage()
	msg.Embeds[0].Title = "test"
	_, err := convertMessageToSheat(msg)
	if err == nil {
		t.Errorf("expected error, but got nil")
	}
	if err.Error() != "即時集計ではありません。" {
		t.Errorf("actual: %s, expected: %s", err.Error(), "即時集計ではありません。")
	}
}

func Test_convertMessageToSheat_error_6v6でない場合(t *testing.T) {
	msg := makeSampleMessage()
	msg.Embeds[0].Title = "即時集計 4v4\nAiZ - Bup"
	_, err := convertMessageToSheat(msg)
	if err == nil {
		t.Errorf("expected error, but got nil")
	}
	if err.Error() != "6v6ではありません。" {
		t.Errorf("actual: %s, expected: %s", err.Error(), "6v6ではありません。")
	}
}

func Test_convertMessageToSheat_error_mogiが終了していない場合(t *testing.T) {
	msg := makeSampleMessage()
	msg.Embeds[0].Description = "`429 : 555(-126)` x@xx"
	_, err := convertMessageToSheat(msg)
	if err == nil {
		t.Errorf("expected error, but got nil")
	}
	if err.Error() != "mogiが終了していません。" {
		t.Errorf("actual: %s, expected: %s", err.Error(), "mogiが終了していません。")
	}
}

func assertResult(t *testing.T, actual result, expected result) {
	if actual.TrackName != expected.TrackName {
		t.Errorf("actual: %s, expected: %s", actual.TrackName, expected.TrackName)
	}
	for i, actualPlace := range actual.Places {
		if actualPlace != expected.Places[i] {
			t.Errorf("%s actual: %d, expected: %d", expected.TrackName, actualPlace, expected.Places[i])
		}
	}
}

func Test_extractPlaces(t *testing.T) {
	actual, err := extractPlaces("`32 : 50(-18)` | `4,5,7,9,10,11`")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := [6]int8{4, 5, 7, 9, 10, 11}
	for i, actualPlace := range actual {
		if actualPlace != expected[i] {
			t.Errorf("actual: %d, expected: %d", actualPlace, expected[i])
		}
	}
}

func Test_toDiscordTimestamp(t *testing.T) {
	actual := toDiscordTimestamp(time.Date(2024, time.February, 9, 15, 47, 8, 614, time.UTC), "f")
	if actual != "<t:1707493628:f>" {
		t.Errorf("actual: %s, expected: %s", actual, "<t:1707493628:f>")
	}
}

/**
 * make sample meassage
 *
 * 元ネタ
 * {"type":"rich","title":"即時集計 6v6\nAiZ - Bup","description":"`429 : 555(-126)`  `@0`","color":13814454,"image":{"url":"https://cdn.discordapp.com/attachments/1163174408597282929/1205540398420729916/result.png?ex=65d8bdfc\u0026is=65c648fc\u0026hm=2983c56713fb246983673a68790317ceb0c038346c9ca884e9d6c704f72f0d77\u0026","proxy_url":"https://media.discordapp.net/attachments/1163174408597282929/1205540398420729916/result.png?ex=65d8bdfc\u0026is=65c648fc\u0026hm=2983c56713fb246983673a68790317ceb0c038346c9ca884e9d6c704f72f0d77\u0026","width":1280,"height":720},"fields":[{"name":"1  - \u003c:bAD:1083984326015848519\u003e アムステルダム","value":"`32 : 50(-18)` | `4,5,7,9,10,11`"},{"name":"2  - \u003c:rWS:968803704042045460\u003e ワリスタ","value":"`43 : 39(+4)` | `2,4,5,6,8,11`"},{"name":"3  - \u003c:MW_:968803703903637544\u003e ワリスノ","value":"`31 : 51(-20)` | `4,5,8,9,10,11`"},{"name":"4  - \u003c:bLAL:1128798442370637954\u003e ロサンゼルス","value":"`37 : 45(-8)` | `3,4,5,8,9,12`"},{"name":"5  - \u003c:rDKJ:968803704130117672\u003e ジャングル","value":"`31 : 51(-20)` | `4,6,7,9,10,11`"},{"name":"6  - \u003c:bDKS:1083984413026684938\u003e スノボ","value":"`30 : 52(-22)` | `3,6,7,9,11,12`"},{"name":"7  - \u003c:rRRd:968816715473510421\u003e 64虹","value":"`46 : 36(+10)` | `1,3,5,6,8,12`"},{"name":"8  - \u003c:rMC:968803704222396446\u003e GBAマリサ","value":"`29 : 53(-24)` | `4,6,7,9,11,12`"},{"name":"9  - \u003c:dHC:968816715192471572\u003e ハイラル","value":"`28 : 54(-26)` | `4,5,8,10,11,12`"},{"name":"10  - \u003c:bRR:1055496691421302794\u003e 7虹","value":"`28 : 54(-26)` | `3,7,8,9,11,12`"},{"name":"11  - \u003c:bLL:1055488205631271044\u003e ロンドン","value":"`52 : 30(+22)` | `1,2,3,7,8,9`"},{"name":"12  - \u003c:unknown:1050033104569499698\u003e しんでん","value":"`42 : 40(+2)` | `1,2,6,8,11,12`"}]}
 */
func makeSampleMessage() *discordgo.Message {
	return &discordgo.Message{
		Timestamp: time.Date(2024, time.February, 9, 15, 47, 8, 614, time.UTC),
		Embeds: []*discordgo.MessageEmbed{{
			Title:       "即時集計 6v6\nAiZ - Bup",
			Description: "`429 : 555(-126)`  `@0`",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "1  - <:bAD:1083984326015848519> アムステルダム",
					Value: "`32 : 50(-18)` | `4,5,7,9,10,11`",
				},
				{
					Name:  "2  - <:rWS:968803704042045460> ワリスタ",
					Value: "`43 : 39(+4)` | `2,4,5,6,8,11`",
				},
				{
					Name:  "3  - <:MW_:968803703903637544> ワリスノ",
					Value: "`31 : 51(-20)` | `4,5,8,9,10,11`",
				},
				{
					Name:  "4  - <:bLAL:1128798442370637954> ロサンゼルス",
					Value: "`37 : 45(-8)` | `3,4,5,8,9,12`",
				},
				{
					Name:  "5  - <:rDKJ:968803704130117672> ジャングル",
					Value: "`31 : 51(-20)` | `4,6,7,9,10,11`",
				},
				{
					Name:  "6  - <:bDKS:1083984413026684938> スノボ",
					Value: "`30 : 52(-22)` | `3,6,7,9,11,12`",
				},
				{
					Name:  "7  - <:rRRd:968816715473510421> 64虹",
					Value: "`46 : 36(+10)` | `1,3,5,6,8,12`",
				},
				{
					Name:  "8  - <:rMC:968803704222396446> GBAマリサ",
					Value: "`29 : 53(-24)` | `4,6,7,9,11,12`",
				},
				{
					Name:  "9  - <:dHC:968816715192471572> ハイラル",
					Value: "`28 : 54(-26)` | `4,5,8,10,11,12`",
				},
				{
					Name:  "10  - <:bRR:1055496691421302794> 7虹",
					Value: "`28 : 54(-26)` | `3,7,8,9,11,12`",
				},
				{
					Name:  "11  - <:bLL:1055488205631271044> ロンドン",
					Value: "`52 : 30(+22)` | `1,2,3,7,8,9`",
				},
				{
					Name:  "12  - <:unknown:1050033104569499698> しんでん",
					Value: "`42 : 40(+2)` | `1,2,6,8,11,12`",
				},
			},
		}},
	}
}
