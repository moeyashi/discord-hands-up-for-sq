package repository

import "context"

type Repository interface {
	GetVersion(ctx context.Context) (*version, error)
	PutSQList(ctx context.Context, guildID string, sqList []string) error
	GetGuild(ctx context.Context, guildID string) (*Guild, error)
	GetSQMembers(ctx context.Context, guildID string, sqTitle string) ([]Member, error)
	PutSQMembers(ctx context.Context, guildID string, sqTitle string, members []Member) error
}

type MemberTypes int

const (
	MemberTypesParticipant MemberTypes = 1
	MemberTypesTemporary   MemberTypes = 2
)

type Member struct {
	UserID     string      `firestore:"userID"`
	UserName   string      `firestore:"userName"`
	MemberType MemberTypes `firestore:"memberType"`
}

type SQ struct {
	Title   string   `firestore:"title"`
	Members []Member `firestore:"members"`
}

type Guild struct {
	SQList []SQ `firestore:"sqList"`
}
