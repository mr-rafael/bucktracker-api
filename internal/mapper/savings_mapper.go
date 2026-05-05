package mapper

import (
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/google/uuid"
)

func ToSavingsResponse(plan domain.SavingsPlan) dto.SavingsResponseParams {
	response := dto.SavingsResponseParams{}

	response.MonthlyInterestRate = multiplierToHighPrecisionPercent(plan.InterestMultiplierM)
	response.TotalInterestEarnings = int(plan.TotalInterestEarnings.Round(0).IntPart())
	response.TotalDeposited = int(plan.TotalDeposited.Round(0).IntPart())
	response.RateOfReturn = plan.RateOfReturn.String()
	response.InflationAdjustedROR = plan.InflationAdjustedROR.String()
	for _, status := range plan.Plan {
		response.Plan = append(response.Plan, dto.SavingsStatus{
			Date:         status.Date,
			Interest:     status.Interest,
			Tax:          status.Tax,
			Contribution: status.Contribution,
			Increase:     status.Increase,
			Capital:      status.Capital,
		})
	}
	return response
}

func ToSaveSavingsResponse(savings db.Saving) dto.SavingsSaveResponseParams {
	return dto.SavingsSaveResponseParams{
		ID:                    savings.ID.String(),
		Name:                  savings.Name,
		StartingCapital:       int(savings.StartingCapital),
		YearlyInterestRate:    savings.YearlyInflationRate.String,
		InterestRateType:      savings.InterestRateType,
		MonthlyContribution:   int(savings.MonthlyContribution),
		DurationYears:         int(savings.DurationYears),
		TaxRate:               savings.TaxRate,
		YearlyInflationRate:   savings.YearlyInflationRate.String,
		StartDate:             savings.StartDate.Time.Format(time.RFC3339),
		MonthlyInterestRate:   savings.MonthlyInterestRate,
		TotalInterestEarnings: int(savings.TotalInterestEarnings),
		RateOfReturn:          savings.RateOfReturn,
		InflationAdjustedROR:  savings.InflationAdjustedRor,
	}
}

func ToSavingsListResponse(rows []db.GetSavingsByUserIDRow) dto.SavingsListResponseParams {
	params := dto.SavingsListResponseParams{}
	for _, row := range rows {
		newRow := dto.SavingsInfo{
			ID:              row.ID.String(),
			Name:            row.Name,
			StartingCapital: int(row.StartingCapital),
		}
		params.Plans = append(params.Plans, newRow)
	}
	return params
}

func ToGetSavingsResponse(plan domain.SavingsPlan) dto.SavedSavingsResponseParams {
	originalParams := dto.OriginalSavingsData{
		StartingCapital:     plan.OriginalData.StartingCapital,
		YearlyInterestRate:  plan.OriginalData.YearlyInterestRate,
		InterestRateType:    plan.OriginalData.InterestRateType,
		MonthlyContribution: plan.OriginalData.MonthlyContribution,
		DurationYears:       plan.OriginalData.DurationYears,
		TaxRate:             plan.OriginalData.TaxRate,
		YearlyInflationRate: plan.OriginalData.YearlyInflationRate,
		StartDate:           plan.OriginalData.StartDate,
	}
	monthlyInterestRate := multiplierToHighPrecisionPercent(plan.InterestMultiplierM)
	calculatedParams := dto.CalculatedSavingsData{
		MonthlyInterestRate:   monthlyInterestRate,
		TotalInterestEarnings: int(plan.TotalInterestEarnings.Round(0).IntPart()),
		RateOfReturn:          plan.RateOfReturn.String(),
		InflationAdjustedROR:  plan.InflationAdjustedROR.String(),
	}
	params := dto.SavedSavingsResponseParams{
		ID:             plan.ID.String(),
		Name:           plan.Name,
		OriginalData:   originalParams,
		CalculatedData: calculatedParams,
	}
	for _, status := range plan.Plan {
		params.Plan = append(params.Plan, dto.SavingsStatus{
			Date:         status.Date,
			Interest:     status.Interest,
			Tax:          status.Tax,
			Contribution: status.Contribution,
			Increase:     status.Increase,
			Capital:      status.Capital,
		})
	}
	return params
}

func ToSavingsInput(input dto.SavingsRequestParams) domain.SavingsInput {
	savings := domain.SavingsInput{
		StartingCapital:     input.StartingCapital,
		YearlyInterestRate:  input.YearlyInterestRate,
		InterestRateType:    input.InterestRateType,
		MonthlyContribution: input.MonthlyContribution,
		DurationYears:       input.DurationYears,
		TaxRate:             input.TaxRate,
		YearlyInflationRate: input.YearlyInflationRate,
		StartDate:           input.StartDate,
	}

	return savings
}

func ToSaveSavingsInput(userId uuid.UUID, input dto.SavingsSaveRequestParams) domain.SaveSavingsInput {
	savings := domain.SaveSavingsInput{
		UserID:              userId,
		PlanName:            input.Name,
		StartingCapital:     input.StartingCapital,
		YearlyInterestRate:  input.YearlyInterestRate,
		InterestRateType:    input.InterestRateType,
		MonthlyContribution: input.MonthlyContribution,
		DurationYears:       input.DurationYears,
		TaxRate:             input.TaxRate,
		YearlyInflationRate: input.YearlyInflationRate,
		StartDate:           input.StartDate,
	}
	return savings
}
