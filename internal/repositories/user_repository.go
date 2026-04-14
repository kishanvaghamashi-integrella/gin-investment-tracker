package repository

import (
	"context"
	"fmt"
	model "gin-investment-tracker/internal/models"
	"gin-investment-tracker/internal/util"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO USERS(name, email, password_hash)
		VALUES($1, $2, $3)
		RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, userId int64) error {
	query := `
		UPDATE users
		SET is_active=FALSE, updated_at=$2
		WHERE id=$1
	`

	res, err := r.db.Exec(ctx, query, userId, time.Now())
	if err != nil {
		slog.Error(err.Error())
		return util.NewInternalError("failed to delete user")
	}

	if res.RowsAffected() == 0 {
		return util.NewNotFoundError(fmt.Sprintf("no user found with id %d", userId))
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password_hash, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = TRUE
	`

	var user model.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, name, email, password_hash, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = TRUE
	`

	var user model.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, util.NewNotFoundError("User not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND is_active = TRUE)`

	var exists bool
	if err := r.db.QueryRow(ctx, query, id).Scan(&exists); err != nil {
		slog.Error("failed to check user existence", "error", err.Error())
		return false, util.NewInternalError("failed to check user existence")
	}

	return exists, nil
}
