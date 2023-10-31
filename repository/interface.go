package repository

import "context"

type Repository interface {
	GetVersion(ctx context.Context) (*version, error)
}