package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

type LoansRepo struct {
	queries *db.Queries
}

func NewLoansRepo(queries *db.Queries) *LoansRepo {
	return &LoansRepo{queries: queries}
}

func (r *LoansRepo) SaveLoanPaymentPlan(ctx context.Context, plan domain.LoanPaymentPlan) (db.Loan, error) {
	loanParams, err := toLoanInsertQueryParams(plan)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Error preparing params for insert query: %v", err)
	}

	queryResult, err := r.queries.CreateLoan(ctx, loanParams)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Failed to save to database: %v", err)
	}

	for _, status := range plan.Plan {
		_, err := r.queries.CreateLoanState(ctx, toLoanStateInsertParams(status, queryResult.ID))
		if err != nil {
			return db.Loan{}, fmt.Errorf("Failed to save loan status to database: %v", err)
		}
	}
	return queryResult, nil
}

func (r *LoansRepo) GetLoanPaymentPlansByUser(ctx context.Context, userID uuid.UUID) ([]db.GetLoansByUserIDRow, error) {
	queryUserID := pgtype.UUID{
		Bytes: userID,
		Valid: true,
	}

	result, err := r.queries.GetLoansByUserID(ctx, queryUserID)
	if err != nil {
		return []db.GetLoansByUserIDRow{}, fmt.Errorf("failed to fetch user's loan payment plans: %v", err)
	}
	return result, nil
}

func (r *LoansRepo) GetLoanByID(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) (domain.LoanPaymentPlan, error) {
	queryLoanID := pgtype.UUID{
		Bytes: loanID,
		Valid: true,
	}

	loanQueryResult, err := r.queries.GetLoan(ctx, toLoanGetParams(loanID, userID))
	if err != nil {
		return domain.LoanPaymentPlan{}, fmt.Errorf("failed to fetch loan pament plan from database: %v", err)
	}
	plan, err := toLoanPaymentPlan(loanQueryResult)

	statesQueryResult, err := r.queries.GetLoanStatesByLoanID(ctx, queryLoanID)
	if err != nil {
		return domain.LoanPaymentPlan{}, fmt.Errorf("failed to fetch loan payment plan rows from database: %v", err)
	}
	for _, state := range statesQueryResult {
		plan.Plan = append(plan.Plan, domain.LoanStatus{
			Date:          state.Date.Time,
			Payment:       decimal.NewFromInt32(state.Payment),
			Interest:      decimal.NewFromInt32(state.Interest),
			OtherPayments: decimal.NewFromInt32(state.OtherPayments),
			Paydown:       decimal.NewFromInt32(state.Paydown),
			Principal:     decimal.NewFromInt32(state.Principal),
		})
	}

	return plan, nil
}

func (r *LoansRepo) GetLoanInitialData(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) (domain.UpdateLoanData, error) {

	loanQueryResult, err := r.queries.GetLoanInitialData(ctx, toInitialLoanDataGetParams(loanID, userID))
	if err != nil {
		return domain.UpdateLoanData{}, fmt.Errorf("failed to fetch loan pament plan from database: %v", err)
	}
	loansInput := domain.LoansInput{
		StartingPrincipal:  int(loanQueryResult.StartingPrincipal),
		YearlyInterestRate: loanQueryResult.YearlyInterestRate,
		MonthlyPayment:     int(loanQueryResult.MonthlyPayment),
		EscrowPayment:      int(loanQueryResult.EscrowPayment),
		StartDate:          loanQueryResult.StartDate.Time.Format(time.RFC3339),
	}
	loanData := domain.UpdateLoanData{
		ID:       loanID,
		Name:     loanQueryResult.Name,
		LoanData: loansInput,
	}

	return loanData, nil
}

func (r *LoansRepo) UpdateLoan(ctx context.Context, plan domain.LoanPaymentPlan) (db.Loan, error) {
	loanParams, err := toLoanUpdateQueryParams(plan)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Error preparing params for insert query: %v", err)
	}

	queryResult, err := r.queries.UpdateLoan(ctx, loanParams)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Failed to update loan on database: %v", err)
	}

	err = r.queries.DeleteLoanStatesByLoanID(ctx, loanParams.ID)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Error deleting old payment plan data: %v", err)
	}

	for _, status := range plan.Plan {
		_, err := r.queries.CreateLoanState(ctx, toLoanStateInsertParams(status, queryResult.ID))
		if err != nil {
			return db.Loan{}, fmt.Errorf("Failed to save loan status to database: %v", err)
		}
	}
	return queryResult, nil
}

func (r *LoansRepo) DeleteLoan(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) error {
	return r.queries.DeleteLoan(ctx, db.DeleteLoanParams(toLoanGetParams(loanID, userID)))
}

func toLoanInsertQueryParams(plan domain.LoanPaymentPlan) (db.CreateLoanParams, error) {
	startDate, err := time.Parse("2006-01-02", plan.OriginalData.StartDate)
	if err != nil {
		return db.CreateLoanParams{}, err
	}
	return db.CreateLoanParams{
		UserID: pgtype.UUID{
			Bytes: plan.UserID,
			Valid: true,
		},
		Name:               plan.Name,
		StartingPrincipal:  int32(plan.OriginalData.StartingPrincipal),
		YearlyInterestRate: plan.OriginalData.YearlyInterestRate,
		MonthlyPayment:     int32(plan.OriginalData.MonthlyPayment),
		EscrowPayment:      int32(plan.OriginalData.EscrowPayment),
		StartDate: pgtype.Timestamptz{
			Time:  startDate,
			Valid: true,
		},
		MonthlyInterestRate: multiplierToPercent(plan.InterestMultiplierM),
		DurationMonths:      int32(plan.DurationMonths),
		TotalExpenditure:    int32(plan.TotalExpenditure.Round(0).IntPart()),
		TotalPaid:           int32(plan.TotalPaid.Round(0).IntPart()),
		CostOfCredit:        plan.CostOfCreditPercent.String(),
	}, nil
}

func toLoanUpdateQueryParams(plan domain.LoanPaymentPlan) (db.UpdateLoanParams, error) {
	startDate, err := time.Parse("2006-01-02", plan.OriginalData.StartDate)
	if err != nil {
		return db.UpdateLoanParams{}, err
	}
	return db.UpdateLoanParams{
		ID: pgtype.UUID{
			Bytes: plan.ID,
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: plan.UserID,
			Valid: true,
		},
		Name:               plan.Name,
		StartingPrincipal:  int32(plan.OriginalData.StartingPrincipal),
		YearlyInterestRate: plan.OriginalData.YearlyInterestRate,
		MonthlyPayment:     int32(plan.OriginalData.MonthlyPayment),
		EscrowPayment:      int32(plan.OriginalData.EscrowPayment),
		StartDate: pgtype.Timestamptz{
			Time:  startDate,
			Valid: true,
		},
		MonthlyInterestRate: multiplierToPercent(plan.InterestMultiplierM),
		DurationMonths:      int32(plan.DurationMonths),
		TotalExpenditure:    int32(plan.TotalExpenditure.Round(0).IntPart()),
		TotalPaid:           int32(plan.TotalPaid.Round(0).IntPart()),
		CostOfCredit:        plan.CostOfCreditPercent.String(),
	}, nil
}

func toLoanStateInsertParams(status domain.LoanStatus, loanID pgtype.UUID) db.CreateLoanStateParams {
	params := db.CreateLoanStateParams{
		LoanID: loanID,
		Date: pgtype.Timestamptz{
			Time:  status.Date,
			Valid: true,
		},
		Payment:       int32(status.Payment.Round(0).IntPart()),
		Interest:      int32(status.Interest.Round(0).IntPart()),
		OtherPayments: int32(status.OtherPayments.Round(0).IntPart()),
		Paydown:       int32(status.Paydown.Round(0).IntPart()),
		Principal:     int32(status.Principal.Round(0).IntPart()),
	}
	return params
}

func toLoanPaymentPlan(queryResult db.Loan) (domain.LoanPaymentPlan, error) {
	originalPlanData := domain.LoansInput{
		StartingPrincipal:  int(queryResult.StartingPrincipal),
		YearlyInterestRate: queryResult.YearlyInterestRate,
		MonthlyPayment:     int(queryResult.MonthlyPayment),
		EscrowPayment:      int(queryResult.EscrowPayment),
		StartDate:          queryResult.StartDate.Time.Format(time.RFC3339),
	}
	costOfCredit, err := decimal.NewFromString(queryResult.CostOfCredit)
	if err != nil {
		return domain.LoanPaymentPlan{}, fmt.Errorf("corrupted cost of credit data for savings plan: %v", err)
	}
	plan := domain.LoanPaymentPlan{
		ID:                  queryResult.ID.Bytes,
		UserID:              queryResult.UserID.Bytes,
		Name:                queryResult.Name,
		OriginalData:        originalPlanData,
		StartingPrincipal:   decimal.NewFromInt32(queryResult.StartingPrincipal),
		InterestMultiplierM: percentToMultiplier(queryResult.MonthlyInterestRate),
		PaymentM:            decimal.NewFromInt32(queryResult.MonthlyPayment),
		EscrowM:             decimal.NewFromInt32(queryResult.EscrowPayment),
		DurationMonths:      int(queryResult.DurationMonths),
		TotalExpenditure:    decimal.NewFromInt32(queryResult.TotalExpenditure),
		TotalPaid:           decimal.NewFromInt32(queryResult.TotalPaid),
		CostOfCreditPercent: costOfCredit,
	}

	return plan, nil
}

func toLoanGetParams(loanID uuid.UUID, userID uuid.UUID) db.GetLoanParams {
	return db.GetLoanParams{
		ID: pgtype.UUID{
			Bytes: loanID,
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: userID,
			Valid: true,
		},
	}
}

func toInitialLoanDataGetParams(loanID uuid.UUID, userID uuid.UUID) db.GetLoanInitialDataParams {
	return db.GetLoanInitialDataParams{
		ID: pgtype.UUID{
			Bytes: loanID,
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: userID,
			Valid: true,
		},
	}
}
