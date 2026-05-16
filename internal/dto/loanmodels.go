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
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	OriginalData   OriginalLoanData   `json:"originalData"`
	CalculatedData CalculatedLoanData `json:"calculatedData"`
	PaymentPlan    []LoanStatus       `json:"paymentPlan"`
}

type OriginalLoanData struct {
	StartingPrincipal  int    `json:"startingPrincipal"`
	YearlyInterestRate string `json:"yearlyInterestRate"`
	MonthlyPayment     int    `json:"monthlyPayment"`
	EscrowPayment      int    `json:"escrowPayment"`
	StartDate          string `json:"startDate"`
}

type CalculatedLoanData struct {
	MonthlyInterestRate string `json:"monthlyInterestRate"`
	DurationMonths      int    `json:"durationMonths"`
	TotalExpenditure    int    `json:"totalExpenditure"`
	TotalPaid           int    `json:"totalPaid"`
	CostOfCredit        string `json:"costOfCredit"`
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
