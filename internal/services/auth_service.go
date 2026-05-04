package service

import (
	"context"
	dto "gin-investment-tracker/internal/dtos"
	model "gin-investment-tracker/internal/models"
	repository "gin-investment-tracker/internal/repositories"
	"gin-investment-tracker/internal/util"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Create(ctx context.Context, req *dto.CreateUserRequest) error {
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

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
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

	token, err := util.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, util.NewInternalError("failed to generate token")
	}

	return &dto.LoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}, nil
}

func (s *AuthService) GoogleLogin(ctx context.Context, userInfo *dto.GoogleUserInfo) (*dto.LoginResponse, error) {
	user, err := s.repo.GetByGoogleID(ctx, userInfo.Sub)
	if err != nil {
		slog.Error("got error in auth service", "error", err.Error())
		return nil, util.NewInternalError("failed to fetch user")
	}

	if user == nil {
		user = &model.User{
			Name:     userInfo.Name,
			Email:    userInfo.Email,
			GoogleID: &userInfo.Sub,
		}
		if err := s.repo.CreateGoogleUser(ctx, user); err != nil {
			return nil, util.NewInternalError("failed to create user")
		}
	}

	token, err := util.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, util.NewInternalError("failed to generate token")
	}

	return &dto.LoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}, nil
}

func (s *AuthService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AuthService) Delete(ctx context.Context, userId int64) error {
	err := s.repo.Delete(ctx, userId)
	return err
}
