package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Mr-Rafael/bucktracker-api/internal/dto"
	"github.com/Mr-Rafael/bucktracker-api/internal/mapper"
	"github.com/Mr-Rafael/bucktracker-api/internal/service"
	"github.com/google/uuid"
)

type SavingsHandler struct {
	savingsService *service.SavingsService
}

func NewSavingsHandler(service *service.SavingsService) *SavingsHandler {
	return &SavingsHandler{savingsService: service}
}

func (handler *SavingsHandler) HandleCalculateSavings(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := dto.SavingsRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
		return
	}

	result, err := handler.savingsService.CalculateSavingsPlan(mapper.ToSavingsInput(reqParams))
	if err != nil {
		var inputErr service.SavingsInputError
		switch {
		case errors.As(err, &inputErr):
			respondWithError(writer, err.Error(), err.Error(), http.StatusBadRequest)
		default:
			respondWithError(writer, err.Error(), err.Error(), http.StatusInternalServerError)
		}
		return
	}
	respondWithJSON(writer, mapper.ToSavingsResponse(result), http.StatusOK)
}

func (handler *SavingsHandler) HandleSaveSavings(writer http.ResponseWriter, request *http.Request) {
	userID, ok := request.Context().Value(userIDKey).(string)
	if !ok {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(request.Body)
	reqParams := dto.SavingsSaveRequestParams{}
	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
		return
	}

	result, err := handler.savingsService.SaveSavingsPlan(context.Background(), mapper.ToSaveSavingsInput(userUUID, reqParams))
	if err != nil {
		var inputErr service.SavingsInputError
		switch {
		case errors.As(err, &inputErr):
			respondWithError(writer, err.Error(), err.Error(), http.StatusBadRequest)
		default:
			respondWithError(writer, err.Error(), err.Error(), http.StatusInternalServerError)
		}
		return
	}

	respondWithJSON(writer, mapper.ToSaveSavingsResponse(result), http.StatusCreated)
}

func (handler *SavingsHandler) HandleListSavings(writer http.ResponseWriter, request *http.Request) {
	userID, ok := request.Context().Value(userIDKey).(string)
	if !ok {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}

	result, err := handler.savingsService.GetSavingsPlansByUser(context.Background(), userUUID)

	respondWithJSON(writer, mapper.ToSavingsListResponse(result), http.StatusOK)
}

func (handler *SavingsHandler) HandleGetSavings(writer http.ResponseWriter, request *http.Request) {
	userID, ok := request.Context().Value(userIDKey).(string)
	if !ok {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}
	planID := request.PathValue("id")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		respondWithErrorCode(writer, "invalid plan ID in URL", http.StatusUnauthorized)
		return
	}

	result, err := handler.savingsService.GetSavingsPlan(context.Background(), planUUID, userUUID)
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("attempt to fetch plan %v by user %v", planUUID, userUUID), http.StatusUnauthorized)
		return
	}

	respondWithJSON(writer, mapper.ToGetSavingsResponse(result), http.StatusOK)
}

func (handler *SavingsHandler) HandleUpdateSavings(writer http.ResponseWriter, request *http.Request) {
	userID, ok := request.Context().Value(userIDKey).(string)
	if !ok {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}
	planID := request.PathValue("id")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		respondWithErrorCode(writer, "invalid plan ID in URL", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(request.Body)
	reqParams := dto.SavingsUpdateRequestParams{}
	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad update savings request", http.StatusBadRequest)
	}

	result, err := handler.savingsService.UpdateSavings(context.Background(), mapper.ToUpdateSavingsInput(planUUID, userUUID, reqParams))
	fmt.Printf("From the service, received the yearly interest rate: %v\n", result.YearlyInterestRate)
	if err != nil {
		var inputErr service.SavingsInputError
		switch {
		case errors.As(err, &inputErr):
			respondWithError(writer, err.Error(), err.Error(), http.StatusBadRequest)
		default:
			respondWithError(writer, err.Error(), err.Error(), http.StatusInternalServerError)
		}
		return
	}

	respondWithJSON(writer, mapper.ToSaveSavingsResponse(result), http.StatusOK)
}

func (handler *SavingsHandler) HandleDeleteSavings(writer http.ResponseWriter, request *http.Request) {
	userID, ok := request.Context().Value(userIDKey).(string)
	if !ok {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}
	planID := request.PathValue("id")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		respondWithErrorCode(writer, "invalid plan ID in URL", http.StatusUnauthorized)
		return
	}
	err = handler.savingsService.DeleteSavingsPlan(context.Background(), planUUID, userUUID)
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("failed attempt to delete savings plan %v by user %v", planID, userID), http.StatusNotFound)
		return
	}
	respondWithCode(writer, http.StatusNoContent)
}
