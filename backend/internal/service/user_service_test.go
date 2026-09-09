package service

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/Mr-Rafael/bucktracker-api/internal/db"
)

func TestRegisterUserAlreadyExists(t *testing.T) {
	mockUsersRepo := &MockUsersRepo{
		CreateUserFunc: func(ctx context.Context, params db.CreateUserParams) (db.User, error) {
			return db.User{}, fmt.Errorf("email field violates unique constraint")
		},
	}
	service := NewUserService(mockUsersRepo)
	ctx := context.Background()

	input := RegisterUserInput{
		Email:    "existing@user.com",
		Password: "password",
		Username: "username",
	}

	_, err := service.RegisterUser(ctx, input)
	if err == nil {
		log.Fatalf("Expected an error to return for existing user, but none did.")
	}
}
