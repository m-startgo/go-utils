package m_time

import (
	"math"
	"strings"
	"testing"
	"time"
)

func TestParseAndFormat(t *testing.T) {
	// RFC3339 string
	s := "2021-01-02T03:04:05Z"
	tt, err := ParseString(s)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}
	if want := "2021-01-02 03:04:05.000000 +00:00"; tt.Format(DefaultToken) != want {
		t.Fatalf("Format mismatch: got %q want %q", tt.Format(DefaultToken), want)
	}

	// numeric string milliseconds
	ms := int64(1609459200000) // 2021-01-01T00:00:00Z
	p, err := ParseString("1609459200000")
	if err != nil {
		t.Fatalf("ParseString numeric failed: %v", err)
	}
	if p.ToTime().UTC().Unix()*1000 != ms/1 { // sanity check: Unix()*1000 == ms
		t.Fatalf("numeric ms parse mismatch: %v", p.ToTime())
	}

	// int64 milliseconds
	p2, err := ParseInt64(1609459200000)
	if err != nil {
		t.Fatalf("ParseInt64 failed: %v", err)
	}
	if p2.ToTime().UTC().Year() != 2021 {
		t.Fatalf("ParseInt64 year mismatch: %v", p2.ToTime())
	}

	// float seconds with fraction
	pf, err := ParseFloat64(1610000000.5)
	if err != nil {
		t.Fatalf("ParseFloat64 failed: %v", err)
	}
	if pf.ToTime().UTC().Nanosecond() != 500000000 {
		t.Fatalf("ParseFloat64 fractional second mismatch: %v", pf.ToTime())
	}

	// generic Parse wrapper
	if _, err := Parse("2021-01-02"); err != nil {
		t.Fatalf("Parse wrapper failed for string: %v", err)
	}
	if _, err := Parse(int64(1610000000)); err != nil {
		t.Fatalf("Parse wrapper failed for int64: %v", err)
	}
	if _, err := Parse(1610000000.5); err != nil {
		t.Fatalf("Parse wrapper failed for float64: %v", err)
	}
}

func TestArithmeticAndDayHelpers(t *testing.T) {
	base := FromTime(time.Date(2021, 3, 14, 15, 9, 26, 0, time.UTC))
	plus1 := base.AddDays(1)
	if !plus1.IsSameDay(FromTime(base.ToTime().Add(24 * time.Hour))) {
		t.Fatalf("AddDays incorrect: %v vs %v", plus1, base)
	}
	plusH := base.AddHours(2)
	if plusH.ToTime().Hour() != 17 {
		t.Fatalf("AddHours incorrect: hour=%d", plusH.ToTime().Hour())
	}

	sod := base.StartOfDay()
	if sod.ToTime().Hour() != 0 || sod.ToTime().Minute() != 0 {
		t.Fatalf("StartOfDay incorrect: %v", sod.ToTime())
	}
	eod := base.EndOfDay()
	if eod.ToTime().Hour() != 23 {
		t.Fatalf("EndOfDay hour incorrect: %v", eod.ToTime())
	}
}

func TestUnixTimestamps(t *testing.T) {
	tt := FromTime(time.Date(2022, 6, 1, 12, 0, 0, 123456000, time.UTC))
	expectedUnix := tt.ToTime().Unix()
	if tt.Unix() != expectedUnix {
		t.Fatalf("Unix() mismatch: got %d expected %d", tt.Unix(), expectedUnix)
	}
	expectedNsec := int64(tt.ToTime().Nanosecond())
	if tt.UnixNano()%1e9 != expectedNsec {
		t.Fatalf("UnixNano fractional mismatch: got %d expected %d", tt.UnixNano()%1e9, expectedNsec)
	}
}

func TestEdgeCases(t *testing.T) {
	// 极小/极大时间戳
	// 1) 32-bit 秒边界之前/之后
	small, err := ParseInt64(-2147483648) // 1970-01-... negative seconds
	if err != nil {
		t.Fatalf("ParseInt64 small failed: %v", err)
	}
	_ = small

	// 2) 纳秒级极大值（18+ 位）
	bigNsec := int64(1e18) // 1e18 纳秒
	big, err := ParseInt64(bigNsec)
	if err != nil {
		t.Fatalf("ParseInt64 big failed: %v", err)
	}
	expectedBig := time.Unix(0, bigNsec).UTC()
	if !big.ToTime().UTC().Equal(expectedBig) {
		t.Fatalf("big timestamp mismatch: got %v expected %v", big.ToTime().UTC(), expectedBig)
	}

	// 时区字符串解析（带时区偏移）
	tzStr := "2021-03-14T02:30:00-05:00" // EST offset -05:00
	tzt, err := ParseString(tzStr)
	if err != nil {
		t.Fatalf("ParseString tz failed: %v", err)
	}
	if tzt.ToTime().Location().String() == "" {
		// location may be empty name, just ensure time parsed and offset applied
		if tzt.ToTime().Hour() != 2 {
			t.Fatalf("tz parse hour mismatch: %v", tzt.ToTime())
		}
	}

	// 夏令时边界（以 America/New_York 为例）
	loc, _ := time.LoadLocation("America/New_York")
	// 2021-03-14 在美国开始 DST，02:00 跳到 03:00
	// 以本地时间构造 2021-03-14 02:30 并观察解析
	// 我们用 FromTime 构造并验证 StartOfDay/EndOfDay 在该时区的行为
	SetDefaultLocation(loc)
	localT := FromTime(time.Date(2021, 3, 14, 2, 30, 0, 0, loc))
	sod := localT.StartOfDay()
	if sod.ToTime().Hour() != 0 {
		t.Fatalf("StartOfDay local incorrect: %v", sod.ToTime())
	}
	// 恢复默认 (UTC)
	SetDefaultLocation(nil)
}

func TestMoreEdgeCases(t *testing.T) {
	// invalid inputs
	if _, err := Parse(nil); err == nil {
		t.Fatalf("expected error for nil Parse")
	}
	if _, err := Parse(struct{}{}); err == nil {
		t.Fatalf("expected error for unsupported type")
	}

	// uint64 overflow
	if _, err := Parse(uint64(math.MaxInt64) + 1); err == nil {
		t.Fatalf("expected error for uint64 overflow")
	}

	// format token edge cases
	tt := FromTime(time.Date(1999, 12, 31, 23, 59, 59, 999000000, time.UTC))
	if s := tt.Format("YYYY/MM/DD HH:mm:ss.SSS"); !strings.Contains(s, "1999/12/31") {
		t.Fatalf("format token failed: %s", s)
	}

	// concurrent SetDefaultLocation
	locNY, _ := time.LoadLocation("America/New_York")
	done := make(chan struct{})
	go func() {
		for i := 0; i < 1000; i++ {
			SetDefaultLocation(locNY)
			SetDefaultLocation(nil)
		}
		close(done)
	}()
	<-done
}
