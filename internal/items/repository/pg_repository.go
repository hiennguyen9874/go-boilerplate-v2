package repository

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/hiennguyen9874/go-boilerplate-v2/ent"
	"github.com/hiennguyen9874/go-boilerplate-v2/ent/item"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/items"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
)

type ItemPgRepo struct {
	client *ent.Client
}

func CreateItemPgRepository(client *ent.Client) items.ItemPgRepository {
	return &ItemPgRepo{client: client}
}

func (r *ItemPgRepo) mapModel(db_obj *ent.Item) *models.Item {
	return &models.Item{
		Id:      db_obj.ID,
		Title:   db_obj.Title,
		OwnerId: db_obj.OwnerID,
	}
}

func (r *ItemPgRepo) mapModels(db_objs []*ent.Item) []*models.Item {
	objs := make([]*models.Item, len(db_objs))
	for i, db_obj := range db_objs {
		objs[i] = r.mapModel(db_obj)
	}
	return objs
}

func (r *ItemPgRepo) Get(ctx context.Context, id uint) (*models.Item, error) {
	db_obj, err := r.client.Item.Query().
		Where(item.ID(id)).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return r.mapModel(db_obj), nil
}

func (r *ItemPgRepo) GetMulti(ctx context.Context, offset, limit int) ([]*models.Item, error) {
	db_objs, err := r.client.Item.Query().
		Order(item.ByID(sql.OrderAsc())).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModels(db_objs), nil
}

func (r *ItemPgRepo) Delete(ctx context.Context, id uint) (*models.Item, error) {
	db_obj, err := r.client.Item.Query().
		Where(item.ID(id)).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	err = r.client.Item.DeleteOne(db_obj).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.mapModel(db_obj), nil
}

func (r *ItemPgRepo) Update(ctx context.Context, id uint, obj_update *models.ItemUpdate) (*models.Item, error) {
	query := r.client.Item.UpdateOneID(id)
	if obj_update.Title != nil {
		query = query.SetTitle(*obj_update.Title)
	}
	if obj_update.Description != nil {
		query = query.SetDescription(*obj_update.Description)
	}
	db_obj, err := query.Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}

func (r *ItemPgRepo) GetMultiByOwnerId(ctx context.Context, ownerId uint, offset, limit int) ([]*models.Item, error) {
	db_objs, err := r.client.Item.Query().
		Where(item.OwnerID(ownerId)).
		Order(item.ByID(sql.OrderAsc())).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModels(db_objs), nil

}

func (r *ItemPgRepo) CreateWithOwner(ctx context.Context, ownerId uint, obj_create *models.ItemCreate) (*models.Item, error) {
	db_obj, err := r.client.Item.
		Create().
		SetTitle(obj_create.Title).
		SetDescription(obj_create.Description).
		SetOwnerID(ownerId).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}

func (r *ItemPgRepo) DeleteWithoutGet(ctx context.Context, id uint) error {
	return r.client.Item.DeleteOneID(id).Exec(ctx)
}
