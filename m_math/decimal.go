package m_math

/*
我想基于 shopspring/decimal 最新版 封装一个易于使用的计算库。
*/

import (
	"fmt"

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

// DivSafe 除法，除数为零时返回错误（非破坏性新增，不替换现有 Div）
func (d Decimal) DivSafe(other Decimal) (Decimal, error) {
	if other.value.IsZero() {
		return Zero, fmt.Errorf("divide by zero")
	}
	return Decimal{value: d.value.Div(other.value)}, nil
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

/* ============= 比较 ================= */

// Cmp 比较两个 Decimal 对象
// 返回值: -1 (小于), 0 (等于), 1 (大于)
func (d Decimal) Cmp(other Decimal) int {
	return d.value.Cmp(other.value)
}

// Equal 判断是否相等
func (d Decimal) Equal(other Decimal) bool {
	return d.Cmp(other) == 0
}

// GreaterThan 判断是否大于
func (d Decimal) GreaterThan(other Decimal) bool {
	return d.Cmp(other) == 1
}

// LessThan 判断是否小于
func (d Decimal) LessThan(other Decimal) bool {
	return d.Cmp(other) == -1
}

// Neg 返回取负值 正数变负数，负数变正数，0 的相反数仍是 0
func (d Decimal) Neg() Decimal {
	return Decimal{value: d.value.Neg()}
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
