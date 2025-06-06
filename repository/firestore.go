package repository

import (
	"context"
	"errors"
	"log"
	"os"

	"cloud.google.com/go/firestore"
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
	client, err := firestore.NewClient(ctx, projectID, GetGoogleDefaultCredentialClientOption())
	if err != nil {
		return nil, err
	}
	return &firestoreRepository{client: client}, nil
}

type version struct {
	Version int64 `firestore:"version"`
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

func (r *firestoreRepository) GetGuild(ctx context.Context, guildID string) (*Guild, error) {
	guild, err := r.getGuildOrCreate(ctx, guildID)
	if err != nil {
		return nil, err
	}
	return guild, nil
}

func (r *firestoreRepository) PutGuildName(ctx context.Context, guildID string, name string) (*Guild, error) {
	guild, err := r.getGuildOrCreate(ctx, guildID)
	if err != nil {
		return nil, err
	}
	guild.Name = name
	_, err = r.getGuildDocRef(guildID).Set(ctx, guild)
	if err != nil {
		return nil, err
	}
	return guild, nil
}

func (r *firestoreRepository) GetSQList(ctx context.Context, guild *Guild) ([]SQ, error) {
	return guild.SQList, nil
}

func (r *firestoreRepository) PutSQList(ctx context.Context, guild *Guild, sqList []SQ) error {
	guild.SQList = sqList
	_, err := r.getGuildDocRef(guild.ID).Set(ctx, guild)
	return err
}

func (r *firestoreRepository) GetMogiList(ctx context.Context, guild *Guild) ([]Mogi, error) {
	return guild.MogiList, nil
}

func (r *firestoreRepository) GetMogi(ctx context.Context, guild *Guild, mogiTitle string) (*Mogi, error) {
	for _, mogi := range guild.MogiList {
		if mogi.Title() == mogiTitle {
			return &mogi, nil
		}
	}
	return nil, errors.New("not found")
}

func (r *firestoreRepository) AppendMogiList(ctx context.Context, guild *Guild, mogi Mogi) error {
	// すでに存在するかチェック
	for _, m := range guild.MogiList {
		if m.Title() == mogi.Title() {
			return nil
		}
	}

	guild.MogiList = append(guild.MogiList, mogi)
	_, err := r.getGuildDocRef(guild.ID).Set(ctx, guild)
	return err
}

func (r *firestoreRepository) DeleteMogi(ctx context.Context, guild *Guild, mogiTitle string) error {
	mogiList := []Mogi{}
	for _, mogi := range guild.MogiList {
		if mogi.Title() != mogiTitle {
			mogiList = append(mogiList, mogi)
		}
	}
	guild.MogiList = mogiList
	_, err := r.getGuildDocRef(guild.ID).Set(ctx, guild)
	return err
}

func (r *firestoreRepository) GetSQMembers(ctx context.Context, guild *Guild, sqTitle string) ([]Member, error) {
	for _, sq := range guild.SQList {
		if sq.Title == sqTitle {
			return sq.Members, nil
		}
	}
	return nil, errors.New("not found")
}

func (r *firestoreRepository) PutSQMembers(ctx context.Context, guild *Guild, sqTitle string, members []Member) error {
	for i, sq := range guild.SQList {
		if sq.Title == sqTitle {
			guild.SQList[i].Members = members
			_, err := r.getGuildDocRef(guild.ID).Set(ctx, guild)
			return err
		}
	}
	return errors.New("not found")
}

func (r *firestoreRepository) GetMogiMembers(ctx context.Context, guild *Guild, mogiTitle string) ([]Member, error) {
	for _, mogi := range guild.MogiList {
		if mogi.Title() == mogiTitle {
			return mogi.Members, nil
		}
	}
	return nil, errors.New("not found")
}

func (r *firestoreRepository) PutMogiMembers(ctx context.Context, guild *Guild, mogiTitle string, members []Member) error {
	for i, mogi := range guild.MogiList {
		if mogi.Title() == mogiTitle {
			guild.MogiList[i].Members = members
			_, err := r.getGuildDocRef(guild.ID).Set(ctx, guild)
			return err
		}
	}
	return errors.New("not found")
}

func (r *firestoreRepository) PutResultsSpreadsheet(ctx context.Context, guild *Guild, spreadsheet string) error {
	guild.Spreadsheet = spreadsheet
	_, err := r.getGuildDocRef(guild.ID).Set(ctx, guild)
	return err
}

func (r *firestoreRepository) getGuildDocRef(guildID string) *firestore.DocumentRef {
	return r.client.Collection(versionCollection).Doc(firestoreVersion).Collection(guildCollection).Doc(guildID)
}

func (r *firestoreRepository) getGuildOrCreate(ctx context.Context, guildID string) (*Guild, error) {
	ref := r.getGuildDocRef(guildID)
	guild, err := ref.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Println("create new guild")
			newGuild := Guild{ID: guildID}
			_, err := ref.Create(ctx, newGuild)
			if err != nil {
				if status.Code(err) == codes.AlreadyExists {
					log.Println("guild already exists")
				} else {
					return nil, err
				}
			}
			return &newGuild, nil
		} else {
			return nil, err
		}
	}

	existsGuild := &Guild{}
	if err := guild.DataTo(existsGuild); err != nil {
		return nil, err
	}

	return existsGuild, nil
}
