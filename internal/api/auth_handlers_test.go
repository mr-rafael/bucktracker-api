package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginOK(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	mockUserID, _ := uuid.NewRandom()
	mockUserPassword := "password"
	mockUserPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(mockUserPassword), bcrypt.DefaultCost)

	mockAuthRepo := &service.MockAuthRepo{
		CreateRefreshTokenFunc: func(ctx context.Context, userID pgtype.UUID, tokenHash string, expDate time.Time) (db.RefreshToken, error) {
			return db.RefreshToken{
				TokenHash: "TOKENHASH",
			}, nil
		},
	}
	mockUsersRepo := &service.MockUsersRepo{
		GetUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
			return db.User{
				ID: pgtype.UUID{
					Bytes: mockUserID,
					Valid: true,
				},
				PasswordHash: string(mockUserPasswordHash),
			}, nil
		},
	}
	service := service.NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)
	handler := NewAuthHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/login",
		strings.NewReader(fmt.Sprintf(`{"password": "%v"}`, mockUserPassword)),
	)
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
