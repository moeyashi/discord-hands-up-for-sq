package usecase

import (
	"sort"
	"time"

	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

// NewSQList は新しいSQリストを既存のSQリストと結合します。
//
// 既存のSQリストからは、現在時刻よりも過去のSQを取り除きます。
// 新しいSQリストからは、現在時刻よりも未来のSQを取り出します。
// 新しいSQリストに既存のSQリストに含まれているSQがある場合、新しいSQリストのSQを優先します。
// 返却されるSQリストは、timestampの昇順でソートされています。
func NewSQList(
	existingSQList []repository.SQ,
	newSQList []repository.SQ,
	now time.Time,
) []repository.SQ {
	filteredExistingSQList := sqListInFuture(existingSQList, now)
	filteredNewSQList := sqListInFuture(newSQList, now)

	newSQMap := map[string]repository.SQ{}
	for _, sq := range filteredExistingSQList {
		newSQMap[sq.ID] = sq
	}
	for _, sq := range filteredNewSQList {
		newSQMap[sq.ID] = sq
	}

	newSQList = []repository.SQ{}
	for _, sq := range newSQMap {
		newSQList = append(newSQList, sq)
	}
	sort.Slice(newSQList, func(i, j int) bool {
		return newSQList[i].Timestamp.Before(newSQList[j].Timestamp)
	})

	return newSQList
}

func sqListInFuture(sqList []repository.SQ, now time.Time) []repository.SQ {
	nowUnix := now.Unix()
	featureSQList := []repository.SQ{}
	for _, sq := range sqList {
		timestampUnix := sq.Timestamp.Unix()
		if nowUnix <= timestampUnix {
			featureSQList = append(featureSQList, sq)
		}
	}
	return featureSQList
}
