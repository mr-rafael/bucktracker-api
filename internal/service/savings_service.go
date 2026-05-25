package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type SavingsInputError struct {
	Message string
}

func (e SavingsInputError) Error() string {
	return e.Message
}

type SavingsService struct {
	savingsRepo SavingsRepository
}

type SavingsRepository interface {
	SaveSavingsPlan(context.Context, domain.SavingsPlan) (db.Saving, error)
	GetSavingsPlansByUser(context.Context, uuid.UUID) ([]db.GetSavingsByUserIDRow, error)
	GetSavingsPlanByID(context.Context, uuid.UUID, uuid.UUID) (domain.SavingsPlan, error)
	GetSavingsInitialData(context.Context, uuid.UUID, uuid.UUID) (domain.UpdateSavingsData, error)
	UpdateSavings(context.Context, domain.SavingsPlan) (db.Saving, error)
	DeleteSavingsPlan(context.Context, uuid.UUID, uuid.UUID) error
}

func NewSavingsService(repo SavingsRepository) *SavingsService {
	return &SavingsService{savingsRepo: repo}
}

const minStartCapCents = "1"
const maxStartCapCents = "1000000000"
const minSavIntRate = "0.0001"
const maxSavIntRate = "1"
const minDurYears = "1"
const maxDurYears = "50"
const minMonthContrib = "0"
const maxMonthContrib = "1000000000"
const minTaxPercent = "0"
const maxTaxPercent = "100"

func (s *SavingsService) CalculateSavingsPlan(input domain.SavingsInput) (domain.SavingsPlan, error) {
	plan, err := initializeSavingsPlan(input, uuid.Nil, "")
	if err != nil {
		return domain.SavingsPlan{}, err
	}

	plan = calculateSavings(plan)

	return plan, nil
}

func (s *SavingsService) SaveSavingsPlan(ctx context.Context, input domain.SaveSavingsInput) (db.Saving, error) {
	plan, err := initializeSavingsPlan(toSavingsInput(input), input.UserID, input.PlanName)
	if err != nil {
		return db.Saving{}, err
	}

	plan = calculateSavings(plan)
	result, err := s.savingsRepo.SaveSavingsPlan(ctx, plan)
	if err != nil {
		return db.Saving{}, err
	}

	return result, nil
}

func (s *SavingsService) GetSavingsPlansByUser(ctx context.Context, input uuid.UUID) ([]db.GetSavingsByUserIDRow, error) {
	result, err := s.savingsRepo.GetSavingsPlansByUser(ctx, input)
	if err != nil {
		return []db.GetSavingsByUserIDRow{}, err
	}
	return result, nil
}

func (s *SavingsService) GetSavingsPlan(ctx context.Context, planID uuid.UUID, userID uuid.UUID) (domain.SavingsPlan, error) {
	result, err := s.savingsRepo.GetSavingsPlanByID(ctx, planID, userID)
	if err != nil {
		return domain.SavingsPlan{}, err
	}
	return result, nil
}

func (s *SavingsService) UpdateSavings(ctx context.Context, input domain.UpdateSavingsInput) (db.Saving, error) {
	originalData, err := s.savingsRepo.GetSavingsInitialData(ctx, input.ID, input.UserID)
	if err != nil {
		return db.Saving{}, fmt.Errorf("Savings plan not found.")
	}
	patchedData := patchSavingsFields(originalData, input)

	plan, err := initializeSavingsPlan(patchedData.SavingsData, input.UserID, patchedData.Name)
	if err != nil {
		return db.Saving{}, fmt.Errorf("failed to initialize the savings plan struct: %v", err)
	}
	plan = calculateSavings(plan)
	plan.ID = input.ID
	result, err := s.savingsRepo.UpdateSavings(ctx, plan)
	if err != nil {
		return db.Saving{}, err
	}

	return result, nil
}

func (s *SavingsService) DeleteSavingsPlan(ctx context.Context, planID uuid.UUID, userID uuid.UUID) error {
	return s.savingsRepo.DeleteSavingsPlan(ctx, planID, userID)
}

func calculateSavings(plan domain.SavingsPlan) domain.SavingsPlan {
	for i := 0; i < int(plan.DurationMonths.IntPart()); i++ {
		state := plan.PassMonth()
		state = plan.GenerateInterest(state)
		state = plan.Contribute(state)
		plan.Plan = append(plan.Plan, state)
	}
	plan.FinalCalculations()
	return plan
}

func initializeSavingsPlan(input domain.SavingsInput, userID uuid.UUID, name string) (domain.SavingsPlan, error) {
	plan := domain.SavingsPlan{}
	aHundred := decimal.NewFromInt(100)

	plan.OriginalData = input
	plan.UserID = userID
	plan.Name = name

	if len(input.InterestRateType) == 0 {
		plan.OriginalData.InterestRateType = "APY"
	}
	if len(input.TaxRate) == 0 {
		plan.OriginalData.TaxRate = "0"
	}
	if len(input.YearlyInflationRate) == 0 {
		plan.OriginalData.YearlyInflationRate = "0"
	}

	startingCapital := decimal.NewFromInt(int64(input.StartingCapital))
	if !decimalIsBetween(startingCapital, minStartCapCents, maxStartCapCents) {
		return domain.SavingsPlan{}, SavingsInputError{Message: fmt.Sprintf("invalid starting amount '%v'. the valid range is 0.01-1,000,000,000", startingCapital.Div(aHundred).Round(2))}
	}
	plan.StartingCapital = startingCapital
	plan.CurrentCapital = startingCapital
	plan.TotalDeposited = startingCapital

	monthlyInterestRate, err := toMonthlyInterestMultiplier(input.YearlyInterestRate, input.InterestRateType)
	if err != nil {
		return domain.SavingsPlan{}, SavingsInputError{Message: fmt.Sprintf("invalid interest rate: %v", input.YearlyInterestRate)}
	}
	if !decimalIsBetween(monthlyInterestRate, minSavIntRate, maxSavIntRate) {
		return domain.SavingsPlan{}, SavingsInputError{Message: "invalid interest rate. The valid range is 0.001-1"}
	}
	plan.InterestMultiplierM = monthlyInterestRate

	durationMonths := decimal.NewFromInt(int64(input.DurationYears)).Mul(decimal.NewFromInt(12))
	if !decimalIsBetween(decimal.NewFromInt32(int32(input.DurationYears)), minDurYears, maxDurYears) {
		return domain.SavingsPlan{}, SavingsInputError{Message: fmt.Sprintf("invalid plan duration. The valid range is %v-%v", minDurYears, maxDurYears)}
	}
	plan.DurationMonths = durationMonths

	monthlyContribution := decimal.NewFromInt(int64(input.MonthlyContribution))
	if !decimalIsBetween(monthlyContribution, minMonthContrib, maxMonthContrib) {
		return domain.SavingsPlan{}, SavingsInputError{Message: "invalid monthly contribution amount. The valid range is 0-1,000,000,000"}
	}
	plan.MonthlyContribution = monthlyContribution

	tax, err := toTaxMultiplier(input.TaxRate)
	if err != nil {
		return domain.SavingsPlan{}, fmt.Errorf("invalid tax rate %v", input.TaxRate)
	}
	if !stringNumberBetween(input.TaxRate, minTaxPercent, maxTaxPercent) {
		return domain.SavingsPlan{}, SavingsInputError{Message: fmt.Sprintf("invalid tax rate '%v'. The valid range is %v-%v%%.", input.TaxRate, minTaxPercent, maxTaxPercent)}
	}
	plan.TaxMultiplierM = tax

	inflation, err := toInflationMultiplier(input.YearlyInflationRate)
	if err != nil {
		return domain.SavingsPlan{}, SavingsInputError{Message: fmt.Sprintf("invalid inflation rate %v", input.YearlyInflationRate)}
	}
	plan.InflationMultiplierY = inflation

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return domain.SavingsPlan{}, SavingsInputError{Message: fmt.Sprintf("invalid start date: %v", input.StartDate)}
	}
	plan.Date = startDate

	return plan, nil
}

func toSavingsInput(input domain.SaveSavingsInput) domain.SavingsInput {
	return domain.SavingsInput{
		StartingCapital:     input.StartingCapital,
		YearlyInterestRate:  input.YearlyInterestRate,
		InterestRateType:    input.InterestRateType,
		MonthlyContribution: input.MonthlyContribution,
		DurationYears:       input.DurationYears,
		TaxRate:             input.TaxRate,
		YearlyInflationRate: input.YearlyInflationRate,
		StartDate:           input.StartDate,
	}
}

func patchSavingsFields(savingsData domain.UpdateSavingsData, patchData domain.UpdateSavingsInput) domain.UpdateSavingsData {
	if patchData.PlanName != nil {
		savingsData.Name = *patchData.PlanName
	}
	if patchData.StartingCapital != nil {
		savingsData.SavingsData.StartingCapital = *patchData.StartingCapital
	}
	if patchData.YearlyInterestRate != nil {
		savingsData.SavingsData.YearlyInterestRate = *patchData.YearlyInterestRate
	}
	if patchData.InterestRateType != nil {
		savingsData.SavingsData.InterestRateType = *patchData.InterestRateType
	}
	if patchData.MonthlyContribution != nil {
		savingsData.SavingsData.MonthlyContribution = *patchData.MonthlyContribution
	}
	if patchData.TaxRate != nil {
		savingsData.SavingsData.TaxRate = *patchData.TaxRate
	}
	if patchData.YearlyInflationRate != nil {
		savingsData.SavingsData.YearlyInflationRate = *patchData.YearlyInflationRate
	}
	if patchData.StartDate != nil {
		savingsData.SavingsData.StartDate = *patchData.StartDate
	}
	return savingsData
}
