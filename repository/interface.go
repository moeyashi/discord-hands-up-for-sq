package repository

import "context"

type Repository interface {
	GetVersion(ctx context.Context) (*version, error)
	PutSQList(ctx context.Context, guildID string, sqList []string) error
	GetGuild(ctx context.Context, guildID string) (*Guild, error)
}

type SQ struct {
	Title string `firestore:"title"`
	Members []string `firestore:"members"`
}

type Guild struct {
	SQList []SQ `firestore:"sqList"`
}