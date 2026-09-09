package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mr-Rafael/bucktracker-api/internal/db"
	"github.com/Mr-Rafael/bucktracker-api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	mockUsersRepo := &service.MockUsersRepo{
		CreateUserFunc: func(ctx context.Context, params db.CreateUserParams) (db.User, error) {
			return db.User{}, nil
		},
	}
	service := service.NewUserService(mockUsersRepo)
	handler := NewUsersHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/savings/calculate",
		strings.NewReader(`{
			"email":    "test@user.com",
			"password": "password",
			"username": "username"
		}`),
	)
	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestCreateUserError(t *testing.T) {
	mockUsersRepo := &service.MockUsersRepo{
		CreateUserFunc: func(ctx context.Context, params db.CreateUserParams) (db.User, error) {
			return db.User{}, fmt.Errorf("Error saving user")
		},
	}
	service := service.NewUserService(mockUsersRepo)
	handler := NewUsersHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/savings/calculate",
		strings.NewReader(`{
			"email":    "test@user.com",
			"password": "password",
			"username": "username"
		}`),
	)
	rr := httptest.NewRecorder()

	handler.CreateUser(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
