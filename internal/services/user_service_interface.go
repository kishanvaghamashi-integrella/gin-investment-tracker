package service

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"
)

type UserServiceInterface interface {
	Create(ctx context.Context, req *dto.CreateUserRequest) error
	Login(ctx context.Context, req *dto.LoginRequest) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, userId int64) error
}
