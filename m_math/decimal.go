package m_math

import (
	"github.com/shopspring/decimal"
)

// Decimal 封装了 shopspring/decimal 库的 Decimal 类型，提供更易用的接口
type Decimal struct {
	value decimal.Decimal
}

// Value 返回底层 shopspring/decimal.Decimal（便于需要原生 API 的场景）
func (d Decimal) Value() decimal.Decimal {
	return d.value
}

// 常量定义
var (
	Zero = Decimal{value: decimal.Zero}
	One  = Decimal{value: decimal.NewFromInt(1)}
)
