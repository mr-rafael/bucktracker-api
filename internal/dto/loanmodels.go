package dto

import "time"

type LoanRequestParams struct {
	StartingPrincipal  int    `json:"startingPrincipal"`
	YearlyInterestRate string `json:"yearlyInterestRate"`
	MonthlyPayment     int    `json:"monthlyPayment"`
	EscrowPayment      int    `json:"escrowPayment"`
	StartDate          string `json:"startDate"`
}

type LoanResponseParams struct {
	DurationMonths      int          `json:"durationMonths"`
	TotalExpenditure    int          `json:"totalExpenditure"`
	TotalPaid           int          `json:"totalPaid"`
	CostOfCreditPercent string       `json:"costOfCreditPercent"`
	Plan                []LoanStatus `json:"plan"`
}

type SavedLoanResponseParams struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	OriginalData       OriginalLoanData     `json:"originalData"`
	DefaultPaymentPlan PaymentPlanSummary   `json:"defaultPaymentPlan"`
	PaymentPlans       []PaymentPlanSummary `json:"paymentPlans"`
}

type OriginalLoanData struct {
	StartingPrincipal  int    `json:"startingPrincipal"`
	YearlyInterestRate string `json:"yearlyInterestRate"`
	MonthlyPayment     int    `json:"monthlyPayment"`
	EscrowPayment      int    `json:"escrowPayment"`
	StartDate          string `json:"startDate"`
}

type LoanStatus struct {
	Date          time.Time `json:"date"`
	Payment       int       `json:"payment"`
	Interest      int       `json:"interest"`
	OtherPayments int       `json:"otherPayments"`
	Paydown       int       `json:"paydown"`
	Principal     int       `json:"principal"`
}

type LoanSaveRequestParams struct {
	Name               string `json:"name"`
	StartingPrincipal  int    `json:"startingPrincipal"`
	YearlyInterestRate string `json:"yearlyInterestRate"`
	MonthlyPayment     int    `json:"monthlyPayment"`
	EscrowPayment      int    `json:"escrowPayment"`
	StartDate          string `json:"startDate"`
}

type PaymentPlanSummary struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	DurationMonths int    `json:"durationMonths"`
}

type LoanCreateResponseParams struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	StartingPrincipal  int                  `json:"startingPrincipal"`
	YearlyInterestRate string               `json:"yearlyInterestRate"`
	DefaultPaymentPlan PaymentPlanSummary   `json:"defaultPaymentPlan"`
	PaymentPlans       []PaymentPlanSummary `json:"paymentPlans"`
}

// LoanSaveResponseParams is used by update responses.
type LoanSaveResponseParams struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	StartingPrincipal   int    `json:"startingPrincipal"`
	YearlyInterestRate  string `json:"yearlyInterestRate"`
	MonthlyPayment      int    `json:"monthlyPayment"`
	EscrowPayment       int    `json:"escrowPayment"`
	StartDate           string `json:"startDate"`
	DurationMonths      int    `json:"durationMonths"`
	TotalExpenditure    int    `json:"totalExpenditure"`
	TotalPaid           int    `json:"totalPaid"`
	CostOfCreditPercent string `json:"costOfCreditPercent"`
}

type LoanUpdateRequestParams struct {
	Name               *string `json:"name"`
	StartingPrincipal  *int    `json:"startingPrincipal"`
	YearlyInterestRate *string `json:"yearlyInterestRate"`
	MonthlyPayment     *int    `json:"monthlyPayment"`
	EscrowPayment      *int    `json:"escrowPayment"`
	StartDate          *string `json:"startDate"`
}

type LoanListResponseParams struct {
	Loans []LoanInfo `json:"loans"`
}

type LoanInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	LoanAmount int    `json:"loanAmount"`
}
