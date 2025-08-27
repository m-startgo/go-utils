package mhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// RequestConfig 请求配置，用于创建 HTTP 请求。
//
// 字段说明：
//   - URL: 请求 URL 地址（必填）
//   - Params: URL 查询参数（可选）
//   - Data: 请求体数据（可选，POST 请求使用）
//   - Header: 请求头信息（可选）
//   - Timeout: 请求超时时间（单位：秒，默认 10）
//   - Retry: 请求重试次数（可选，默认 0）
//   - RetryDelay: 重试延迟时间（单位：秒，默认 1）
type RequestConfig struct {
	URL        string            `json:"url"`
	Params     map[string]any    `json:"params"`
	Data       map[string]any    `json:"data"`
	Header     map[string]string `json:"header"`
	Timeout    int               `json:"timeout"`
	Retry      int               `json:"retry"`
	RetryDelay int               `json:"retryDelay"`
}

// Request 是一个简单的 HTTP 客户端封装，绑定了配置和 http.Client。
type Request struct {
	config *RequestConfig
	client *http.Client
}

// Fetch 根据 RequestConfig 创建并返回 *Request。
//
// 注意：如果配置无效会返回错误。函数会为未设置的 Timeout/RetryDelay 填充默认值。
func Fetch(config RequestConfig) (*Request, error) {
	// 验证必填参数
	if config.URL == "" {
		return nil, errors.New("请求 URL 不能为空")
	}

	// 设置默认值（与注释一致，默认 10 秒）
	if config.Timeout <= 0 {
		config.Timeout = 10
	}
	if config.RetryDelay <= 0 {
		config.RetryDelay = 1
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	return &Request{
		config: &config,
		client: client,
	}, nil
}

/*
Get 发送GET请求

返回值：
  - *Response: 响应信息
  - error: 请求错误
*/
// Get 发送 GET 请求并返回响应。
func (r *Request) Get() (*Response, error) {
	return r.sendRequest("GET", nil)
}

// Response 封装 HTTP 响应信息。
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

/*
Post 发送POST请求

返回值：
  - *Response: 响应信息
  - error: 请求错误
*/
// Post 发送 POST 请求，自动根据 Content-Type 编码请求体（支持 JSON 与 form）。
func (r *Request) Post() (*Response, error) {
	var contentType string
	if r.config.Header != nil {
		contentType = r.config.Header["Content-Type"]
	}

	var bodyBytes []byte
	// 根据 Content-Type 选择合适的序列化方式
	if contentType == "application/x-www-form-urlencoded" {
		formData := url.Values{}
		if r.config.Data != nil {
			for key, value := range r.config.Data {
				formData.Add(key, fmt.Sprintf("%v", value))
			}
		}
		bodyBytes = []byte(formData.Encode())
	} else {
		// 默认 JSON 处理
		var jsonData []byte
		var err error
		if r.config.Data != nil {
			jsonData, err = json.Marshal(r.config.Data)
			if err != nil {
				return nil, fmt.Errorf("请求数据 JSON 序列化失败: %w", err)
			}
		} else {
			jsonData = []byte("{}")
		}
		bodyBytes = jsonData
	}

	return r.sendRequestBytes("POST", bodyBytes)
}

// sendRequestBytes 是 sendRequest 的变体，接受已序列化的请求体字节切片，
// 以便在重试时可以为每次请求重新创建 io.Reader，避免 body 被消费后不可重用的问题。
func (r *Request) sendRequestBytes(method string, bodyBytes []byte) (*Response, error) {
	var bodyBuf *bytes.Buffer
	if bodyBytes != nil {
		bodyBuf = bytes.NewBuffer(bodyBytes)
	}
	return r.sendRequest(method, bodyBuf)
}

/*
sendRequest 发送HTTP请求的内部实现

参数：
  - method: HTTP方法(GET/POST等)
  - body: 请求体

返回值：
  - *Response: 响应信息
  - error: 请求错误
*/
func (r *Request) sendRequest(method string, body *bytes.Buffer) (*Response, error) {
	// 解析 URL
	parsedUrl, urlErr := url.Parse(r.config.URL)
	if urlErr != nil {
		return nil, fmt.Errorf("URL 解析失败: %w", urlErr)
	}

	// 添加查询参数
	query := parsedUrl.Query()
	if r.config.Params != nil {
		for key, value := range r.config.Params {
			query.Add(key, fmt.Sprintf("%v", value))
		}
	}
	parsedUrl.RawQuery = query.Encode()

	// 为了支持重试并且保证请求体在每次重试都可用，需要在每次请求前重新创建 *http.Request。
	maxRetries := r.config.Retry

	var lastErr error
	var resp *http.Response

	// 将 body 的内容保存在字节切片中，供每次重试时创建新的 reader
	var bodyBytes []byte
	if body != nil {
		bodyBytes = body.Bytes()
	}

	for i := 0; i <= maxRetries; i++ {
		// 为每次请求创建新的请求实例
		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, reqErr := http.NewRequest(method, parsedUrl.String(), reqBody)
		if reqErr != nil {
			return nil, fmt.Errorf("创建请求失败: %w", reqErr)
		}

		// 设置请求头
		if r.config.Header != nil {
			for key, value := range r.config.Header {
				req.Header.Set(key, value)
			}
		}

		// 设置默认 Content-Type（仅在 POST 且未设置时）
		if method == "POST" && req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
		}

		resp, lastErr = r.client.Do(req)
		// 如果没有网络错误并且状态码小于 500 则认为成功（不再重试）
		if lastErr == nil && resp != nil && resp.StatusCode < 500 {
			break
		}

		// 如果收到响应且准备重试，记得关闭 body 以免泄露连接
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
			resp = nil
		}

		// 重试延迟
		if i < maxRetries {
			time.Sleep(time.Duration(r.config.RetryDelay) * time.Second)
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("请求执行失败: %w", lastErr)
	}
	if resp == nil {
		return nil, errors.New("未获得有效的响应对象")
	}

	// 确保在退出前关闭响应体
	defer resp.Body.Close()

	// 读取响应体
	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", readErr)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       respBody,
	}, nil
}
