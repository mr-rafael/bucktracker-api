package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/Mr-Rafael/finance-calculator/internal/mapper"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(service *service.AuthService) AuthHandler {
	return AuthHandler{authService: service}
}

func (handler *AuthHandler) Login(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := dto.UserLoginRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("received bad login request: %v", err), http.StatusBadRequest)
		return
	}
	result, err := handler.authService.Login(context.Background(), mapper.ToLoginInput(reqParams))
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("failed login attempt for user '%v': %v", reqParams.Email, err), http.StatusUnauthorized)
		return
	}
	secure := os.Getenv("ENV") == "production"
	http.SetCookie(writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    result.RefreshToken,
		HttpOnly: true,
		Secure:   secure,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 7,
	})

	respondWithJSON(writer, mapper.ToLoginResponse(result), http.StatusOK)
}

func (handler *AuthHandler) Refresh(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("refresh_token")
	if err != nil {
		respondWithError(writer, "request is missing refresh token", "missing refresh token", http.StatusUnauthorized)
		return
	}

	result, err := handler.authService.Refresh(context.Background(), (mapper.ToRefreshInput(cookie.Value)))
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("Error refreshing token: %v", err), http.StatusUnauthorized)
	}

	respondWithJSON(writer, mapper.ToRefreshResponse(result), http.StatusOK)
}
