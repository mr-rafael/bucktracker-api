package repository

import (
	"context"
	"log"
	"testing"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewUsersRepo(queries)
	testEmail := "unit@test.com"

	insertParams := db.CreateUserParams{
		Email:        testEmail,
		PasswordHash: "password",
		Username:     "Unit",
	}

	got, err := repo.CreateUser(ctx, insertParams)
	if err != nil {
		log.Fatalf("Failed to insert user: %v", err)
	}
	err = repo.DeleteUser(ctx, got.ID)
	if err != nil {
		log.Fatalf("Failed to delete user after insertion: %v", err)
	}

	want := db.User{
		Email: testEmail,
	}

	if got.Email != want.Email {
		log.Fatalf("The inserted email (%v) and expected email (%v) did not match.", got.Email, want.Email)
	}
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewUsersRepo(queries)
	testEmail := "test@mail.com"
	testID := "af38df43-3ced-4869-9930-93a0fa0cf1e0"

	got, err := repo.GetUserByEmail(ctx, testEmail)
	if err != nil {
		log.Fatalf("Failed to get user from database: %v", err)
	}
	testUUID, err := uuid.Parse(testID)
	if err != nil {
		log.Fatalf("Failed to parse test user's UUID: %v", testUUID)
	}
	want := db.User{
		ID: pgtype.UUID{
			Bytes: testUUID,
			Valid: true,
		},
	}

	if got.ID.Bytes != want.ID.Bytes {
		log.Fatalf("The expected UUID (%v) did not match the fetched one (%v)", want.ID.Bytes, got.ID.Bytes)
	}
}

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewUsersRepo(queries)
	testEmail := "test@mail.com"
	testID := "af38df43-3ced-4869-9930-93a0fa0cf1e0"
	testUUID, err := uuid.Parse(testID)
	if err != nil {
		log.Fatalf("Failed to parse test user's UUID: %v", testUUID)
	}
	getUserParam := pgtype.UUID{
		Bytes: testUUID,
		Valid: true,
	}

	got, err := repo.GetUserByID(ctx, getUserParam)
	if err != nil {
		log.Fatalf("Failed to get user from database: %v", err)
	}
	want := db.User{
		Email: testEmail,
	}

	if got.Email != want.Email {
		log.Fatalf("The expected email (%v) did not match the fetched one (%v)", want.ID.Bytes, got.ID.Bytes)
	}
}

func CreateTestUserIfNotExists() (db.User, error) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	testEmail := "unit@test.com"

	insertParams := db.CreateUserParams{
		Email:        testEmail,
		PasswordHash: "password",
		Username:     "Unit",
	}

	user, err := queries.GetUserByEmail(ctx, testEmail)
	if err != nil {
		return queries.CreateUser(ctx, insertParams)

	}
	return user, nil
}

func DeleteTestUser() {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	testEmail := "unit@test.com"

	user, err := queries.GetUserByEmail(ctx, testEmail)
	if err == nil {
		queries.DeleteUser(ctx, user.ID)
		return
	}
}
