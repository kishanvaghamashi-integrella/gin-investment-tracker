package repository

import (
	"context"
	model "gin-investment-tracker/internal/models"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, userId int64) error
	ExistsByID(ctx context.Context, id int64) (bool, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
}
