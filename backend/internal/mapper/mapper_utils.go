package mapper

import "github.com/shopspring/decimal"

func multiplierToHighPrecisionPercent(mult decimal.Decimal) string {
	oneHundred := decimal.NewFromInt(100)
	return mult.Mul(oneHundred).Round(10).String()
}
