package m_math

import (
	"fmt"
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
func TestEdgeCases(t *testing.T) {
	// 测试非常大的数
	large := NewFromFloat(1e20)
	small := NewFromFloat(1e-20)
	result := large.Mul(small)
	expected := NewFromFloat(1)
	if !result.Equal(expected) {
		t.Errorf("Edge case multiplication failed: expected %v, got %v", expected, result)
	}

	// 测试零值运算
	zero := Zero
	result = zero.Add(zero)
	if !result.Equal(zero) {
		t.Errorf("Zero addition failed: expected %v, got %v", zero, result)
	}

	// 测试负数运算
	negA := NewFromFloat(-1.23)
	negB := NewFromFloat(-4.56)
	result = negA.Add(negB)
	expected = NewFromFloat(-5.79)
	if !result.Equal(expected) {
		t.Errorf("Negative addition failed: expected %v, got %v", expected, result)
	}

	// 测试精度边界
	pi := NewFromFloat(math.Pi)
	result = pi.Round(10)
	expectedStr := fmt.Sprintf("%.10f", math.Pi)
	expected, _ = NewFromString(expectedStr)
	if !result.Equal(expected) {
		t.Errorf("Precision boundary test failed: expected %v, got %v", expected, result)
	}
}
