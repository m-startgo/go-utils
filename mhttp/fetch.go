package mhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// FetchOptions 请求选项
type FetchOptions struct {
	URL        string
	Data       []byte
	DataMap    map[string]any
	Params     map[string]string
	Headers    map[string]string
	Timeout    int    // seconds
	Retry      int    // 重试次数
	RetryDelay int    // 重试次数延迟 seconds
	Method     string // 允许值：GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS（不区分大小写，会在 Do 中规范化为大写）
	// MaxBodySize 限制读取响应体的最大字节数，0 表示不限制
	MaxBodySize int64
}

// Fetch 请求封装
type Fetch struct {
	opts FetchOptions
}

// package-level transport 用于复用连接池
var defaultTransport = &http.Transport{}

// NewFetch 创建一个 Fetch 实例
func NewFetch(opts FetchOptions) *Fetch {
	return &Fetch{opts: opts}
}

// Get 发起 GET 请求，并返回响应 body
// Do 发起请求，使用 FetchOptions.Method，要求 Method 非空并且为标准 HTTP 方法
// 调用示例：
//
//	res, err := NewFetch(FetchOptions{URL: "https://...", Method: http.MethodPost, DataMap: m}).Do()
func (f *Fetch) Do() ([]byte, error) {
	opts := f.opts
	if opts.Method == "" {
		return nil, errors.New("empty Method")
	}
	m := strings.ToUpper(opts.Method)
	switch m {
	case http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions:
		// ok
	default:
		return nil, fmt.Errorf("invalid method: %s", opts.Method)
	}
	opts.Method = m
	return f.do(opts)
}

// do 执行请求，并支持重试
func (f *Fetch) do(opts FetchOptions) ([]byte, error) {
	// 保护性检查
	if opts.URL == "" {
		return nil, errors.New("empty URL")
	}

	// 构造 URL 和 params
	u, err := url.Parse(opts.URL)
	if err != nil {
		return nil, err
	}
	if len(opts.Params) > 0 {
		q := u.Query()
		for k, v := range opts.Params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	// body 构造
	var rawBody []byte
	if opts.Method != http.MethodGet {
		if len(opts.Data) > 0 {
			rawBody = opts.Data
		} else if opts.DataMap != nil {
			jb, jerr := json.Marshal(opts.DataMap)
			if jerr != nil {
				return nil, jerr
			}
			rawBody = jb
			// 设置默认 Content-Type（修改本地 opts 副本，不影响调用方）
			if opts.Headers == nil {
				opts.Headers = map[string]string{"Content-Type": "application/json"}
			} else {
				if _, ok := opts.Headers["Content-Type"]; !ok {
					opts.Headers["Content-Type"] = "application/json"
				}
			}
		}
	}

	// 超时时间
	tout := 30
	if opts.Timeout > 0 {
		tout = opts.Timeout
	}

	// 重试参数
	retry := opts.Retry
	if retry < 0 {
		retry = 0
	}
	retryDelay := opts.RetryDelay
	if retryDelay <= 0 {
		retryDelay = 1
	}

	client := &http.Client{Transport: defaultTransport, Timeout: time.Duration(tout) * time.Second}

	var lastErr error

	for attempt := 0; attempt <= retry; attempt++ {
		// context with timeout for this attempt
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(tout)*time.Second)
		defer cancel()

		var bodyReader io.Reader
		if rawBody != nil {
			bodyReader = bytes.NewReader(rawBody)
		}

		req, err := http.NewRequestWithContext(ctx, opts.Method, u.String(), bodyReader)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		// headers
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request error: %w", err)
			// 网络/超时类错误重试
			if attempt < retry {
				time.Sleep(time.Duration(retryDelay) * time.Second)
				continue
			}
			return nil, lastErr
		}

		// ensure body closed
		var respBody []byte
		func() {
			defer resp.Body.Close()
			var reader io.Reader = resp.Body
			if opts.MaxBodySize > 0 {
				reader = io.LimitReader(resp.Body, opts.MaxBodySize)
			}
			b, rerr := io.ReadAll(reader)
			if rerr != nil {
				lastErr = fmt.Errorf("read body: %w", rerr)
				return
			}
			respBody = b
		}()

		// 判断状态码
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			// 对 5xx 做重试，对 4xx 一般不重试
			lastErr = fmt.Errorf("http status %d: %s", resp.StatusCode, string(respBody))
			if resp.StatusCode >= 500 && attempt < retry {
				time.Sleep(time.Duration(retryDelay) * time.Second)
				continue
			}
			return nil, lastErr
		}

		// 成功
		return respBody, nil
	}

	return nil, lastErr
}
