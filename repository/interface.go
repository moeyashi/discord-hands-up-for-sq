package repository

import (
	"context"
	"time"
)

type Repository interface {
	GetVersion(ctx context.Context) (*version, error)
	GetGuild(ctx context.Context, guildID string) (*Guild, error)
	GetSQList(ctx context.Context, guild *Guild) ([]SQ, error)
	PutSQList(ctx context.Context, guild *Guild, sqList []SQ) error
	GetMogiList(ctx context.Context, guild *Guild) ([]Mogi, error)
	GetMogi(ctx context.Context, guild *Guild, mogiTitle string) (*Mogi, error)
	AppendMogiList(ctx context.Context, guild *Guild, mogi Mogi) error
	DeleteMogi(ctx context.Context, guild *Guild, mogiTitle string) error
	GetSQMembers(ctx context.Context, guild *Guild, sqTitle string) ([]Member, error)
	PutSQMembers(ctx context.Context, guild *Guild, sqTitle string, members []Member) error
	GetMogiMembers(ctx context.Context, guild *Guild, mogiTitle string) ([]Member, error)
	PutMogiMembers(ctx context.Context, guild *Guild, mogiTitle string, members []Member) error
	PutResultsSpreadsheet(ctx context.Context, guild *Guild, spreadsheet string) error
}

type LoungeRepository interface {
	GetLoungeName(ctx context.Context, userID string) (*getLoungeNameResponse, error)
}

type getLoungeNameResponse struct {
	Name string `json:"name"`
}

type MemberTypes int

const (
	MemberTypesParticipant MemberTypes = 1
	MemberTypesTemporary   MemberTypes = 2
	MemberTypesSub         MemberTypes = 3
)

type Member struct {
	UserID     string      `firestore:"userID"`
	UserName   string      `firestore:"userName"`
	MemberType MemberTypes `firestore:"memberType"`
}

type SQ struct {
	ID        string    `firestore:"id"`
	Title     string    `firestore:"title"`
	Members   []Member  `firestore:"members"`
	Format    string    `firestore:"format"`
	Timestamp time.Time `firestore:"timestamp"`
}

type Mogi struct {
	Timestamp time.Time `firestore:"timestamp"`
	Members   []Member  `firestore:"members"`
}

type Guild struct {
	ID          string `firestore:"id"`
	SQList      []SQ   `firestore:"sqList"`
	Spreadsheet string `firestore:"spreadsheet"`
	MogiList    []Mogi `firestore:"mogiList"`
}

func (mogi Mogi) Title() string {
	return mogi.Timestamp.Format("01月02日 15時")
}

func MakeMogi(now time.Time, month, date, hour int64) *Mogi {
	year := nextYear(now, month, date)
	mogiTimestamp := time.Date(year, time.Month(month-1), int(date), int(hour), 0, 0, 0, time.Local)
	return &Mogi{
		Timestamp: mogiTimestamp,
	}
}

// 次にその月日を持つ年を返す
func nextYear(now time.Time, month, date int64) int {
	if month-1 < int64(now.Month()) || (month-1 == int64(now.Month()) && date < int64(now.Day())) {
		return now.Year() + 1
	}
	return now.Year()
}
