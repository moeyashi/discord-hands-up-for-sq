package repository

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/firestore"
)

type firestoreRepository struct {
	client *firestore.Client
}

func New(ctx context.Context) (Repository, error) {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &firestoreRepository{client: client}, nil
}

type version struct {
	Version    int64  `firestore:"version"`
}

func (r *firestoreRepository) GetVersion(ctx context.Context) (*version, error) {
	versions, err := r.client.Collection("v").OrderBy("version", firestore.Desc).Limit(1).Documents(ctx).GetAll()
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
