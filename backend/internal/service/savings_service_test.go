package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/Mr-Rafael/bucktracker-api/internal/db"
	"github.com/Mr-Rafael/bucktracker-api/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestCalculateSavingsPlanDurationTooShort(t *testing.T) {
	mockSavingsRepo := &MockSavingsRepo{}
	service := NewSavingsService(mockSavingsRepo)

	input := domain.SavingsInput{
		StartingCapital:     10000000,
		YearlyInterestRate:  "5",
		InterestRateType:    "APY",
		MonthlyContribution: 10000,
		DurationYears:       0,
		TaxRate:             "0",
		YearlyInflationRate: "0",
		StartDate:           "1970-01-01",
	}

	_, err := service.CalculateSavingsPlan(input)
	if err == nil {
		log.Fatalf("Expected an error when calculating savings for an invalid duration (0 years), but none was returned.")
	}
}

func TestCalculateSavingsPlanMinDuration(t *testing.T) {
	mockSavingsRepo := &MockSavingsRepo{}
	service := NewSavingsService(mockSavingsRepo)

	input := domain.SavingsInput{
		StartingCapital:     10000000,
		YearlyInterestRate:  "5",
		InterestRateType:    "APY",
		MonthlyContribution: 10000,
		DurationYears:       1,
		TaxRate:             "0",
		YearlyInflationRate: "0",
		StartDate:           "1970-01-01",
	}

	got, err := service.CalculateSavingsPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the Savings Plan: %v", err)
	}

	rateOfReturn, _ := decimal.NewFromString("4.97")
	want := domain.SavingsPlan{
		CurrentCapital:        decimal.NewFromInt32(10622726),
		TotalDeposited:        decimal.NewFromInt32(10120000),
		TotalInterestEarnings: decimal.NewFromInt32(502726),
		RateOfReturn:          rateOfReturn,
	}

	if got.CurrentCapital.Round(0).Compare(want.CurrentCapital) != 0 {
		log.Fatalf("Expected calculated principal (%v) to match expected principal (%v), but it didn't.", got.CurrentCapital.Round(0), want.CurrentCapital)
	}

	if got.TotalDeposited.Round(0).Compare(want.TotalDeposited) != 0 {
		log.Fatalf("Expected total deposited (%v) to match expected total (%v), but it didn't.", got.TotalDeposited.Round(0), want.TotalDeposited)
	}

	if got.TotalInterestEarnings.Round(0).Compare(want.TotalInterestEarnings) != 0 {
		log.Fatalf("Expected total interest earnings (%v) to match expected earnings (%v), but they didn't.", got.TotalInterestEarnings.Round(0), want.TotalInterestEarnings)
	}

	if got.RateOfReturn.Round(2).Compare(want.RateOfReturn) != 0 {
		log.Fatalf("Expected rate of return (%v) to match expected rate (%v), but it didn't.", got.RateOfReturn.Round(2), want.RateOfReturn)
	}
}

func TestCalculateSavingsWithTaxAndInflation(t *testing.T) {
	mockSavingsRepo := &MockSavingsRepo{}
	service := NewSavingsService(mockSavingsRepo)

	input := domain.SavingsInput{
		StartingCapital:     10000000,
		YearlyInterestRate:  "5",
		InterestRateType:    "APY",
		MonthlyContribution: 10000,
		DurationYears:       5,
		TaxRate:             "5",
		YearlyInflationRate: "6",
		StartDate:           "1970-01-01",
	}

	got, err := service.CalculateSavingsPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the Savings Plan: %v", err)
	}

	rateOfReturn, _ := decimal.NewFromString("25.3")
	want := domain.SavingsPlan{
		CurrentCapital:        decimal.NewFromInt32(13282311),
		TotalDeposited:        decimal.NewFromInt32(10600000),
		TotalInterestEarnings: decimal.NewFromInt32(2682311),
		RateOfReturn:          rateOfReturn,
	}

	if got.CurrentCapital.Round(0).Compare(want.CurrentCapital) != 0 {
		log.Fatalf("Expected calculated principal (%v) to match expected principal (%v) at the end of loan, but it didn't.", got.CurrentCapital.Round(0), want.CurrentCapital)
	}

	if got.TotalDeposited.Round(0).Compare(want.TotalDeposited) != 0 {
		log.Fatalf("Expected total deposited (%v) to match expected total (%v), but it didn't.", got.TotalDeposited.Round(0), want.TotalDeposited)
	}

	if got.TotalInterestEarnings.Round(0).Compare(want.TotalInterestEarnings) != 0 {
		log.Fatalf("Expected total interest earnings (%v) to match expected earnings (%v), but they didn't.", got.TotalInterestEarnings.Round(0), want.TotalInterestEarnings)
	}

	if got.RateOfReturn.Round(2).Compare(want.RateOfReturn) != 0 {
		log.Fatalf("Expected rate of return (%v) to match expected rate (%v), but it didn't.", got.RateOfReturn.Round(2), want.RateOfReturn)
	}
}

func TestCalculateSavingsMaxDuration(t *testing.T) {
	mockSavingsRepo := &MockSavingsRepo{}
	service := NewSavingsService(mockSavingsRepo)

	input := domain.SavingsInput{
		StartingCapital:     100000,
		YearlyInterestRate:  "5",
		InterestRateType:    "APY",
		MonthlyContribution: 1000,
		DurationYears:       50,
		TaxRate:             "5",
		YearlyInflationRate: "6",
		StartDate:           "1970-01-01",
	}

	got, err := service.CalculateSavingsPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the Savings Plan: %v", err)
	}

	rateOfReturn, _ := decimal.NewFromString("382.88")
	want := domain.SavingsPlan{
		CurrentCapital:        decimal.NewFromInt32(3380157),
		TotalDeposited:        decimal.NewFromInt32(700000),
		TotalInterestEarnings: decimal.NewFromInt32(2680157),
		RateOfReturn:          rateOfReturn,
	}

	if got.CurrentCapital.Round(0).Compare(want.CurrentCapital) != 0 {
		log.Fatalf("Expected calculated principal (%v) to match expected principal (%v) at the end of loan, but it didn't.", got.CurrentCapital.Round(0), want.CurrentCapital)
	}

	if got.TotalDeposited.Round(0).Compare(want.TotalDeposited) != 0 {
		log.Fatalf("Expected total deposited (%v) to match expected total (%v), but it didn't.", got.TotalDeposited.Round(0), want.TotalDeposited)
	}

	if got.TotalInterestEarnings.Round(0).Compare(want.TotalInterestEarnings) != 0 {
		log.Fatalf("Expected total interest earnings (%v) to match expected earnings (%v), but they didn't.", got.TotalInterestEarnings.Round(0), want.TotalInterestEarnings)
	}

	if got.RateOfReturn.Round(2).Compare(want.RateOfReturn) != 0 {
		log.Fatalf("Expected rate of return (%v) to match expected rate (%v), but it didn't.", got.RateOfReturn.Round(2), want.RateOfReturn)
	}
}

func TestCalculateSavingsTermTooLong(t *testing.T) {
	mockSavingsRepo := &MockSavingsRepo{}
	service := NewSavingsService(mockSavingsRepo)

	input := domain.SavingsInput{
		StartingCapital:     100000,
		YearlyInterestRate:  "5",
		InterestRateType:    "APY",
		MonthlyContribution: 1000,
		DurationYears:       51,
		TaxRate:             "5",
		YearlyInflationRate: "6",
		StartDate:           "1970-01-01",
	}

	_, err := service.CalculateSavingsPlan(input)
	if err == nil {
		log.Fatalf("Expected an error when calculating savings for an invalid duration (50+ years), but none was returned.")
	}
}

func TestCalculateStartingPrincipalTooLow(t *testing.T) {
	mockSavingsRepo := &MockSavingsRepo{}
	service := NewSavingsService(mockSavingsRepo)

	input := domain.SavingsInput{
		StartingCapital:     0,
		YearlyInterestRate:  "5",
		InterestRateType:    "APY",
		MonthlyContribution: 1000,
		DurationYears:       51,
		TaxRate:             "5",
		YearlyInflationRate: "6",
		StartDate:           "1970-01-01",
	}

	_, err := service.CalculateSavingsPlan(input)
	if err == nil {
		log.Fatalf("Expected an error when calculating savings for an invalid starting principal (<1 cent), but none was returned.")
	}
}

func TestCalculateSavingsAPRInterestRate(t *testing.T) {
	mockSavingsRepo := &MockSavingsRepo{}
	service := NewSavingsService(mockSavingsRepo)

	input := domain.SavingsInput{
		StartingCapital:     500000,
		YearlyInterestRate:  "3.75",
		InterestRateType:    "APR",
		MonthlyContribution: 500000,
		DurationYears:       10,
		TaxRate:             "5",
		YearlyInflationRate: "0",
		StartDate:           "1970-01-01",
	}

	got, err := service.CalculateSavingsPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the Savings Plan: %v", err)
	}

	rateOfReturn, _ := decimal.NewFromString("20.11")
	want := domain.SavingsPlan{
		CurrentCapital:        decimal.NewFromInt32(72664943),
		TotalDeposited:        decimal.NewFromInt32(60500000),
		TotalInterestEarnings: decimal.NewFromInt32(12164943),
		RateOfReturn:          rateOfReturn,
	}

	if got.CurrentCapital.Round(0).Compare(want.CurrentCapital) != 0 {
		log.Fatalf("Expected calculated principal (%v) to match expected principal (%v) at the end of loan, but it didn't.", got.CurrentCapital.Round(0), want.CurrentCapital)
	}

	if got.TotalDeposited.Round(0).Compare(want.TotalDeposited) != 0 {
		log.Fatalf("Expected total deposited (%v) to match expected total (%v), but it didn't.", got.TotalDeposited.Round(0), want.TotalDeposited)
	}

	if got.TotalInterestEarnings.Round(0).Compare(want.TotalInterestEarnings) != 0 {
		log.Fatalf("Expected total interest earnings (%v) to match expected earnings (%v), but they didn't.", got.TotalInterestEarnings.Round(0), want.TotalInterestEarnings)
	}

	if got.RateOfReturn.Round(2).Compare(want.RateOfReturn) != 0 {
		log.Fatalf("Expected rate of return (%v) to match expected rate (%v), but it didn't.", got.RateOfReturn.Round(2), want.RateOfReturn)
	}
}

func TestCalculateSavingsZeroContribution(t *testing.T) {
	mockSavingsRepo := &MockSavingsRepo{}
	service := NewSavingsService(mockSavingsRepo)

	input := domain.SavingsInput{
		StartingCapital:     500000,
		YearlyInterestRate:  "3.75",
		InterestRateType:    "APY",
		MonthlyContribution: 0,
		DurationYears:       3,
		TaxRate:             "12",
		YearlyInflationRate: "0",
		StartDate:           "1970-01-01",
	}

	got, err := service.CalculateSavingsPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the Savings Plan: %v", err)
	}

	rateOfReturn, _ := decimal.NewFromString("10.21")
	want := domain.SavingsPlan{
		CurrentCapital:        decimal.NewFromInt32(551044),
		TotalDeposited:        decimal.NewFromInt32(500000),
		TotalInterestEarnings: decimal.NewFromInt32(51044),
		RateOfReturn:          rateOfReturn,
	}

	if got.CurrentCapital.Round(0).Compare(want.CurrentCapital) != 0 {
		log.Fatalf("Expected calculated principal (%v) to match expected principal (%v) at the end of loan, but it didn't.", got.CurrentCapital.Round(0), want.CurrentCapital)
	}

	if got.TotalDeposited.Round(0).Compare(want.TotalDeposited) != 0 {
		log.Fatalf("Expected total deposited (%v) to match expected total (%v), but it didn't.", got.TotalDeposited.Round(0), want.TotalDeposited)
	}

	if got.TotalInterestEarnings.Round(0).Compare(want.TotalInterestEarnings) != 0 {
		log.Fatalf("Expected total interest earnings (%v) to match expected earnings (%v), but they didn't.", got.TotalInterestEarnings.Round(0), want.TotalInterestEarnings)
	}

	if got.RateOfReturn.Round(2).Compare(want.RateOfReturn) != 0 {
		log.Fatalf("Expected rate of return (%v) to match expected rate (%v), but it didn't.", got.RateOfReturn.Round(2), want.RateOfReturn)
	}
}

func TestSaveSavingsPlan(t *testing.T) {
	mockUserID := uuid.Nil
	mockSavingsRepo := &MockSavingsRepo{
		SaveSavingsPlanFunc: func(ctx context.Context, plan domain.SavingsPlan) (db.Saving, error) {
			return db.Saving{
				TotalInterestEarnings: 2347213,
				TotalDeposited:        68000000,
				RateOfReturn:          "3.45",
			}, nil
		},
	}
	service := NewSavingsService(mockSavingsRepo)
	ctx := context.Background()

	input := domain.SaveSavingsInput{
		UserID:              mockUserID,
		PlanName:            "test",
		StartingCapital:     50000000,
		YearlyInterestRate:  "4.25",
		InterestRateType:    "APY",
		MonthlyContribution: 1500000,
		DurationYears:       1,
		TaxRate:             "5",
		YearlyInflationRate: "6",
		StartDate:           "1970-01-01",
	}

	got, err := service.SaveSavingsPlan(ctx, input)
	if err != nil {
		log.Fatalf("Error calculating the Savings Plan: %v", err)
	}

	want := db.Saving{
		TotalInterestEarnings: 2347213,
		TotalDeposited:        68000000,
		RateOfReturn:          "3.45",
	}

	if got.TotalInterestEarnings != want.TotalInterestEarnings {
		log.Fatalf("Expected saved total deposited (%v) to match expected total deposited (%v), but it didn't.", got.TotalInterestEarnings, want.TotalInterestEarnings)
	}

	if got.TotalDeposited != want.TotalDeposited {
		log.Fatalf("Expected saved interest earnings (%v) to match expected earnings (%v), but they didn't.", got.TotalDeposited, want.TotalDeposited)
	}

	if got.RateOfReturn != want.RateOfReturn {
		log.Fatalf("Expected saved rate of return (%v) to match expected rate (%v), but it didn't.", got.RateOfReturn, want.RateOfReturn)
	}
}

func TestUpdateSavings(t *testing.T) {
	originalName := "Original Name"
	updatedName := "Updated Name"
	originalCapital := 400000
	updatedCapital := 500000
	originalInterest := "4.2"
	updatedInterest := "4.75"
	mockSavingsRepo := &MockSavingsRepo{
		GetSavingsInitialDataFunc: func(ctx context.Context, planID uuid.UUID, userID uuid.UUID) (domain.UpdateSavingsData, error) {
			return domain.UpdateSavingsData{
				ID:   planID,
				Name: originalName,
				SavingsData: domain.SavingsInput{
					StartingCapital:     originalCapital,
					YearlyInterestRate:  originalInterest,
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
	service := NewSavingsService(mockSavingsRepo)
	ctx := context.Background()

	input := domain.UpdateSavingsInput{
		PlanName:           &updatedName,
		StartingCapital:    &updatedCapital,
		YearlyInterestRate: &updatedInterest,
	}

	got, err := service.UpdateSavings(ctx, input)
	if err != nil {
		log.Fatalf("Error updating the savings plan: %v", err)
	}

	want := db.Saving{
		Name:               updatedName,
		StartingCapital:    int32(updatedCapital),
		YearlyInterestRate: updatedInterest,
	}

	if want.Name != got.Name {
		log.Fatalf("Updated loan name returned (%v) did not match the expected one (%v).", got.Name, want.Name)
	}
	if want.StartingCapital != got.StartingCapital {
		log.Fatalf("Updated capital returned (%v cents) did not match the expected one (%v cents).", got.StartingCapital, want.StartingCapital)
	}
	if want.YearlyInterestRate != got.YearlyInterestRate {
		log.Fatalf("Updated interest rate returned (%v) did not match the expected one (%v).", got.YearlyInterestRate, want.YearlyInterestRate)
	}
}

func TestGetSavingsPlansByUserNoPlans(t *testing.T) {
	mockUserID := uuid.Nil
	mockSavingsRepo := &MockSavingsRepo{
		GetSavingsPlansByUserFunc: func(ctx context.Context, id uuid.UUID) ([]db.GetSavingsByUserIDRow, error) {
			return nil, nil
		},
	}
	service := NewSavingsService(mockSavingsRepo)
	ctx := context.Background()

	got, err := service.GetSavingsPlansByUser(ctx, mockUserID)
	if err != nil {
		log.Fatalf("Error fetching savings plans for user: %v", err)
	}

	if len(got) > 0 {
		log.Fatalf("Expected the list of savings plan rows to come back empty, but it didn't.")
	}
}

func TestGetSavingsPlanNotFound(t *testing.T) {
	mockUserID := uuid.Nil
	mockSavingsRepo := &MockSavingsRepo{
		GetSavingsPlansByUserFunc: func(ctx context.Context, id uuid.UUID) ([]db.GetSavingsByUserIDRow, error) {
			return nil, fmt.Errorf("savings plan not found")
		},
	}
	service := NewSavingsService(mockSavingsRepo)
	ctx := context.Background()

	_, err := service.GetSavingsPlansByUser(ctx, mockUserID)
	if !strings.Contains(err.Error(), "savings plan not found") {
		log.Fatalf("Expected an error log due to the savings plan not being found, but the function didn't return it: %v", err)
	}
}
