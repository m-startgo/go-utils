package m_time

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	now := New()
	if now == nil {
		t.Error("New() returned nil")
		return
	}
	// 检查时间是否在当前时间附近（1秒内）
	if now.tm.After(time.Now().Add(1*time.Second)) || now.tm.Before(time.Now().Add(-1*time.Second)) {
		t.Error("New() returned time not close to current time")
	}
}

func TestNewFromTime(t *testing.T) {
	tm := time.Date(2023, 10, 15, 13, 45, 26, 0, time.UTC)
	t1 := NewFromTime(tm)
	if t1 == nil {
		t.Error("NewFromTime() returned nil")
	} else if !t1.tm.Equal(tm) {
		t.Error("NewFromTime() did not preserve the input time")
	}
}

func TestNewFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			input:    "2023-10-15",
			expected: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			input:    "2023-10-15 13:45:26",
			expected: time.Date(2023, 10, 15, 13, 45, 26, 0, time.UTC),
			wantErr:  false,
		},
		{
			input:    "invalid-date",
			expected: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t1, err := NewFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !t1.tm.Equal(tt.expected) {
				t.Errorf("NewFromString() = %v, expected %v", t1.tm, tt.expected)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	t1 := NewFromTime(time.Date(2023, 10, 15, 13, 45, 26, 0, time.UTC))
	formatted := t1.Format("2006-01-02 15:04:05")
	if formatted != "2023-10-15 13:45:26" {
		t.Errorf("Format() returned %s, expected 2023-10-15 13:45:26", formatted)
	}
}

func TestIn(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("Skipping test due to unable to load location")
	}
	t1 := NewFromTime(time.Date(2023, 10, 15, 13, 45, 26, 0, time.UTC))
	t2 := t1.In(loc)
	if t2 == nil {
		t.Error("In() returned nil")
	}
	// 纽约时间比UTC晚4小时（夏令时）或5小时（冬令时）
	// 这里简化处理，只检查是否转换了时区
	if t2 != nil && t2.tm.Location() != loc {
		t.Error("In() did not change the location")
	}
}

func TestAdd(t *testing.T) {
	t1 := NewFromTime(time.Date(2023, 10, 15, 13, 45, 26, 0, time.UTC))
	t2 := t1.Add(24 * time.Hour)
	expected := time.Date(2023, 10, 16, 13, 45, 26, 0, time.UTC)
	if !t2.tm.Equal(expected) {
		t.Errorf("Add() returned %v, expected %v", t2.tm, expected)
	}
}

func TestSubtract(t *testing.T) {
	t1 := NewFromTime(time.Date(2023, 10, 15, 13, 45, 26, 0, time.UTC))
	t2 := t1.Subtract(24 * time.Hour)
	expected := time.Date(2023, 10, 14, 13, 45, 26, 0, time.UTC)
	if !t2.tm.Equal(expected) {
		t.Errorf("Subtract() returned %v, expected %v", t2.tm, expected)
	}
}

func TestStartOf(t *testing.T) {
	t1 := NewFromTime(time.Date(2023, 10, 15, 13, 45, 26, 0, time.UTC))
	t2 := t1.StartOf("day")
	expected := time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC)
	if !t2.tm.Equal(expected) {
		t.Errorf("StartOf() returned %v, expected %v", t2.tm, expected)
	}
}

func TestEndOf(t *testing.T) {
	t1 := NewFromTime(time.Date(2023, 10, 15, 13, 45, 26, 0, time.UTC))
	t2 := t1.EndOf("day")
	expected := time.Date(2023, 10, 15, 23, 59, 59, 999999999, time.UTC)
	if !t2.tm.Equal(expected) {
		t.Errorf("EndOf() returned %v, expected %v", t2.tm, expected)
	}
}

func TestDaysInMonth(t *testing.T) {
	// 测试2月（非闰年）
	t1 := NewFromTime(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC))
	if days := t1.DaysInMonth(); days != 28 {
		t.Errorf("DaysInMonth() for Feb 2023 returned %d, expected 28", days)
	}
	// 测试2月（闰年）
	t2 := NewFromTime(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC))
	if days := t2.DaysInMonth(); days != 29 {
		t.Errorf("DaysInMonth() for Feb 2024 returned %d, expected 29", days)
	}
	// 测试大月
	t3 := NewFromTime(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	if days := t3.DaysInMonth(); days != 31 {
		t.Errorf("DaysInMonth() for Jan 2023 returned %d, expected 31", days)
	}
	// 测试小月
	t4 := NewFromTime(time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC))
	if days := t4.DaysInMonth(); days != 30 {
		t.Errorf("DaysInMonth() for Apr 2023 returned %d, expected 30", days)
	}
}

func TestDiff(t *testing.T) {
	t1 := NewFromTime(time.Date(2023, 10, 15, 13, 45, 26, 0, time.UTC))
	t2 := NewFromTime(time.Date(2023, 9, 15, 13, 45, 26, 0, time.UTC))

	// 测试天数差
	diffDays := t1.Diff(t2, "day")
	expectedDays := 30
	if diffDays.IntPart() != int64(expectedDays) {
		t.Errorf("Diff() days returned %v, expected %d", diffDays, expectedDays)
	}

	// 测试月数差
	diffMonths := t1.Diff(t2, "month")
	expectedMonths := 1
	if diffMonths.IntPart() != int64(expectedMonths) {
		t.Errorf("Diff() months returned %v, expected %d", diffMonths, expectedMonths)
	}
}

func TestUnix(t *testing.T) {
	tm := time.Date(2023, 10, 15, 13, 45, 26, 0, time.UTC)
	t1 := NewFromTime(tm)
	if unix := t1.Unix(); unix != tm.Unix() {
		t.Errorf("Unix() returned %d, expected %d", unix, tm.Unix())
	}
}

func TestUnixMilli(t *testing.T) {
	tm := time.Date(2023, 10, 15, 13, 45, 26, 123000000, time.UTC)
	t1 := NewFromTime(tm)
	if unixMilli := t1.UnixMilli(); unixMilli != tm.UnixMilli() {
		t.Errorf("UnixMilli() returned %d, expected %d", unixMilli, tm.UnixMilli())
	}
}

func TestUnixMicro(t *testing.T) {
	tm := time.Date(2023, 10, 15, 13, 45, 26, 123456000, time.UTC)
	t1 := NewFromTime(tm)
	if unixMicro := t1.UnixMicro(); unixMicro != tm.UnixMicro() {
		t.Errorf("UnixMicro() returned %d, expected %d", unixMicro, tm.UnixMicro())
	}
}

func TestUnixNano(t *testing.T) {
	tm := time.Date(2023, 10, 15, 13, 45, 26, 123456789, time.UTC)
	t1 := NewFromTime(tm)
	if unixNano := t1.UnixNano(); unixNano != tm.UnixNano() {
		t.Errorf("UnixNano() returned %d, expected %d", unixNano, tm.UnixNano())
	}
}

func TestTime(t *testing.T) {
	tm := time.Date(2023, 10, 15, 13, 45, 26, 0, time.UTC)
	t1 := NewFromTime(tm)
	if t1.Time() != tm {
		t.Error("Time() did not return the correct time")
	}
}

func TestString(t *testing.T) {
	tm := time.Date(2023, 10, 15, 13, 45, 26, 0, time.UTC)
	t1 := NewFromTime(tm)
	expected := tm.String()
	if str := t1.String(); str != expected {
		t.Errorf("String() returned %s, expected %s", str, expected)
	}
}

func TestNormalizeUnit(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"y", "year"},
		{"yr", "year"},
		{"years", "year"},
		{"M", "month"},
		{"mon", "month"},
		{"months", "month"},
		{"d", "day"},
		{"days", "day"},
		{"h", "hour"},
		{"hours", "hour"},
		{"m", "minute"},
		{"min", "minute"},
		{"minutes", "minute"},
		{"s", "second"},
		{"sec", "second"},
		{"seconds", "second"},
		{"ms", "millisecond"},
		{"milliseconds", "millisecond"},
		{"us", "microsecond"},
		{"µs", "microsecond"},
		{"microseconds", "microsecond"},
		{"ns", "nanosecond"},
		{"nanoseconds", "nanosecond"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := normalizeUnit(tt.input); got != tt.expected {
				t.Errorf("normalizeUnit(%q) = %q, expected %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMaxInt(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{3, 5, 5},
		{-1, -5, -1},
		{0, 0, 0},
		{100, 50, 100},
	}

	for _, tt := range tests {
		t.Run(testName(tt.a, tt.b), func(t *testing.T) {
			if got := maxInt(tt.a, tt.b); got != tt.expected {
				t.Errorf("maxInt(%d, %d) = %d, expected %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// 为测试名称格式化添加辅助函数
func testName(a, b int) string {
	return fmt.Sprintf("%d_%d", a, b)
}
