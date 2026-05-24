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

type InputError struct {
	Message string
}

func (e InputError) Error() string {
	return e.Message
}

type LoansService struct {
	loansRepo LoansRepository
}

type LoansRepository interface {
	SaveLoanPaymentPlan(context.Context, domain.LoanPaymentPlan) (db.Loan, error)
	GetLoanPaymentPlansByUser(context.Context, uuid.UUID) ([]db.GetLoansByUserIDRow, error)
	GetLoanByID(context.Context, uuid.UUID, uuid.UUID) (domain.LoanPaymentPlan, error)
	GetLoanInitialData(context.Context, uuid.UUID, uuid.UUID) (domain.UpdateLoanData, error)
	UpdateLoan(context.Context, domain.LoanPaymentPlan) (db.Loan, error)
	DeleteLoan(context.Context, uuid.UUID, uuid.UUID) error
}

func NewLoansService(repo LoansRepository) *LoansService {
	return &LoansService{loansRepo: repo}
}

const minLoanCents = "1"
const maxLoanCents = "100000000000"
const minInterestRate = "0"
const maxInterestRate = "100"
const minMonthlyPaymentCents = "1"
const maxMonthlyPaymentCents = "100000000000"
const minEscrowCents = "0"
const maxEscrowCents = "100000000000"
const maxPaymentYears = 30

func (s *LoansService) CalculateLoanPaymentPlan(input domain.LoansInput) (domain.LoanPaymentPlan, error) {
	plan, err := initializePaymentPlan(input, uuid.Nil, "")
	if err != nil {
		return domain.LoanPaymentPlan{}, err
	}

	plan, err = calculatePaymentPlan(plan)
	if err != nil {
		return domain.LoanPaymentPlan{}, fmt.Errorf("Error calculating payment plan: %v", err)
	}

	return plan, nil
}

func (s *LoansService) SaveLoanPaymentPlan(ctx context.Context, input domain.SaveLoanInput) (db.Loan, error) {
	plan, err := initializePaymentPlan(saveInputToLoanInput(input), input.UserID, input.LoanName)
	if err != nil {
		return db.Loan{}, err
	}

	plan, err = calculatePaymentPlan(plan)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Error calculating payment plan: %v", err)
	}
	result, err := s.loansRepo.SaveLoanPaymentPlan(ctx, plan)
	if err != nil {
		return db.Loan{}, err
	}

	return result, nil
}

func (s *LoansService) GetLoansByUser(ctx context.Context, input uuid.UUID) ([]db.GetLoansByUserIDRow, error) {
	result, err := s.loansRepo.GetLoanPaymentPlansByUser(ctx, input)
	if err != nil {
		return []db.GetLoansByUserIDRow{}, err
	}
	return result, nil
}

func (s *LoansService) GetLoan(ctx context.Context, planID uuid.UUID, userID uuid.UUID) (domain.LoanPaymentPlan, error) {
	result, err := s.loansRepo.GetLoanByID(ctx, planID, userID)
	if err != nil {
		return domain.LoanPaymentPlan{}, err
	}
	return result, nil
}

func (s *LoansService) UpdateLoan(ctx context.Context, input domain.UpdateLoanInput) (db.Loan, error) {
	originalData, err := s.loansRepo.GetLoanInitialData(ctx, input.ID, input.UserID)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Loan not found.")
	}
	patchedData := patchLoanFields(originalData, input)

	plan, err := initializePaymentPlan(patchedData.LoanData, input.UserID, patchedData.Name)
	if err != nil {
		return db.Loan{}, err
	}
	plan, err = calculatePaymentPlan(plan)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Error calculating payment plan: %v", err)
	}
	plan.ID = input.ID
	result, err := s.loansRepo.UpdateLoan(ctx, plan)
	if err != nil {
		return db.Loan{}, err
	}

	return result, nil
}

func (s *LoansService) DeleteLoan(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) error {
	return s.loansRepo.DeleteLoan(ctx, loanID, userID)
}

func calculatePaymentPlan(plan domain.LoanPaymentPlan) (domain.LoanPaymentPlan, error) {
	i := 0
	for plan.CurrentPrincipal.Compare(decimal.Zero) == 1 {
		i++
		if i > maxPaymentYears*12 {
			remainder := plan.CurrentPrincipal.Div(decimal.NewFromInt(100)).Round(2).String()
			return domain.LoanPaymentPlan{}, fmt.Errorf("The payment plan exceeds the year limit (%v years), with a remaining %v to pay", maxPaymentYears, remainder)
		}
		payment := plan.PassMonth()
		payment = plan.GenerateInterest(payment)
		payment = plan.ChargeEscrow(payment)
		payment = plan.MakePayment(payment)
		plan.Plan = append(plan.Plan, payment)
	}
	plan.FinalCalculations()

	return plan, nil
}

func initializePaymentPlan(input domain.LoansInput, userID uuid.UUID, name string) (domain.LoanPaymentPlan, error) {
	plan := domain.LoanPaymentPlan{}
	oneHundred := decimal.NewFromInt(100)

	plan.OriginalData = input
	plan.UserID = userID
	plan.Name = name

	startingPrincipal := decimal.NewFromInt(int64(input.StartingPrincipal))
	if !decimalIsBetween(startingPrincipal, minLoanCents, maxLoanCents) {
		return domain.LoanPaymentPlan{}, InputError{Message: fmt.Sprintf("invalid starting principal: '%v'. the accepted range is 0.01 - 1,000,000,000", startingPrincipal.Div(oneHundred).Round(2))}
	}
	plan.StartingPrincipal = startingPrincipal
	plan.CurrentPrincipal = startingPrincipal

	monthlyInterestRate, err := getMonthlyAPRMultiplier(input.YearlyInterestRate)
	if !stringNumberBetween(input.YearlyInterestRate, minInterestRate, maxInterestRate) {
		return domain.LoanPaymentPlan{}, InputError{Message: fmt.Sprintf("invalid interest rate: '%v'. the accepted range is 0%% - 100%%", input.YearlyInterestRate)}
	}
	if err != nil {
		return domain.LoanPaymentPlan{}, InputError{Message: fmt.Sprintf("invalid interest rate: '%v'", input.YearlyInterestRate)}
	}
	plan.InterestMultiplierM = monthlyInterestRate

	monthlyPayment := decimal.NewFromInt(int64(input.MonthlyPayment))
	if !decimalIsBetween(monthlyPayment, minMonthlyPaymentCents, maxMonthlyPaymentCents) {
		return domain.LoanPaymentPlan{}, InputError{Message: fmt.Sprintf("invalid monthly payments: '%v'. the accepted range is 0.01 - 1,000,000,000", monthlyPayment.Div(oneHundred).Round(2))}
	}
	plan.PaymentM = monthlyPayment

	escrow := decimal.NewFromInt(int64(input.EscrowPayment))
	if !decimalIsBetween(escrow, minEscrowCents, maxEscrowCents) {
		return domain.LoanPaymentPlan{}, InputError{Message: fmt.Sprintf("invalid escrow payment: '%v'. the accepted range is 0.01 - 1,000,000,000", escrow.Div(oneHundred).Round(2))}
	}
	plan.EscrowM = escrow

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return domain.LoanPaymentPlan{}, InputError{Message: fmt.Sprintf("invalid start date: %v", input.StartDate)}
	}
	plan.Date = startDate

	err = checkIfEnoughMonthlyPayment(plan)
	if err != nil {
		return domain.LoanPaymentPlan{}, err
	}

	return plan, nil
}

func saveInputToLoanInput(input domain.SaveLoanInput) domain.LoansInput {
	return domain.LoansInput{
		StartingPrincipal:  input.StartingPrincipal,
		YearlyInterestRate: input.YearlyInterestRate,
		MonthlyPayment:     input.MonthlyPayment,
		EscrowPayment:      input.EscrowPayment,
		StartDate:          input.StartDate,
	}

}

func checkIfEnoughMonthlyPayment(plan domain.LoanPaymentPlan) error {
	firstMonthInterest := plan.StartingPrincipal.Mul(plan.InterestMultiplierM)
	minPayment := firstMonthInterest.Add(plan.EscrowM)
	aHundred := decimal.NewFromInt32(100)

	if plan.PaymentM.Compare(minPayment) != 1 {
		return fmt.Errorf("The monthly payment is not enough to cover interest and escrow payment for the first month (total $%v). Please enter a higher monthly payment.", minPayment.Div(aHundred).Round(2).String())
	}
	return nil
}

func patchLoanFields(loanData domain.UpdateLoanData, patchData domain.UpdateLoanInput) domain.UpdateLoanData {
	if patchData.LoanName != nil {
		loanData.Name = *patchData.LoanName
	}
	if patchData.StartingPrincipal != nil {
		loanData.LoanData.StartingPrincipal = *patchData.StartingPrincipal
	}
	if patchData.YearlyInterestRate != nil {
		loanData.LoanData.YearlyInterestRate = *patchData.YearlyInterestRate
	}
	if patchData.MonthlyPayment != nil {
		loanData.LoanData.MonthlyPayment = *patchData.MonthlyPayment
	}
	if patchData.EscrowPayment != nil {
		loanData.LoanData.EscrowPayment = *patchData.EscrowPayment
	}
	if patchData.StartDate != nil {
		loanData.LoanData.StartDate = *patchData.StartDate
	}
	return loanData
}
