package mapper

import (
	"fmt"
	"time"

	"github.com/Mr-Rafael/bucktracker-api/internal/db"
	"github.com/Mr-Rafael/bucktracker-api/internal/domain"
	"github.com/Mr-Rafael/bucktracker-api/internal/dto"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func ToLoanResponse(loan domain.Loan) dto.LoanResponseParams {
	response := dto.LoanResponseParams{}
	if loan.DefaultPaymentPlan == nil {
		return response
	}

	plan := loan.DefaultPaymentPlan
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

func ToCreateLoanResponse(loan domain.Loan) dto.LoanCreateResponseParams {
	response := dto.LoanCreateResponseParams{
		ID:                 loan.ID.String(),
		Name:               loan.Name,
		StartingPrincipal:  loan.OriginalData.StartingPrincipal,
		YearlyInterestRate: loan.OriginalData.YearlyInterestRate,
		PaymentPlans:       []dto.PaymentPlanSummary{},
	}
	if loan.DefaultPaymentPlan != nil {
		summary := dto.PaymentPlanSummary{
			ID:             loan.DefaultPaymentPlan.ID.String(),
			Name:           loan.DefaultPaymentPlan.Name,
			DurationMonths: loan.DefaultPaymentPlan.DurationMonths,
		}
		response.DefaultPaymentPlan = summary
		response.PaymentPlans = append(response.PaymentPlans, summary)
	}
	return response
}

func ToSaveLoanResponse(loan domain.Loan) dto.LoanSaveResponseParams {
	response := dto.LoanSaveResponseParams{
		ID:                 loan.ID.String(),
		Name:               loan.Name,
		StartingPrincipal:  loan.OriginalData.StartingPrincipal,
		YearlyInterestRate: loan.OriginalData.YearlyInterestRate,
		MonthlyPayment:     loan.OriginalData.MonthlyPayment,
		EscrowPayment:      loan.OriginalData.EscrowPayment,
		StartDate:          loan.OriginalData.StartDate,
	}
	if loan.DefaultPaymentPlan != nil {
		response.DurationMonths = loan.DefaultPaymentPlan.DurationMonths
		response.TotalExpenditure = int(loan.DefaultPaymentPlan.TotalExpenditure.Round(0).IntPart())
		response.TotalPaid = int(loan.DefaultPaymentPlan.TotalPaid.Round(0).IntPart())
		response.CostOfCreditPercent = loan.DefaultPaymentPlan.CostOfCreditPercent.String()
	}
	return response
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

func ToPaymentPlanDetailResponse(plan domain.LoanPaymentPlan) dto.PaymentPlanDetailResponseParams {
	response := dto.PaymentPlanDetailResponseParams{
		ID:                   plan.ID.String(),
		Name:                 plan.Name,
		DurationMonths:       plan.DurationMonths,
		TotalExpenditure:     int(plan.TotalExpenditure.Round(0).IntPart()),
		TotalPaid:            int(plan.TotalPaid.Round(0).IntPart()),
		CostOfCredit:         plan.CostOfCreditPercent.String(),
		PrincipalPayments:    []dto.PrincipalPaymentParams{},
		PaymentPlanBreakdown: []dto.LoanStatus{},
	}
	for _, payment := range plan.PrincipalPayments {
		response.PrincipalPayments = append(response.PrincipalPayments, dto.PrincipalPaymentParams{
			Date:   payment.Date.Format("2006-01-02"),
			Amount: int(payment.AmountPaid.Round(0).IntPart()),
		})
	}
	for _, status := range plan.Plan {
		response.PaymentPlanBreakdown = append(response.PaymentPlanBreakdown, dto.LoanStatus{
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

func ToGetLoanResponse(loan domain.Loan) dto.SavedLoanResponseParams {
	params := dto.SavedLoanResponseParams{
		ID:   loan.ID.String(),
		Name: loan.Name,
		OriginalData: dto.OriginalLoanData{
			StartingPrincipal:  loan.OriginalData.StartingPrincipal,
			YearlyInterestRate: loan.OriginalData.YearlyInterestRate,
			MonthlyPayment:     loan.OriginalData.MonthlyPayment,
			EscrowPayment:      loan.OriginalData.EscrowPayment,
			StartDate:          loan.OriginalData.StartDate,
		},
		PaymentPlans: []dto.PaymentPlanSummary{},
	}
	if loan.DefaultPaymentPlan != nil {
		params.DefaultPaymentPlan = dto.PaymentPlanSummary{
			ID:             loan.DefaultPaymentPlan.ID.String(),
			Name:           loan.DefaultPaymentPlan.Name,
			DurationMonths: loan.DefaultPaymentPlan.DurationMonths,
		}
	}
	if len(loan.PaymentPlans) > 0 {
		for _, plan := range loan.PaymentPlans {
			params.PaymentPlans = append(params.PaymentPlans, dto.PaymentPlanSummary{
				ID:             plan.ID.String(),
				Name:           plan.Name,
				DurationMonths: plan.DurationMonths,
			})
		}
	} else if loan.DefaultPaymentPlan != nil {
		params.PaymentPlans = append(params.PaymentPlans, params.DefaultPaymentPlan)
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

func ToCreatePaymentPlanInput(loanID uuid.UUID, userID uuid.UUID, input dto.CreatePaymentPlanRequestParams) (domain.CreatePaymentPlanInput, error) {
	principalPayments := make([]domain.PrincipalPayment, 0, len(input.PrincipalPayments))
	for _, payment := range input.PrincipalPayments {
		date, err := time.Parse("2006-01-02", payment.Date)
		if err != nil {
			return domain.CreatePaymentPlanInput{}, fmt.Errorf("invalid principal payment date: %v", payment.Date)
		}
		principalPayments = append(principalPayments, domain.PrincipalPayment{
			AmountPaid: decimal.NewFromInt(int64(payment.Amount)),
			Date:       date,
		})
	}
	return domain.CreatePaymentPlanInput{
		LoanID:            loanID,
		UserID:            userID,
		Name:              input.Name,
		PrincipalPayments: principalPayments,
	}, nil
}
