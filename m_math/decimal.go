package m_math

import (
	"fmt"

	"github.com/shopspring/decimal"
)

/*

示例用法

链式：
NewFromString("1.23").Add(NewFromFloat(2.0)).Mul(MustFromString("3")).StringFixed(2)

聚合：
Sum(a, b, c).Div(NewFromInt(3)).Round(2)

*/

// Decimal 是对 shopspring/decimal 的轻量封装，提供链式算术和常用工具函数。
type Decimal struct {
	d decimal.Decimal
}

// --------------- 构造器 ----------------

// Zero 返回 0
func Zero() Decimal { return Decimal{decimal.Zero} }

// NewFromFloat 从 float64 创建
func NewFromFloat(f float64) Decimal { return Decimal{decimal.NewFromFloat(f)} }

// NewFromInt 从 int64 创建
func NewFromInt(i int64) Decimal { return Decimal{decimal.NewFromInt(i)} }

// NewFromString 解析字符串，失败返回 error
func NewFromString(s string) (Decimal, error) {
	d, err := decimal.NewFromString(s)
	return Decimal{d}, err
}

// MustFromString 解析字符串，失败会 panic（便于初始化常量）
func MustFromString(s string) Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return Decimal{d}
}

// FromDecimal 从原生 decimal 创建
func FromDecimal(d decimal.Decimal) Decimal {
	return Decimal{d}
}

// --------------- 基本输出 ----------------

func (x Decimal) Decimal() decimal.Decimal {
	return x.d
}

func (x Decimal) String() string {
	return x.d.String()
}

// StringFixed 指定位数（四舍五入）输出
func (x Decimal) StringFixed(places int32) string {
	return x.d.StringFixed(places)
}

// Float64 转 float64（可能有精度丢失）
func (x Decimal) Float64() float64 {
	f, _ := x.d.Float64()
	return f
}

// Int64 转 int64（向零截断）
func (x Decimal) Int64() int64 {
	return x.d.IntPart()
}

// --------------- 算术（返回新 Decimal，支持链式） ----------------

func (x Decimal) Add(y Decimal) Decimal {
	return Decimal{x.d.Add(y.d)}
}

func (x Decimal) Sub(y Decimal) Decimal {
	return Decimal{x.d.Sub(y.d)}
}

func (x Decimal) Mul(y Decimal) Decimal {
	return Decimal{x.d.Mul(y.d)}
}

func (x Decimal) Div(y Decimal) Decimal {
	return Decimal{x.d.Div(y.d)}
}

// Round 四舍五入到小数位 places
func (x Decimal) Round(places int32) Decimal {
	return Decimal{x.d.Round(places)}
}

// Floor 向下取整
func (x Decimal) Floor() Decimal {
	return Decimal{x.d.Floor()}
}

// Ceil 向上取整
func (x Decimal) Ceil() Decimal {
	return Decimal{x.d.Ceil()}
}

// Abs 绝对值
func (x Decimal) Abs() Decimal {
	return Decimal{x.d.Abs()}
}

// --------------- 比较 ----------------

func (x Decimal) Equal(y Decimal) bool {
	return x.d.Equal(y.d)
}

func (x Decimal) GreaterThan(y Decimal) bool {
	return x.d.GreaterThan(y.d)
}

func (x Decimal) LessThan(y Decimal) bool {
	return x.d.LessThan(y.d)
}

func (x Decimal) Cmp(y Decimal) int {
	return x.d.Cmp(y.d)
}

func (x Decimal) IsZero() bool {
	return x.d.IsZero()
}

// --------------- 聚合辅助 ----------------

// Sum 对一组 Decimal 求和
func Sum(ds ...Decimal) Decimal {
	s := decimal.Zero
	for _, v := range ds {
		s = s.Add(v.d)
	}
	return Decimal{s}
}

// Avg 对一组 Decimal 求平均（长度为0返回0）
func Avg(ds ...Decimal) Decimal {
	if len(ds) == 0 {
		return Zero()
	}
	sum := Sum(ds...)
	return Decimal{sum.d.Div(decimal.NewFromInt(int64(len(ds))))}
}

// Min 返回最小值（长度为0会 panic）
func Min(ds ...Decimal) Decimal {
	if len(ds) == 0 {
		panic("Min requires at least one value")
	}
	min := ds[0].d
	for _, v := range ds[1:] {
		if v.d.LessThan(min) {
			min = v.d
		}
	}
	return Decimal{min}
}

// Max 返回最大值（长度为0会 panic）
func Max(ds ...Decimal) Decimal {
	if len(ds) == 0 {
		panic("Max requires at least one value")
	}
	max := ds[0].d
	for _, v := range ds[1:] {
		if v.d.GreaterThan(max) {
			max = v.d
		}
	}
	return Decimal{max}
}

// --------------- 方便的格式化辅助 ----------------

// FormatCurrency 按固定小数位格式化（例如 2 位），可选前缀
func FormatCurrency(x Decimal, places int32, prefix string) string {
	return fmt.Sprintf("%s%s", prefix, x.d.StringFixed(places))
}
