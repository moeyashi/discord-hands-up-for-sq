package repository

import "context"

type Repository interface {
	GetVersion(ctx context.Context) (*version, error)
	GetGuild(ctx context.Context, guildID string) (*Guild, error)
	GetSQList(ctx context.Context, guild *Guild) ([]SQ, error)
	PutSQList(ctx context.Context, guild *Guild, sqList []string) error
	GetSQMembers(ctx context.Context, guild *Guild, sqTitle string) ([]Member, error)
	PutSQMembers(ctx context.Context, guild *Guild, sqTitle string, members []Member) error
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
	Title   string   `firestore:"title"`
	Members []Member `firestore:"members"`
}

type Guild struct {
	ID     string `firestore:"id"`
	SQList []SQ   `firestore:"sqList"`
}
