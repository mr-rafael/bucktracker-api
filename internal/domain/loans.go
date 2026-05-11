package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LoansInput struct {
	StartingPrincipal  int
	YearlyInterestRate string
	MonthlyPayment     int
	EscrowPayment      int
	StartDate          string
}

type SaveLoanInput struct {
	UserID             uuid.UUID
	LoanName           string
	StartingPrincipal  int
	YearlyInterestRate string
	MonthlyPayment     int
	EscrowPayment      int
	StartDate          string
}

type UpdateLoanInput struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	LoanName           *string
	StartingPrincipal  *int
	YearlyInterestRate *string
	MonthlyPayment     *int
	EscrowPayment      *int
	StartDate          *string
}

type LoanPaymentPlan struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	Name                string
	OriginalData        LoansInput
	StartingPrincipal   decimal.Decimal
	CurrentPrincipal    decimal.Decimal
	InterestMultiplierM decimal.Decimal
	PaymentM            decimal.Decimal
	EscrowM             decimal.Decimal
	Date                time.Time
	DurationMonths      int
	TotalExpenditure    decimal.Decimal
	TotalPaid           decimal.Decimal
	CostOfCreditPercent decimal.Decimal
	Plan                []LoanStatus
}

type LoanStatus struct {
	Date          time.Time
	Payment       decimal.Decimal
	Interest      decimal.Decimal
	OtherPayments decimal.Decimal
	Paydown       decimal.Decimal
	Principal     decimal.Decimal
}

func (p *LoanPaymentPlan) PassMonth() LoanStatus {
	p.Date = p.Date.AddDate(0, 1, 0)
	p.DurationMonths += 1
	return LoanStatus{
		Date: p.Date,
	}
}

func (p *LoanPaymentPlan) GenerateInterest(s LoanStatus) LoanStatus {
	interest := p.CurrentPrincipal.Mul(p.InterestMultiplierM)
	p.TotalExpenditure = p.TotalExpenditure.Add(interest)

	s.Interest = interest
	return s
}

func (p *LoanPaymentPlan) ChargeEscrow(s LoanStatus) LoanStatus {
	p.TotalExpenditure = p.TotalExpenditure.Add(p.EscrowM)

	s.OtherPayments = p.EscrowM
	return s
}

func (p *LoanPaymentPlan) MakePayment(s LoanStatus) LoanStatus {
	paydown := p.PaymentM.Sub(s.Interest).Sub(s.OtherPayments)
	if p.CurrentPrincipal.Cmp(paydown) == -1 {
		payment := p.CurrentPrincipal.Add(s.Interest).Add(s.OtherPayments)
		p.TotalPaid = p.TotalPaid.Add(payment)
		s.Payment = payment
		s.Paydown = p.CurrentPrincipal
		p.CurrentPrincipal = decimal.Zero
		s.Principal = p.CurrentPrincipal
	} else {
		p.TotalPaid = p.TotalPaid.Add(p.PaymentM)
		s.Payment = p.PaymentM
		s.Paydown = paydown
		p.CurrentPrincipal = p.CurrentPrincipal.Sub(paydown)
		s.Principal = p.CurrentPrincipal
	}

	return s
}

func (p *LoanPaymentPlan) FinalCalculations() {
	p.CostOfCreditPercent = p.TotalPaid.Div(p.StartingPrincipal)
}
