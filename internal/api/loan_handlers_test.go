package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCalculateLoan(t *testing.T) {
	mockLoansRepo := &service.MockLoansRepo{}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/loans/calculate",
		strings.NewReader(`{
			"startingPrincipal": 10000000,
			"yearlyInterestRate": "5",
			"monthlyPayment": 1500000,
			"escrowPayment": 10000,
			"startDate": "1970-01-01"
		}`),
	)
	rr := httptest.NewRecorder()

	handler.HandleCalculateLoan(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestCalculateBadRequest(t *testing.T) {
	mockLoansRepo := &service.MockLoansRepo{}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/loans/calculate",
		strings.NewReader(`{
			"principal": 10000000,
			"interestRate": "5",
			"payment": 1500000,
			"escrow": 10000,
			"startDate": "1970-01-01"
		}`),
	)
	rr := httptest.NewRecorder()

	handler.HandleCalculateLoan(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSaveLoan(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockLoansRepo{}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/loans/save",
		strings.NewReader(`{
   			"name": "Test 2",
			"startingPrincipal": 10000000,
			"yearlyInterestRate": "5",
			"monthlyPayment": 900076,
			"escrowPayment": 10000,
			"startDate": "1970-01-01"
		}`),
	)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleSaveLoan(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestSaveLoanBadRequest(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockLoansRepo{}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/loans/save",
		strings.NewReader(`{
			"startingPrincipal": 10000000,
			"interestRate": "5",
			"monthlyPayment": 1500000,
			"escrowPayment": 10000
		}`),
	)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleSaveLoan(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestListLoans(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockLoansRepo{}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/app/loans/calculate",
		nil,
	)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleListLoans(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestListLoansUnauthorized(t *testing.T) {
	mockLoansRepo := &service.MockLoansRepo{}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/app/loans/calculate",
		nil,
	)
	rr := httptest.NewRecorder()

	handler.HandleListLoans(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetLoan(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()
	mockLoanID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockLoansRepo{}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/app/loans/get/%v", mockLoanID.String()),
		nil,
	)
	req.SetPathValue("id", mockLoanID.String())
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleGetLoan(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestGetLoanUnauthorized(t *testing.T) {
	mockLoansRepo := &service.MockLoansRepo{}
	mockLoanID, _ := uuid.NewRandom()
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/app/loans/get/%v", mockLoanID.String()),
		nil,
	)
	rr := httptest.NewRecorder()

	handler.HandleGetLoan(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateLoan(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()
	mockLoanID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockLoansRepo{
		GetLoanInitialDataFunc: func(ctx context.Context, planID uuid.UUID, userID uuid.UUID) (domain.UpdateLoanData, error) {
			return domain.UpdateLoanData{
				ID:   planID,
				Name: "originalName",
				LoanData: domain.LoansInput{
					StartingPrincipal:  10000,
					YearlyInterestRate: "5",
					MonthlyPayment:     1000,
					EscrowPayment:      100,
					StartDate:          "1970-01-01",
				},
			}, nil
		},
		UpdateLoanFunc: func(ctx context.Context, plan domain.LoanPaymentPlan) (db.Loan, error) {
			return db.Loan{
				Name:               plan.Name,
				StartingPrincipal:  int32(plan.OriginalData.StartingPrincipal),
				YearlyInterestRate: plan.OriginalData.YearlyInterestRate,
			}, nil
		},
	}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/app/loans",
		strings.NewReader(`{
			"interestRate": "5"
		}`),
	)
	req.SetPathValue("id", mockLoanID.String())
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleUpdateLoan(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateLoanUnauthorized(t *testing.T) {
	mockLoanID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockLoansRepo{}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("/app/loans/update/%v", mockLoanID.String()),
		nil,
	)
	rr := httptest.NewRecorder()

	handler.HandleUpdateLoan(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDeleteLoan(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()
	mockLoanID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockLoansRepo{}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("/app/loans/%v", mockLoanID.String()),
		nil,
	)
	req.SetPathValue("id", mockLoanID.String())
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleDeleteLoan(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteLoanUnauthorized(t *testing.T) {
	mockLoanID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockLoansRepo{}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("/app/loans/%v", mockLoanID.String()),
		nil,
	)
	req.SetPathValue("id", mockLoanID.String())
	rr := httptest.NewRecorder()

	handler.HandleDeleteLoan(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}
