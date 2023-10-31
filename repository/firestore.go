package repository

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
)

type firestoreRepository struct {
	client *firestore.Client
}

func New(ctx context.Context) (*firestoreRepository, error) {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &firestoreRepository{client: client}, nil
}

func (r *firestoreRepository) GetVersion(ctx context.Context) (string, error) {
	versions, err := r.client.Collection("v").OrderBy("id", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	if err != nil || len(versions) == 0 {
		return "", err
	}
	return versions[0].Ref.ID, nil
}
