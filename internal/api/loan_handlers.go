package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/Mr-Rafael/finance-calculator/internal/mapper"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
	"github.com/google/uuid"
)

type LoanHandler struct {
	loanService *service.LoansService
}

func NewLoanHandler(service *service.LoansService) *LoanHandler {
	return &LoanHandler{loanService: service}
}

func (handler *LoanHandler) HandleCalculateLoan(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := dto.LoanRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
		return
	}

	result, err := handler.loanService.CalculateLoanPaymentPlan(mapper.ToLoanInput(reqParams))
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error calculating loan payment plan: %v", err), fmt.Sprintf("Error calculating loan payment plan: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(writer, mapper.ToLoanResponse(result), http.StatusOK)
}

func (handler *LoanHandler) HandleSaveLoan(writer http.ResponseWriter, request *http.Request) {
	userID := request.Context().Value(userIDKey).(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(request.Body)
	reqParams := dto.LoanSaveRequestParams{}
	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
		return
	}

	result, err := handler.loanService.SaveLoanPaymentPlan(context.Background(), mapper.ToSaveLoanInput(userUUID, reqParams))
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error saving the plan: %v", err), fmt.Sprintf("Error saving the plan: %v", err), http.StatusInternalServerError)
		return
	}

	respondWithJSON(writer, mapper.ToSaveLoanResponse(result), http.StatusCreated)
}

func (handler *LoanHandler) HandleListLoans(writer http.ResponseWriter, request *http.Request) {
	userID := request.Context().Value(userIDKey).(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}

	result, err := handler.loanService.GetLoansByUser(context.Background(), userUUID)

	respondWithJSON(writer, mapper.ToLoanListResponse(result), http.StatusOK)
}

func (handler *LoanHandler) HandleGetLoan(writer http.ResponseWriter, request *http.Request) {
	userID := request.Context().Value(userIDKey).(string)
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

	result, err := handler.loanService.GetLoan(context.Background(), planUUID, userUUID)
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("attempt to fetch loan %v by user %v", planUUID, userUUID), http.StatusUnauthorized)
		return
	}

	respondWithJSON(writer, mapper.ToGetLoanResponse(result), http.StatusOK)
}

func (handler *LoanHandler) HandleUpdateLoan(writer http.ResponseWriter, request *http.Request) {
	userID := request.Context().Value(userIDKey).(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}
	planID := request.PathValue("id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		respondWithErrorCode(writer, "invalid plan ID", http.StatusNotFound)
	}

	decoder := json.NewDecoder(request.Body)
	reqParams := dto.LoanUpdateRequestParams{}
	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad update loan request", http.StatusBadRequest)
		return
	}

	result, err := handler.loanService.UpdateLoan(context.Background(), mapper.ToUpdateLoanInput(planUUID, userUUID, reqParams))
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error saving the plan: %v", err), fmt.Sprintf("Error saving the plan: %v", err), http.StatusInternalServerError)
		return
	}

	respondWithJSON(writer, mapper.ToSaveLoanResponse(result), http.StatusCreated)
}

func (handler *LoanHandler) HandleDeleteLoan(writer http.ResponseWriter, request *http.Request) {
	userID := request.Context().Value(userIDKey).(string)
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
	err = handler.loanService.DeleteLoan(context.Background(), planUUID, userUUID)
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("failed attempt to delete loan %v by user %v", planID, userID), http.StatusUnauthorized)
	}
	respondWithCode(writer, http.StatusNoContent)
}
