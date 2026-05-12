package mapper

import (
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/google/uuid"
)

func ToLoanResponse(plan domain.LoanPaymentPlan) dto.LoanResponseParams {
	response := dto.LoanResponseParams{}

	response.DurationMonths = plan.DurationMonths
	response.TotalExpenditure = int(plan.TotalExpenditure.Round(0).IntPart())
	response.TotalPaid = int(plan.TotalPaid.Round(0).IntPart())
	response.CostOfCreditPercent = plan.CostOfCreditPercent.Round(2).String()
	for _, status := range plan.Plan {
		response.Plan = append(response.Plan, dto.LoanStatus{
			Date:          status.Date,
			Payment:       int(status.Payment.Round(0).IntPart()),
			Interest:      int(status.Interest.Round(0).IntPart()),
			OtherPayments: int(status.OtherPayments.Round(0).IntPart()),
			Paydown:       int(status.Paydown.Round(0).IntPart()),
			Principal:     int(status.Principal.Round(0).IntPart()),
		})
	}
	return response
}

func ToSaveLoanResponse(loan db.Loan) dto.LoanSaveResponseParams {
	return dto.LoanSaveResponseParams{
		ID:                  loan.ID.String(),
		Name:                loan.Name,
		StartingPrincipal:   int(loan.StartingPrincipal),
		YearlyInterestRate:  loan.YearlyInterestRate,
		MonthlyPayment:      int(loan.MonthlyPayment),
		EscrowPayment:       int(loan.EscrowPayment),
		StartDate:           loan.StartDate.Time.Format(time.RFC3339),
		DurationMonths:      int(loan.DurationMonths),
		TotalExpenditure:    int(loan.TotalExpenditure),
		TotalPaid:           int(loan.TotalPaid),
		CostOfCreditPercent: loan.CostOfCredit,
	}
}

func ToLoanListResponse(rows []db.GetLoansByUserIDRow) dto.LoanListResponseParams {
	params := dto.LoanListResponseParams{}
	for _, row := range rows {
		newRow := dto.LoanInfo{
			ID:         row.ID.String(),
			Name:       row.Name,
			LoanAmount: int(row.StartingPrincipal),
		}
		params.Loans = append(params.Loans, newRow)
	}
	return params
}

func ToGetLoanResponse(plan domain.LoanPaymentPlan) dto.SavedLoanResponseParams {
	originalParams := dto.OriginalLoanData{
		StartingPrincipal:  plan.OriginalData.StartingPrincipal,
		YearlyInterestRate: plan.OriginalData.YearlyInterestRate,
		MonthlyPayment:     plan.OriginalData.MonthlyPayment,
		EscrowPayment:      plan.OriginalData.EscrowPayment,
		StartDate:          plan.OriginalData.StartDate,
	}
	monthlyInterestRate := multiplierToHighPrecisionPercent(plan.InterestMultiplierM)
	calculatedParams := dto.CalculatedLoanData{
		MonthlyInterestRate: monthlyInterestRate,
		DurationMonths:      plan.DurationMonths,
		TotalExpenditure:    int(plan.TotalExpenditure.Round(0).IntPart()),
		TotalPaid:           int(plan.TotalPaid.Round(0).IntPart()),
		CostOfCredit:        plan.CostOfCreditPercent.String(),
	}
	params := dto.SavedLoanResponseParams{
		ID:             plan.ID.String(),
		Name:           plan.Name,
		OriginalData:   originalParams,
		CalculatedData: calculatedParams,
	}
	for _, status := range plan.Plan {
		params.PaymentPlan = append(params.PaymentPlan, dto.LoanStatus{
			Date:          status.Date,
			Payment:       int(status.Payment.Round(0).IntPart()),
			Interest:      int(status.Interest.Round(0).IntPart()),
			OtherPayments: int(status.OtherPayments.Round(0).IntPart()),
			Paydown:       int(status.Paydown.Round(0).IntPart()),
			Principal:     int(status.Principal.Round(0).IntPart()),
		})
	}
	return params
}

func ToLoanInput(input dto.LoanRequestParams) domain.LoansInput {
	loan := domain.LoansInput{
		StartingPrincipal:  input.StartingPrincipal,
		YearlyInterestRate: input.YearlyInterestRate,
		MonthlyPayment:     input.MonthlyPayment,
		EscrowPayment:      input.EscrowPayment,
		StartDate:          input.StartDate,
	}

	return loan
}

func ToSaveLoanInput(userId uuid.UUID, input dto.LoanSaveRequestParams) domain.SaveLoanInput {
	loan := domain.SaveLoanInput{
		UserID:             userId,
		LoanName:           input.Name,
		StartingPrincipal:  input.StartingPrincipal,
		YearlyInterestRate: input.YearlyInterestRate,
		MonthlyPayment:     input.MonthlyPayment,
		EscrowPayment:      input.EscrowPayment,
		StartDate:          input.StartDate,
	}
	return loan
}

func ToUpdateLoanInput(loanID uuid.UUID, userId uuid.UUID, input dto.LoanUpdateRequestParams) domain.UpdateLoanInput {
	loan := domain.UpdateLoanInput{
		ID:                 loanID,
		UserID:             userId,
		LoanName:           input.Name,
		StartingPrincipal:  input.StartingPrincipal,
		YearlyInterestRate: input.YearlyInterestRate,
		MonthlyPayment:     input.MonthlyPayment,
		EscrowPayment:      input.EscrowPayment,
		StartDate:          input.StartDate,
	}
	return loan
}
