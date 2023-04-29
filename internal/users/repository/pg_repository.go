package repository

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/hiennguyen9874/go-boilerplate-v2/ent"
	"github.com/hiennguyen9874/go-boilerplate-v2/ent/user"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/users"
)

type UserPgRepo struct {
	client *ent.Client
}

func CreateUserPgRepository(client *ent.Client) users.UserPgRepository {
	return &UserPgRepo{client: client}
}

func (r *UserPgRepo) mapModel(db_obj *ent.User) *models.User {
	return &models.User{
		Id:                 db_obj.ID,
		CreateTime:         db_obj.CreateTime,
		UpdateTime:         db_obj.UpdateTime,
		Name:               db_obj.Name,
		Email:              db_obj.Email,
		Password:           db_obj.Password,
		IsActive:           db_obj.IsActive,
		IsSuperUser:        db_obj.IsSuperUser,
		Verified:           db_obj.Verified,
		VerificationCode:   db_obj.VerificationCode,
		PasswordResetToken: db_obj.PasswordResetToken,
		PasswordResetAt:    db_obj.PasswordResetAt,
	}
}

func (r *UserPgRepo) mapModels(db_objs []*ent.User) []*models.User {
	objs := make([]*models.User, len(db_objs))
	for i, db_obj := range db_objs {
		objs[i] = r.mapModel(db_obj)
	}
	return objs
}

func (r *UserPgRepo) Get(ctx context.Context, id uint) (*models.User, error) {
	db_obj, err := r.client.User.Query().
		Where(user.ID(id)).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	return r.mapModel(db_obj), nil
}

func (r *UserPgRepo) GetMulti(ctx context.Context, limit, offset int) ([]*models.User, error) {
	db_objs, err := r.client.User.Query().
		Order(user.ByID(sql.OrderAsc())).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModels(db_objs), nil
}

func (r *UserPgRepo) Create(ctx context.Context, obj_in *models.UserCreate) (*models.User, error) {
	db_obj, err := r.client.User.
		Create().
		SetName(obj_in.Name).
		SetEmail(obj_in.Email).
		SetPassword(obj_in.Password).
		SetNillableIsActive(obj_in.IsActive).
		SetNillableIsSuperUser(obj_in.IsSuperUser).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}

func (r *UserPgRepo) Delete(ctx context.Context, id uint) (*models.User, error) {
	db_obj, err := r.client.User.Query().
		Where(user.ID(id)).
		Only(ctx)

	if err != nil {
		return nil, err
	}

	err = r.client.User.DeleteOne(db_obj).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.mapModel(db_obj), nil
}

func (r *UserPgRepo) Update(ctx context.Context, id uint, obj_update *models.UserUpdate) (*models.User, error) {
	query := r.client.User.UpdateOneID(id)
	if obj_update.Name != nil {
		query = query.SetName(*obj_update.Name)
	}
	if obj_update.Email != nil {
		query = query.SetEmail(*obj_update.Email)
	}
	if obj_update.Password != nil {
		query = query.SetPassword(*obj_update.Password)
	}
	db_obj, err := query.
		SetNillableIsActive(obj_update.IsActive).
		SetNillableIsSuperUser(obj_update.IsSuperUser).
		SetNillableVerified(obj_update.Verified).
		SetNillableVerificationCode(obj_update.VerificationCode).
		SetNillablePasswordResetToken(obj_update.PasswordResetToken).
		SetNillablePasswordResetAt(obj_update.PasswordResetAt).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}

func (r *UserPgRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	db_obj, err := r.client.User.Query().
		Where(user.Email(email)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}

func (r *UserPgRepo) UpdatePassword(ctx context.Context, id uint, newPassword string) (*models.User, error) {
	db_obj, err := r.client.User.UpdateOneID(id).
		SetPassword(newPassword).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}

func (r *UserPgRepo) UpdateVerificationCode(ctx context.Context, id uint, newVerificationCode string) (*models.User, error) {
	db_obj, err := r.client.User.UpdateOneID(id).
		SetVerificationCode(newVerificationCode).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}

func (r *UserPgRepo) UpdateVerification(
	ctx context.Context,
	id uint,
	newVerificationCode string,
	newVerified bool,
) (*models.User, error) {
	db_obj, err := r.client.User.UpdateOneID(id).
		SetVerified(newVerified).
		SetVerificationCode(newVerificationCode).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}

func (r *UserPgRepo) GetByVerificationCode(
	ctx context.Context,
	verificationCode string,
) (*models.User, error) {
	db_obj, err := r.client.User.Query().
		Where(user.VerificationCode(verificationCode)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}

func (r *UserPgRepo) UpdatePasswordReset(
	ctx context.Context,
	id uint,
	passwordResetToken string,
	passwordResetAt time.Time,
) (*models.User, error) {
	db_obj, err := r.client.User.UpdateOneID(id).
		SetPasswordResetToken(passwordResetToken).
		SetPasswordResetAt(passwordResetAt).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}

func (r *UserPgRepo) GetByResetToken(ctx context.Context, resetToken string) (*models.User, error) {
	db_obj, err := r.client.User.Query().
		Where(user.PasswordResetToken(resetToken)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}

func (r *UserPgRepo) GetByResetTokenResetAt(
	ctx context.Context,
	resetToken string,
	resetAt time.Time,
) (*models.User, error) {
	db_obj, err := r.client.User.Query().
		Where(
			user.PasswordResetToken(resetToken),
			user.PasswordResetAtGTE(resetAt),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}

func (r *UserPgRepo) UpdatePasswordResetToken(
	ctx context.Context,
	id uint,
	newPassword string,
	resetToken string,
) (*models.User, error) {
	db_obj, err := r.client.User.UpdateOneID(id).
		SetPassword(newPassword).
		SetPasswordResetToken(resetToken).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapModel(db_obj), nil
}
