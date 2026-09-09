package service

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Mr-Rafael/bucktracker-api/internal/auth"
	"github.com/Mr-Rafael/bucktracker-api/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestValidateAccessToken(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	mockAuthRepo := &MockAuthRepo{}
	mockUsersRepo := &MockUsersRepo{}
	service := NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)

	want := "001"
	claims := auth.AccessClaims{
		UserID: want,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "savings-app",
			Subject:   want,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(mockAccessSecret))
	if err != nil {
		log.Fatalf("Failed to generate test access token: %v", err)
	}

	got, err := service.ValidateAccessToken(signedToken)
	if err != nil {
		log.Fatalf("Failed to validate the access token: %v", err)
	}

	if want != got {
		log.Fatalf("The expected User ID (%v) did not match the one obtained from validation (%v)", want, got)
	}
}

func TestValidateExpiredToken(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	mockAuthRepo := &MockAuthRepo{
		CreateRefreshTokenFunc: nil,
	}
	mockUsersRepo := &MockUsersRepo{
		CreateUserFunc: nil,
	}
	service := NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)

	userID := "001"
	claims := auth.AccessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "savings-app",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(mockAccessSecret))
	if err != nil {
		log.Fatalf("Failed to generate test access token: %v", err)
	}

	_, err = service.ValidateAccessToken(signedToken)
	if err == nil {
		log.Fatalf("Validate function returned the expired token as valid.")
	}
}

func TestValidateInvalidSignature(t *testing.T) {
	mockAccessSecret := "ACCESS"
	badAccessSecret := "NOACCESS"
	mockRefreshSecret := "REFRESH"
	mockAuthRepo := &MockAuthRepo{
		CreateRefreshTokenFunc: nil,
	}
	mockUsersRepo := &MockUsersRepo{
		CreateUserFunc: nil,
	}
	service := NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)

	userID := "001"
	claims := auth.AccessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "savings-app",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(badAccessSecret))
	if err != nil {
		log.Fatalf("Failed to generate test access token: %v", err)
	}

	_, err = service.ValidateAccessToken(signedToken)
	if err == nil {
		log.Fatalf("Validate function wrongly signed token as valid.")
	}
}

func TestLogin(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	ctx := context.Background()
	mockUserID := uuid.Nil

	mockAuthRepo := &MockAuthRepo{
		CreateRefreshTokenFunc: func(ctx context.Context, userID pgtype.UUID, tokenHash string, expDate time.Time) (db.RefreshToken, error) {
			return db.RefreshToken{
				TokenHash: "TOKENHASH",
			}, nil
		},
	}
	mockUsersRepo := &MockUsersRepo{
		GetUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
			return db.User{
				ID: pgtype.UUID{
					Bytes: mockUserID,
					Valid: true,
				},
				PasswordHash: "$2a$10$olKeSVnknIIssUqv85e5wuH3dTMgNjjX1OClqan2TTpVe2tWoHIea",
			}, nil
		},
	}
	service := NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)

	input := LoginInput{
		Email:    "test2@mail.com",
		Password: "password",
	}

	got, err := service.Login(ctx, input)
	if err != nil {
		log.Fatalf("Failed to log in with the test user: %v", err)
	}

	tokenUserID, err := service.ValidateAccessToken(got.AccessToken)
	if err != nil {
		log.Fatalf("Login function returned an invalid access token: %v", err)
	}
	if tokenUserID != mockUserID.String() {
		log.Fatalf("Login function returned an access token with the incorrect User ID: %v. As a string: %v", tokenUserID, string(tokenUserID))
	}

	tokenUserID, err = service.ValidateRefreshToken(got.RefreshToken)
	if err != nil {
		log.Fatalf("Login function returned an invalid refresh token |%v|: %v", got.RefreshToken, err)
	}
	if tokenUserID != mockUserID.String() {
		log.Fatalf("Login function returned a refresh token with the incorrect User ID: %v", tokenUserID)
	}
}

func TestLoginBadUsername(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	ctx := context.Background()

	mockAuthRepo := &MockAuthRepo{
		CreateRefreshTokenFunc: func(ctx context.Context, userID pgtype.UUID, tokenHash string, expDate time.Time) (db.RefreshToken, error) {
			return db.RefreshToken{
				TokenHash: "TOKENHASH",
			}, nil
		},
	}
	mockUsersRepo := &MockUsersRepo{
		GetUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
			return db.User{}, fmt.Errorf("User not found.")
		},
	}
	service := NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)

	input := LoginInput{
		Email:    "test2@mail.com",
		Password: "password",
	}

	_, err := service.Login(ctx, input)
	if err == nil {
		log.Fatalf("The login should have failed because the user doesn't exist, but it did not.")
	}
}

func TestLoginBadPassword(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	ctx := context.Background()
	mockUserID := uuid.Nil

	mockAuthRepo := &MockAuthRepo{
		CreateRefreshTokenFunc: func(ctx context.Context, userID pgtype.UUID, tokenHash string, expDate time.Time) (db.RefreshToken, error) {
			return db.RefreshToken{
				TokenHash: "TOKENHASH",
			}, nil
		},
	}
	mockUsersRepo := &MockUsersRepo{
		GetUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
			return db.User{
				ID: pgtype.UUID{
					Bytes: mockUserID,
					Valid: true,
				},
				PasswordHash: "$2a$10$olKeSVnknIIssUqv85e5wuH3dTMgNjjX1OClqan2TTpVe2tWoHIeb",
			}, nil
		},
	}
	service := NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)

	input := LoginInput{
		Email:    "test2@mail.com",
		Password: "password",
	}

	_, err := service.Login(ctx, input)
	if err == nil {
		log.Fatalf("The login should have failed due to password mismatch, but it did not.")
	}
}

func TestRefresh(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	ctx := context.Background()
	mockUserID := uuid.Nil
	want := mockUserID.String()

	mockAuthRepo := &MockAuthRepo{
		GetTokenByHashFunc: func(ctx context.Context, hash string) (db.RefreshToken, error) {
			return db.RefreshToken{
				ID: pgtype.UUID{
					Bytes: mockUserID,
					Valid: true,
				},
				UserID: pgtype.UUID{
					Bytes: mockUserID,
					Valid: true,
				},
				Revoked: pgtype.Bool{
					Bool:  false,
					Valid: true,
				},
				ExpiresAt: pgtype.Timestamptz{
					Time:  time.Now().Add(1 * time.Minute),
					Valid: true,
				},
			}, nil
		},
	}
	mockUsersRepo := &MockUsersRepo{}

	service := NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)

	input := RefreshInput{
		RefreshToken: "refreshtoken",
	}
	tokenInfo, err := service.Refresh(ctx, input)
	if err != nil {
		log.Fatalf("There was an error generating the refresh token: %v", err)
	}
	got, err := service.ValidateAccessToken(tokenInfo.AccessToken)
	if err != nil {
		log.Fatalf("The generated access token was invalid: %v", err)
	}
	if got != want {
		log.Fatalf("Expected user ID | %v |, but got |%v|.", want, got)
	}
}

func TestRefreshExpiredToken(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	ctx := context.Background()
	mockUserID := uuid.Nil

	mockAuthRepo := &MockAuthRepo{
		GetTokenByHashFunc: func(ctx context.Context, hash string) (db.RefreshToken, error) {
			return db.RefreshToken{
				ID: pgtype.UUID{
					Bytes: mockUserID,
					Valid: true,
				},
				UserID: pgtype.UUID{
					Bytes: mockUserID,
					Valid: true,
				},
				Revoked: pgtype.Bool{
					Bool:  false,
					Valid: true,
				},
				ExpiresAt: pgtype.Timestamptz{
					Time:  time.Now().Add(-1 * time.Minute),
					Valid: true,
				},
			}, nil
		},
	}
	mockUsersRepo := &MockUsersRepo{}

	service := NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)

	input := RefreshInput{
		RefreshToken: "refreshtoken",
	}
	_, err := service.Refresh(ctx, input)
	if err == nil {
		log.Fatal("Expected the refresh function to fail with an expired refresh token, but it didn't.")
	}
}

func TestRefreshNotFound(t *testing.T) {
	mockAccessSecret := "ACCESS"
	mockRefreshSecret := "REFRESH"
	ctx := context.Background()

	mockAuthRepo := &MockAuthRepo{
		GetTokenByHashFunc: func(ctx context.Context, hash string) (db.RefreshToken, error) {
			return db.RefreshToken{}, fmt.Errorf("token not found.")
		},
	}
	mockUsersRepo := &MockUsersRepo{}

	service := NewAuthService(mockAuthRepo, mockUsersRepo, mockAccessSecret, mockRefreshSecret)

	input := RefreshInput{
		RefreshToken: "refreshtoken",
	}
	_, err := service.Refresh(ctx, input)
	if err == nil {
		log.Fatal("Expected the refresh function to fail when the token is not found, but it didn't.")
	}

}
