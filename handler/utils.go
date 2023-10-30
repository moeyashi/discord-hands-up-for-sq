package handler

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

func createHandsUpCommandsForBetweenNowAndTomorrow(sqInfo string, now time.Time) []string {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Println(err)
		return []string{}
	}
	dayAfterTomorrow := now.AddDate(0, 0, 2)
	dayAfterTomorrow = time.Date(dayAfterTomorrow.Year(), dayAfterTomorrow.Month(), dayAfterTomorrow.Day(), 0, 0, 0, 0, dayAfterTomorrow.Location())
	nowUnix := now.Unix()
	dayAfterTomorrowUnix := dayAfterTomorrow.Unix()
	re := regexp.MustCompile("`#(\\d+)` \\*\\*(\\dv\\d):\\*\\* <t:(\\d+):F>")
	results := re.FindAllStringSubmatch(sqInfo, -1)
	commands := []string{}
	for _, subMatches := range results {
		timestamp, err := strconv.ParseInt(subMatches[3], 10, 64)
		if err != nil {
			log.Println(err)
			return []string{}
		}
		if nowUnix <= timestamp && timestamp < dayAfterTomorrowUnix {
			hourContent := time.Unix(timestamp, 0).In(jst).Format("2日15:04")
			mogiFormat := subMatches[2]
			command := fmt.Sprintf("/hands-up set hour:%s %s number:12", hourContent, mogiFormat)
			commands = append(commands, command)
		}
	}
	return commands
}
