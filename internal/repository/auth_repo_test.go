package repository

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestCreateRefreshToken(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	authRepo := NewAuthRepo(queries)

	userData, err := CreateTestUserIfNotExists()

	got, err := authRepo.CreateRefreshToken(ctx, userData.ID, "test_token", time.Now())
	if err != nil {
		log.Fatalf("Error saving the refresh token in database: %v", err)
	}
	want := db.RefreshToken{
		UserID: userData.ID,
	}

	if got.UserID.Bytes != want.UserID.Bytes {
		log.Fatalf("Saved (%v) and expected (%v) user IDs did not match.", got.UserID.Bytes, want.UserID.Bytes)
	}

	DeleteTestUser()
}

func TestGetTokenByHash(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewAuthRepo(queries)

	tokenHash := "1d0c6a19d3602a7608a1a2218671c1d444ca415fc685cb9182c2584e9ce395b6"

	got, err := repo.GetTokenByHash(ctx, tokenHash)
	if err != nil {
		log.Fatalf("Error getting token from database: %v", err)
	}

	testTokenID, err := uuid.Parse("aabfe0ac-3e13-4744-8a31-4073b69caa68")
	if err != nil {
		log.Fatalf("failed to parse the test token uuid: %v", err)
	}
	want := db.RefreshToken{
		ID: pgtype.UUID{
			Bytes: testTokenID,
			Valid: true,
		},
	}

	if got.ID.Bytes != want.ID.Bytes {
		log.Fatalf("Read (%v) and expected (%v) user IDs did not match.", got.UserID.Bytes, want.UserID.Bytes)
	}
}

func TestRevokeToken(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewAuthRepo(queries)

	test_user_id, err := uuid.Parse("af38df43-3ced-4869-9930-93a0fa0cf1e0")
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	user := pgtype.UUID{
		Bytes: test_user_id,
		Valid: true,
	}
	newToken, err := repo.CreateRefreshToken(ctx, user, "test_token", time.Now())
	if err != nil {
		log.Fatalf("Error saving the refresh token in database: %v", err)
	}

	err = repo.RevokeTokenByUserID(ctx, newToken.UserID)
	if err != nil {
		log.Fatalf("Error revoking refresh token: %v", err)
	}

	got, err := repo.GetTokenByHash(ctx, newToken.TokenHash)
	if err != nil {
		log.Fatalf("Error getting revoked token: %v", err)
	}

	want := db.RefreshToken{
		Revoked: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
	}

	if got.ID.Valid != want.ID.Valid {

	}
}
