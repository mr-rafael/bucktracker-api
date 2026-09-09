package repository

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Mr-Rafael/bucktracker-api/internal/db"
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

	testUser, err := CreateTestUserIfNotExists()
	if err != nil {
		log.Fatalf("Failed to create a test user for unit test: %v", err)
		return
	}
	tokenData, err := CreateTestTokenIfNotExists(testUser.ID)
	if err != nil {
		log.Fatalf("Failed to create test refresh token: %v", err)
		DeleteTestUser()
		return
	}
	got, err := repo.GetTokenByHash(ctx, tokenData.TokenHash)
	if err != nil {
		log.Fatalf("Error getting token from database: %v", err)
	}
	want := db.RefreshToken{
		ID: pgtype.UUID{
			Bytes: tokenData.ID.Bytes,
			Valid: true,
		},
	}

	if got.ID.Bytes != want.ID.Bytes {
		log.Fatalf("The obtained (%v) and expected (%v) user IDs did not match.", got.UserID.Bytes, want.UserID.Bytes)
	}

	DeleteTestUser()
}

func TestRevokeToken(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewAuthRepo(queries)

	testUser, err := CreateTestUserIfNotExists()
	if err != nil {
		log.Fatalf("Failed to create a test user for the unit test: %v", err)
		return
	}

	newToken, err := repo.CreateRefreshToken(ctx, testUser.ID, "test_token", time.Now())
	if err != nil {
		log.Fatalf("Error saving the refresh token in database: %v", err)
		return
	}

	err = repo.RevokeTokenByUserID(ctx, newToken.UserID)
	if err != nil {
		log.Fatalf("Error revoking refresh token: %v", err)
		return
	}

	got, err := repo.GetTokenByHash(ctx, newToken.TokenHash)
	if err != nil {
		log.Fatalf("Error getting revoked token: %v", err)
		return
	}

	want := db.RefreshToken{
		Revoked: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
	}

	if got.Revoked.Bool != want.Revoked.Bool {
		log.Fatalf("Expected the token to be invalid (%v), but it was valid (%v).", want.Revoked.Bool, got.Revoked.Bool)
	}

	DeleteTestUser()
}

func CreateTestTokenIfNotExists(userID pgtype.UUID) (db.RefreshToken, error) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	testHash := "TESTHASH"

	insertParams := db.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: testHash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(5 * time.Minute),
			Valid: true,
		},
		Revoked: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
	}

	token, err := queries.GetTokenByHash(ctx, testHash)
	if err != nil {
		return queries.CreateRefreshToken(ctx, insertParams)
	}
	return token, nil
}
