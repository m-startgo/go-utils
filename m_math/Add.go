package m_math

import (
	"github.com/shopspring/decimal"
)

func toDec(s string) decimal.Decimal {
	n, _ := decimal.NewFromString(s)
	return n
}

// a+b
func Add(a, b string) string {
	n := toDec(a).Add(toDec(b))
	return n.String()
}
