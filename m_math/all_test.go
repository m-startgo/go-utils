package m_math

import (
	"math"
	"testing"
)

// 测试除零情况
var zero = Zero

// 测试 Decimal 的四则运算
func TestDecimalArithmetic(t *testing.T) {
	// 测试加法
	a := NewFromFloat(1.23)
	b := NewFromFloat(4.56)
	expected := NewFromFloat(5.79)
	result := a.Add(b)
	if !result.Equal(expected) {
		t.Errorf("Add failed: expected %v, got %v", expected, result)
	}

	// 测试减法
	expected = NewFromFloat(-3.33)
	result = a.Sub(b)
	if !result.Equal(expected) {
		t.Errorf("Sub failed: expected %v, got %v", expected, result)
	}

	// 测试乘法
	expected = NewFromFloat(5.6088)
	result = a.Mul(b)
	if !result.Equal(expected) {
		t.Errorf("Mul failed: expected %v, got %v", expected, result)
	}

	// 测试除法
	expected = NewFromFloat(0.2697368421052632)
	result = a.Div(b)
	if !result.Equal(expected) {
		t.Errorf("Div failed: expected %v, got %v", expected, result)
	}

	result = a.Div(zero)
	if !result.Equal(zero) {
		t.Errorf("Div by zero failed: expected %v, got %v", zero, result)
	}

	// 测试安全除法
	_, err := a.DivSafe(zero)
	if err == nil {
		t.Error("DivSafe should return error when dividing by zero")
	}
}

// 测试百分比计算
func TestPercentageCalculations(t *testing.T) {
	// 测试 Pct
	a := NewFromFloat(50)
	b := NewFromFloat(100)
	expected := NewFromFloat(50)
	result := Pct(a, b)
	if !result.Equal(expected) {
		t.Errorf("Pct failed: expected %v, got %v", expected, result)
	}

	// 测试 PctN
	expected = NewFromFloat(50).Round(2)
	result = PctN(a, b, 2)
	if !result.Equal(expected) {
		t.Errorf("PctN failed: expected %v, got %v", expected, result)
	}

	// 测试 ChgPct
	expected = NewFromFloat(25)
	result = ChgPct(NewFromFloat(125), b)
	if !result.Equal(expected) {
		t.Errorf("ChgPct failed: expected %v, got %v", expected, result)
	}

	// 测试 ChgPctN
	expected = NewFromFloat(25).Round(2)
	result = ChgPctN(NewFromFloat(125), b, 2)
	if !result.Equal(expected) {
		t.Errorf("ChgPctN failed: expected %v, got %v", expected, result)
	}

	// 测试除零情况
	result = Pct(a, zero)
	if !result.Equal(zero) {
		t.Errorf("Pct with zero divisor failed: expected %v, got %v", zero, result)
	}

	result = ChgPct(a, zero)
	if !result.Equal(zero) {
		t.Errorf("ChgPct with zero divisor failed: expected %v, got %v", zero, result)
	}
}

// 测试取值操作
func TestDecimalValueOperations(t *testing.T) {
	// 测试 Neg
	a := NewFromFloat(1.23)
	expected := NewFromFloat(-1.23)
	result := a.Neg()
	if !result.Equal(expected) {
		t.Errorf("Neg failed: expected %v, got %v", expected, result)
	}

	// 测试 Round
	expected = NewFromFloat(1.2)
	result = a.Round(1)
	if !result.Equal(expected) {
		t.Errorf("Round failed: expected %v, got %v", expected, result)
	}

	// 测试 Truncate
	b := NewFromFloat(1.29)
	expected = NewFromFloat(1.2)
	result = b.Truncate(1)
	if !result.Equal(expected) {
		t.Errorf("Truncate failed: expected %v, got %v", expected, result)
	}

	// 测试 Abs
	a = NewFromFloat(-1.23)
	expected = NewFromFloat(1.23)
	result = a.Abs()
	if !result.Equal(expected) {
		t.Errorf("Abs failed: expected %v, got %v", expected, result)
	}
}

// 测试比较操作
func TestDecimalComparison(t *testing.T) {
	a := NewFromFloat(1.23)
	b := NewFromFloat(4.56)
	zero := Zero

	// 测试 Cmp
	if a.Cmp(b) != -1 {
		t.Error("Cmp failed: expected -1")
	}
	if a.Cmp(a) != 0 {
		t.Error("Cmp failed: expected 0")
	}
	if b.Cmp(a) != 1 {
		t.Error("Cmp failed: expected 1")
	}

	// 测试 Equal
	if !a.Equal(a) {
		t.Error("Equal failed: expected true")
	}
	if a.Equal(b) {
		t.Error("Equal failed: expected false")
	}

	// 测试 GreaterThan
	if !b.GreaterThan(a) {
		t.Error("GreaterThan failed: expected true")
	}
	if a.GreaterThan(b) {
		t.Error("GreaterThan failed: expected false")
	}

	// 测试 LessThan
	if !a.LessThan(b) {
		t.Error("LessThan failed: expected true")
	}
	if b.LessThan(a) {
		t.Error("LessThan failed: expected false")
	}

	// 测试 IsZero
	if !zero.IsZero() {
		t.Error("IsZero failed: expected true")
	}
	if a.IsZero() {
		t.Error("IsZero failed: expected false")
	}

	// 测试 IsPositive
	if !a.IsPositive() {
		t.Error("IsPositive failed: expected true")
	}
	if zero.IsPositive() {
		t.Error("IsPositive failed: expected false")
	}

	// 测试 IsNegative
	c := NewFromFloat(-1.23)
	if !c.IsNegative() {
		t.Error("IsNegative failed: expected true")
	}
	if zero.IsNegative() {
		t.Error("IsNegative failed: expected false")
	}

	// 测试 Precision
	d := NewFromFloat(12.345)
	if d.Precision() != 3 {
		t.Error("Precision failed: expected 3")
	}
	e := NewFromInt(100)
	if e.Precision() != 0 {
		t.Error("Precision failed: expected 0")
	}
}

// 测试 Decimal 创建函数
func TestDecimalCreation(t *testing.T) {
	// 测试 NewFromFloat
	a := NewFromFloat(1.23)
	if a.String() != "1.23" {
		t.Errorf("NewFromFloat failed: expected 1.23, got %s", a.String())
	}

	// 测试 NewFromInt
	b := NewFromInt(123)
	if b.String() != "123" {
		t.Errorf("NewFromInt failed: expected 123, got %s", b.String())
	}

	// 测试 NewFromString
	c, err := NewFromString("1.23")
	if err != nil {
		t.Errorf("NewFromString failed: %v", err)
	}
	if c.String() != "1.23" {
		t.Errorf("NewFromString failed: expected 1.23, got %s", c.String())
	}

	// 测试 NewFromString 错误情况
	_, err = NewFromString("invalid")
	if err == nil {
		t.Error("NewFromString should return error for invalid string")
	}
}

// 测试 Sum 和 Mean 函数
func TestSumAndMean(t *testing.T) {
	// 测试 Sum
	a := NewFromFloat(1.23)
	b := NewFromFloat(4.56)
	c := NewFromFloat(7.89)
	expected := NewFromFloat(13.68)
	result := Sum(a, b, c)
	if !result.Equal(expected) {
		t.Errorf("Sum failed: expected %v, got %v", expected, result)
	}

	// 测试空参数的 Sum
	result = Sum()
	if !result.Equal(Zero) {
		t.Errorf("Sum with no args failed: expected %v, got %v", Zero, result)
	}

	// 测试 Mean
	expected = NewFromFloat(4.56)
	result = Mean(a, b, c)
	if !result.Equal(expected) {
		t.Errorf("Mean failed: expected %v, got %v", expected, result)
	}

	// 测试空参数的 Mean
	result = Mean()
	if !result.Equal(Zero) {
		t.Errorf("Mean with no args failed: expected %v, got %v", Zero, result)
	}
}

// 测试随机数生成函数
func TestRandomFunctions(t *testing.T) {
	// 测试 RandIntN
	for i := 0; i < 10; i++ {
		n := RandIntN(10)
		if n < 0 || n >= 10 {
			t.Errorf("RandIntN failed: got %d, expected in range [0, 10)", n)
		}
	}

	// 测试 RandIntN 边界情况
	n := RandIntN(0)
	if n != 0 {
		t.Errorf("RandIntN with 0 failed: expected 0, got %d", n)
	}

	// 测试 RandInt
	for i := 0; i < 10; i++ {
		n := RandInt(5, 10)
		if n < 5 || n > 10 {
			t.Errorf("RandInt failed: got %d, expected in range [5, 10]", n)
		}
	}

	// 测试 RandInt 参数交换
	n = RandInt(10, 5)
	if n < 5 || n > 10 {
		t.Errorf("RandInt with swapped params failed: got %d, expected in range [5, 10]", n)
	}

	// 测试 RandInt 相等参数
	n = RandInt(5, 5)
	if n != 5 {
		t.Errorf("RandInt with equal params failed: expected 5, got %d", n)
	}
}

// 测试转换函数
func TestTransformFunctions(t *testing.T) {
	// 测试 String
	a := NewFromFloat(1.23)
	if a.String() != "1.23" {
		t.Errorf("String failed: expected 1.23, got %s", a.String())
	}

	// 测试 Float64
	f := a.Float64()
	if math.Abs(f-1.23) > 1e-10 {
		t.Errorf("Float64 failed: expected 1.23, got %f", f)
	}

	// 测试 IntPart
	b := NewFromFloat(123.456)
	if b.IntPart() != 123 {
		t.Errorf("IntPart failed: expected 123, got %d", b.IntPart())
	}
}

// 测试边界条件和特殊情况

func TestDecimalFractionOperations(t *testing.T) {
	// 0.1 + 0.2 = 0.3 (use strings to avoid float binary inaccuracies)
	a, err := NewFromString("0.1")
	if err != nil {
		t.Fatalf("NewFromString failed: %v", err)
	}
	b, err := NewFromString("0.2")
	if err != nil {
		t.Fatalf("NewFromString failed: %v", err)
	}
	expected, err := NewFromString("0.3")
	if err != nil {
		t.Fatalf("NewFromString failed: %v", err)
	}
	res := a.Add(b)
	if !res.Equal(expected) {
		t.Errorf("0.1 + 0.2 failed: expected %v, got %v", expected, res)
	}

	// 0.1 * 0.2 = 0.02
	expected, _ = NewFromString("0.02")
	res = a.Mul(b)
	if !res.Equal(expected) {
		t.Errorf("0.1 * 0.2 failed: expected %v, got %v", expected, res)
	}
}

func TestDivisionPrecisionAndRounding(t *testing.T) {
	one, _ := NewFromString("1")
	three, _ := NewFromString("3")

	// 1 / 3 rounded to 5 decimals => 0.33333
	div := one.Div(three).Round(5)
	expected, _ := NewFromString("0.33333")
	if !div.Equal(expected) {
		t.Errorf("1/3 Round(5) failed: expected %v, got %v", expected, div)
	}

	// DivSafe should return value and nil error for non-zero divisor
	v, err := one.DivSafe(three)
	if err != nil {
		t.Errorf("DivSafe unexpected error: %v", err)
	}
	if v.Cmp(one.Div(three)) != 0 {
		t.Errorf("DivSafe returned different result: expected %v, got %v", one.Div(three), v)
	}

	// DivSafe by zero should error
	_, err = one.DivSafe(Zero)
	if err == nil {
		t.Errorf("DivSafe should error when dividing by zero")
	}
}

func TestPctAndChgPctWithDecimals(t *testing.T) {
	// Pct(2.5, 5) => (2.5/5)*100 = 50
	a, _ := NewFromString("2.5")
	b, _ := NewFromString("5")
	expected, _ := NewFromString("50")
	res := Pct(a, b)
	if !res.Equal(expected) {
		t.Errorf("Pct with decimals failed: expected %v, got %v", expected, res)
	}

	// PctN with rounding: Pct(1,3) ~ 33.33333 -> round 4 => 33.3333
	num, _ := NewFromString("1")
	den, _ := NewFromString("3")
	res = PctN(num, den, 4)
	expected, _ = NewFromString("33.3333")
	if !res.Equal(expected) {
		t.Errorf("PctN rounding failed: expected %v, got %v", expected, res)
	}

	// ChgPct decimals: from 120 to 100 => 20%
	from, _ := NewFromString("120")
	to, _ := NewFromString("100")
	res = ChgPct(from, to)
	expected, _ = NewFromString("20")
	if !res.Equal(expected) {
		t.Errorf("ChgPct with decimals failed: expected %v, got %v", expected, res)
	}

	// ChgPctN rounding
	res = ChgPctN(from, to, 1)
	expected, _ = NewFromString("20.0")
	if !res.Equal(expected) {
		t.Errorf("ChgPctN rounding failed: expected %v, got %v", expected, res)
	}
}

func TestArithmeticWithRepeatingDecimals(t *testing.T) {
	// 0.33333 * 3 = 0.99999 -> Round(5) should be 0.99999; adding small adjusts to 1.0 after rounding if needed
	v, _ := NewFromString("0.33333")
	three, _ := NewFromString("3")
	res := v.Mul(three).Round(5)
	expected, _ := NewFromString("0.99999")
	if !res.Equal(expected) {
		t.Errorf("Repeating decimal multiplication failed: expected %v, got %v", expected, res)
	}

	// If we use higher precision 1/3 * 3 -> 1
	one, _ := NewFromString("1")
	third := one.Div(three).Mul(three) // 1/3
	third = third.Round(2)

	if !third.Equal(one) {
		t.Errorf("1/3 * 3 failed: expected %v, got %v", one, res)
	}
}

func TestSmallDecimalEdgeCases(t *testing.T) {
	// Very small numbers multiplication
	a, _ := NewFromString("0.0003")
	b, _ := NewFromString("0.0002")
	expected, _ := NewFromString("0.00000006")
	res := a.Mul(b)
	if !res.Equal(expected) {
		t.Errorf("Small decimal multiplication failed: expected %v, got %v", expected, res)
	}

	// Addition of complementary small decimals
	x, _ := NewFromString("0.00000005")
	y, _ := NewFromString("0.00000001")
	expected, _ = NewFromString("0.00000006")
	if !x.Add(y).Equal(expected) {
		t.Errorf("Small decimal addition failed: expected %v, got %v", expected, x.Add(y))
	}
}

func TestAggregateDecimalOperations(t *testing.T) {
	a, _ := NewFromString("1.11")
	b, _ := NewFromString("2.22")
	c, _ := NewFromString("3.33")

	// Sum
	expected, _ := NewFromString("6.66")
	if !Sum(a, b, c).Equal(expected) {
		t.Errorf("Sum with decimals failed: expected %v, got %v", expected, Sum(a, b, c))
	}

	// Mean
	expected, _ = NewFromString("2.22")
	if !Mean(a, b, c).Equal(expected) {
		t.Errorf("Mean with decimals failed: expected %v, got %v", expected, Mean(a, b, c))
	}

	// Sum with mixed precision
	d, _ := NewFromString("1.2345")
	e, _ := NewFromString("2.0")
	expected, _ = NewFromString("3.2345")
	if !Sum(d, e).Equal(expected) {
		t.Errorf("Sum mixed precision failed: expected %v, got %v", expected, Sum(d, e))
	}

	// Precision detection for string-created decimals with trailing zeros
	p, _ := NewFromString("1.2300")
	if p.Precision() != 4 {
		t.Errorf("Precision detection failed: expected 4, got %d", p.Precision())
	}
}
