package mapper

import (
	"testing"

	"github.com/Mr-Rafael/bucktracker-api/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestToGetSavingsResponseIncludesTotalDeposited(t *testing.T) {
	planID := uuid.MustParse("a5b140fd-583c-4edb-8668-e7ee986d2a37")
	plan := domain.SavingsPlan{
		ID:   planID,
		Name: "test",
		OriginalData: domain.SavingsInput{
			StartingCapital:     700000,
			YearlyInterestRate:  "4.75",
			InterestRateType:    "APY",
			MonthlyContribution: 15000,
			DurationYears:       1,
			TaxRate:             "0",
			YearlyInflationRate: "0",
			StartDate:           "2026-02-01",
		},
		InterestMultiplierM:   decimal.RequireFromString("0.003874684992"),
		TotalInterestEarnings: decimal.NewFromInt(37136),
		TotalDeposited:        decimal.NewFromInt(880000),
		RateOfReturn:          decimal.RequireFromString("4.22"),
		InflationAdjustedROR:  decimal.RequireFromString("4.22"),
	}

	got := ToGetSavingsResponse(plan)

	require.Equal(t, 880000, got.CalculatedData.TotalDeposited)
	require.Equal(t, 37136, got.CalculatedData.TotalInterestEarnings)
	require.Equal(t, "4.22", got.CalculatedData.RateOfReturn)
}
