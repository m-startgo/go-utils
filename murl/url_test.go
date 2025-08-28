package murl

import (
	"strings"
	"testing"
)

// 对比辅助函数
func eq(t *testing.T, name string, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got=%q want=%q", name, got, want)
	}
}

func TestParse_AddsHTTPWhenNoScheme(t *testing.T) {
	raw := "example.com/path?a=1#f"
	u, err := Parse(raw)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	eq(t, "Scheme", u.Scheme(), "http")
	eq(t, "Hostname", u.Hostname(), "example.com")
	eq(t, "Host", u.Host(), "example.com")
	eq(t, "Path", u.Path(), "/path")
	eq(t, "QueryValue a", u.QueryValue("a"), "1")
	eq(t, "Fragment", u.Fragment(), "f")

	// Parse 会把原始输入补全为 http://...
	if !strings.HasPrefix(u.Raw(), "http://") {
		t.Fatalf("Raw not prefixed with http://: %s", u.Raw())
	}
}

func TestParse_EmptyInput_Error(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Fatalf("expected error for empty input")
	}
	if !strings.Contains(err.Error(), "err:murl.Parse|empty input") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMustParse_PanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("MustParse did not panic on bad input")
		}
	}()
	// 传入空字符串，MustParse 应该 panic
	_ = MustParse("")
}

func TestUserInfo_Port_Hostname(t *testing.T) {
	u := MustParse("http://user:pass@example.com:8080/path")
	eq(t, "User", u.User(), "user:pass")
	eq(t, "Port", u.Port(), "8080")
	eq(t, "Host", u.Host(), "example.com:8080")
	eq(t, "Hostname", u.Hostname(), "example.com")
}

func TestUserInfo_UsernameOnly(t *testing.T) {
	u := MustParse("http://alice@example.com")
	eq(t, "User", u.User(), "alice")
}

func TestNilReceiver_IsSafe(t *testing.T) {
	var u *URL

	eq(t, "Raw", u.Raw(), "")
	eq(t, "Scheme", u.Scheme(), "")
	eq(t, "Host", u.Host(), "")
	eq(t, "Hostname", u.Hostname(), "")
	eq(t, "Port", u.Port(), "")
	eq(t, "Path", u.Path(), "")
	if q := u.Query(); len(q) != 0 {
		t.Fatalf("expected empty query for nil receiver, got %v", q)
	}
	eq(t, "QueryValue", u.QueryValue("x"), "")
	eq(t, "Fragment", u.Fragment(), "")
	eq(t, "User", u.User(), "")
	eq(t, "String", u.String(), "")
}
