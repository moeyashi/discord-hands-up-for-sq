package repository

import (
	"context"
	"errors"
	"log"
	"os"
	"slices"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const firestoreVersion = "1"

const versionCollection = "v"
const guildCollection = "guilds"

type firestoreRepository struct {
	client *firestore.Client
}

func New(ctx context.Context) (Repository, error) {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	var client *firestore.Client
	credentialJSON := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	var err error
	if credentialJSON != "" {
		client, err = firestore.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(credentialJSON)))
	} else {
		client, err = firestore.NewClient(ctx, projectID)
	}
	if err != nil {
		return nil, err
	}
	return &firestoreRepository{client: client}, nil
}

type version struct {
	Version    int64  `firestore:"version"`
}

func (r *firestoreRepository) GetVersion(ctx context.Context) (*version, error) {
	versions, err := r.client.Collection(versionCollection).OrderBy("version", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, errors.New("version not found")
	}
	var v version
	err = versions[0].DataTo(&v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *firestoreRepository) PutSQList(ctx context.Context, guildID string, sqList []string) error {
	guild, err := r.getGuildOrCreate(ctx, guildID)
	if err != nil {
		return err
	}
	log.Println("guild.SQList:", guild.SQList)

	addedSQTitle := []string{}
	newSQList := []SQ{}
	// すでにfirestoreにあるものはそのまま残す
	for _, sq := range guild.SQList {
		log.Println("sq.Title:", sq.Title)
		if slices.Contains(sqList, sq.Title) {
			newSQList = append(newSQList, sq)
			addedSQTitle = append(addedSQTitle, sq.Title)
		}
	}
	log.Println("addedSQTitle:", addedSQTitle)
	// 新規追加
	for _, title := range sqList {
		if !slices.Contains(addedSQTitle, title) {
			log.Println("title:", title)
			newSQList = append(newSQList, SQ{Title: title})
		}
	}
	
	guild.SQList = newSQList
	_, err = r.getGuildDocRef(guildID).Set(ctx, guild)
	return err
}

func (r *firestoreRepository) GetGuild(ctx context.Context, guildID string) (*Guild, error) {
	guild, err := r.getGuildOrCreate(ctx, guildID)
	if err != nil {
		return nil, err
	}
	return guild, nil
}

func (r *firestoreRepository) getGuildDocRef(guildID string) *firestore.DocumentRef {
	return r.client.Collection(versionCollection).Doc(firestoreVersion).Collection(guildCollection).Doc(guildID)
}

func (r *firestoreRepository) getGuildOrCreate(ctx context.Context, guildID string) (*Guild, error) {
	ref := r.getGuildDocRef(guildID)
	guild, err := ref.Get(ctx)
	if err != nil {
		log.Println("create new guild")
		if status.Code(err) == codes.NotFound {
			newGuild := Guild{}
			_, err := ref.Create(ctx, newGuild)
			if err != nil {
				return nil, err
			}
			return &newGuild, nil
		} else {
			return nil, err
		}
	}

	existsGuild := &Guild{}
	if err := guild.DataTo(existsGuild); err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("existsGuild:", existsGuild)

	return existsGuild, nil
}