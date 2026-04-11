package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CodebyTecs/wishlist-service/internal/domain"
)

type WishlistItemRepository interface {
	Create(ctx context.Context, userID, wishlistID string, item domain.WishlistItem) (domain.WishlistItem, error)
	ListByWishlistAndUser(ctx context.Context, userID, wishlistID string) ([]domain.WishlistItem, error)
	GetByIDAndWishlistAndUser(ctx context.Context, userID, wishlistID, itemID string) (domain.WishlistItem, error)
	UpdateByIDAndWishlistAndUser(ctx context.Context, userID, wishlistID, itemID string, update domain.WishlistItemUpdate) (domain.WishlistItem, error)
	DeleteByIDAndWishlistAndUser(ctx context.Context, userID, wishlistID, itemID string) error
	ListByPublicToken(ctx context.Context, publicToken string) ([]domain.WishlistItem, error)
	ReserveByPublicTokenAndItemID(ctx context.Context, publicToken, itemID string) error
}

type PostgresWishlistItemRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresWishlistItemRepository(pool *pgxpool.Pool) *PostgresWishlistItemRepository {
	return &PostgresWishlistItemRepository{pool: pool}
}

func (r *PostgresWishlistItemRepository) Create(ctx context.Context, userID, wishlistID string, item domain.WishlistItem) (domain.WishlistItem, error) {
	const query = `
		INSERT INTO wishlist_items (id, wishlist_id, name, description, product_link, priority, is_reserved, reserved_at, created_at)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9
		WHERE EXISTS (
			SELECT 1
			FROM wishlists w
			WHERE w.id = $2 AND w.user_id = $10
		)
		RETURNING id, wishlist_id, name, description, product_link, priority, is_reserved, reserved_at, created_at
	`

	var saved domain.WishlistItem
	var description sql.NullString
	var productLink sql.NullString
	var reservedAt sql.NullTime
	err := r.pool.QueryRow(
		ctx,
		query,
		item.ID,
		wishlistID,
		item.Name,
		nullIfEmpty(item.Description),
		nullIfEmpty(item.ProductLink),
		item.Priority,
		item.IsReserved,
		nil,
		item.CreatedAt,
		userID,
	).Scan(
		&saved.ID,
		&saved.WishlistID,
		&saved.Name,
		&description,
		&productLink,
		&saved.Priority,
		&saved.IsReserved,
		&reservedAt,
		&saved.CreatedAt,
	)
	if err == nil {
		saved.Description = nullableString(description)
		saved.ProductLink = nullableString(productLink)
		saved.ReservedAt = nullableTime(reservedAt)
		return saved, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.WishlistItem{}, domain.ErrNotFound
	}

	return domain.WishlistItem{}, err
}

func (r *PostgresWishlistItemRepository) ListByWishlistAndUser(ctx context.Context, userID, wishlistID string) ([]domain.WishlistItem, error) {
	const query = `
		SELECT wi.id, wi.wishlist_id, wi.name, wi.description, wi.product_link, wi.priority, wi.is_reserved, wi.reserved_at, wi.created_at
		FROM wishlist_items wi
		INNER JOIN wishlists w ON w.id = wi.wishlist_id
		WHERE wi.wishlist_id = $1 AND w.user_id = $2
		ORDER BY wi.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, wishlistID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.WishlistItem, 0)
	for rows.Next() {
		var item domain.WishlistItem
		var description sql.NullString
		var productLink sql.NullString
		var reservedAt sql.NullTime
		if err := rows.Scan(
			&item.ID,
			&item.WishlistID,
			&item.Name,
			&description,
			&productLink,
			&item.Priority,
			&item.IsReserved,
			&reservedAt,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.Description = nullableString(description)
		item.ProductLink = nullableString(productLink)
		item.ReservedAt = nullableTime(reservedAt)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *PostgresWishlistItemRepository) GetByIDAndWishlistAndUser(ctx context.Context, userID, wishlistID, itemID string) (domain.WishlistItem, error) {
	const query = `
		SELECT wi.id, wi.wishlist_id, wi.name, wi.description, wi.product_link, wi.priority, wi.is_reserved, wi.reserved_at, wi.created_at
		FROM wishlist_items wi
		INNER JOIN wishlists w ON w.id = wi.wishlist_id
		WHERE wi.id = $1 AND wi.wishlist_id = $2 AND w.user_id = $3
		LIMIT 1
	`

	var item domain.WishlistItem
	var description sql.NullString
	var productLink sql.NullString
	var reservedAt sql.NullTime
	err := r.pool.QueryRow(ctx, query, itemID, wishlistID, userID).Scan(
		&item.ID,
		&item.WishlistID,
		&item.Name,
		&description,
		&productLink,
		&item.Priority,
		&item.IsReserved,
		&reservedAt,
		&item.CreatedAt,
	)
	if err == nil {
		item.Description = nullableString(description)
		item.ProductLink = nullableString(productLink)
		item.ReservedAt = nullableTime(reservedAt)
		return item, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.WishlistItem{}, domain.ErrNotFound
	}

	return domain.WishlistItem{}, err
}

func (r *PostgresWishlistItemRepository) UpdateByIDAndWishlistAndUser(ctx context.Context, userID, wishlistID, itemID string, update domain.WishlistItemUpdate) (domain.WishlistItem, error) {
	const query = `
		UPDATE wishlist_items wi
		SET
			name = CASE WHEN $4 THEN $5 ELSE wi.name END,
			description = CASE WHEN $6 THEN $7 ELSE wi.description END,
			product_link = CASE WHEN $8 THEN $9 ELSE wi.product_link END,
			priority = CASE WHEN $10 THEN $11 ELSE wi.priority END
		FROM wishlists w
		WHERE wi.id = $1
			AND wi.wishlist_id = $2
			AND w.id = wi.wishlist_id
			AND w.user_id = $3
		RETURNING wi.id, wi.wishlist_id, wi.name, wi.description, wi.product_link, wi.priority, wi.is_reserved, wi.reserved_at, wi.created_at
	`

	var item domain.WishlistItem
	var description sql.NullString
	var productLink sql.NullString
	var reservedAt sql.NullTime
	err := r.pool.QueryRow(
		ctx,
		query,
		itemID,
		wishlistID,
		userID,
		update.UpdateName,
		update.Name,
		update.UpdateDescription,
		nullIfEmpty(update.Description),
		update.UpdateProductLink,
		nullIfEmpty(update.ProductLink),
		update.UpdatePriority,
		update.Priority,
	).Scan(
		&item.ID,
		&item.WishlistID,
		&item.Name,
		&description,
		&productLink,
		&item.Priority,
		&item.IsReserved,
		&reservedAt,
		&item.CreatedAt,
	)
	if err == nil {
		item.Description = nullableString(description)
		item.ProductLink = nullableString(productLink)
		item.ReservedAt = nullableTime(reservedAt)
		return item, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.WishlistItem{}, domain.ErrNotFound
	}

	return domain.WishlistItem{}, err
}

func (r *PostgresWishlistItemRepository) DeleteByIDAndWishlistAndUser(ctx context.Context, userID, wishlistID, itemID string) error {
	const query = `
		DELETE FROM wishlist_items wi
		USING wishlists w
		WHERE wi.id = $1
			AND wi.wishlist_id = $2
			AND w.id = wi.wishlist_id
			AND w.user_id = $3
	`

	result, err := r.pool.Exec(ctx, query, itemID, wishlistID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *PostgresWishlistItemRepository) ListByPublicToken(ctx context.Context, publicToken string) ([]domain.WishlistItem, error) {
	const query = `
		SELECT wi.id, wi.wishlist_id, wi.name, wi.description, wi.product_link, wi.priority, wi.is_reserved, wi.reserved_at, wi.created_at
		FROM wishlist_items wi
		INNER JOIN wishlists w ON w.id = wi.wishlist_id
		WHERE w.public_token = $1
		ORDER BY wi.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, publicToken)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.WishlistItem, 0)
	for rows.Next() {
		var item domain.WishlistItem
		var description sql.NullString
		var productLink sql.NullString
		var reservedAt sql.NullTime
		if err := rows.Scan(
			&item.ID,
			&item.WishlistID,
			&item.Name,
			&description,
			&productLink,
			&item.Priority,
			&item.IsReserved,
			&reservedAt,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.Description = nullableString(description)
		item.ProductLink = nullableString(productLink)
		item.ReservedAt = nullableTime(reservedAt)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *PostgresWishlistItemRepository) ReserveByPublicTokenAndItemID(ctx context.Context, publicToken, itemID string) error {
	const reserveQuery = `
		UPDATE wishlist_items wi
		SET is_reserved = TRUE, reserved_at = NOW()
		FROM wishlists w
		WHERE wi.id = $1
			AND w.id = wi.wishlist_id
			AND w.public_token = $2
			AND wi.is_reserved = FALSE
	`

	result, err := r.pool.Exec(ctx, reserveQuery, itemID, publicToken)
	if err != nil {
		return err
	}
	if result.RowsAffected() > 0 {
		return nil
	}

	const checkQuery = `
		SELECT wi.is_reserved
		FROM wishlist_items wi
		INNER JOIN wishlists w ON w.id = wi.wishlist_id
		WHERE wi.id = $1 AND w.public_token = $2
		LIMIT 1
	`

	var isReserved bool
	err = r.pool.QueryRow(ctx, checkQuery, itemID, publicToken).Scan(&isReserved)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	if err != nil {
		return err
	}
	if isReserved {
		return domain.ErrAlreadyReserved
	}

	return domain.ErrConflict
}

func nullableString(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

func nullableTime(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}
