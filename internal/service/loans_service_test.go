package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/google/uuid"
)

func TestCalculateLoanPaymentPlan(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  10000000,
		YearlyInterestRate: "5",
		MonthlyPayment:     900076,
		EscrowPayment:      10000,
		StartDate:          "1970-01-01",
	}

	got, err := service.CalculateLoanPaymentPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the loan payment plan: %v", err)
	}

	want := domain.LoanPaymentPlan{
		DurationMonths: 12,
	}

	if got.DurationMonths != want.DurationMonths {
		log.Fatalf("Expected a duration of %v months, but got %v.", want.DurationMonths, got.DurationMonths)
	}
}

func TestCalculateLoanMaxTerm(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)
	maxLoanTerm := 360

	input := domain.LoansInput{
		StartingPrincipal:  10000000,
		YearlyInterestRate: "5",
		MonthlyPayment:     63683,
		EscrowPayment:      10000,
		StartDate:          "1970-01-01",
	}

	got, err := service.CalculateLoanPaymentPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the loan payment plan.")
	}

	want := domain.LoanPaymentPlan{
		DurationMonths: maxLoanTerm,
	}

	if got.DurationMonths != want.DurationMonths {
		log.Fatalf("Expected a duration of %v months, but got %v.", want.DurationMonths, got.DurationMonths)
	}
}

func TestCalculateLoanMinTerm(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)
	minLoanTerm := 1

	input := domain.LoansInput{
		StartingPrincipal:  1000000,
		YearlyInterestRate: "5",
		MonthlyPayment:     1014167,
		EscrowPayment:      10000,
		StartDate:          "1970-01-01",
	}

	got, err := service.CalculateLoanPaymentPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the loan payment plan.")
	}

	want := domain.LoanPaymentPlan{
		DurationMonths: minLoanTerm,
	}

	if got.DurationMonths != want.DurationMonths {
		log.Fatalf("Expected a duration of %v month, but got %v.", want.DurationMonths, got.DurationMonths)
	}
}

func TestCalculateLoanTermTooLong(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  10000000,
		YearlyInterestRate: "5",
		MonthlyPayment:     63682,
		EscrowPayment:      10000,
		StartDate:          "1970-01-01",
	}

	_, err := service.CalculateLoanPaymentPlan(input)
	if err == nil {
		log.Fatalf("Expected the loan calcuation to fail due to the term being longer than 360 months, but it didn't.")
	}
}

func TestCalculateZeroInterestAndEscrow(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  10000000,
		YearlyInterestRate: "0",
		MonthlyPayment:     1000000,
		EscrowPayment:      0,
		StartDate:          "1970-01-01",
	}

	got, err := service.CalculateLoanPaymentPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the loan payment plan.")
	}

	want := domain.LoanPaymentPlan{
		DurationMonths: 10,
	}

	if got.DurationMonths != want.DurationMonths {
		log.Fatalf("Expected a duration of %v months, but got %v.", want.DurationMonths, got.DurationMonths)
	}
}

func TestCalculateLoanPrincipalTooHigh(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  100000000001,
		YearlyInterestRate: "1",
		MonthlyPayment:     1,
		EscrowPayment:      0,
		StartDate:          "1970-01-01",
	}

	_, err := service.CalculateLoanPaymentPlan(input)
	if err == nil {
		log.Fatalf("Expected the loan calculation to fail due to starting principal being larger than the accepted amount (100000000000), but it didn't.")
	}
}

func TestCalculateLoanInterestTooHigh(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  100,
		YearlyInterestRate: "101",
		MonthlyPayment:     1,
		EscrowPayment:      0,
		StartDate:          "1970-01-01",
	}

	_, err := service.CalculateLoanPaymentPlan(input)
	if err == nil {
		log.Fatalf("Expected the loan calculation to fail due to the interest rate being higher than valid percent (100%% yearly), but it didn't.")
	}
}

func TestCalculateLoanMonthlyPaymentTooLow(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  10000000,
		YearlyInterestRate: "5",
		MonthlyPayment:     51666,
		EscrowPayment:      10000,
		StartDate:          "1970-01-01",
	}

	_, err := service.CalculateLoanPaymentPlan(input)
	if !strings.Contains(err.Error(), "not enough to cover interest and escrow payment") {
		log.Fatalf("Expected the loan calculation to fail due to the monthly payment (%v cents) not even covering interest and escrow, but it didn't.", input.MonthlyPayment)
	}
}

func TestCalculateLoanEscrowTooHigh(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  100,
		YearlyInterestRate: "1",
		MonthlyPayment:     1,
		EscrowPayment:      100000000001,
		StartDate:          "1970-01-01",
	}

	_, err := service.CalculateLoanPaymentPlan(input)
	if err == nil {
		log.Fatalf("Expected the loan calculation to fail due to escrow payment being higher than the valid amount (100000000000 cents), but it didn't.")
	}
}

func TestCalculateLoanInvalidDateFormat(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  100,
		YearlyInterestRate: "1",
		MonthlyPayment:     1,
		EscrowPayment:      100000000001,
		StartDate:          "01/01/1970",
	}

	_, err := service.CalculateLoanPaymentPlan(input)
	if err == nil {
		log.Fatalf("Expected the loan calculation to fail due to invalid start date format, but it didn't.")
	}
}

func TestSaveLoanPaymentPlan(t *testing.T) {
	mockUserID := uuid.Nil
	mockLoansRepo := &MockLoansRepo{
		SaveLoanPaymentPlanFunc: func(ctx context.Context, plan domain.LoanPaymentPlan) (db.Loan, error) {
			return db.Loan{
				DurationMonths:   int32(plan.DurationMonths),
				TotalPaid:        int32(plan.TotalPaid.Round(0).IntPart()),
				TotalExpenditure: int32(plan.TotalExpenditure.Round(0).IntPart()),
				CostOfCredit:     plan.CostOfCreditPercent.String(),
			}, nil
		},
	}
	service := NewLoansService(mockLoansRepo)
	ctx := context.Background()

	input := domain.SaveLoanInput{
		UserID:             mockUserID,
		LoanName:           "test",
		StartingPrincipal:  10000000,
		YearlyInterestRate: "5",
		MonthlyPayment:     900076,
		EscrowPayment:      10000,
		StartDate:          "1970-01-01",
	}

	want := db.Loan{
		DurationMonths:   12,
		TotalPaid:        10383416,
		TotalExpenditure: 383416,
		CostOfCredit:     "1.0383416398261762",
	}

	got, err := service.SaveLoanPaymentPlan(ctx, input)
	if err != nil {
		log.Fatalf("Error saving the loan payment plan: %v", err)
	}

	if want.DurationMonths != got.DurationMonths {
		log.Fatalf("Expected the duration in months saved on database (%v) to match the expected ones (%v), but they didn't.", got.DurationMonths, want.DurationMonths)
	}
	if want.TotalPaid != got.TotalPaid {
		log.Fatalf("Expected the total paid saved on database (%v cents) to match the expected one (%v cents), but it didn't.", got.TotalPaid, want.TotalPaid)
	}
	if want.TotalExpenditure != got.TotalExpenditure {
		log.Fatalf("Expected the total expenditure saved on database (%v cents) to match the expected one (%v cents), but it didn't.", got.TotalExpenditure, want.TotalExpenditure)
	}
	if want.CostOfCredit != got.CostOfCredit {
		log.Fatalf("Expected the cost of credit saved on database (%v%%) to match the expected one (%v%%), but it didn't.", got.CostOfCredit, want.CostOfCredit)
	}
}
func TestUpdateLoan(t *testing.T) {
	originalName := "Original Name"
	updatedName := "Updated Name"
	originalPrincipal := 10000
	updatedPrincipal := 15000
	originalInterest := "5"
	updatedInterest := "4"
	mockLoansRepo := &MockLoansRepo{
		GetLoanInitialDataFunc: func(ctx context.Context, planID uuid.UUID, userID uuid.UUID) (domain.UpdateLoanData, error) {
			return domain.UpdateLoanData{
				ID:   planID,
				Name: originalName,
				LoanData: domain.LoansInput{
					StartingPrincipal:  originalPrincipal,
					YearlyInterestRate: originalInterest,
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
	service := NewLoansService(mockLoansRepo)
	ctx := context.Background()

	input := domain.UpdateLoanInput{
		LoanName:           &updatedName,
		StartingPrincipal:  &updatedPrincipal,
		YearlyInterestRate: &updatedInterest,
	}

	got, err := service.UpdateLoan(ctx, input)
	if err != nil {
		log.Fatalf("Error updating the loan payment plan: %v", err)
	}

	want := db.Loan{
		Name:               updatedName,
		StartingPrincipal:  int32(updatedPrincipal),
		YearlyInterestRate: updatedInterest,
	}

	if want.Name != got.Name {
		log.Fatalf("Updated loan name returned (%v) did not match the expected one (%v).", got.Name, want.Name)
	}
	if want.StartingPrincipal != got.StartingPrincipal {
		log.Fatalf("Updated principal returned (%v cents) did not match the expected one (%v cents).", got.StartingPrincipal, want.StartingPrincipal)
	}
	if want.YearlyInterestRate != got.YearlyInterestRate {
		log.Fatalf("Updated interest rate returned (%v) did not match the expected one (%v).", got.YearlyInterestRate, want.YearlyInterestRate)
	}
}

func TestGetLoansByUserNoLoans(t *testing.T) {
	mockUserID := uuid.Nil
	mockLoansRepo := &MockLoansRepo{
		GetLoanPaymentPlansByUserFunc: func(ctx context.Context, userID uuid.UUID) ([]db.GetLoansByUserIDRow, error) {
			return nil, nil
		},
	}
	service := NewLoansService(mockLoansRepo)
	ctx := context.Background()

	got, err := service.GetLoansByUser(ctx, mockUserID)
	if err != nil {
		log.Fatalf("Expected list loans to return no error when no loans were found for user, but it did return an error: %v", err)
	}
	if len(got) > 0 {
		log.Fatalf("Expected list of loans to be null or empty, but it came with results.")
	}
}

func TestGetLoanNotFound(t *testing.T) {
	mockUserID := uuid.Nil
	mockLoanID := uuid.Nil
	mockLoansRepo := &MockLoansRepo{
		GetLoanByIDFunc: func(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) (domain.LoanPaymentPlan, error) {
			return domain.LoanPaymentPlan{}, fmt.Errorf("loan not found")
		},
	}
	service := NewLoansService(mockLoansRepo)
	ctx := context.Background()

	_, err := service.loansRepo.GetLoanByID(ctx, mockLoanID, mockUserID)
	if !strings.Contains(err.Error(), "loan not found") {
		log.Fatalf("Expected an error log due to the loan not being found, but the function didn't return it: %v", err)
	}
}
