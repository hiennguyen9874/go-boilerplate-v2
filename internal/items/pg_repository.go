package items

import (
	"context"

	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
)

type ItemPgRepository interface {
	Get(ctx context.Context, id uint) (*models.Item, error)
	GetMulti(ctx context.Context, offset, limit int) ([]*models.Item, error)
	Delete(ctx context.Context, id uint) (*models.Item, error)
	Update(ctx context.Context, id uint, obj_update *models.ItemUpdate) (*models.Item, error)
	GetMultiByOwnerId(ctx context.Context, ownerId uint, offset, limit int) ([]*models.Item, error)
	CreateWithOwner(ctx context.Context, ownerId uint, obj_create *models.ItemCreate) (*models.Item, error)
	DeleteWithoutGet(ctx context.Context, id uint) error
}
