package usecase_test

import (
	"testing"
	"time"

	"github.com/moeyashi/discord-hands-up-for-sq/repository"
	"github.com/moeyashi/discord-hands-up-for-sq/usecase"
)

func Test_NewSQList_通常ケース_既存のSQがすべて終了した後に新しいSQ情報が来た場合_新しいSQ情報のみ返却される(t *testing.T) {
	now := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	existingSQList := []repository.SQ{
		{
			ID:        "1",
			Title:     "31日23:59 2v2",
			Format:    "2v2",
			Timestamp: time.Date(2020, 12, 31, 23, 59, 59, 0, time.UTC),
			Members: []repository.Member{
				{UserID: "1", UserName: "user1", MemberType: repository.MemberTypesParticipant},
			},
		},
	}
	newSQList := []repository.SQ{
		{
			ID:        "2",
			Title:     "1日00:00 2v2",
			Format:    "2v2",
			Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	actual := usecase.NewSQList(existingSQList, newSQList, now)

	if len(actual) != 1 {
		t.Fatalf("actual = %v, want %v", actual, 1)
	}
	assertSQ(t, actual[0], newSQList[0])
}

func Test_NewSQList_既存のSQがすべて終了する前に新しいSQ情報が来た場合_既存のSQの未終了のSQと新しいSQが返却される(t *testing.T) {
	now := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	existingSQList := []repository.SQ{
		{
			ID:        "1",
			Title:     "31日23:59 2v2",
			Format:    "2v2",
			Timestamp: time.Date(2020, 12, 31, 23, 59, 59, 0, time.UTC),
			Members: []repository.Member{
				{UserID: "1", UserName: "user1", MemberType: repository.MemberTypesParticipant},
			},
		},
		{
			ID:        "2",
			Title:     "1日00:00 2v2",
			Format:    "2v2",
			Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			Members: []repository.Member{
				{UserID: "2", UserName: "user2", MemberType: repository.MemberTypesParticipant},
			},
		},
	}
	newSQList := []repository.SQ{
		{
			ID:        "3",
			Title:     "2日00:00 2v2",
			Format:    "2v2",
			Timestamp: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	actual := usecase.NewSQList(existingSQList, newSQList, now)

	if len(actual) != 2 {
		t.Fatalf("actual = %v, want %v", actual, 2)
	}
	assertSQ(t, actual[0], existingSQList[1])
	assertSQ(t, actual[1], newSQList[0])
}

func Test_NewSQList_新しいSQではなくSQの修正が発生した場合_修正後で置き換わる(t *testing.T) {
	now := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	existingSQList := []repository.SQ{
		{
			ID:        "1",
			Title:     "2日00:00 2v2",
			Format:    "2v2",
			Timestamp: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
			Members: []repository.Member{
				{UserID: "1", UserName: "user1", MemberType: repository.MemberTypesParticipant},
			},
		},
		{
			ID:        "2",
			Title:     "2日10:00 3v3",
			Format:    "3v23",
			Timestamp: time.Date(2021, 1, 2, 10, 0, 0, 0, time.UTC),
			Members: []repository.Member{
				{UserID: "2", UserName: "user2", MemberType: repository.MemberTypesParticipant},
			},
		},
	}
	newSQList := []repository.SQ{
		{
			ID:        "1",
			Title:     "2日01:00 2v2",
			Format:    "2v2",
			Timestamp: time.Date(2021, 1, 2, 1, 0, 0, 0, time.UTC),
		},
	}

	actual := usecase.NewSQList(existingSQList, newSQList, now)

	if len(actual) != 2 {
		t.Fatalf("actual = %v, want %v", actual, 2)
	}
	assertSQ(t, actual[0], repository.SQ{
		ID:        newSQList[0].ID,
		Title:     newSQList[0].Title,
		Format:    newSQList[0].Format,
		Timestamp: newSQList[0].Timestamp,
		Members:   existingSQList[0].Members,
	})
	assertSQ(t, actual[1], existingSQList[1])
}

func assertSQ(t *testing.T, actual repository.SQ, expected repository.SQ) {
	if actual.ID == expected.ID &&
		actual.Title == expected.Title &&
		actual.Format == expected.Format &&
		actual.Timestamp.Equal(expected.Timestamp) &&
		equalsMembers(actual.Members, expected.Members) {
		return
	}
	t.Errorf("actual = %v, want %v", actual, expected)
}

func equalsMembers(actual []repository.Member, expected []repository.Member) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i := range actual {
		if actual[i].UserID != expected[i].UserID ||
			actual[i].UserName != expected[i].UserName ||
			actual[i].MemberType != expected[i].MemberType {
			return false
		}
	}
	return true
}
