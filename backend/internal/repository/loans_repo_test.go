package repository

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Mr-Rafael/bucktracker-api/internal/db"
	"github.com/Mr-Rafael/bucktracker-api/internal/domain"
	"github.com/Mr-Rafael/bucktracker-api/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func testLoanParams(userID uuid.UUID) domain.Loan {
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
	return domain.Loan{
		ID:                  uuid.Nil,
		UserID:              userID,
		Name:                "test",
		OriginalData:        domain.LoansInput(originalData),
		StartingPrincipal:   decimal.Zero,
		CurrentPrincipal:    decimal.Zero,
		InterestMultiplierM: decimal.Zero,
		PaymentM:            decimal.Zero,
		EscrowM:             decimal.Zero,
		Date:                time.Now(),
		DefaultPaymentPlan: &domain.LoanPaymentPlan{
			DurationMonths:      0,
			TotalExpenditure:    decimal.Zero,
			TotalPaid:           decimal.Zero,
			CostOfCreditPercent: decimal.Zero,
			Plan:                []domain.LoanStatus{status},
		},
	}
}

func TestSaveLoanPaymentPlan(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewLoansRepo(queries)

	testUser, err := CreateTestUserIfNotExists()
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	params := testLoanParams(testUser.ID.Bytes)

	got, err := repo.SaveLoanPaymentPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving loan in database: %v", err)
	}

	want := db.Loan{
		UserID: pgtype.UUID{
			Bytes: testUser.ID.Bytes,
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

	testUser, err := CreateTestUserIfNotExists()
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	params := testLoanParams(testUser.ID.Bytes)

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
			Bytes: testUser.ID.Bytes,
			Valid: true,
		},
	}

	if got.UserID != want.UserID.Bytes {
		log.Fatalf("The created loan and the retrieved loan didn't match")
	}
}

func TestGetPaymentPlanByID(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewLoansRepo(queries)

	testUser, err := CreateTestUserIfNotExists()
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	params := testLoanParams(testUser.ID.Bytes)
	params.DefaultPaymentPlan.Name = "Default Payment Plan"
	params.DefaultPaymentPlan.DurationMonths = 1
	params.DefaultPaymentPlan.TotalExpenditure = decimal.NewFromInt(100)
	params.DefaultPaymentPlan.TotalPaid = decimal.NewFromInt(200)
	params.DefaultPaymentPlan.CostOfCreditPercent = decimal.NewFromInt(5)

	savedLoan, err := repo.SaveLoanPaymentPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving loan in database: %v", err)
	}

	got, err := repo.GetPaymentPlanByID(ctx, savedLoan.ID.Bytes, savedLoan.DefaultPaymentPlan.Bytes, savedLoan.UserID.Bytes)
	if err != nil {
		log.Fatalf("Error getting payment plan from database: %v", err)
	}

	if got.ID != savedLoan.DefaultPaymentPlan.Bytes {
		log.Fatalf("Payment plan IDs did not match")
	}
	if got.Name != "Default Payment Plan" {
		log.Fatalf("Expected payment plan name Default Payment Plan, got %v", got.Name)
	}
	if got.DurationMonths != 1 {
		log.Fatalf("Expected duration months 1, got %v", got.DurationMonths)
	}
	if len(got.Plan) != 1 {
		log.Fatalf("Expected 1 loan state in payment plan breakdown, got %v", len(got.Plan))
	}
}

func TestGetLoansByUser(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewLoansRepo(queries)

	testUser, err := CreateTestUserIfNotExists()
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}
	loansBefore, err := repo.queries.GetLoansByUserID(ctx, testUser.ID)
	if err != nil {
		log.Fatalf("Error fetching loans before adding new one.")
	}
	want := len(loansBefore) + 1

	params := testLoanParams(testUser.ID.Bytes)
	_, err = repo.SaveLoanPaymentPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving loan in database: %v", err)
	}
	loansAfter, err := repo.queries.GetLoansByUserID(ctx, testUser.ID)
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

	testUser, err := CreateTestUserIfNotExists()
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	params := testLoanParams(testUser.ID.Bytes)

	result, err := repo.SaveLoanPaymentPlan(ctx, params)

	updatedName := "updatedLoanTest"
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
		log.Fatalf("Loan interest rate returned from the database (%v) doesn't match the expected one (%v).", got.YearlyInterestRate, want.YearlyInterestRate)
	}
}

func TestDeleteLoan(t *testing.T) {
	ctx := context.Background()
	queries := initializeQueries(ctx)
	repo := NewLoansRepo(queries)

	testUser, err := CreateTestUserIfNotExists()
	if err != nil {
		log.Fatalf("failed to parse the test user uuid: %v", err)
	}

	params := testLoanParams(testUser.ID.Bytes)
	loanInfo, err := repo.SaveLoanPaymentPlan(ctx, params)
	if err != nil {
		log.Fatalf("Error saving loan in database: %v", err)
	}
	err = repo.DeleteLoan(ctx, loanInfo.ID.Bytes, loanInfo.UserID.Bytes)
	if err != nil {
		log.Fatalf("Error deleting loan: %v", err)
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
