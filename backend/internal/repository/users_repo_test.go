package repository

import (
	"context"
	"log"
	"testing"

	"github.com/Mr-Rafael/bucktracker-api/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewUsersRepo(queries)
	testEmail := "create@unit.test"

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
	testUser, err := CreateTestUserIfNotExists()
	if err != nil {
		log.Fatalf("Failed to create a test user.")
		return
	}

	got, err := repo.GetUserByEmail(ctx, testUser.Email)
	if err != nil {
		log.Fatalf("Failed to get user from database: %v", err)
	}
	want := db.User{
		ID: pgtype.UUID{
			Bytes: testUser.ID.Bytes,
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
	testEmail := "unit@test.com"
	testUser, err := CreateTestUserIfNotExists()
	if err != nil {
		log.Fatalf("Failed to create a test user.")
		return
	}
	testUUID, err := uuid.Parse(testUser.ID.String())
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

	user, _ := queries.GetUserByEmail(ctx, testEmail)
	if user.ID.Bytes == uuid.Nil || len(user.ID.Bytes) <= 0 {
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
