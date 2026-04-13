package service

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"
	repository "gin-investment-tracker/internal/repositories"
	"gin-investment-tracker/internal/util"

	"github.com/jackc/pgx/v5"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req *dto.CreateUserRequest) error {
	newHashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: newHashedPassword,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, util.NewBadRequestError("invalid email or password")
		}
		return nil, util.NewInternalError("failed to process login")
	}

	if !util.CheckPassword(user.PasswordHash, req.Password) {
		return nil, util.NewBadRequestError("invalid email or password")
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	return user, err
}

func (s *UserService) Delete(ctx context.Context, userId int64) error {
	err := s.repo.Delete(ctx, userId)
	return err
}
