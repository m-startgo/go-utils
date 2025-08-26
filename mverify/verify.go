package mverify

import (
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// IsEmail 判断字符串是否为合法的邮箱格式（简单校验）。
// 输入: 任意字符串
// 输出: true 表示看起来像邮箱，false 表示不是
// 注意: 此函数做轻量级格式校验，不连接网络或验证域名存在性。
func IsEmail(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// 使用更严格的邮箱正则：只允许小写字母/数字开头，local 部分仅允许下划线、点、数字、小写字母和连字符，域名为小写，TLD 长度 2-4
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(s)
}

// IsMobile 判断字符串是否为中国大陆手机号（支持可选的 +86 前缀）。
// 规则：可带 "+86" 或 "86" 前缀，主号码为 1 开头且第二位为 3-9 的 11 位数字。
func IsMobile(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// 移除空格和短线
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	// 支持 +86 或 86 前缀
	if strings.HasPrefix(s, "+86") {
		s = s[3:]
	} else if strings.HasPrefix(s, "86") {
		s = s[2:]
	}
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(s)
}

// IsURL 简单判断是否为合法 URL（要求包含 scheme 和 host）。
func IsURL(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

// IsNumeric 判断字符串是否为整数或小数（可带负号）。
func IsNumeric(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	re := regexp.MustCompile(`^[+-]?(?:\d+)(?:\.\d+)?$`)
	return re.MatchString(s)
}

// IsIDCard 判断中国身份证号（15 位或 18 位）的基本合法性。
// 对 18 位身份证会校验最后一位校验码；15 位仅做纯数字长度校验。
func IsIDCard(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// 15 位全数字
	re15 := regexp.MustCompile(`^\d{15}$`)
	if re15.MatchString(s) {
		return true
	}
	// 18 位，前 17 位数字，最后一位数字或 X/x
	re18 := regexp.MustCompile(`^\d{17}[0-9Xx]$`)
	if !re18.MatchString(s) {
		return false
	}
	// 校验最后一位
	id17 := s[:17]
	expected := calcIDCardChecksum(id17)
	last := strings.ToUpper(s[17:])
	return last == expected
}

// calcIDCardChecksum 根据前 17 位数字计算第 18 位校验码（返回字符串，如 "X" 或 "0"）。
func calcIDCardChecksum(id17 string) string {
	// 加权因子
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	// 校验码映射
	mapping := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
	sum := 0
	for i, ch := range id17 {
		digit := int(ch - '0')
		sum += digit * weights[i]
	}
	mod := sum % 11
	return mapping[mod]
}

// IsIPv4 判断字符串是否为合法的 IPv4 地址（点分十进制，0.0.0.0 - 255.255.255.255）。
func IsIPv4(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}
	return ip.To4() != nil
}

// IsPort 判断字符串是否为合法端口（0-65535）。
func IsPort(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return n >= 0 && n <= 65535
}
