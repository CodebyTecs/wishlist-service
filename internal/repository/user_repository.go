package repository

import (
	"context"
	"errors"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id string) (domain.User, error)
}

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user domain.User) error {
	const query = `
		INSERT INTO users (id, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.pool.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.CreatedAt)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrAlreadyExists
	}

	return err
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == nil {
		return user, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrNotFound
	}

	return domain.User{}, err
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	const query = `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == nil {
		return user, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrNotFound
	}

	return domain.User{}, err
}
