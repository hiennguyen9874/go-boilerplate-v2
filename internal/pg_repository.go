package internal

import (
	"context"

	"github.com/google/uuid"
)

type PgRepository[M any] interface {
	Create(ctx context.Context, exp *M) (*M, error)
	Get(ctx context.Context, id uuid.UUID) (*M, error)
	GetMulti(ctx context.Context, limit, offset int) ([]*M, error)
	Delete(ctx context.Context, id uuid.UUID) (*M, error)
	Update(ctx context.Context, exp *M, values map[string]interface{}) (*M, error)
}
