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

type LoansRepo struct {
	queries *db.Queries
}

func NewLoansRepo(queries *db.Queries) *LoansRepo {
	return &LoansRepo{queries: queries}
}

func (r *LoansRepo) SaveLoanPaymentPlan(ctx context.Context, loan domain.Loan) (db.Loan, error) {
	if loan.DefaultPaymentPlan == nil {
		return db.Loan{}, fmt.Errorf("loan has no default payment plan")
	}

	loanParams, err := toLoanInsertQueryParams(loan)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Error preparing params for insert query: %v", err)
	}

	savedLoan, err := r.queries.CreateLoan(ctx, loanParams)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Failed to save to database: %v", err)
	}

	paymentPlan, err := r.queries.CreatePaymentPlan(ctx, toPaymentPlanInsertParams(*loan.DefaultPaymentPlan, savedLoan.ID))
	if err != nil {
		return db.Loan{}, fmt.Errorf("Failed to save payment plan to database: %v", err)
	}

	updateParams, err := toLoanUpdateQueryParams(loan, savedLoan.ID, paymentPlan.ID)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Error preparing params for update query: %v", err)
	}
	savedLoan, err = r.queries.UpdateLoan(ctx, updateParams)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Failed to set default payment plan: %v", err)
	}

	for _, status := range loan.DefaultPaymentPlan.Plan {
		_, err := r.queries.CreateLoanState(ctx, toLoanStateInsertParams(status, paymentPlan.ID))
		if err != nil {
			return db.Loan{}, fmt.Errorf("Failed to save loan status to database: %v", err)
		}
	}

	for _, principalPayment := range loan.DefaultPaymentPlan.PrincipalPayments {
		_, err := r.queries.CreatePrincipalPayment(ctx, toPrincipalPaymentInsertParams(principalPayment, paymentPlan.ID))
		if err != nil {
			return db.Loan{}, fmt.Errorf("Failed to save principal payment to database: %v", err)
		}
	}

	return savedLoan, nil
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

func (r *LoansRepo) GetPaymentPlanByID(ctx context.Context, loanID uuid.UUID, paymentPlanID uuid.UUID, userID uuid.UUID) (domain.LoanPaymentPlan, error) {
	_, err := r.queries.GetLoan(ctx, toLoanGetParams(loanID, userID))
	if err != nil {
		return domain.LoanPaymentPlan{}, fmt.Errorf("failed to fetch loan from database: %v", err)
	}

	paymentPlan, err := r.queries.GetPaymentPlanByIDAndLoanID(ctx, db.GetPaymentPlanByIDAndLoanIDParams{
		ID: pgtype.UUID{
			Bytes: paymentPlanID,
			Valid: true,
		},
		LoanID: pgtype.UUID{
			Bytes: loanID,
			Valid: true,
		},
	})
	if err != nil {
		return domain.LoanPaymentPlan{}, fmt.Errorf("failed to fetch payment plan from database: %v", err)
	}

	costOfCredit, err := decimal.NewFromString(paymentPlan.CostOfCredit)
	if err != nil {
		return domain.LoanPaymentPlan{}, fmt.Errorf("corrupted cost of credit data for loan payment plan: %v", err)
	}

	plan := domain.LoanPaymentPlan{
		ID:                  paymentPlan.ID.Bytes,
		Name:                paymentPlan.Name,
		DurationMonths:      int(paymentPlan.DurationMonths),
		TotalExpenditure:    decimal.NewFromInt32(paymentPlan.TotalExpenditure),
		TotalPaid:           decimal.NewFromInt32(paymentPlan.TotalPaid),
		CostOfCreditPercent: costOfCredit,
	}

	states, err := r.queries.GetLoanStatesByPaymentPlanID(ctx, paymentPlan.ID)
	if err != nil {
		return domain.LoanPaymentPlan{}, fmt.Errorf("failed to fetch loan states from database: %v", err)
	}
	for _, state := range states {
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

func (r *LoansRepo) GetLoanByID(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) (domain.Loan, error) {
	loanQueryResult, err := r.queries.GetLoan(ctx, toLoanGetParams(loanID, userID))
	if err != nil {
		return domain.Loan{}, fmt.Errorf("failed to fetch loan from database: %v", err)
	}
	if !loanQueryResult.DefaultPaymentPlan.Valid {
		return domain.Loan{}, fmt.Errorf("loan has no default payment plan")
	}

	paymentPlan, err := r.queries.GetPaymentPlan(ctx, loanQueryResult.DefaultPaymentPlan)
	if err != nil {
		return domain.Loan{}, fmt.Errorf("failed to fetch payment plan from database: %v", err)
	}

	loan, err := toDomainLoan(loanQueryResult, paymentPlan)
	if err != nil {
		return domain.Loan{}, err
	}

	allPlans, err := r.queries.GetPaymentPlansByLoanID(ctx, loanQueryResult.ID)
	if err != nil {
		return domain.Loan{}, fmt.Errorf("failed to fetch payment plans from database: %v", err)
	}
	for _, plan := range allPlans {
		loan.PaymentPlans = append(loan.PaymentPlans, domain.LoanPaymentPlan{
			ID:             plan.ID.Bytes,
			Name:           plan.Name,
			DurationMonths: int(plan.DurationMonths),
		})
	}

	return loan, nil
}

func (r *LoansRepo) GetLoanInitialData(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) (domain.UpdateLoanData, error) {
	loanQueryResult, err := r.queries.GetLoanInitialData(ctx, toInitialLoanDataGetParams(loanID, userID))
	if err != nil {
		return domain.UpdateLoanData{}, fmt.Errorf("failed to fetch loan from database: %v", err)
	}
	loansInput := domain.LoansInput{
		StartingPrincipal:  int(loanQueryResult.StartingPrincipal),
		YearlyInterestRate: loanQueryResult.YearlyInterestRate,
		MonthlyPayment:     int(loanQueryResult.MonthlyPayment),
		EscrowPayment:      int(loanQueryResult.EscrowPayment),
		StartDate:          loanQueryResult.StartDate.Time.Format("2006-01-02"),
	}
	loanData := domain.UpdateLoanData{
		ID:       loanID,
		Name:     loanQueryResult.Name,
		LoanData: loansInput,
	}

	return loanData, nil
}

func (r *LoansRepo) UpdateLoan(ctx context.Context, loan domain.Loan) (db.Loan, error) {
	if loan.DefaultPaymentPlan == nil {
		return db.Loan{}, fmt.Errorf("loan has no default payment plan")
	}

	existingLoan, err := r.queries.GetLoan(ctx, toLoanGetParams(loan.ID, loan.UserID))
	if err != nil {
		return db.Loan{}, fmt.Errorf("Failed to fetch loan for update: %v", err)
	}
	if !existingLoan.DefaultPaymentPlan.Valid {
		return db.Loan{}, fmt.Errorf("loan has no default payment plan")
	}

	loanParams, err := toLoanUpdateQueryParams(loan, existingLoan.ID, existingLoan.DefaultPaymentPlan)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Error preparing params for update query: %v", err)
	}

	queryResult, err := r.queries.UpdateLoan(ctx, loanParams)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Failed to update loan on database: %v", err)
	}

	_, err = r.queries.UpdatePaymentPlan(ctx, toPaymentPlanUpdateParams(*loan.DefaultPaymentPlan, existingLoan.DefaultPaymentPlan))
	if err != nil {
		return db.Loan{}, fmt.Errorf("Failed to update payment plan on database: %v", err)
	}

	err = r.queries.DeleteLoanStatesByPaymentPlanID(ctx, existingLoan.DefaultPaymentPlan)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Error deleting old payment plan data: %v", err)
	}

	for _, status := range loan.DefaultPaymentPlan.Plan {
		_, err := r.queries.CreateLoanState(ctx, toLoanStateInsertParams(status, existingLoan.DefaultPaymentPlan))
		if err != nil {
			return db.Loan{}, fmt.Errorf("Failed to save loan status to database: %v", err)
		}
	}

	err = r.queries.DeletePrincipalPaymentsByPaymentPlanID(ctx, existingLoan.DefaultPaymentPlan)
	if err != nil {
		return db.Loan{}, fmt.Errorf("Error deleting old principal payments: %v", err)
	}

	for _, principalPayment := range loan.DefaultPaymentPlan.PrincipalPayments {
		_, err := r.queries.CreatePrincipalPayment(ctx, toPrincipalPaymentInsertParams(principalPayment, existingLoan.DefaultPaymentPlan))
		if err != nil {
			return db.Loan{}, fmt.Errorf("Failed to save principal payment to database: %v", err)
		}
	}

	return queryResult, nil
}

func (r *LoansRepo) DeleteLoan(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) error {
	existingLoan, err := r.queries.GetLoan(ctx, toLoanGetParams(loanID, userID))
	if err != nil {
		return fmt.Errorf("Not found.")
	}

	// Clear the circular FK before deleting so payment_plans can cascade.
	clearParams := db.UpdateLoanParams{
		ID:                  existingLoan.ID,
		UserID:              existingLoan.UserID,
		Name:                existingLoan.Name,
		StartingPrincipal:   existingLoan.StartingPrincipal,
		YearlyInterestRate:  existingLoan.YearlyInterestRate,
		MonthlyPayment:      existingLoan.MonthlyPayment,
		EscrowPayment:       existingLoan.EscrowPayment,
		StartDate:           existingLoan.StartDate,
		MonthlyInterestRate: existingLoan.MonthlyInterestRate,
		DefaultPaymentPlan:  pgtype.UUID{Valid: false},
	}
	_, err = r.queries.UpdateLoan(ctx, clearParams)
	if err != nil {
		return fmt.Errorf("Failed to clear default payment plan before delete: %v", err)
	}

	rows, err := r.queries.DeleteLoan(ctx, db.DeleteLoanParams{
		ID:     existingLoan.ID,
		UserID: existingLoan.UserID,
	})
	if err != nil || rows <= 0 {
		return fmt.Errorf("Not found.")
	}
	return nil
}

func toLoanInsertQueryParams(loan domain.Loan) (db.CreateLoanParams, error) {
	startDate, err := time.Parse("2006-01-02", loan.OriginalData.StartDate)
	if err != nil {
		return db.CreateLoanParams{}, err
	}
	return db.CreateLoanParams{
		UserID: pgtype.UUID{
			Bytes: loan.UserID,
			Valid: true,
		},
		Name:               loan.Name,
		StartingPrincipal:  int32(loan.OriginalData.StartingPrincipal),
		YearlyInterestRate: loan.OriginalData.YearlyInterestRate,
		MonthlyPayment:     int32(loan.OriginalData.MonthlyPayment),
		EscrowPayment:      int32(loan.OriginalData.EscrowPayment),
		StartDate: pgtype.Timestamptz{
			Time:  startDate,
			Valid: true,
		},
		MonthlyInterestRate: multiplierToPercent(loan.InterestMultiplierM),
	}, nil
}

func toLoanUpdateQueryParams(loan domain.Loan, loanID pgtype.UUID, defaultPaymentPlanID pgtype.UUID) (db.UpdateLoanParams, error) {
	startDate, err := time.Parse("2006-01-02", loan.OriginalData.StartDate)
	if err != nil {
		return db.UpdateLoanParams{}, err
	}
	return db.UpdateLoanParams{
		ID: loanID,
		UserID: pgtype.UUID{
			Bytes: loan.UserID,
			Valid: true,
		},
		Name:               loan.Name,
		StartingPrincipal:  int32(loan.OriginalData.StartingPrincipal),
		YearlyInterestRate: loan.OriginalData.YearlyInterestRate,
		MonthlyPayment:     int32(loan.OriginalData.MonthlyPayment),
		EscrowPayment:      int32(loan.OriginalData.EscrowPayment),
		StartDate: pgtype.Timestamptz{
			Time:  startDate,
			Valid: true,
		},
		MonthlyInterestRate: multiplierToPercent(loan.InterestMultiplierM),
		DefaultPaymentPlan:  defaultPaymentPlanID,
	}, nil
}

const defaultPaymentPlanName = "Default Payment Plan"

func paymentPlanNameOrDefault(name string) string {
	if name == "" {
		return defaultPaymentPlanName
	}
	return name
}

func toPaymentPlanInsertParams(plan domain.LoanPaymentPlan, loanID pgtype.UUID) db.CreatePaymentPlanParams {
	return db.CreatePaymentPlanParams{
		LoanID:           loanID,
		Name:             paymentPlanNameOrDefault(plan.Name),
		DurationMonths:   int32(plan.DurationMonths),
		TotalExpenditure: int32(plan.TotalExpenditure.Round(0).IntPart()),
		TotalPaid:        int32(plan.TotalPaid.Round(0).IntPart()),
		CostOfCredit:     plan.CostOfCreditPercent.String(),
	}
}

func toPaymentPlanUpdateParams(plan domain.LoanPaymentPlan, paymentPlanID pgtype.UUID) db.UpdatePaymentPlanParams {
	return db.UpdatePaymentPlanParams{
		ID:               paymentPlanID,
		Name:             paymentPlanNameOrDefault(plan.Name),
		DurationMonths:   int32(plan.DurationMonths),
		TotalExpenditure: int32(plan.TotalExpenditure.Round(0).IntPart()),
		TotalPaid:        int32(plan.TotalPaid.Round(0).IntPart()),
		CostOfCredit:     plan.CostOfCreditPercent.String(),
	}
}

func toLoanStateInsertParams(status domain.LoanStatus, paymentPlanID pgtype.UUID) db.CreateLoanStateParams {
	return db.CreateLoanStateParams{
		PaymentPlanID: paymentPlanID,
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
}

func toPrincipalPaymentInsertParams(payment domain.PrincipalPayment, paymentPlanID pgtype.UUID) db.CreatePrincipalPaymentParams {
	return db.CreatePrincipalPaymentParams{
		PaymentPlanID: paymentPlanID,
		AmountPaid:    int32(payment.AmountPaid.Round(0).IntPart()),
		Date: pgtype.Timestamptz{
			Time:  payment.Date,
			Valid: true,
		},
	}
}

func toDomainLoan(loan db.Loan, paymentPlan db.PaymentPlan) (domain.Loan, error) {
	originalPlanData := domain.LoansInput{
		StartingPrincipal:  int(loan.StartingPrincipal),
		YearlyInterestRate: loan.YearlyInterestRate,
		MonthlyPayment:     int(loan.MonthlyPayment),
		EscrowPayment:      int(loan.EscrowPayment),
		StartDate:          loan.StartDate.Time.Format(time.RFC3339),
	}
	costOfCredit, err := decimal.NewFromString(paymentPlan.CostOfCredit)
	if err != nil {
		return domain.Loan{}, fmt.Errorf("corrupted cost of credit data for loan payment plan: %v", err)
	}
	return domain.Loan{
		ID:                  loan.ID.Bytes,
		UserID:              loan.UserID.Bytes,
		Name:                loan.Name,
		OriginalData:        originalPlanData,
		StartingPrincipal:   decimal.NewFromInt32(loan.StartingPrincipal),
		InterestMultiplierM: percentToMultiplier(loan.MonthlyInterestRate),
		PaymentM:            decimal.NewFromInt32(loan.MonthlyPayment),
		EscrowM:             decimal.NewFromInt32(loan.EscrowPayment),
		DefaultPaymentPlan: &domain.LoanPaymentPlan{
			ID:                  paymentPlan.ID.Bytes,
			Name:                paymentPlan.Name,
			DurationMonths:      int(paymentPlan.DurationMonths),
			TotalExpenditure:    decimal.NewFromInt32(paymentPlan.TotalExpenditure),
			TotalPaid:           decimal.NewFromInt32(paymentPlan.TotalPaid),
			CostOfCreditPercent: costOfCredit,
		},
	}, nil
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
