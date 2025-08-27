package mhttp

/*
我想基于 github.com/gocolly/colly/v2 封装一个请求库


可以像如下这样调用

resData,err:= NewFetch({
	URL: "http://www.test.com",
	Data: []byte,
	DataMap: map[string]any,
	Params: map[string]string,
	Headers: map[string]string,
	Timeout: 10,
	Retry: 3,
	RetryDelay: 2,
}).Get()

传递了 Data 则会忽略 DataMap
Params 则代表在 URL 后面拼接的参数

上述的参数名字和 key 名 你都可以随意修改，我只是给你一个例子，
要支持 Get 和 Post 两种请求方式

*/

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gocolly/colly/v2"
)

// FetchOptions 请求选项
type FetchOptions struct {
	URL        string
	Data       []byte
	DataMap    map[string]any
	Params     map[string]string
	Headers    map[string]string
	Timeout    int // seconds
	Retry      int
	RetryDelay int // seconds
	Method     string
}

// Fetch 请求封装
type Fetch struct {
	opts FetchOptions
}

// NewFetch 创建一个 Fetch 实例
// 示例:
//
//	f := NewFetch(FetchOptions{URL: "http://example.com", Timeout: 10}).Get()
func NewFetch(opts FetchOptions) *Fetch {
	return &Fetch{opts: opts}
}

// Get 发起 GET 请求，并返回响应 body
func (f *Fetch) Get() ([]byte, error) {
	f.opts.Method = http.MethodGet
	return f.do()
}

// Post 发起 POST 请求，并返回响应 body
func (f *Fetch) Post() ([]byte, error) {
	f.opts.Method = http.MethodPost
	return f.do()
}

// do 执行请求，并支持重试
func (f *Fetch) do() ([]byte, error) {
	opts := f.opts

	if opts.URL == "" {
		return nil, errors.New("err:mhttp.Fetch|do|empty URL")
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

	// body 构造：保存原始字节切片，避免传递 typed-nil 到接口引起 panic
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
			// 如果用户未显式设置 Content-Type，则设置为 application/json
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

	var lastErr error
	var respBody []byte

	for attempt := 0; attempt <= retry; attempt++ {
		// 每次都创建新的 collector，避免状态污染
		c := colly.NewCollector()
		c.SetRequestTimeout(time.Duration(tout) * time.Second)

		respBody = nil
		c.OnResponse(func(r *colly.Response) {
			// 保存最后一次响应
			respBody = make([]byte, len(r.Body))
			copy(respBody, r.Body)
		})

		// 转换 headers
		hdr := http.Header{}
		for k, v := range opts.Headers {
			hdr.Set(k, v)
		}

		var body io.Reader
		if rawBody != nil {
			// 为每次请求创建新的 reader，接口为 io.Reader，避免 typed-nil
			body = bytes.NewReader(rawBody)
		}

		err = c.Request(opts.Method, u.String(), body, nil, hdr)
		if err == nil && respBody != nil {
			return respBody, nil
		}

		if err != nil {
			lastErr = err
		} else {
			lastErr = errors.New("err:mhttp.Fetch|do|empty response")
		}

		// 如果还有重试机会，等待后重试
		if attempt < retry {
			time.Sleep(time.Duration(retryDelay) * time.Second)
			continue
		}
	}

	return nil, lastErr
}
