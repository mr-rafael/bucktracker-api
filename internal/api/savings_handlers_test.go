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

func TestCalculateSavings(t *testing.T) {
	mockSavingsRepo := &service.MockSavingsRepo{}
	service := service.NewSavingsService(mockSavingsRepo)
	handler := NewSavingsHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/savings/calculate",
		strings.NewReader(`{
			"startingCapital": 10000000,
			"yearlyInterestRate": "5",
			"monthlyContribution": 10000,
			"durationYears": 1,
			"startDate": "2026-02-01"
		}`),
	)
	rr := httptest.NewRecorder()

	handler.HandleCalculateSavings(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestCalculateSavingsBadRequest(t *testing.T) {
	mockSavingsRepo := &service.MockSavingsRepo{}
	service := service.NewSavingsService(mockSavingsRepo)
	handler := NewSavingsHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/savings/calculate",
		strings.NewReader(`{
			"capital": 10000000,
			"yearlyInterestRate": "5",
			"monthlyContribution": 10000,
			"durationYears": 1,
			"startDate": "2026-02-01"
		}`),
	)
	rr := httptest.NewRecorder()

	handler.HandleCalculateSavings(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSaveSavings(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()

	mockSavingsRepo := &service.MockSavingsRepo{}
	service := service.NewSavingsService(mockSavingsRepo)
	handler := NewSavingsHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/savings/save",
		strings.NewReader(`{
		    "name": "test",
			"startingCapital": 700000,
			"yearlyInterestRate": "4.75",
			"monthlyContribution": 15000,
			"durationYears": 1,
			"startDate": "2026-02-01"
		}`),
	)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleSaveSavings(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestSaveSavingsBadRequest(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()

	mockSavingsRepo := &service.MockSavingsRepo{}
	service := service.NewSavingsService(mockSavingsRepo)
	handler := NewSavingsHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/savings/save",
		strings.NewReader(`{
		    "savingsName": "test",
			"badfieldname": 700000,
			"yearlyInterestRate": "4.75",
			"monthlyContribution": 15000,
			"durationYears": 1,
			"startDate": "2026-02-01"
		}`),
	)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleSaveSavings(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestListSavings(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()

	mockSavingsRepo := &service.MockSavingsRepo{}
	service := service.NewSavingsService(mockSavingsRepo)
	handler := NewSavingsHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/app/savings/list",
		nil,
	)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleListSavings(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestListSavingsUnauthorized(t *testing.T) {
	mockSavingsRepo := &service.MockSavingsRepo{}
	service := service.NewSavingsService(mockSavingsRepo)
	handler := NewSavingsHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/app/savings/list",
		nil,
	)
	rr := httptest.NewRecorder()

	handler.HandleListSavings(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetSavings(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()
	mockLoanID, _ := uuid.NewRandom()

	mockSavingsRepo := &service.MockSavingsRepo{}
	service := service.NewSavingsService(mockSavingsRepo)
	handler := NewSavingsHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/app/savings/get/%v", mockLoanID.String()),
		nil,
	)
	req.SetPathValue("id", mockLoanID.String())
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleGetSavings(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestGetSavingsUnauthorized(t *testing.T) {
	mockLoanID, _ := uuid.NewRandom()

	mockSavingsRepo := &service.MockSavingsRepo{}
	service := service.NewSavingsService(mockSavingsRepo)
	handler := NewSavingsHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/app/savings/get/%v", mockLoanID.String()),
		nil,
	)
	req.SetPathValue("id", mockLoanID.String())
	rr := httptest.NewRecorder()

	handler.HandleGetSavings(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUpdateSavings(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()
	mockSavingsID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockSavingsRepo{
		GetSavingsInitialDataFunc: func(ctx context.Context, planID uuid.UUID, userID uuid.UUID) (domain.UpdateSavingsData, error) {
			return domain.UpdateSavingsData{
				ID:   planID,
				Name: "originalName",
				SavingsData: domain.SavingsInput{
					StartingCapital:     10000,
					YearlyInterestRate:  "5",
					InterestRateType:    "APY",
					MonthlyContribution: 100,
					DurationYears:       1,
					TaxRate:             "12",
					YearlyInflationRate: "8",
					StartDate:           "1970-01-01",
				},
			}, nil
		},
		UpdateSavingsFunc: func(ctx context.Context, plan domain.SavingsPlan) (db.Saving, error) {
			return db.Saving{
				Name:               plan.Name,
				StartingCapital:    int32(plan.OriginalData.StartingCapital),
				YearlyInterestRate: plan.OriginalData.YearlyInterestRate,
			}, nil
		},
	}
	service := service.NewSavingsService(mockLoansRepo)
	handler := NewSavingsHandler(service)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/app/savings",
		strings.NewReader(`{
			"interestRate": "5"
		}`),
	)
	req.SetPathValue("id", mockSavingsID.String())
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleUpdateSavings(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateSavingsUnauthorized(t *testing.T) {
	mockSavingsID, _ := uuid.NewRandom()

	mockSavingsRepo := &service.MockSavingsRepo{}
	service := service.NewSavingsService(mockSavingsRepo)
	handler := NewSavingsHandler(service)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/app/savings",
		strings.NewReader(`{
			"interestRate": "5"
		}`),
	)
	req.SetPathValue("id", mockSavingsID.String())
	rr := httptest.NewRecorder()

	handler.HandleUpdateSavings(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDeleteSavings(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()
	mockSavingsID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockSavingsRepo{}
	service := service.NewSavingsService(mockLoansRepo)
	handler := NewSavingsHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("/app/savings/%v", mockSavingsID.String()),
		nil,
	)
	req.SetPathValue("id", mockSavingsID.String())
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleDeleteSavings(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteSavingsUnauthorized(t *testing.T) {
	mockSavingsID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockSavingsRepo{}
	service := service.NewSavingsService(mockLoansRepo)
	handler := NewSavingsHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("/app/savings/%v", mockSavingsID.String()),
		nil,
	)
	req.SetPathValue("id", mockSavingsID.String())
	rr := httptest.NewRecorder()

	handler.HandleDeleteSavings(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}
