package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Mr-Rafael/bucktracker-api/internal/db"
	"github.com/Mr-Rafael/bucktracker-api/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type SavingsRepo struct {
	queries *db.Queries
}

func NewSavingsRepo(queries *db.Queries) *SavingsRepo {
	return &SavingsRepo{queries: queries}
}

func (r *SavingsRepo) SaveSavingsPlan(ctx context.Context, plan domain.SavingsPlan) (db.Saving, error) {
	savingsParams, err := toSavingsInsertQueryParams(plan)
	if err != nil {
		return db.Saving{}, fmt.Errorf("Error preparing params for insert query: %v", err)
	}

	queryResult, err := r.queries.CreateSavings(ctx, savingsParams)
	if err != nil {
		return db.Saving{}, fmt.Errorf("Failed to save to database: %v", err)
	}

	for _, status := range plan.Plan {
		_, err := r.queries.CreateSavingsState(ctx, toSavingsStateInsertParams(status, queryResult.ID))
		if err != nil {
			return db.Saving{}, fmt.Errorf("Failed to save savings status to database: %v", err)
		}
	}
	return queryResult, nil
}

func (r *SavingsRepo) GetSavingsPlansByUser(ctx context.Context, userID uuid.UUID) ([]db.GetSavingsByUserIDRow, error) {
	queryUserID := pgtype.UUID{
		Bytes: userID,
		Valid: true,
	}

	result, err := r.queries.GetSavingsByUserID(ctx, queryUserID)
	if err != nil {
		return []db.GetSavingsByUserIDRow{}, fmt.Errorf("failed to fetch user's savings plans: %v", err)
	}
	return result, nil
}

func (r *SavingsRepo) GetSavingsPlanByID(ctx context.Context, planID uuid.UUID, userID uuid.UUID) (domain.SavingsPlan, error) {
	querySavingsID := pgtype.UUID{
		Bytes: planID,
		Valid: true,
	}

	savingsQueryResult, err := r.queries.GetSavings(ctx, toSavingsGetParams(planID, userID))
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("failed to fetch savings plan from database: %v", err)
	}
	plan, err := toSavingsPlan(savingsQueryResult)

	statesQueryResult, err := r.queries.GetSavingsStatesBySavingsID(ctx, querySavingsID)
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("failed to fetch savings plan rows from database: %v", err)
	}
	for _, state := range statesQueryResult {
		plan.Plan = append(plan.Plan, domain.SavingsStatus{
			Date:         state.Date.Time,
			Interest:     int(state.Interest),
			Tax:          int(state.Tax),
			Contribution: int(state.Contribution),
			Increase:     int(state.Increase),
			Capital:      int(state.Capital),
		})
	}

	return plan, nil
}

func (r *SavingsRepo) GetSavingsInitialData(ctx context.Context, savingsID uuid.UUID, userID uuid.UUID) (domain.UpdateSavingsData, error) {
	savingsQueryResult, err := r.queries.GetSavingsInitialData(ctx, toInitialSavingsDataGetParams(savingsID, userID))
	if err != nil {
		return domain.UpdateSavingsData{}, fmt.Errorf("failed to fetch savings plan original data: %v", err)
	}
	savingsInput := domain.SavingsInput{
		StartingCapital:     int(savingsQueryResult.StartingCapital),
		YearlyInterestRate:  savingsQueryResult.YearlyInterestRate,
		InterestRateType:    savingsQueryResult.InterestRateType,
		MonthlyContribution: int(savingsQueryResult.MonthlyContribution),
		DurationYears:       int(savingsQueryResult.DurationYears),
		TaxRate:             savingsQueryResult.TaxRate,
		YearlyInflationRate: savingsQueryResult.YearlyInflationRate.String,
		StartDate:           savingsQueryResult.StartDate.Time.Format("2006-01-02"),
	}
	savingsData := domain.UpdateSavingsData{
		ID:          savingsID,
		Name:        savingsQueryResult.Name,
		SavingsData: savingsInput,
	}

	return savingsData, nil
}

func (r *SavingsRepo) UpdateSavings(ctx context.Context, plan domain.SavingsPlan) (db.Saving, error) {
	savingsParams, err := toSavingsUpdateQueryParams(plan)
	if err != nil {
		return db.Saving{}, fmt.Errorf("Error preparing params for update query: %v", err)
	}

	queryResult, err := r.queries.UpdateSavings(ctx, savingsParams)
	if err != nil {
		return db.Saving{}, fmt.Errorf("Failed to update savings plan on database: %v", err)
	}

	err = r.queries.DeleteSavingsStatesBySavingsID(ctx, savingsParams.ID)
	if err != nil {
		return db.Saving{}, fmt.Errorf("Error deleting old savings plan data: %v", err)
	}

	for _, status := range plan.Plan {
		_, err := r.queries.CreateSavingsState(ctx, toSavingsStateInsertParams(status, queryResult.ID))
		if err != nil {
			return db.Saving{}, fmt.Errorf("Failed to update savings status in database: %v", err)
		}
	}
	return queryResult, nil
}

func (r *SavingsRepo) DeleteSavingsPlan(ctx context.Context, planID uuid.UUID, userID uuid.UUID) error {
	rows, err := r.queries.DeleteSavings(ctx, db.DeleteSavingsParams(toSavingsGetParams(planID, userID)))
	if err != nil || rows <= 0 {
		return fmt.Errorf("Not found.")
	}
	return nil
}

func toSavingsInsertQueryParams(plan domain.SavingsPlan) (db.CreateSavingsParams, error) {
	startDate, err := time.Parse("2006-01-02", plan.OriginalData.StartDate)
	if err != nil {
		return db.CreateSavingsParams{}, err
	}
	return db.CreateSavingsParams{
		UserID: pgtype.UUID{
			Bytes: plan.UserID,
			Valid: true,
		},
		Name:                plan.Name,
		StartingCapital:     int32(plan.OriginalData.StartingCapital),
		YearlyInterestRate:  plan.OriginalData.YearlyInterestRate,
		InterestRateType:    plan.OriginalData.InterestRateType,
		MonthlyContribution: int32(plan.OriginalData.MonthlyContribution),
		DurationYears:       int32(plan.OriginalData.DurationYears),
		TaxRate:             plan.OriginalData.TaxRate,
		YearlyInflationRate: pgtype.Text{
			String: plan.OriginalData.YearlyInflationRate,
			Valid:  true,
		},
		StartDate: pgtype.Date{
			Time:  startDate,
			Valid: true,
		},
		MonthlyInterestRate:   multiplierToPercent(plan.InterestMultiplierM),
		TotalDeposited:        int32(plan.TotalDeposited.Round(0).IntPart()),
		TotalInterestEarnings: int32(plan.TotalInterestEarnings.Round(0).IntPart()),
		RateOfReturn:          plan.RateOfReturn.String(),
		InflationAdjustedRor:  plan.InflationAdjustedROR.String(),
	}, nil
}

func toSavingsStateInsertParams(status domain.SavingsStatus, savingsID pgtype.UUID) db.CreateSavingsStateParams {
	params := db.CreateSavingsStateParams{
		SavingsID: savingsID,
		Date: pgtype.Date{
			Time:  status.Date,
			Valid: true,
		},
		Interest:     int32(status.Interest),
		Tax:          int32(status.Tax),
		Contribution: int32(status.Contribution),
		Increase:     int32(status.Increase),
		Capital:      int32(status.Capital),
	}
	return params
}

func toSavingsUpdateQueryParams(plan domain.SavingsPlan) (db.UpdateSavingsParams, error) {
	startDate, err := time.Parse("2006-01-02", plan.OriginalData.StartDate)
	if err != nil {
		return db.UpdateSavingsParams{}, err
	}
	return db.UpdateSavingsParams{
		ID: pgtype.UUID{
			Bytes: plan.ID,
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: plan.UserID,
			Valid: true,
		},
		Name:                plan.Name,
		StartingCapital:     int32(plan.OriginalData.StartingCapital),
		YearlyInterestRate:  plan.OriginalData.YearlyInterestRate,
		InterestRateType:    plan.OriginalData.InterestRateType,
		MonthlyContribution: int32(plan.MonthlyContribution.Round(0).IntPart()),
		DurationYears:       int32(plan.DurationMonths.Div(decimal.NewFromInt32(12)).IntPart()),
		TaxRate:             plan.OriginalData.TaxRate,
		YearlyInflationRate: pgtype.Text{
			String: plan.OriginalData.YearlyInflationRate,
			Valid:  true,
		},
		StartDate: pgtype.Date{
			Time:  startDate,
			Valid: true,
		},
		MonthlyInterestRate:   multiplierToPercent(plan.InterestMultiplierM),
		TotalDeposited:        int32(plan.TotalDeposited.Round(0).IntPart()),
		TotalInterestEarnings: int32(plan.TotalInterestEarnings.Round(0).IntPart()),
		RateOfReturn:          plan.RateOfReturn.String(),
		InflationAdjustedRor:  plan.InflationAdjustedROR.String(),
	}, nil
}

func multiplierToPercent(mult decimal.Decimal) string {
	oneHundred := decimal.NewFromInt(100)
	return mult.Mul(oneHundred).String()
}

func percentToMultiplier(p string) decimal.Decimal {
	oneHundred := decimal.NewFromInt(100)
	decimalP, err := decimal.NewFromString(p)
	if err != nil {
		return decimal.Zero
	}
	return decimalP.Div(oneHundred)
}

func toSavingsPlan(queryResult db.Saving) (domain.SavingsPlan, error) {
	originalPlanData := domain.SavingsInput{
		StartingCapital:     int(queryResult.StartingCapital),
		YearlyInterestRate:  queryResult.YearlyInterestRate,
		InterestRateType:    queryResult.InterestRateType,
		MonthlyContribution: int(queryResult.MonthlyContribution),
		DurationYears:       int(queryResult.DurationYears),
		TaxRate:             queryResult.TaxRate,
		YearlyInflationRate: queryResult.YearlyInflationRate.String,
		StartDate:           queryResult.StartDate.Time.Format("2006-01-02"),
	}
	rateOfReturn, err := decimal.NewFromString(queryResult.RateOfReturn)
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("corrupted rate of return data for savings plan: %v", err)
	}
	inflationAdjustedReturn, err := decimal.NewFromString(queryResult.InflationAdjustedRor)
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("corrupted inflation rate of return data for savings plan: %v", err)
	}
	plan := domain.SavingsPlan{
		ID:                    queryResult.ID.Bytes,
		UserID:                queryResult.UserID.Bytes,
		Name:                  queryResult.Name,
		OriginalData:          originalPlanData,
		StartingCapital:       decimal.NewFromInt32(queryResult.StartingCapital),
		MonthlyContribution:   decimal.NewFromInt32(queryResult.MonthlyContribution),
		DurationMonths:        decimal.NewFromInt32(queryResult.DurationYears).Mul(decimal.NewFromInt(12)),
		InterestMultiplierM:   percentToMultiplier(queryResult.MonthlyInterestRate),
		TotalInterestEarnings: decimal.NewFromInt32(queryResult.TotalInterestEarnings),
		TotalDeposited:        decimal.NewFromInt32(queryResult.TotalDeposited),
		RateOfReturn:          rateOfReturn,
		InflationAdjustedROR:  inflationAdjustedReturn,
	}

	return plan, nil
}

func toSavingsGetParams(savingsID uuid.UUID, userID uuid.UUID) db.GetSavingsParams {
	return db.GetSavingsParams{
		ID: pgtype.UUID{
			Bytes: savingsID,
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: userID,
			Valid: true,
		},
	}
}

func toInitialSavingsDataGetParams(savingsID uuid.UUID, userID uuid.UUID) db.GetSavingsInitialDataParams {
	return db.GetSavingsInitialDataParams{
		ID: pgtype.UUID{
			Bytes: savingsID,
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: userID,
			Valid: true,
		},
	}
}
