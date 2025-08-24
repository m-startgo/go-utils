package m_math

import "fmt"

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

// DivSafe 除法，除数为零时返回错误
func (d Decimal) DivSafe(other Decimal) (Decimal, error) {
	if other.value.IsZero() {
		return Zero, fmt.Errorf("divide by zero")
	}
	return Decimal{value: d.value.Div(other.value)}, nil
}

/* =============百分比计算================= */

// (a/b)*100，除数为零时返回 Zero
func Pct(a, b Decimal) Decimal {
	if b.IsZero() {
		return Zero
	}
	return a.Div(b).Mul(NewFromInt(100))
}

// (a/b)*100 并保留指定小数位数，除数为零时返回 Zero
func PctN(a, b Decimal, places int32) Decimal {
	return Pct(a, b).Round(places)
}

// ((a-b)/b)*100，除数为零时返回 Zero
func ChgPct(a, b Decimal) Decimal {
	if b.IsZero() {
		return Zero
	}
	return a.Sub(b).Div(b).Mul(NewFromInt(100))
}

// ((a-b)/b)*100 并保留指定小数位数，除数为零时返回 Zero
func ChgPctN(a, b Decimal, places int32) Decimal {
	return ChgPct(a, b).Round(places)
}

/* =============取值================= */

// Neg 返回取负值 正数变负数，负数变正数，0 的相反数仍是 0
func (d Decimal) Neg() Decimal {
	return Decimal{value: d.value.Neg()}
}

// Round 四舍五入到指定小数位数
func (d Decimal) Round(places int32) Decimal {
	return Decimal{value: d.value.Round(places)}
}

// Truncate 截断到指定小数位数（不四舍五入）
func (d Decimal) Truncate(places int32) Decimal {
	return Decimal{value: d.value.Truncate(places)}
}

// Abs 返回绝对值
func (d Decimal) Abs() Decimal {
	return Decimal{value: d.value.Abs()}
}
