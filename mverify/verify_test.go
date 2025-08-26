package mverify

import "testing"

func TestIsEmail(t *testing.T) {
	cases := map[string]bool{
		"test@example.com":      true,
		"a.b_tag@sub.domain.cn": true,
		"invalid@":              false,
		"":                      false,
	}
	for s, want := range cases {
		if got := IsEmail(s); got != want {
			t.Fatalf("IsEmail(%q) = %v, want %v", s, got, want)
		}
	}
}

func TestIsMobile(t *testing.T) {
	cases := map[string]bool{
		"13800138000":    true,
		"+8613800138000": true,
		"8613800138000":  true,
		"123":            false,
		"":               false,
	}
	for s, want := range cases {
		if got := IsMobile(s); got != want {
			t.Fatalf("IsMobile(%q) = %v, want %v", s, got, want)
		}
	}
}

func TestIsURL(t *testing.T) {
	cases := map[string]bool{
		"https://example.com":        true,
		"http://localhost:8080/path": true,
		"//example.com":              false,
		"not a url":                  false,
	}
	for s, want := range cases {
		if got := IsURL(s); got != want {
			t.Fatalf("IsURL(%q) = %v, want %v", s, got, want)
		}
	}
}

func TestIsNumeric(t *testing.T) {
	cases := map[string]bool{
		"123":   true,
		"-1.23": true,
		"+0.5":  true,
		"1.":    false,
		"":      false,
	}
	for s, want := range cases {
		if got := IsNumeric(s); got != want {
			t.Fatalf("IsNumeric(%q) = %v, want %v", s, got, want)
		}
	}
}

func TestIsIDCard(t *testing.T) {
	// 18 位示例（随机示例，校验位需正确）
	// 例子: 11010519491231002X 是常见测试样例
	cases := map[string]bool{
		"11010519491231002X": true,
		"110105194912310021": false,
		"130503670401001":    true, // 15 位
		"":                   false,
	}
	for s, want := range cases {
		if got := IsIDCard(s); got != want {
			t.Fatalf("IsIDCard(%q) = %v, want %v", s, got, want)
		}
	}
}

func TestIsIPv4(t *testing.T) {
	cases := map[string]bool{
		"127.0.0.1":       true,
		"192.168.1.1":     true,
		"255.255.255.255": true,
		"256.0.0.1":       false,
		"::1":             false,
		"":                false,
	}
	for s, want := range cases {
		if got := IsIPv4(s); got != want {
			t.Fatalf("IsIPv4(%q) = %v, want %v", s, got, want)
		}
	}
}

func TestIsPort(t *testing.T) {
	cases := map[string]bool{
		"0":     true,
		"22":    true,
		"8080":  true,
		"65535": true,
		"65536": false,
		"-1":    false,
		"abc":   false,
		"":      false,
	}
	for s, want := range cases {
		if got := IsPort(s); got != want {
			t.Fatalf("IsPort(%q) = %v, want %v", s, got, want)
		}
	}
}
