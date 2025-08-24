package m_math

/*
我想基于 shopspring/decimal 最新版 封装一个易于使用的计算库。
*/

import (
	"github.com/shopspring/decimal"
)

// Decimal 封装了 shopspring/decimal 库的 Decimal 类型，提供更易用的接口
type Decimal struct {
	value decimal.Decimal
}

// 常量定义
var (
	Zero = Decimal{value: decimal.Zero}
	One  = Decimal{value: decimal.NewFromInt(1)}
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

/* =============四则运算================= */

// Add 加法运算
func (d Decimal) Add(other Decimal) Decimal {
	return Decimal{value: d.value.Add(other.value)}
}

// Sub 减法运算
func (d Decimal) Sub(other Decimal) Decimal {
	return Decimal{value: d.value.Sub(other.value)}
}

// Mul 乘法运算
func (d Decimal) Mul(other Decimal) Decimal {
	return Decimal{value: d.value.Mul(other.value)}
}

// Div 除法运算，除数为零时返回零值
func (d Decimal) Div(other Decimal) Decimal {
	if other.value.IsZero() {
		// 可以选择返回零值，也可以返回 d 本身，或 panic，建议返回零值
		return Zero
	}
	return Decimal{value: d.value.Div(other.value)}
}

/* ============= 转换为标准对象 ================= */

// String 返回 Decimal 对象的字符串表示
func (d Decimal) String() string {
	return d.value.String()
}

// Float64 将 Decimal 对象转换为 float64，自动保留 float64 支持的最大有效位数
func (d Decimal) Float64() float64 {
	const maxFloat64Digits = 17 // float64 最大有效数字位数
	rounded := d.value.Round(maxFloat64Digits)
	f, _ := rounded.Float64()
	return f
}

// IntPart 获取 Decimal 对象的整数部分
func (d Decimal) IntPart() int64 {
	return d.value.IntPart()
}

// Cmp 比较两个 Decimal 对象
// 返回值: -1 (小于), 0 (等于), 1 (大于)
func (d Decimal) Cmp(other Decimal) int {
	return d.value.Cmp(other.value)
}

// Round 四舍五入到指定小数位数
func (d Decimal) Round(places int32) Decimal {
	return Decimal{value: d.value.Round(places)}
}

// Abs 返回绝对值
func (d Decimal) Abs() Decimal {
	return Decimal{value: d.value.Abs()}
}

// IsZero 判断是否为零值
func (d Decimal) IsZero() bool {
	return d.value.IsZero()
}

// IsPositive 判断是否为正值
func (d Decimal) IsPositive() bool {
	return d.value.IsPositive()
}

// IsNegative 判断是否为负值
func (d Decimal) IsNegative() bool {
	return d.value.IsNegative()
}
