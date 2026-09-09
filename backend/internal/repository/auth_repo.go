package repository

import (
	"context"
	"time"

	"github.com/Mr-Rafael/bucktracker-api/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthRepo struct {
	queries *db.Queries
}

func NewAuthRepo(queries *db.Queries) *AuthRepo {
	return &AuthRepo{queries: queries}
}

func (r *AuthRepo) CreateRefreshToken(ctx context.Context, userID pgtype.UUID, tokenHash string, expDate time.Time) (db.RefreshToken, error) {
	params := ToRefreshTokenCreateParams(userID, tokenHash, expDate)
	return r.queries.CreateRefreshToken(ctx, params)
}

func (r *AuthRepo) GetTokenByHash(ctx context.Context, hash string) (db.RefreshToken, error) {
	return r.queries.GetTokenByHash(ctx, hash)
}

func (r *AuthRepo) RevokeTokenByUserID(ctx context.Context, id pgtype.UUID) error {
	return r.queries.RevokeTokenByUserID(ctx, id)
}

func ToRefreshTokenCreateParams(user pgtype.UUID, tokenHash string, expDate time.Time) db.CreateRefreshTokenParams {
	return db.CreateRefreshTokenParams{
		UserID:    user,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expDate,
			Valid: true,
		},
		Revoked: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
	}
}
