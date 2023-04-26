package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/items"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/repository"
	"gorm.io/gorm"
)

type ItemPgRepo struct {
	repository.PgRepo[models.Item]
}

func CreateItemPgRepository(db *gorm.DB) items.ItemPgRepository {
	return &ItemPgRepo{
		PgRepo: repository.CreatePgRepo[models.Item](db),
	}
}

func (r *ItemPgRepo) GetMultiByOwnerId(ctx context.Context, ownerId uuid.UUID, limit, offset int) ([]*models.Item, error) {
	var objs []*models.Item
	r.DB.WithContext(ctx).Where("owner_id = ?", ownerId.String()).Limit(limit).Offset(offset).Find(&objs)
	return objs, nil
}

func (r *ItemPgRepo) CreateWithOwner(ctx context.Context, ownerId uuid.UUID, exp *models.Item) (*models.Item, error) {
	exp.OwnerId = ownerId
	if result := r.DB.WithContext(ctx).Create(exp); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *ItemPgRepo) DeleteWithoutGet(ctx context.Context, id uuid.UUID) error {
	if result := r.DB.WithContext(ctx).Delete(&models.Item{}, "id = ?", id.String()); result.Error != nil {
		return result.Error
	}
	return nil
}
