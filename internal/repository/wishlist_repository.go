package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
)

type WishlistRepository interface {
	Create(ctx context.Context, wishlist domain.Wishlist) (domain.Wishlist, error)
	ListByUserID(ctx context.Context, userID string) ([]domain.Wishlist, error)
	GetByIDAndUserID(ctx context.Context, wishlistID, userID string) (domain.Wishlist, error)
	GetByPublicToken(ctx context.Context, publicToken string) (domain.Wishlist, error)
	UpdateByIDAndUserID(ctx context.Context, wishlistID, userID string, update domain.WishlistUpdate) (domain.Wishlist, error)
	DeleteByIDAndUserID(ctx context.Context, wishlistID, userID string) error
}

type PostgresWishlistRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresWishlistRepository(pool *pgxpool.Pool) *PostgresWishlistRepository {
	return &PostgresWishlistRepository{pool: pool}
}

func (r *PostgresWishlistRepository) Create(ctx context.Context, wishlist domain.Wishlist) (domain.Wishlist, error) {
	const query = `
		INSERT INTO wishlists (id, user_id, name, description, event_date, public_token, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, name, description, event_date, public_token, created_at
	`

	var saved domain.Wishlist
	var description sql.NullString
	err := r.pool.QueryRow(
		ctx,
		query,
		wishlist.ID,
		wishlist.UserID,
		wishlist.Name,
		nullIfEmpty(wishlist.Description),
		wishlist.EventDate,
		wishlist.PublicToken,
		wishlist.CreatedAt,
	).Scan(
		&saved.ID,
		&saved.UserID,
		&saved.Name,
		&description,
		&saved.EventDate,
		&saved.PublicToken,
		&saved.CreatedAt,
	)
	if err == nil {
		saved.Description = descriptionToString(description)
		return saved, nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.Wishlist{}, domain.ErrConflict
	}

	return domain.Wishlist{}, err
}

func (r *PostgresWishlistRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Wishlist, error) {
	const query = `
		SELECT id, user_id, name, description, event_date, public_token, created_at
		FROM wishlists
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wishlists := make([]domain.Wishlist, 0)
	for rows.Next() {
		var w domain.Wishlist
		var description sql.NullString
		if err := rows.Scan(&w.ID, &w.UserID, &w.Name, &description, &w.EventDate, &w.PublicToken, &w.CreatedAt); err != nil {
			return nil, err
		}
		w.Description = descriptionToString(description)
		wishlists = append(wishlists, w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wishlists, nil
}

func (r *PostgresWishlistRepository) GetByIDAndUserID(ctx context.Context, wishlistID, userID string) (domain.Wishlist, error) {
	const query = `
		SELECT id, user_id, name, description, event_date, public_token, created_at
		FROM wishlists
		WHERE id = $1 AND user_id = $2
		LIMIT 1
	`

	var w domain.Wishlist
	var description sql.NullString
	err := r.pool.QueryRow(ctx, query, wishlistID, userID).Scan(
		&w.ID,
		&w.UserID,
		&w.Name,
		&description,
		&w.EventDate,
		&w.PublicToken,
		&w.CreatedAt,
	)
	if err == nil {
		w.Description = descriptionToString(description)
		return w, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Wishlist{}, domain.ErrNotFound
	}

	return domain.Wishlist{}, err
}

func (r *PostgresWishlistRepository) UpdateByIDAndUserID(ctx context.Context, wishlistID, userID string, update domain.WishlistUpdate) (domain.Wishlist, error) {
	const query = `
		UPDATE wishlists
		SET
			name = CASE WHEN $3 THEN $4 ELSE name END,
			description = CASE WHEN $5 THEN $6 ELSE description END,
			event_date = CASE WHEN $7 THEN $8 ELSE event_date END
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, name, description, event_date, public_token, created_at
	`

	var updated domain.Wishlist
	var description sql.NullString
	err := r.pool.QueryRow(
		ctx,
		query,
		wishlistID,
		userID,
		update.UpdateName,
		update.Name,
		update.UpdateDescription,
		nullIfEmpty(update.Description),
		update.UpdateEventDate,
		update.EventDate,
	).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.Name,
		&description,
		&updated.EventDate,
		&updated.PublicToken,
		&updated.CreatedAt,
	)
	if err == nil {
		updated.Description = descriptionToString(description)
		return updated, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Wishlist{}, domain.ErrNotFound
	}

	return domain.Wishlist{}, err
}

func (r *PostgresWishlistRepository) GetByPublicToken(ctx context.Context, publicToken string) (domain.Wishlist, error) {
	const query = `
		SELECT id, user_id, name, description, event_date, public_token, created_at
		FROM wishlists
		WHERE public_token = $1
		LIMIT 1
	`

	var w domain.Wishlist
	var description sql.NullString
	err := r.pool.QueryRow(ctx, query, publicToken).Scan(
		&w.ID,
		&w.UserID,
		&w.Name,
		&description,
		&w.EventDate,
		&w.PublicToken,
		&w.CreatedAt,
	)
	if err == nil {
		w.Description = descriptionToString(description)
		return w, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Wishlist{}, domain.ErrNotFound
	}

	return domain.Wishlist{}, err
}

func (r *PostgresWishlistRepository) DeleteByIDAndUserID(ctx context.Context, wishlistID, userID string) error {
	const query = `
		DELETE FROM wishlists
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.pool.Exec(ctx, query, wishlistID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func descriptionToString(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
