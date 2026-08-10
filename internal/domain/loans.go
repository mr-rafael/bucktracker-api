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

type UpdateLoanData struct {
	ID       uuid.UUID
	Name     string
	LoanData LoansInput
}

type Loan struct {
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
	DefaultPaymentPlan  *LoanPaymentPlan
	PaymentPlans        []LoanPaymentPlan
}

type LoanPaymentPlan struct {
	ID                  uuid.UUID
	Name                string
	DurationMonths      int
	TotalExpenditure    decimal.Decimal
	TotalPaid           decimal.Decimal
	CostOfCreditPercent decimal.Decimal
	Plan                []LoanStatus
	PrincipalPayments   []PrincipalPayment
}

type PrincipalPayment struct {
	AmountPaid decimal.Decimal
	Date       time.Time
}

type LoanStatus struct {
	Date          time.Time
	Payment       decimal.Decimal
	Interest      decimal.Decimal
	OtherPayments decimal.Decimal
	Paydown       decimal.Decimal
	Principal     decimal.Decimal
}

func (l *Loan) ensureDefaultPaymentPlan() *LoanPaymentPlan {
	if l.DefaultPaymentPlan == nil {
		l.DefaultPaymentPlan = &LoanPaymentPlan{
			Name: "Default Payment Plan",
		}
	}
	if l.DefaultPaymentPlan.Name == "" {
		l.DefaultPaymentPlan.Name = "Default Payment Plan"
	}
	return l.DefaultPaymentPlan
}

func (l *Loan) PassMonth() LoanStatus {
	l.Date = l.Date.AddDate(0, 1, 0)
	plan := l.ensureDefaultPaymentPlan()
	plan.DurationMonths += 1
	return LoanStatus{
		Date: l.Date,
	}
}

func (l *Loan) GenerateInterest(s LoanStatus) LoanStatus {
	interest := l.CurrentPrincipal.Mul(l.InterestMultiplierM)
	plan := l.ensureDefaultPaymentPlan()
	plan.TotalExpenditure = plan.TotalExpenditure.Add(interest)

	s.Interest = interest
	return s
}

func (l *Loan) ChargeEscrow(s LoanStatus) LoanStatus {
	plan := l.ensureDefaultPaymentPlan()
	plan.TotalExpenditure = plan.TotalExpenditure.Add(l.EscrowM)

	s.OtherPayments = l.EscrowM
	return s
}

func (l *Loan) MakePayment(s LoanStatus) LoanStatus {
	plan := l.ensureDefaultPaymentPlan()
	paydown := l.PaymentM.Sub(s.Interest).Sub(s.OtherPayments)
	if l.CurrentPrincipal.Cmp(paydown) == -1 {
		payment := l.CurrentPrincipal.Add(s.Interest).Add(s.OtherPayments)
		plan.TotalPaid = plan.TotalPaid.Add(payment)
		s.Payment = payment
		s.Paydown = l.CurrentPrincipal
		l.CurrentPrincipal = decimal.Zero
		s.Principal = l.CurrentPrincipal
	} else {
		plan.TotalPaid = plan.TotalPaid.Add(l.PaymentM)
		s.Payment = l.PaymentM
		s.Paydown = paydown
		l.CurrentPrincipal = l.CurrentPrincipal.Sub(paydown)
		s.Principal = l.CurrentPrincipal
	}

	return s
}

func (l *Loan) FinalCalculations() {
	plan := l.ensureDefaultPaymentPlan()
	one := decimal.NewFromInt32(1)
	oneHundred := decimal.NewFromInt32(100)
	plan.CostOfCreditPercent = plan.TotalPaid.Div(l.StartingPrincipal).Sub(one).Mul(oneHundred)
}
