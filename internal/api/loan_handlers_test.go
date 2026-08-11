package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Mr-Rafael/bucktracker-api/internal/db"
	"github.com/Mr-Rafael/bucktracker-api/internal/domain"
	"github.com/Mr-Rafael/bucktracker-api/internal/dto"
	"github.com/Mr-Rafael/bucktracker-api/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
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
	mockLoanID, _ := uuid.NewRandom()
	mockPlanID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockLoansRepo{
		SaveLoanPaymentPlanFunc: func(ctx context.Context, loan domain.Loan) (db.Loan, error) {
			require.NotNil(t, loan.DefaultPaymentPlan)
			require.Equal(t, "Default Payment Plan", loan.DefaultPaymentPlan.Name)
			require.Empty(t, loan.DefaultPaymentPlan.PrincipalPayments)
			return db.Loan{
				ID: pgtype.UUID{Bytes: mockLoanID, Valid: true},
				DefaultPaymentPlan: pgtype.UUID{
					Bytes: mockPlanID,
					Valid: true,
				},
			}, nil
		},
	}
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

	var body dto.LoanCreateResponseParams
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&body))
	require.Equal(t, mockLoanID.String(), body.ID)
	require.Equal(t, "Test 2", body.Name)
	require.Equal(t, 10000000, body.StartingPrincipal)
	require.Equal(t, "5", body.YearlyInterestRate)
	require.Equal(t, mockPlanID.String(), body.DefaultPaymentPlan.ID)
	require.Equal(t, "Default Payment Plan", body.DefaultPaymentPlan.Name)
	require.Equal(t, 12, body.DefaultPaymentPlan.DurationMonths)
	require.Len(t, body.PaymentPlans, 1)
	require.Equal(t, body.DefaultPaymentPlan, body.PaymentPlans[0])
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

func TestGetPaymentPlan(t *testing.T) {
	mockUserID, _ := uuid.NewRandom()
	mockLoanID, _ := uuid.NewRandom()
	mockPlanID, _ := uuid.NewRandom()
	stateDate := time.Date(1970, 2, 1, 0, 0, 0, 0, time.UTC)

	mockLoansRepo := &service.MockLoansRepo{
		GetPaymentPlanByIDFunc: func(ctx context.Context, loanID uuid.UUID, paymentPlanID uuid.UUID, userID uuid.UUID) (domain.LoanPaymentPlan, error) {
			require.Equal(t, mockLoanID, loanID)
			require.Equal(t, mockPlanID, paymentPlanID)
			require.Equal(t, mockUserID, userID)
			return domain.LoanPaymentPlan{
				ID:                  mockPlanID,
				Name:                "Default Payment Plan",
				DurationMonths:      12,
				TotalExpenditure:    decimal.NewFromInt(120000),
				TotalPaid:           decimal.NewFromInt(10800000),
				CostOfCreditPercent: decimal.RequireFromString("8.00"),
				Plan: []domain.LoanStatus{
					{
						Date:          stateDate,
						Payment:       decimal.NewFromInt(900076),
						Interest:      decimal.NewFromInt(41667),
						OtherPayments: decimal.NewFromInt(10000),
						Paydown:       decimal.NewFromInt(848409),
						Principal:     decimal.NewFromInt(9151591),
					},
				},
			}, nil
		},
	}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/app/loans/%v/payment-plans/%v", mockLoanID.String(), mockPlanID.String()),
		nil,
	)
	req.SetPathValue("loanId", mockLoanID.String())
	req.SetPathValue("paymentPlanId", mockPlanID.String())
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), userIDKey, mockUserID.String())

	handler.HandleGetPaymentPlan(rr, req.WithContext(ctx))

	require.Equal(t, http.StatusOK, rr.Code)

	var body dto.PaymentPlanDetailResponseParams
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&body))
	require.Equal(t, mockPlanID.String(), body.ID)
	require.Equal(t, "Default Payment Plan", body.Name)
	require.Equal(t, 12, body.DurationMonths)
	require.Equal(t, 120000, body.TotalExpenditure)
	require.Equal(t, 10800000, body.TotalPaid)
	require.Equal(t, "8", body.CostOfCredit)
	require.Len(t, body.PaymentPlanBreakdown, 1)
	require.Equal(t, 900076, body.PaymentPlanBreakdown[0].Payment)
	require.Equal(t, 41667, body.PaymentPlanBreakdown[0].Interest)
	require.Equal(t, 10000, body.PaymentPlanBreakdown[0].OtherPayments)
	require.Equal(t, 848409, body.PaymentPlanBreakdown[0].Paydown)
	require.Equal(t, 9151591, body.PaymentPlanBreakdown[0].Principal)
}

func TestGetPaymentPlanUnauthorized(t *testing.T) {
	mockLoanID, _ := uuid.NewRandom()
	mockPlanID, _ := uuid.NewRandom()

	mockLoansRepo := &service.MockLoansRepo{}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/app/loans/%v/payment-plans/%v", mockLoanID.String(), mockPlanID.String()),
		nil,
	)
	req.SetPathValue("loanId", mockLoanID.String())
	req.SetPathValue("paymentPlanId", mockPlanID.String())
	rr := httptest.NewRecorder()

	handler.HandleGetPaymentPlan(rr, req)

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
		GetLoanByIDFunc: func(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) (domain.Loan, error) {
			planID, _ := uuid.NewRandom()
			return domain.Loan{
				ID:   loanID,
				Name: "originalName",
				OriginalData: domain.LoansInput{
					StartingPrincipal:  10000,
					YearlyInterestRate: "5",
					MonthlyPayment:     1000,
					EscrowPayment:      100,
					StartDate:          "1970-01-01",
				},
				DefaultPaymentPlan: &domain.LoanPaymentPlan{
					ID:             planID,
					Name:           "Default Payment Plan",
					DurationMonths: 12,
				},
				PaymentPlans: []domain.LoanPaymentPlan{
					{
						ID:             planID,
						Name:           "Default Payment Plan",
						DurationMonths: 12,
					},
				},
			}, nil
		},
		UpdateLoanFunc: func(ctx context.Context, loan domain.Loan) (db.Loan, error) {
			return db.Loan{
				ID: pgtype.UUID{
					Bytes: loan.ID,
					Valid: true,
				},
				Name:               loan.Name,
				StartingPrincipal:  int32(loan.OriginalData.StartingPrincipal),
				YearlyInterestRate: loan.OriginalData.YearlyInterestRate,
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
