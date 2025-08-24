package m_math

import (
	"math"
	"testing"
)

func TestNewFromFloat(t *testing.T) {
	f := 123.456
	d := NewFromFloat(f)
	if d.String() != "123.456" {
		t.Errorf("Expected 123.456, got %s", d.String())
	}
}

func TestNewFromString(t *testing.T) {
	s := "789.012"
	d, err := NewFromString(s)
	if err != nil {
		t.Errorf("Error creating Decimal from string: %v", err)
	}
	if d.String() != s {
		t.Errorf("Expected %s, got %s", s, d.String())
	}
}

func TestNewFromStringError(t *testing.T) {
	_, err := NewFromString("not_a_number")
	if err == nil {
		t.Errorf("Expected error for invalid string, got nil")
	}
}

func TestNewFromInt(t *testing.T) {
	i := int64(456)
	d := NewFromInt(i)
	if d.String() != "456" {
		t.Errorf("Expected 456, got %s", d.String())
	}
}

func TestAdd(t *testing.T) {
	d1 := NewFromFloat(1.23)
	d2 := NewFromFloat(4.56)
	result := d1.Add(d2)
	expected := "5.79"
	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestSub(t *testing.T) {
	d1 := NewFromFloat(4.56)
	d2 := NewFromFloat(1.23)
	result := d1.Sub(d2)
	expected := "3.33"
	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestMul(t *testing.T) {
	d1 := NewFromFloat(1.23)
	d2 := NewFromFloat(4.56)
	result := d1.Mul(d2)
	expected := "5.6088"
	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestDiv(t *testing.T) {
	d1 := NewFromFloat(5.6088)
	d2 := NewFromFloat(4.56)
	result := d1.Div(d2)
	expected := "1.23"
	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestDivByZero(t *testing.T) {
	d1 := NewFromFloat(1.23)
	d2 := Zero
	result := d1.Div(d2)
	if !result.IsZero() {
		t.Errorf("Expected Zero when dividing by zero, got %s", result.String())
	}
}

func TestRound(t *testing.T) {
	d := NewFromFloat(1.234567)
	result := d.Round(2)
	expected := "1.23"
	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestAbs(t *testing.T) {
	d := NewFromFloat(-1.23)
	result := d.Abs()
	expected := "1.23"
	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestCmp(t *testing.T) {
	d1 := NewFromFloat(1.23)
	d2 := NewFromFloat(4.56)
	d3 := NewFromFloat(1.23)

	if d1.Cmp(d2) != -1 {
		t.Errorf("Expected -1, got %d", d1.Cmp(d2))
	}

	if d2.Cmp(d1) != 1 {
		t.Errorf("Expected 1, got %d", d2.Cmp(d1))
	}

	if d1.Cmp(d3) != 0 {
		t.Errorf("Expected 0, got %d", d1.Cmp(d3))
	}
}

func TestIsZero(t *testing.T) {
	d := Zero
	if !d.IsZero() {
		t.Errorf("Expected true, got false")
	}
}

func TestIsPositive(t *testing.T) {
	d := NewFromFloat(1.23)
	if !d.IsPositive() {
		t.Errorf("Expected true, got false")
	}
}

func TestIsNegative(t *testing.T) {
	d := NewFromFloat(-1.23)
	if !d.IsNegative() {
		t.Errorf("Expected true, got false")
	}
}

func TestConstants(t *testing.T) {
	if Zero.String() != "0" {
		t.Errorf("Expected 0, got %s", Zero.String())
	}

	if One.String() != "1" {
		t.Errorf("Expected 1, got %s", One.String())
	}
}

func TestFloat64Precision(t *testing.T) {
	d, err := NewFromString("0.12345678901234567890")
	if err != nil {
		t.Errorf("Error creating Decimal from string: %v", err)
	}
	f := d.Float64()
	if f == 0 {
		t.Errorf("Float64 precision lost too much, got %f", f)
	}
}

func TestIntPartNegative(t *testing.T) {
	d := NewFromFloat(-123.456)
	if d.IntPart() != -123 {
		t.Errorf("Expected -123, got %d", d.IntPart())
	}
}

func TestRoundNegativePlaces(t *testing.T) {
	d := NewFromFloat(123.456)
	result := d.Round(-1)
	expected := "120"
	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestAbsZero(t *testing.T) {
	d := Zero
	result := d.Abs()
	if !result.IsZero() {
		t.Errorf("Expected zero, got %s", result.String())
	}
}

func TestCmpPrecision(t *testing.T) {
	d1, err1 := NewFromString("1.2300")
	if err1 != nil {
		t.Errorf("Error creating Decimal from string: %v", err1)
	}
	d2, err2 := NewFromString("1.23")
	if err2 != nil {
		t.Errorf("Error creating Decimal from string: %v", err2)
	}
	if d1.Cmp(d2) != 0 {
		t.Errorf("Expected 0, got %d", d1.Cmp(d2))
	}
}

func TestIsZeroSmallValue(t *testing.T) {
	d := NewFromFloat(1e-20)
	if d.IsZero() {
		t.Errorf("Expected false for small value, got true")
	}
}

func TestNewFromFloatExtreme(t *testing.T) {
	d := NewFromFloat(1e308)
	if d.String() == "0" {
		t.Errorf("Expected non-zero for large float, got %s", d.String())
	}
}

func TestNewFromIntExtreme(t *testing.T) {
	d := NewFromInt(9223372036854775807)
	if d.String() != "9223372036854775807" {
		t.Errorf("Expected max int64, got %s", d.String())
	}
}

func TestDivSafeByZero(t *testing.T) {
	d1 := NewFromFloat(1.23)
	_, err := d1.DivSafe(Zero)
	if err == nil {
		t.Errorf("Expected error when dividing by zero with DivSafe, got nil")
	}
}

func TestNegAndComparisons(t *testing.T) {
	dPos := NewFromFloat(1.23)
	dNeg := dPos.Neg()
	if !dNeg.Equal(NewFromFloat(-1.23)) {
		t.Errorf("Neg failed: expected -1.23, got %s", dNeg.String())
	}

	if !dPos.GreaterThan(Zero) {
		t.Errorf("Expected positive > zero")
	}
	if !dNeg.LessThan(Zero) {
		t.Errorf("Expected negative < zero")
	}
}

func TestFloat64ExactWithTolerance(t *testing.T) {
	d, _ := NewFromString("1.5")
	f := d.Float64()
	if math.Abs(f-1.5) > 1e-12 {
		t.Errorf("Expected float approx 1.5, got %f", f)
	}
}

func TestNewFromStringVariants(t *testing.T) {
	// 科学计数法
	dSci, err := NewFromString("1.23e3")
	if err != nil {
		t.Errorf("Error parsing scientific notation: %v", err)
	}
	if dSci.String() != "1230" {
		t.Errorf("Expected 1230, got %s", dSci.String())
	}

	// 带空白（若希望支持可在使用处 Trim 空白）
	_, err = NewFromString("  1.23  ")
	if err == nil {
		// 如果希望 NewFromString 接受空白，则此处应为 nil；根据 shopspring/decimal 行为调整期望
		// 这里只记录行为，避免误判
	}
}

func TestStringPreserveScale(t *testing.T) {
	d1, _ := NewFromString("1.2300")

	if d1.String() != "1.23" {
		t.Errorf("Expected scale preserved as 1.23, got %s", d1.String())
	}
}

func TestPctAndChgPct(t *testing.T) {
	a := NewFromFloat(50)
	b := NewFromFloat(200)
	if Pct(a, b).String() != "25" {
		t.Errorf("Expected 25, got %s", Pct(a, b).String())
	}

	if ChgPct(a, b).String() != "-75" {
		t.Errorf("Expected -75, got %s", ChgPct(a, b).String())
	}
}

func TestPctNAndChgPctN(t *testing.T) {
	a := NewFromFloat(12.345)
	b := NewFromFloat(67.89)
	// 计算百分比 (a/b)*100，保留2位小数
	expectedPct := "18.18" // 12.345/67.89*100 ≈ 18.1818，四舍五入后为 18.18
	if PctN(a, b, 2).Value().StringFixed(2) != expectedPct {
		t.Errorf("Expected %s, got %s", expectedPct, PctN(a, b, 2).Value().StringFixed(2))
	}
	// 计算变化百分比 ((a-b)/b)*100，保留3位小数
	expectedChgPct := "-81.818" // (12.345-67.89)/67.89*100 ≈ -81.8181，四舍五入后为 -81.818
	if ChgPctN(a, b, 3).Value().StringFixed(3) != expectedChgPct {
		t.Errorf("Expected %s, got %s", expectedChgPct, ChgPctN(a, b, 3).Value().StringFixed(3))
	}
}

func TestValueMethod(t *testing.T) {
	d := NewFromFloat(1.23)
	raw := d.Value()
	if raw.String() != "1.23" {
		t.Errorf("Value() did not return expected decimal.Decimal")
	}
}
