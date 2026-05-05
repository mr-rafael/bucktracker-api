package service

import (
	"log"
	"testing"

	"github.com/Mr-Rafael/finance-calculator/internal/domain"
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
		log.Fatalf("Expected an error when calculating a loan for an invalid duration (0 years), but none was returned.")
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
