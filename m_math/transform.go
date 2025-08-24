package m_math

import (
	"github.com/shopspring/decimal"
)

// NewFromFloat 从 float64 创建一个 Decimal 对象
func NewFromFloat(value float64) Decimal {
	return Decimal{value: decimal.NewFromFloat(value)}
}

// NewFromString 从字符串创建一个 Decimal 对象
func NewFromString(value string) (Decimal, error) {
	v, err := decimal.NewFromString(value)
	return Decimal{value: v}, err
}

// NewFromInt 从 int64 创建一个 Decimal 对象
func NewFromInt(value int64) Decimal {
	return Decimal{value: decimal.NewFromInt(value)}
}

/* ============= 转换为标准对象 ================= */

// String 返回 Decimal 对象的字符串表示
func (d Decimal) String() string {
	return d.value.String()
}

// Float64 将 Decimal 对象转换为 float64，自动保留 float64 支持的最大有效位数
func (d Decimal) Float64() float64 {
	const maxFloat64Digits = 15 // float64 建议有效数字 15 位
	rounded := d.value.Round(maxFloat64Digits)
	f, _ := rounded.Float64()
	return f
}

// IntPart 获取 Decimal 对象的整数部分
func (d Decimal) IntPart() int64 {
	return d.value.IntPart()
}
