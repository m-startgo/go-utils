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

/*
请求配置结构体，用于创建HTTP请求

字段说明：
  - Url: 请求URL地址（必填）
  - Params: URL查询参数（可选）
  - Data: 请求体数据（可选，POST请求使用）
  - Header: 请求头信息（可选）
  - Timeout: 请求超时时间（可选，单位：秒，默认10秒）
  - Retry: 请求重试次数（可选，默认0次不重试）
  - RetryDelay: 重试延迟时间（可选，单位：秒，默认1秒）
*/
type RequestConfig struct {
	Url        string            `json:"url"`
	Params     map[string]any    `json:"params"`
	Data       map[string]any    `json:"data"`
	Header     map[string]string `json:"header"`
	Timeout    int               `json:"timeout"`
	Retry      int               `json:"retry"`
	RetryDelay int               `json:"retryDelay"`
}

/*
HTTP请求客户端结构体
*/
type Request struct {
	config *RequestConfig
	client *http.Client
}

/*
Fetch 创建一个新的HTTP请求实例

参数：
  - config: 请求配置信息

返回值：
  - *Request: 请求实例
  - error: 配置验证错误
*/
func Fetch(config RequestConfig) (*Request, error) {
	// 验证必填参数
	if config.Url == "" {
		return nil, errors.New("请求URL不能为空")
	}

	// 设置默认值
	if config.Timeout <= 0 {
		config.Timeout = 30
	}
	if config.RetryDelay <= 0 {
		config.RetryDelay = 1
	}

	// 创建HTTP客户端
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
func (r *Request) Get() (*Response, error) {
	return r.sendRequest("GET", nil)
}

/*
Response 封装HTTP响应信息

字段说明：
  - StatusCode: HTTP状态码
  - Header: 响应头信息
  - Body: 响应体内容
*/
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
func (r *Request) Post() (*Response, error) {
	var contentType string
	if r.config.Header != nil {
		contentType = r.config.Header["Content-Type"]
	}
	var body *bytes.Buffer

	// 根据Content-Type选择合适的序列化方式
	if contentType == "application/x-www-form-urlencoded" {
		formData := url.Values{}
		if r.config.Data != nil {
			for key, value := range r.config.Data {
				formData.Add(key, fmt.Sprintf("%v", value))
			}
		}
		body = bytes.NewBufferString(formData.Encode())
	} else {
		// 默认JSON处理
		var jsonData []byte
		var err error
		if r.config.Data != nil {
			jsonData, err = json.Marshal(r.config.Data)
			if err != nil {
				return nil, fmt.Errorf("请求数据JSON序列化失败: %w", err)
			}
		} else {
			jsonData = []byte("{}")
		}
		body = bytes.NewBuffer(jsonData)
	}
	return r.sendRequest("POST", body)
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
	// 解析URL
	parsedUrl, urlErr := url.Parse(r.config.Url)
	if urlErr != nil {
		return nil, fmt.Errorf("URL解析失败: %w", urlErr)
	}

	// 添加查询参数
	query := parsedUrl.Query()
	if r.config.Params != nil {
		for key, value := range r.config.Params {
			query.Add(key, fmt.Sprintf("%v", value))
		}
	}
	parsedUrl.RawQuery = query.Encode()

	// 创建请求
	var reqBody io.Reader
	if body != nil {
		reqBody = body
	} else {
		reqBody = nil
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

	// 设置默认Content-Type
	if method == "POST" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	// 实现重试机制
	maxRetries := r.config.Retry
	var resp *http.Response
	var err error

	for i := 0; i <= maxRetries; i++ {
		resp, err = r.client.Do(req)
		if err == nil && resp != nil && resp.StatusCode < 500 {
			break
		}

		// 重试延迟
		if i < maxRetries {
			time.Sleep(time.Duration(r.config.RetryDelay) * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("请求执行失败: %w", err)
	}
	if resp == nil {
		return nil, errors.New("未获得有效的响应对象")
	}

	// 读取响应体并自动关闭
	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		if resp.Body != nil {
			resp.Body.Close()
		}
		return nil, fmt.Errorf("读取响应体失败: %w", readErr)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       respBody,
	}, nil
}

/*
使用示例：

package main

import (
	"fmt"
	"xxn-aihf/utils"
)

func main() {
	// 创建请求配置
	config := utils.RequestConfig{
		Url: "https://example.com/api/resource",
		Params: map[string]any{
			"id":     12345,
			"page":   1,
			"filter": "active",
		},
		Data: map[string]any{
			"key": "value",
		},
		Header: map[string]string{
			"Authorization": "Bearer your_access_token",
			"User-Agent":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		},
		Timeout:    10,
		Retry:      3,
		RetryDelay: 2,
	}

	// 创建请求实例
	request, err := utils.Fetch(config)
	if err != nil {
		fmt.Printf("配置错误: %v\n", err)
		return
	}

	// 发送GET请求
	resp, err := request.Get()
	if err != nil {
		fmt.Printf("GET请求失败: %v\n", err)
		return
	}
	fmt.Printf("GET请求成功: 状态码=%d, 响应体=%s\n", resp.StatusCode, string(resp.Body))

	// 发送POST请求
	resp, err = request.Post()
	if err != nil {
		fmt.Printf("POST请求失败: %v\n", err)
		return
	}
	fmt.Printf("POST请求成功: 状态码=%d, 响应体=%s\n", resp.StatusCode, string(resp.Body))
}
*/
