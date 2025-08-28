package murl

/*

我想封装一个 url 库
它可以吧任意字符串 解析成 net/url 格式

然后可以方便的提取 里面的 每一个部分

*/

import (
	"fmt"
	"net/url"
	"strings"
)

// URL 是对标准库 net/url.URL 的轻量封装，保留原始输入并提供便捷访问方法。
//
// 用法示例：
//
//  u, err := murl.Parse("example.com/path?a=1#f")
//  if err != nil {
//      // 处理错误
//  }
//  host := u.Hostname()
//  port := u.Port()

type URL struct {
	raw string
	u   *url.URL
}

// Parse 将任意字符串解析为 *murl.URL。
//
// 如果输入缺少 scheme（例如 "example.com/path"），函数会尝试在前面加上 "http://" 再解析。
// 返回值：解析成功返回封装后的 *URL；解析失败返回非 nil 的错误。
// 可能的错误包括 net/url 提供的解析错误。
func Parse(raw string) (res *URL, resErr error) {
	res = &URL{}
	resErr = nil

	if raw == "" {
		resErr = fmt.Errorf("err:murl.Parse|empty input|input is empty")
		return
	}

	// 首先尝试原始解析
	u, err := url.Parse(raw)
	if err != nil {
		// 如果原始解析失败且看起来像是缺少 scheme，则尝试补全 http://
		if !strings.Contains(raw, "://") {
			try := "http://" + raw
			u2, err2 := url.Parse(try)
			if err2 != nil {
				resErr = fmt.Errorf("err:murl.Parse|parse failed|%v", err2)
				return
			}
			u = u2
			raw = try
		} else {
			resErr = fmt.Errorf("err:murl.Parse|parse failed|%v", err)
			return
		}
	} else {
		// 如果解析成功但没有 scheme，补 http
		if u.Scheme == "" && !strings.Contains(raw, "://") {
			try := "http://" + raw
			u2, err2 := url.Parse(try)
			if err2 == nil {
				u = u2
				raw = try
			}
		}
	}

	res.raw = raw
	res.u = u
	return
}

// MustParse 与 Parse 类似，但在解析失败时直接 panic，适合在 init 或测试中快速失败。
func MustParse(raw string) *URL {
	u, err := Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

// Raw 返回原始（可能被补全 scheme 后的）输入字符串。
func (x *URL) Raw() string {
	if x == nil {
		return ""
	}
	return x.raw
}

// URL 返回内部的 *url.URL 指针，直接访问 net/url 的所有方法。
func (x *URL) URL() *url.URL {
	if x == nil {
		return &url.URL{}
	}
	return x.u
}

// Scheme 返回 URL 的 scheme，例如 "http" 或 "https"。
func (x *URL) Scheme() string {
	if x == nil || x.u == nil {
		return ""
	}
	return x.u.Scheme
}

// Host 返回包含端口的主机部分，例如 "example.com:8080"。
func (x *URL) Host() string {
	if x == nil || x.u == nil {
		return ""
	}
	return x.u.Host
}

// Hostname 返回不带端口的主机名，例如 "example.com"。
func (x *URL) Hostname() string {
	if x == nil || x.u == nil {
		return ""
	}
	return x.u.Hostname()
}

// Port 返回端口号（如果有），否则返回空字符串。
func (x *URL) Port() string {
	if x == nil || x.u == nil {
		return ""
	}
	return x.u.Port()
}

// Path 返回 URL 的 path 部分（不包含 query/hash）。
func (x *URL) Path() string {
	if x == nil || x.u == nil {
		return ""
	}
	return x.u.Path
}

// Query 返回解析后的 query 参数集合（url.Values）。
func (x *URL) Query() url.Values {
	if x == nil || x.u == nil {
		return url.Values{}
	}
	return x.u.Query()
}

// QueryValue 返回指定 query key 的第一个值，若不存在则返回空字符串。
func (x *URL) QueryValue(key string) string {
	if x == nil || x.u == nil {
		return ""
	}
	return x.u.Query().Get(key)
}

// Fragment 返回 URL 的 fragment（hash）部分，不包含 '#'。
func (x *URL) Fragment() string {
	if x == nil || x.u == nil {
		return ""
	}
	return x.u.Fragment
}

// User 返回用户信息（username[:password]）的可读表示，若无用户信息返回空字符串。
func (x *URL) User() string {
	if x == nil || x.u == nil || x.u.User == nil {
		return ""
	}
	username := x.u.User.Username()
	if p, ok := x.u.User.Password(); ok && p != "" {
		return username + ":" + p
	}
	return username
}

// String 返回标准化的 URL 字符串表示。
func (x *URL) String() string {
	if x == nil || x.u == nil {
		return ""
	}
	return x.u.String()
}
