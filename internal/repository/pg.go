package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hiennguyen9874/go-boilerplate/internal"
	"gorm.io/gorm"
)

type PgRepo[M any] struct {
	DB *gorm.DB
}

func CreatePgRepo[M any](db *gorm.DB) PgRepo[M] {
	return PgRepo[M]{DB: db}
}

func CreatePgRepository[M any](db *gorm.DB) internal.PgRepository[M] {
	return &PgRepo[M]{DB: db}
}

func (r *PgRepo[M]) Get(ctx context.Context, id uuid.UUID) (*M, error) {
	var obj *M
	if result := r.DB.WithContext(ctx).First(&obj, "id = ?", id.String()); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *PgRepo[M]) GetMulti(ctx context.Context, limit, offset int) ([]*M, error) {
	var objs []*M
	r.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&objs)
	return objs, nil
}

func (r *PgRepo[M]) Create(ctx context.Context, exp *M) (*M, error) {
	if result := r.DB.WithContext(ctx).Create(exp); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *PgRepo[M]) Delete(ctx context.Context, id uuid.UUID) (*M, error) {
	obj, err := r.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	if result := r.DB.WithContext(ctx).Delete(&obj, "id = ?", id.String()); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *PgRepo[M]) Update(ctx context.Context, exp *M, values map[string]interface{}) (*M, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Updates(values); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}
