package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/hiennguyen9874/go-boilerplate/config"
	"github.com/hiennguyen9874/go-boilerplate/internal/items"
	"github.com/hiennguyen9874/go-boilerplate/internal/models"
	"github.com/hiennguyen9874/go-boilerplate/internal/usecase"
	"github.com/hiennguyen9874/go-boilerplate/pkg/logger"
)

type itemUseCase struct {
	usecase.UseCase[models.Item]
	pgRepo items.ItemPgRepository
}

func CreateItemUseCaseI(
	pgRepo items.ItemPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) items.ItemUseCaseI {
	return &itemUseCase{
		UseCase: usecase.CreateUseCase[models.Item](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *itemUseCase) GetMultiByOwnerId(ctx context.Context, ownerId uuid.UUID, limit, offset int) ([]*models.Item, error) {
	return u.pgRepo.GetMultiByOwnerId(ctx, ownerId, limit, offset)
}

func (u *itemUseCase) CreateWithOwner(ctx context.Context, ownerId uuid.UUID, exp *models.Item) (*models.Item, error) {
	return u.pgRepo.CreateWithOwner(ctx, ownerId, exp)
}

func (u *itemUseCase) DeleteWithoutGet(ctx context.Context, id uuid.UUID) error {
	return u.pgRepo.DeleteWithoutGet(ctx, id)
}
