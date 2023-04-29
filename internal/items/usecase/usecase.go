package usecase

import (
	"context"

	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/items"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/logger"
)

type itemUseCase struct {
	pgRepo items.ItemPgRepository
	cfg    *config.Config
	logger logger.Logger
}

func CreateItemUseCase(
	pgRepo items.ItemPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) items.ItemUseCase {
	return &itemUseCase{
		pgRepo: pgRepo,
		cfg:    cfg,
		logger: logger,
	}
}

func (u *itemUseCase) CreateWithOwner(ctx context.Context, ownerId uint, obj_create *models.ItemCreate) (*models.Item, error) {
	return u.pgRepo.CreateWithOwner(ctx, ownerId, obj_create)
}

func (u *itemUseCase) Get(ctx context.Context, id uint) (*models.Item, error) {
	return u.pgRepo.Get(ctx, id)
}

func (u *itemUseCase) GetMulti(ctx context.Context, offset, limit int) ([]*models.Item, error) {
	return u.pgRepo.GetMulti(ctx, offset, limit)
}

func (u *itemUseCase) Delete(ctx context.Context, id uint) (*models.Item, error) {
	return u.pgRepo.Delete(ctx, id)
}

func (u *itemUseCase) Update(ctx context.Context, id uint, obj_update *models.ItemUpdate) (*models.Item, error) {
	return u.pgRepo.Update(ctx, id, obj_update)
}

func (u *itemUseCase) GetMultiByOwnerId(ctx context.Context, ownerId uint, limit, offset int) ([]*models.Item, error) {
	return u.pgRepo.GetMultiByOwnerId(ctx, ownerId, limit, offset)
}

func (u *itemUseCase) DeleteWithoutGet(ctx context.Context, id uint) error {
	return u.pgRepo.DeleteWithoutGet(ctx, id)
}
