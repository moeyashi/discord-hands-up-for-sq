package handler

import (
	"testing"
	"time"
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

func Test_createHandsUpCommandsInFuture_今日明日のイベントが取得できる(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	result := createHandsUpCommandsInFuture(sampleSQInfo, time.Date(2023, 10, 27, 0, 0, 0, 0, jst))
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
		t.Errorf("len(result) = %d, want %d", len(result), len(expected))
	}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("result[%d] = %s, want %s", i, v, expected[i])
		}
	}
}

func Test_createHandsUpCommandsInFuture_未来のイベントがない場合_空のsliceを返却する(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	tests := []struct {
		name string
		now time.Time
	}{
		{
			name: "nowが2023-10-29 09:00:01の場合",
			now: time.Date(2023, 10, 29, 9, 0, 1, 0, jst),
		},
		{
			name: "nowが2023-10-30の場合",
			now: time.Date(2023, 10, 30, 0, 0, 0, 0, jst),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createHandsUpCommandsInFuture(sampleSQInfo, tt.now)
			expected := []string{}
			if len(result) != 0 {
				t.Errorf("len(result) = %d, want %d", len(result), len(expected))
			}
		})
	}
}
