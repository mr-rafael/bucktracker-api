package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type SavingsInput struct {
	StartingCapital     int
	YearlyInterestRate  string
	InterestRateType    string
	MonthlyContribution int
	DurationYears       int
	TaxRate             string
	YearlyInflationRate string
	StartDate           string
}

type SaveSavingsInput struct {
	UserID              uuid.UUID
	PlanName            string
	StartingCapital     int
	YearlyInterestRate  string
	InterestRateType    string
	MonthlyContribution int
	DurationYears       int
	TaxRate             string
	YearlyInflationRate string
	StartDate           string
}

type UpdateSavingsInput struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	PlanName            *string
	StartingCapital     *int
	YearlyInterestRate  *string
	InterestRateType    *string
	MonthlyContribution *int
	DurationYears       *int
	TaxRate             *string
	YearlyInflationRate *string
	StartDate           *string
}

type UpdateSavingsData struct {
	ID          uuid.UUID
	Name        string
	SavingsData SavingsInput
}

type SavingsPlan struct {
	ID                    uuid.UUID
	UserID                uuid.UUID
	Name                  string
	OriginalData          SavingsInput
	StartingCapital       decimal.Decimal
	CurrentCapital        decimal.Decimal
	MonthlyContribution   decimal.Decimal
	DurationMonths        decimal.Decimal
	TaxMultiplierM        decimal.Decimal
	InflationMultiplierY  decimal.Decimal
	Date                  time.Time
	InterestMultiplierM   decimal.Decimal
	TotalInterestEarnings decimal.Decimal
	TotalDeposited        decimal.Decimal
	RateOfReturn          decimal.Decimal
	InflationAdjustedROR  decimal.Decimal
	Plan                  []SavingsStatus
}

type SavingsStatus struct {
	Date         time.Time
	Interest     int
	Tax          int
	Contribution int
	Increase     int
	Capital      int
}

func (p *SavingsPlan) PassMonth() SavingsStatus {
	p.Date = p.Date.AddDate(0, 1, 0)
	return SavingsStatus{
		Date: p.Date,
	}
}

func (p *SavingsPlan) GenerateInterest(s SavingsStatus) SavingsStatus {
	interest := p.CurrentCapital.Mul(p.InterestMultiplierM)
	tax := interest.Mul(p.TaxMultiplierM)
	earnings := interest.Sub(tax)
	p.TotalInterestEarnings = p.TotalInterestEarnings.Add(earnings)
	p.CurrentCapital = p.CurrentCapital.Add(earnings)

	s.Interest = int(interest.Round(0).IntPart())
	s.Tax = int(tax.Round(0).IntPart())
	s.Increase = int(earnings.Round(0).IntPart())
	return s
}

func (p *SavingsPlan) Contribute(s SavingsStatus) SavingsStatus {
	p.CurrentCapital = p.CurrentCapital.Add(p.MonthlyContribution)
	p.TotalDeposited = p.TotalDeposited.Add(p.MonthlyContribution)

	s.Increase = s.Increase + int(p.MonthlyContribution.Round(0).IntPart())
	s.Contribution = int(p.MonthlyContribution.Round(0).IntPart())
	s.Capital = int(p.CurrentCapital.Round(0).IntPart())
	return s
}

func (p *SavingsPlan) FinalCalculations() {
	oneHundred := decimal.NewFromInt(100)
	returnRate := p.TotalInterestEarnings.Div(p.TotalDeposited)
	p.RateOfReturn = returnRate.Mul(oneHundred).Round(2)
	totalInflation := p.InflationMultiplierY.Pow(p.DurationMonths.Div(decimal.NewFromInt(12)))
	p.InflationAdjustedROR = returnRate.Div(totalInflation).Mul(oneHundred).Round(2)
}
