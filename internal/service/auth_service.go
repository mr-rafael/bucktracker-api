package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/auth"
	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo      AuthRepository
	usersRepo     UsersRepository
	accessSecret  string
	refreshSecret string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginInfo struct {
	ID           pgtype.UUID
	Email        string
	UserName     string
	AccessToken  string
	RefreshToken string
}

type RefreshInput struct {
	RefreshToken string
}

type RefreshInfo struct {
	AccessToken string
}

type AuthRepository interface {
	CreateRefreshToken(context.Context, pgtype.UUID, string, time.Time) (db.RefreshToken, error)
	GetTokenByHash(context.Context, string) (db.RefreshToken, error)
	RevokeTokenByUserID(ctx context.Context, id pgtype.UUID) error
}

func NewAuthService(authRepo AuthRepository, usersRepo UsersRepository, accessSecret string, refreshSecret string) *AuthService {
	return &AuthService{
		authRepo:      authRepo,
		usersRepo:     usersRepo,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (s *AuthService) ValidateAccessToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(s.accessSecret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid user id in token")
	}

	return userID, nil
}

func (s *AuthService) ValidateRefreshToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(s.refreshSecret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid user id in token")
	}

	return userID, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (LoginInfo, error) {
	userInfo, err := s.usersRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("failed login attempt for user: %v. user not found", input.Email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.PasswordHash), []byte(input.Password))
	if err != nil {
		return LoginInfo{}, fmt.Errorf("password hash mismatch.")
	}

	accessToken, err := auth.GenerateAccessToken(userInfo.ID.String(), s.accessSecret)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("error generating access token: %v", err)
	}

	refreshToken, expDate, err := auth.GenerateRefreshToken(userInfo.ID.String(), s.refreshSecret)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("error generating refresh token: %v", err)
	}
	refTokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(refreshToken)))

	_, err = s.authRepo.CreateRefreshToken(ctx, userInfo.ID, refTokenHash, expDate)
	if err != nil {
		return LoginInfo{}, fmt.Errorf("error storing the refresh token: %v", err)
	}

	return LoginInfo{
		ID:           userInfo.ID,
		Email:        userInfo.Email,
		UserName:     userInfo.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, input RefreshInput) (RefreshInfo, error) {
	refreshTokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(input.RefreshToken)))

	tokenData, err := s.authRepo.GetTokenByHash(ctx, refreshTokenHash)
	if err != nil {
		return RefreshInfo{}, fmt.Errorf("failed to find refresh token for user '%v' in database: %v", tokenData.UserID, err)
	}

	if tokenData.Revoked.Bool {
		return RefreshInfo{}, fmt.Errorf("attempt to refresh with revoked token for user: %v", tokenData.UserID)
	}

	if tokenData.ExpiresAt.Time.Before(time.Now()) {
		return RefreshInfo{}, fmt.Errorf("refresh attempt with expired token.")
	}

	accessToken, err := auth.GenerateAccessToken(tokenData.UserID.String(), s.accessSecret)
	if err != nil {
		return RefreshInfo{}, fmt.Errorf("error generating access token: %v", err)
	}

	return RefreshInfo{AccessToken: accessToken}, nil
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

func ToLoginInfoModel(dbUser db.User) User {
	return User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt,
	}
}
