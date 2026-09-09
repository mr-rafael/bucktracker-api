package service

import "github.com/shopspring/decimal"

func getMonthlyAPYMultiplier(APY string) (decimal.Decimal, error) {
	decimalInterest, err := decimal.NewFromString(APY)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	decimalInterest = decimalInterest.Div(decimal.NewFromInt(100)).Add(decimal.NewFromInt(1))
	decimalInterest = decimalInterest.Pow(decimal.NewFromFloat(1.0 / 12.0)).Sub(decimal.NewFromInt(1))
	return decimalInterest, nil
}

func getMonthlyAPRMultiplier(APR string) (decimal.Decimal, error) {
	decimalInterest, err := decimal.NewFromString(APR)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	decimalInterest = decimalInterest.Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12))
	return decimalInterest, nil
}

func toTaxMultiplier(rate string) (decimal.Decimal, error) {
	if rate == "" {
		return decimal.NewFromInt(0), nil
	}
	decimalTax, err := decimal.NewFromString(rate)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	decimalTax = decimalTax.Div(decimal.NewFromInt(100))
	return decimalTax, nil
}

func toInflationMultiplier(rate string) (decimal.Decimal, error) {
	if rate == "" {
		return decimal.NewFromInt(1), nil
	}
	decimalInflation, err := decimal.NewFromString(rate)
	if err != nil {
		return decimal.NewFromInt(1), err
	}
	decimalInflation = decimalInflation.Div(decimal.NewFromInt(100)).Add(decimal.NewFromInt(1))
	return decimalInflation, nil
}

func getReturnPercent(rate decimal.Decimal) string {
	returnPercent := rate.Sub(decimal.NewFromInt(1)).Mul(decimal.NewFromInt(100))
	return returnPercent.Round(2).String()
}

func decimalIsBetween(dec decimal.Decimal, lower string, upper string) bool {
	lowerDecimal, _ := decimal.NewFromString(lower)
	upperDecimal, _ := decimal.NewFromString(upper)
	isAboveLower := dec.Compare(lowerDecimal) >= 0
	isBelowUpper := upperDecimal.Compare(dec) >= 0
	return isAboveLower && isBelowUpper
}

func stringNumberBetween(number string, lower string, upper string) bool {
	numDecimal, _ := decimal.NewFromString(number)
	lowerDecimal, _ := decimal.NewFromString(lower)
	upperDecimal, _ := decimal.NewFromString(upper)
	isAboveLower := numDecimal.Compare(lowerDecimal) >= 0
	isBelowUpper := upperDecimal.Compare(numDecimal) >= 0
	return isAboveLower && isBelowUpper
}

func toMonthlyInterestMultiplier(yearlyInterestRate string, interestRateType string) (decimal.Decimal, error) {
	if interestRateType == "APR" {
		return getMonthlyAPRMultiplier(yearlyInterestRate)
	}
	return getMonthlyAPYMultiplier(yearlyInterestRate)
}
