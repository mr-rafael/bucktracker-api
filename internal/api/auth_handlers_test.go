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

func TestLoginIncorrectPassword(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	mockUserID, _ := uuid.NewRandom()
	mockUserPassword := "password"
	mockUserPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(mockUserPassword), bcrypt.DefaultCost)
	badPassword := "badpass"

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
		strings.NewReader(fmt.Sprintf(`{"password": "%v"}`, badPassword)),
	)
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestLoginUserNotFound(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	mockUserPassword := "password"

	mockAuthRepo := &service.MockAuthRepo{
		CreateRefreshTokenFunc: func(ctx context.Context, userID pgtype.UUID, tokenHash string, expDate time.Time) (db.RefreshToken, error) {
			return db.RefreshToken{
				TokenHash: "TOKENHASH",
			}, nil
		},
	}
	mockUsersRepo := &service.MockUsersRepo{
		GetUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
			return db.User{}, fmt.Errorf("Not found")
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

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestRefreshOK(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	mockRefreshToken := "test"
	mockUserID := uuid.New()

	mockAuthRepo := &service.MockAuthRepo{
		GetTokenByHashFunc: func(ctx context.Context, tokenHash string) (db.RefreshToken, error) {
			return db.RefreshToken{
				UserID: pgtype.UUID{
					Bytes: mockUserID,
					Valid: true,
				},
				TokenHash: "TOKENHASH",
				ExpiresAt: pgtype.Timestamptz{
					Time:  time.Now().Add(5 * time.Minute),
					Valid: true,
				},
				Revoked: pgtype.Bool{
					Bool:  false,
					Valid: true,
				},
			}, nil
		},
	}
	mockUsersRepo := &service.MockUsersRepo{}
	service := service.NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)
	handler := NewAuthHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/refresh",
		strings.NewReader(fmt.Sprintf(`{"refresh_token": "%v"}`, mockRefreshToken)),
	)
	cookie := http.Cookie{
		Name:  "refresh_token",
		Value: mockRefreshToken,
	}
	req.AddCookie(&cookie)
	rr := httptest.NewRecorder()

	handler.Refresh(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestRefreshExpired(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	mockRefreshToken := "test"
	mockUserID := uuid.New()

	mockAuthRepo := &service.MockAuthRepo{
		GetTokenByHashFunc: func(ctx context.Context, tokenHash string) (db.RefreshToken, error) {
			return db.RefreshToken{
				UserID: pgtype.UUID{
					Bytes: mockUserID,
					Valid: true,
				},
				TokenHash: "TOKENHASH",
				ExpiresAt: pgtype.Timestamptz{
					Time:  time.Now().Add(-5 * time.Minute),
					Valid: true,
				},
				Revoked: pgtype.Bool{
					Bool:  true,
					Valid: true,
				},
			}, nil
		},
	}
	mockUsersRepo := &service.MockUsersRepo{}
	service := service.NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)
	handler := NewAuthHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/refresh",
		nil,
	)
	cookie := http.Cookie{
		Name:  "refresh_token",
		Value: mockRefreshToken,
	}
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()

	handler.Refresh(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}
