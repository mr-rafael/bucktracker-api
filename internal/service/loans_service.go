package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Mr-Rafael/bucktracker-api/internal/db"
	"github.com/Mr-Rafael/bucktracker-api/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LoanInputError struct {
	Message string
}

func (e LoanInputError) Error() string {
	return e.Message
}

type LoansService struct {
	loansRepo LoansRepository
}

type LoansRepository interface {
	SaveLoanPaymentPlan(context.Context, domain.Loan) (db.Loan, error)
	GetLoanPaymentPlansByUser(context.Context, uuid.UUID) ([]db.GetLoansByUserIDRow, error)
	GetLoanByID(context.Context, uuid.UUID, uuid.UUID) (domain.Loan, error)
	GetPaymentPlanByID(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) (domain.LoanPaymentPlan, error)
	CreatePaymentPlanForLoan(context.Context, uuid.UUID, uuid.UUID, domain.LoanPaymentPlan) (domain.LoanPaymentPlan, error)
	GetLoanInitialData(context.Context, uuid.UUID, uuid.UUID) (domain.UpdateLoanData, error)
	UpdateLoan(context.Context, domain.Loan) (db.Loan, error)
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
const defaultPaymentPlanName = "Default Payment Plan"

func (s *LoansService) CalculateLoanPaymentPlan(input domain.LoansInput) (domain.Loan, error) {
	loan, err := initializeLoan(input, uuid.Nil, "")
	if err != nil {
		return domain.Loan{}, err
	}

	loan, err = calculatePaymentPlan(loan, nil)
	if err != nil {
		return domain.Loan{}, fmt.Errorf("Error calculating payment plan: %v", err)
	}

	return loan, nil
}

func (s *LoansService) SaveLoanPaymentPlan(ctx context.Context, input domain.SaveLoanInput) (domain.Loan, error) {
	loan, err := initializeLoan(saveInputToLoanInput(input), input.UserID, input.LoanName)
	if err != nil {
		return domain.Loan{}, err
	}

	loan, err = calculatePaymentPlan(loan, nil)
	if err != nil {
		return domain.Loan{}, fmt.Errorf("Error calculating payment plan: %v", err)
	}
	// Newly created loans always start with a default plan and no principal payments.
	if loan.DefaultPaymentPlan != nil {
		loan.DefaultPaymentPlan.PrincipalPayments = nil
	}
	result, err := s.loansRepo.SaveLoanPaymentPlan(ctx, loan)
	if err != nil {
		return domain.Loan{}, err
	}

	loan.ID = result.ID.Bytes
	if result.DefaultPaymentPlan.Valid && loan.DefaultPaymentPlan != nil {
		loan.DefaultPaymentPlan.ID = result.DefaultPaymentPlan.Bytes
	}

	return loan, nil
}

func (s *LoansService) GetLoansByUser(ctx context.Context, input uuid.UUID) ([]db.GetLoansByUserIDRow, error) {
	result, err := s.loansRepo.GetLoanPaymentPlansByUser(ctx, input)
	if err != nil {
		return []db.GetLoansByUserIDRow{}, err
	}
	return result, nil
}

func (s *LoansService) GetLoan(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) (domain.Loan, error) {
	result, err := s.loansRepo.GetLoanByID(ctx, loanID, userID)
	if err != nil {
		return domain.Loan{}, err
	}
	return result, nil
}

func (s *LoansService) GetPaymentPlan(ctx context.Context, loanID uuid.UUID, paymentPlanID uuid.UUID, userID uuid.UUID) (domain.LoanPaymentPlan, error) {
	result, err := s.loansRepo.GetPaymentPlanByID(ctx, loanID, paymentPlanID, userID)
	if err != nil {
		return domain.LoanPaymentPlan{}, err
	}
	return result, nil
}

func (s *LoansService) CreatePaymentPlan(ctx context.Context, input domain.CreatePaymentPlanInput) (domain.LoanPaymentPlan, error) {
	originalData, err := s.loansRepo.GetLoanInitialData(ctx, input.LoanID, input.UserID)
	if err != nil {
		return domain.LoanPaymentPlan{}, fmt.Errorf("Loan not found.")
	}

	loan, err := initializeLoan(originalData.LoanData, input.UserID, originalData.Name)
	if err != nil {
		return domain.LoanPaymentPlan{}, err
	}

	paymentPlan := &domain.LoanPaymentPlan{
		Name:              input.Name,
		PrincipalPayments: input.PrincipalPayments,
	}

	loan, err = calculatePaymentPlan(loan, paymentPlan)
	if err != nil {
		return domain.LoanPaymentPlan{}, fmt.Errorf("Error calculating payment plan: %v", err)
	}

	created, err := s.loansRepo.CreatePaymentPlanForLoan(ctx, input.LoanID, input.UserID, *loan.DefaultPaymentPlan)
	if err != nil {
		return domain.LoanPaymentPlan{}, err
	}
	return created, nil
}

func (s *LoansService) UpdateLoan(ctx context.Context, input domain.UpdateLoanInput) (domain.Loan, error) {
	originalData, err := s.loansRepo.GetLoanInitialData(ctx, input.ID, input.UserID)
	if err != nil {
		return domain.Loan{}, fmt.Errorf("Loan not found.")
	}
	patchedData := patchLoanFields(originalData, input)

	existingLoan, err := s.loansRepo.GetLoanByID(ctx, input.ID, input.UserID)
	if err != nil {
		return domain.Loan{}, fmt.Errorf("Loan not found.")
	}

	loan, err := initializeLoan(patchedData.LoanData, input.UserID, patchedData.Name)
	if err != nil {
		return domain.Loan{}, err
	}
	if existingLoan.DefaultPaymentPlan != nil {
		loan.DefaultPaymentPlan = &domain.LoanPaymentPlan{
			Name: existingLoan.DefaultPaymentPlan.Name,
		}
	}
	loan, err = calculatePaymentPlan(loan, nil)
	if err != nil {
		return domain.Loan{}, fmt.Errorf("Error calculating payment plan: %v", err)
	}
	loan.ID = input.ID
	result, err := s.loansRepo.UpdateLoan(ctx, loan)
	if err != nil {
		return domain.Loan{}, err
	}

	loan.ID = result.ID.Bytes
	if result.DefaultPaymentPlan.Valid && loan.DefaultPaymentPlan != nil {
		loan.DefaultPaymentPlan.ID = result.DefaultPaymentPlan.Bytes
	}

	return loan, nil
}

func (s *LoansService) DeleteLoan(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) error {
	return s.loansRepo.DeleteLoan(ctx, loanID, userID)
}

func calculatePaymentPlan(loan domain.Loan, paymentPlan *domain.LoanPaymentPlan) (domain.Loan, error) {
	if paymentPlan == nil {
		planName := defaultPaymentPlanName
		if loan.DefaultPaymentPlan != nil && loan.DefaultPaymentPlan.Name != "" {
			planName = loan.DefaultPaymentPlan.Name
		}
		loan.DefaultPaymentPlan = &domain.LoanPaymentPlan{
			Name: planName,
		}
	} else {
		loan.DefaultPaymentPlan = paymentPlan
	}

	i := 0
	for loan.CurrentPrincipal.Compare(decimal.Zero) == 1 {
		i++
		if i > maxPaymentYears*12 {
			remainder := loan.CurrentPrincipal.Div(decimal.NewFromInt(100)).Round(2).String()
			return domain.Loan{}, fmt.Errorf("The payment plan exceeds the year limit (%v years), with a remaining %v to pay", maxPaymentYears, remainder)
		}
		payment := loan.PassMonth()
		payment = loan.GenerateInterest(payment)
		payment = loan.ChargeEscrow(payment)
		payment = loan.MakePayment(payment)
		loan.DefaultPaymentPlan.Plan = append(loan.DefaultPaymentPlan.Plan, payment)
	}
	loan.FinalCalculations()

	return loan, nil
}

func initializeLoan(input domain.LoansInput, userID uuid.UUID, name string) (domain.Loan, error) {
	loan := domain.Loan{}
	oneHundred := decimal.NewFromInt(100)

	loan.OriginalData = input
	loan.UserID = userID
	loan.Name = name

	startingPrincipal := decimal.NewFromInt(int64(input.StartingPrincipal))
	if !decimalIsBetween(startingPrincipal, minLoanCents, maxLoanCents) {
		return domain.Loan{}, LoanInputError{Message: fmt.Sprintf("invalid starting principal: '%v'. the accepted range is 0.01 - 1,000,000,000", startingPrincipal.Div(oneHundred).Round(2))}
	}
	loan.StartingPrincipal = startingPrincipal
	loan.CurrentPrincipal = startingPrincipal

	monthlyInterestRate, err := getMonthlyAPRMultiplier(input.YearlyInterestRate)
	if !stringNumberBetween(input.YearlyInterestRate, minInterestRate, maxInterestRate) {
		return domain.Loan{}, LoanInputError{Message: fmt.Sprintf("invalid interest rate: '%v'. the accepted range is 0%% - 100%%", input.YearlyInterestRate)}
	}
	if err != nil {
		return domain.Loan{}, LoanInputError{Message: fmt.Sprintf("invalid interest rate: '%v'", input.YearlyInterestRate)}
	}
	loan.InterestMultiplierM = monthlyInterestRate

	monthlyPayment := decimal.NewFromInt(int64(input.MonthlyPayment))
	if !decimalIsBetween(monthlyPayment, minMonthlyPaymentCents, maxMonthlyPaymentCents) {
		return domain.Loan{}, LoanInputError{Message: fmt.Sprintf("invalid monthly payments: '%v'. the accepted range is 0.01 - 1,000,000,000", monthlyPayment.Div(oneHundred).Round(2))}
	}
	loan.PaymentM = monthlyPayment

	escrow := decimal.NewFromInt(int64(input.EscrowPayment))
	if !decimalIsBetween(escrow, minEscrowCents, maxEscrowCents) {
		return domain.Loan{}, LoanInputError{Message: fmt.Sprintf("invalid escrow payment: '%v'. the accepted range is 0.01 - 1,000,000,000", escrow.Div(oneHundred).Round(2))}
	}
	loan.EscrowM = escrow

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return domain.Loan{}, LoanInputError{Message: fmt.Sprintf("invalid start date: %v", input.StartDate)}
	}
	loan.Date = startDate

	err = checkIfEnoughMonthlyPayment(loan)
	if err != nil {
		return domain.Loan{}, err
	}

	return loan, nil
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

func checkIfEnoughMonthlyPayment(loan domain.Loan) error {
	firstMonthInterest := loan.StartingPrincipal.Mul(loan.InterestMultiplierM)
	minPayment := firstMonthInterest.Add(loan.EscrowM)
	aHundred := decimal.NewFromInt32(100)

	if loan.PaymentM.Compare(minPayment) != 1 {
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
