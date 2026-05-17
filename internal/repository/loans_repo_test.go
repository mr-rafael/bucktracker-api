package repository

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func TestSaveLoanPaymentPlan(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewLoansRepo(queries)

	test_user_id, err := uuid.Parse("af38df43-3ced-4869-9930-93a0fa0cf1e0")
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	originalData := dto.LoanRequestParams{
		StartingPrincipal:  0,
		YearlyInterestRate: "0",
		MonthlyPayment:     0,
		EscrowPayment:      0,
		StartDate:          "1970-01-01",
	}
	params := domain.LoanPaymentPlan{
		ID:                  uuid.Nil,
		UserID:              test_user_id,
		Name:                "test",
		OriginalData:        domain.LoansInput(originalData),
		StartingPrincipal:   decimal.Zero,
		CurrentPrincipal:    decimal.Zero,
		InterestMultiplierM: decimal.Zero,
		PaymentM:            decimal.Zero,
		EscrowM:             decimal.Zero,
		Date:                time.Now(),
		DurationMonths:      0,
		TotalExpenditure:    decimal.Zero,
		TotalPaid:           decimal.Zero,
		CostOfCreditPercent: decimal.Zero,
	}
	status := domain.LoanStatus{
		Date:          time.Now(),
		Payment:       decimal.Zero,
		Interest:      decimal.Zero,
		OtherPayments: decimal.Zero,
		Paydown:       decimal.Zero,
		Principal:     decimal.Zero,
	}
	params.Plan = append(params.Plan, status)

	got, err := repo.SaveLoanPaymentPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving loan in database: %v", err)
	}

	want := db.Loan{
		UserID: pgtype.UUID{
			Bytes: test_user_id,
			Valid: true,
		},
	}

	if got.UserID.Bytes != want.UserID.Bytes {
		log.Fatalf("Saved (%v) and expected (%v) user IDs did not match.", got.UserID.Bytes, want.UserID.Bytes)
	}
}

func TestGetLoanPaymentPlan(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewLoansRepo(queries)

	test_user_id, err := uuid.Parse("af38df43-3ced-4869-9930-93a0fa0cf1e0")
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	originalData := dto.LoanRequestParams{
		StartingPrincipal:  0,
		YearlyInterestRate: "0",
		MonthlyPayment:     0,
		EscrowPayment:      0,
		StartDate:          "1970-01-01",
	}
	status := domain.LoanStatus{
		Date:          time.Now(),
		Payment:       decimal.Zero,
		Interest:      decimal.Zero,
		OtherPayments: decimal.Zero,
		Paydown:       decimal.Zero,
		Principal:     decimal.Zero,
	}
	params := domain.LoanPaymentPlan{
		ID:                  uuid.Nil,
		UserID:              test_user_id,
		Name:                "test",
		OriginalData:        domain.LoansInput(originalData),
		StartingPrincipal:   decimal.Zero,
		CurrentPrincipal:    decimal.Zero,
		InterestMultiplierM: decimal.Zero,
		PaymentM:            decimal.Zero,
		EscrowM:             decimal.Zero,
		Date:                time.Now(),
		DurationMonths:      0,
		TotalExpenditure:    decimal.Zero,
		TotalPaid:           decimal.Zero,
		CostOfCreditPercent: decimal.Zero,
	}
	params.Plan = append(params.Plan, status)

	plan, err := repo.SaveLoanPaymentPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving loan in database: %v", err)
	}

	got, err := repo.GetLoanByID(ctx, plan.ID.Bytes, plan.UserID.Bytes)
	if err != nil {
		log.Fatalf("Error getting loan from database: %v", err)
	}

	want := db.Loan{
		UserID: pgtype.UUID{
			Bytes: test_user_id,
			Valid: true,
		},
	}

	if got.UserID != want.UserID.Bytes {
		log.Fatalf("The created loan and the retrieved loan didn't match")
	}
}

func TestGetLoansByUser(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewLoansRepo(queries)
	test_user_id, err := uuid.Parse("af38df43-3ced-4869-9930-93a0fa0cf1e0")
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}
	userUUID := pgtype.UUID{
		Bytes: test_user_id,
		Valid: true,
	}
	loansBefore, err := repo.queries.GetLoansByUserID(ctx, userUUID)
	if err != nil {
		log.Fatalf("Error fetching loans before adding new one.")
	}
	want := len(loansBefore) + 1

	originalData := dto.LoanRequestParams{
		StartingPrincipal:  0,
		YearlyInterestRate: "0",
		MonthlyPayment:     0,
		EscrowPayment:      0,
		StartDate:          "1970-01-01",
	}
	params := domain.LoanPaymentPlan{
		ID:                  uuid.Nil,
		UserID:              test_user_id,
		Name:                "test",
		OriginalData:        domain.LoansInput(originalData),
		StartingPrincipal:   decimal.Zero,
		CurrentPrincipal:    decimal.Zero,
		InterestMultiplierM: decimal.Zero,
		PaymentM:            decimal.Zero,
		EscrowM:             decimal.Zero,
		Date:                time.Now(),
		DurationMonths:      0,
		TotalExpenditure:    decimal.Zero,
		TotalPaid:           decimal.Zero,
		CostOfCreditPercent: decimal.Zero,
	}
	_, err = repo.SaveLoanPaymentPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving loan in database: %v", err)
	}
	loansAfter, err := repo.queries.GetLoansByUserID(ctx, userUUID)
	if err != nil {
		log.Fatalf("Error fetching loans after adding new one.")
	}
	got := len(loansAfter)

	if want != got {
		log.Fatalf("The number of loans before insert (%v) didn't match the number of loans after (%v)", want, got)
	}
}

func TestUpdateLoan(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewLoansRepo(queries)

	test_user_id, err := uuid.Parse("af38df43-3ced-4869-9930-93a0fa0cf1e0")
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	originalData := dto.LoanRequestParams{
		StartingPrincipal:  0,
		YearlyInterestRate: "0",
		MonthlyPayment:     0,
		EscrowPayment:      0,
		StartDate:          "1970-01-01",
	}
	params := domain.LoanPaymentPlan{
		ID:                  uuid.Nil,
		UserID:              test_user_id,
		Name:                "test",
		OriginalData:        domain.LoansInput(originalData),
		StartingPrincipal:   decimal.Zero,
		CurrentPrincipal:    decimal.Zero,
		InterestMultiplierM: decimal.Zero,
		PaymentM:            decimal.Zero,
		EscrowM:             decimal.Zero,
		Date:                time.Now(),
		DurationMonths:      0,
		TotalExpenditure:    decimal.Zero,
		TotalPaid:           decimal.Zero,
		CostOfCreditPercent: decimal.Zero,
	}
	status := domain.LoanStatus{
		Date:          time.Now(),
		Payment:       decimal.Zero,
		Interest:      decimal.Zero,
		OtherPayments: decimal.Zero,
		Paydown:       decimal.Zero,
		Principal:     decimal.Zero,
	}
	params.Plan = append(params.Plan, status)

	result, err := repo.SaveLoanPaymentPlan(ctx, params)

	updatedName := "updatedTest"
	updatedPrincipal := 100
	updatedInterest := "1.05"

	params.ID = result.ID.Bytes
	params.Name = updatedName
	params.OriginalData.StartingPrincipal = updatedPrincipal
	params.OriginalData.YearlyInterestRate = updatedInterest

	got, err := repo.UpdateLoan(ctx, params)

	want := db.Loan{
		Name:               updatedName,
		StartingPrincipal:  int32(updatedPrincipal),
		YearlyInterestRate: updatedInterest,
	}

	if got.Name != want.Name {
		log.Fatalf("Loan name returned from the database (%v) doesn't match the expected one (%v).", got.Name, want.Name)
	}
	if got.StartingPrincipal != want.StartingPrincipal {
		log.Fatalf("Loan starting principal returned from the database (%v) doesn't match the expected one (%v).", got.StartingPrincipal, want.StartingPrincipal)
	}
	if got.YearlyInterestRate != want.YearlyInterestRate {
		log.Fatalf("Loan interest rate returned from the database (%v) doesn't match the expected one (%v).", got.StartingPrincipal, want.StartingPrincipal)
	}
}

func TestDeleteLoan(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewLoansRepo(queries)
	test_user_id, err := uuid.Parse("af38df43-3ced-4869-9930-93a0fa0cf1e0")
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	originalData := dto.LoanRequestParams{
		StartingPrincipal:  0,
		YearlyInterestRate: "0",
		MonthlyPayment:     0,
		EscrowPayment:      0,
		StartDate:          "1970-01-01",
	}
	params := domain.LoanPaymentPlan{
		ID:                  uuid.Nil,
		UserID:              test_user_id,
		Name:                "test",
		OriginalData:        domain.LoansInput(originalData),
		StartingPrincipal:   decimal.Zero,
		CurrentPrincipal:    decimal.Zero,
		InterestMultiplierM: decimal.Zero,
		PaymentM:            decimal.Zero,
		EscrowM:             decimal.Zero,
		Date:                time.Now(),
		DurationMonths:      0,
		TotalExpenditure:    decimal.Zero,
		TotalPaid:           decimal.Zero,
		CostOfCreditPercent: decimal.Zero,
	}
	loanInfo, err := repo.SaveLoanPaymentPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving loan in database: %v", err)
	}
	deleteParams := db.DeleteLoanParams{
		ID:     loanInfo.ID,
		UserID: loanInfo.UserID,
	}
	err = repo.queries.DeleteLoan(ctx, deleteParams)
	if err != nil {
		log.Fatalf("Error deleting loan.")
	}

	getParams := db.GetLoanParams{
		ID:     loanInfo.ID,
		UserID: loanInfo.UserID,
	}

	_, got := repo.queries.GetLoan(ctx, getParams)

	if got == nil {
		log.Fatalf("The loan was not deleted.")
	}
}
