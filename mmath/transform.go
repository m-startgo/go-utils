package m_math

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
