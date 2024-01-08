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
	GetSQMembers(ctx context.Context, guild *Guild, sqTitle string) ([]Member, error)
	PutSQMembers(ctx context.Context, guild *Guild, sqTitle string, members []Member) error
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

type Guild struct {
	ID     string `firestore:"id"`
	SQList []SQ   `firestore:"sqList"`
}
