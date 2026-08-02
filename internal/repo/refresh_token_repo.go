package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/murashi19/koda-b8-backend/internal/models"
)

type RefreshTokenRepo struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepo(db *pgxpool.Pool) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		db: db,
	}
}

func (r *RefreshTokenRepo) Create(ctx context.Context, data *models.RefreshToken) (*models.RefreshToken, error) {
	query := `
		INSERT INTO refresh_tokens (
			user_id,
			token_hash,
			expires_at,
			user_agent,
			ip_address
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`

	return oneRow[models.RefreshToken](
		ctx,
		r.db,
		query,
		data.UserID,
		data.TokenHash,
		data.ExpiresAt,
		data.UserAgent,
		data.IPAddress,
	)
}

func (r *RefreshTokenRepo) GetByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	query := `
		SELECT *
		FROM refresh_tokens
		WHERE token_hash = $1
		LIMIT 1
	`

	return oneRow[models.RefreshToken](
		ctx,
		r.db,
		query,
		tokenHash,
	)
}

func (r *RefreshTokenRepo) Revoke(ctx context.Context, id int64) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE id = $1
	`

	rows, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *RefreshTokenRepo) DeleteExpired(ctx context.Context) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < NOW()
	`

	_, err := r.db.Exec(ctx, query)

	return err
}
