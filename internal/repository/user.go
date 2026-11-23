package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"tg-reminder/internal/domain"
)

type UserRepository struct {
	connectionsPool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{connectionsPool: pool}
}

func (u *UserRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := u.connectionsPool.Exec(ctx, `INSERT INTO users (tg_id, username, age, time_zone)
	 VALUES ($1, $2, $3, $4)`,
		user.TgID, user.Name, user.Age, user.TimeZone)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return err
}

func (u *UserRepository) GetUserByTgID(ctx context.Context, tgID int64) (*domain.User, error) {
	user := &domain.User{}
	err := u.connectionsPool.QueryRow(ctx, "select id, tg_id, username, age, time_zone from users where tg_id = $1;", tgID).Scan(
		&user.ID, &user.TgID, &user.Name, &user.Age, &user.TimeZone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by tg_id: %w", err)
	}
	return user, nil
}
