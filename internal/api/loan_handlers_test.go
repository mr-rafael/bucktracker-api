package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mr-Rafael/finance-calculator/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCalculate(t *testing.T) {
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

func TestListLoansNoUserID(t *testing.T) {
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
