package service

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"
)

type AuthServiceInterface interface {
	Create(ctx context.Context, req *dto.CreateUserRequest) error
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	GoogleLogin(ctx context.Context, userInfo *dto.GoogleUserInfo) (*dto.LoginResponse, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, userId int64) error
}
