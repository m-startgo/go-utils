package mmath

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

// Precision 返回小数点后的有效位数（小数位数）
// 例如 12.345 -> 3， 100 -> 0
func (d Decimal) Precision() int {
	exp := d.value.Exponent()
	if exp < 0 {
		return -int(exp)
	}
	return 0
}
