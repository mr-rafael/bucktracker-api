package service

import (
	"context"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type MockAuthRepo struct {
	CreateRefreshTokenFunc  func(context.Context, pgtype.UUID, string, time.Time) (db.RefreshToken, error)
	GetTokenByHashFunc      func(context.Context, string) (db.RefreshToken, error)
	RevokeTokenByUserIDFunc func(context.Context, pgtype.UUID) error
}

type MockUsersRepo struct {
	CreateUserFunc     func(context.Context, db.CreateUserParams) (db.User, error)
	GetUserByEmailFunc func(context.Context, string) (db.User, error)
	GetUserByIDFunc    func(context.Context, pgtype.UUID) (db.User, error)
	DeleteUserFunc     func(context.Context, pgtype.UUID) error
}

type MockLoansRepo struct {
	SaveLoanPaymentPlanFunc       func(context.Context, domain.LoanPaymentPlan) (db.Loan, error)
	GetLoanPaymentPlansByUserFunc func(context.Context, uuid.UUID) ([]db.GetLoansByUserIDRow, error)
	GetLoanByIDFunc               func(context.Context, uuid.UUID, uuid.UUID) (domain.LoanPaymentPlan, error)
	GetLoanInitialDataFunc        func(context.Context, uuid.UUID, uuid.UUID) (domain.LoansInput, error)
	UpdateLoanFunc                func(context.Context, domain.LoanPaymentPlan) (db.Loan, error)
	DeleteLoanFunc                func(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) error
}

type MockSavingsRepo struct {
	SaveSavingsPlanFunc       func(context.Context, domain.SavingsPlan) (db.Saving, error)
	GetSavingsPlansByUserFunc func(context.Context, uuid.UUID) ([]db.GetSavingsByUserIDRow, error)
	GetSavingsPlanByIDFunc    func(context.Context, uuid.UUID, uuid.UUID) (domain.SavingsPlan, error)
	DeleteSavingsPlanFunc     func(context.Context, uuid.UUID, uuid.UUID) error
}

func (m *MockAuthRepo) CreateRefreshToken(ctx context.Context, userID pgtype.UUID, tokenHash string, expDate time.Time) (db.RefreshToken, error) {
	if m.CreateRefreshTokenFunc != nil {
		return m.CreateRefreshTokenFunc(ctx, userID, tokenHash, expDate)
	}
	return db.RefreshToken{}, nil
}

func (m *MockAuthRepo) GetTokenByHash(ctx context.Context, hash string) (db.RefreshToken, error) {
	if m.GetTokenByHashFunc != nil {
		return m.GetTokenByHashFunc(ctx, hash)
	}
	return db.RefreshToken{}, nil
}

func (m *MockAuthRepo) RevokeTokenByUserID(ctx context.Context, id pgtype.UUID) error {
	if m.RevokeTokenByUserIDFunc != nil {
		return m.RevokeTokenByUserIDFunc(ctx, id)
	}
	return nil
}

func (m *MockUsersRepo) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, params)
	}
	return db.User{}, nil
}

func (m *MockUsersRepo) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return db.User{}, nil
}

func (m *MockUsersRepo) GetUserByID(ctx context.Context, id pgtype.UUID) (db.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return db.User{}, nil
}

func (m *MockUsersRepo) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(ctx, id)
	}
	return nil
}

func (m *MockLoansRepo) SaveLoanPaymentPlan(ctx context.Context, plan domain.LoanPaymentPlan) (db.Loan, error) {
	if m.SaveLoanPaymentPlanFunc != nil {
		return m.SaveLoanPaymentPlanFunc(ctx, plan)
	}
	return db.Loan{}, nil
}

func (m *MockLoansRepo) GetLoanPaymentPlansByUser(ctx context.Context, id uuid.UUID) ([]db.GetLoansByUserIDRow, error) {
	if m.GetLoanPaymentPlansByUserFunc != nil {
		return m.GetLoanPaymentPlansByUserFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockLoansRepo) GetLoanByID(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) (domain.LoanPaymentPlan, error) {
	if m.GetLoanByIDFunc != nil {
		return m.GetLoanByIDFunc(ctx, loanID, userID)
	}
	return domain.LoanPaymentPlan{}, nil
}

func (m *MockLoansRepo) GetLoanInitialData(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) (domain.LoansInput, error) {
	if m.GetLoanByIDFunc != nil {
		return m.GetLoanInitialDataFunc(ctx, loanID, userID)
	}
	return domain.LoansInput{}, nil
}

func (m *MockLoansRepo) UpdateLoan(ctx context.Context, plan domain.LoanPaymentPlan) (db.Loan, error) {
	if m.SaveLoanPaymentPlanFunc != nil {
		return m.UpdateLoanFunc(ctx, plan)
	}
	return db.Loan{}, nil
}

func (m *MockLoansRepo) DeleteLoan(ctx context.Context, loanID uuid.UUID, userID uuid.UUID) error {
	if m.GetLoanByIDFunc != nil {
		return m.DeleteLoanFunc(ctx, loanID, userID)
	}
	return nil
}

func (m *MockSavingsRepo) SaveSavingsPlan(ctx context.Context, plan domain.SavingsPlan) (db.Saving, error) {
	if m.SaveSavingsPlanFunc != nil {
		return m.SaveSavingsPlanFunc(ctx, plan)
	}
	return db.Saving{}, nil
}

func (m *MockSavingsRepo) GetSavingsPlansByUser(ctx context.Context, userID uuid.UUID) ([]db.GetSavingsByUserIDRow, error) {
	if m.GetSavingsPlansByUserFunc != nil {
		return m.GetSavingsPlansByUserFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockSavingsRepo) GetSavingsPlanByID(ctx context.Context, planID uuid.UUID, userID uuid.UUID) (domain.SavingsPlan, error) {
	if m.GetSavingsPlanByIDFunc != nil {
		return m.GetSavingsPlanByIDFunc(ctx, planID, userID)
	}
	return domain.SavingsPlan{}, nil
}

func (m *MockSavingsRepo) DeleteSavingsPlan(ctx context.Context, planID uuid.UUID, userID uuid.UUID) error {
	if m.DeleteSavingsPlanFunc != nil {
		return m.DeleteSavingsPlanFunc(ctx, planID, userID)
	}
	return nil
}
