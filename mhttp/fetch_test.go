package mhttp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// go test -v -run Test_mo7
// func Test_mo7Get(t *testing.T) {
// 	type ResType struct {
// 		Code     int    `json:"code"`
// 		Question string `json:"question"`
// 		Answer   string `json:"answer"`
// 	}

// 	question := "mo7cc"

// 	res, err := NewFetch(FetchOptions{
// 		URL: "https://uapis.cn/api/v1/answerbook/ask",
// 		Params: map[string]string{
// 			"question": question,
// 		},
// 		Method: "GET",
// 	}).Do()
// 	if err != nil {
// 		t.Fatalf("请求出错 error: %v", err)
// 	}
// 	var data ResType
// 	if err := json.Unmarshal(res, &data); err != nil {
// 		t.Fatalf("解析结果出错 error: %v", err)
// 	}
// 	if data.Code != 200 {
// 		t.Fatalf("请求失败 error: %s", string(res))
// 	}
// 	if data.Question != question {
// 		t.Fatalf("请求结果不符合预期 error: %s", data.Question)
// 	}
// 	t.Logf("Get响应结果: %+v", data)
// }

// go test -v -run Test_mo7Post
// func Test_mo7Post(t *testing.T) {
// 	type ResType struct {
// 		Code     int    `json:"code"`
// 		Question string `json:"question"`
// 		Answer   string `json:"answer"`
// 	}
// 	question := "mo7cc"
// 	res, err := NewFetch(FetchOptions{
// 		URL: "https://uapis.cn/api/v1/answerbook/ask",
// 		DataMap: map[string]any{
// 			"question": question,
// 		},
// 		Method: "POST",
// 	}).Do()
// 	if err != nil {
// 		t.Fatalf("请求出错 error: %v", err)
// 	}
// 	var data ResType
// 	if err := json.Unmarshal(res, &data); err != nil {
// 		t.Fatalf("解析结果出错 error: %v", err)
// 	}
// 	if data.Code != 200 {
// 		t.Fatalf("请求失败 error: %s", string(res))
// 	}
// 	if data.Question != question {
// 		t.Fatalf("请求结果不符合预期 error: %s", data.Question)
// 	}
// 	t.Logf("Post响应结果: %+v", data)
// }

func TestDo_Get_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello-get"))
	}))
	defer srv.Close()

	res, err := NewFetch(FetchOptions{URL: srv.URL, Method: http.MethodGet, Timeout: 5}).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(res) != "hello-get" {
		t.Fatalf("unexpected body: %s", string(res))
	}
}

func TestDo_Post_JSON(t *testing.T) {
	var gotContentType string
	var gotBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		body, _ := io.ReadAll(r.Body)
		gotBody = body
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	data := map[string]any{"k": "v", "n": 1}
	res, err := NewFetch(FetchOptions{URL: srv.URL, Method: http.MethodPost, DataMap: data, Timeout: 5}).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expect, _ := json.Marshal(data)
	if string(res) != string(expect) {
		t.Fatalf("unexpected response body: %s", string(res))
	}
	if !strings.Contains(gotContentType, "application/json") {
		t.Fatalf("expected content-type application/json, got %q", gotContentType)
	}
	if string(gotBody) != string(expect) {
		t.Fatalf("server got unexpected body: %s", string(gotBody))
	}
}

func TestDo_Params_Append(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.URL.RawQuery))
	}))
	defer srv.Close()

	res, err := NewFetch(FetchOptions{URL: srv.URL, Method: http.MethodGet, Params: map[string]string{"a": "1", "b": "x"}, Timeout: 5}).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	qs := string(res)
	if !strings.Contains(qs, "a=1") || !strings.Contains(qs, "b=x") {
		t.Fatalf("unexpected query string: %s", qs)
	}
}

func TestDo_InvalidMethod_Error(t *testing.T) {
	_, err := NewFetch(FetchOptions{URL: "http://example.com", Method: "BAD"}).Do()
	if err == nil || !strings.Contains(err.Error(), "invalid method") {
		t.Fatalf("expected invalid method error, got: %v", err)
	}
}

func TestDo_EmptyURL_Error(t *testing.T) {
	_, err := NewFetch(FetchOptions{Method: http.MethodGet}).Do()
	if err == nil || !strings.Contains(err.Error(), "empty URL") {
		t.Fatalf("expected empty URL error, got: %v", err)
	}
}

func TestDo_StatusCodeNon2xx_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("server-error"))
	}))
	defer srv.Close()

	_, err := NewFetch(FetchOptions{URL: srv.URL, Method: http.MethodGet, Timeout: 5}).Do()
	if err == nil || !strings.Contains(err.Error(), "http status 500") {
		t.Fatalf("expected http status 500 error, got: %v", err)
	}
	if err != nil && !strings.Contains(err.Error(), "server-error") {
		t.Fatalf("expected error to include response body, got: %v", err)
	}
}

func TestDo_MaxBodySize_Truncate(t *testing.T) {
	big := strings.Repeat("a", 1024)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(big))
	}))
	defer srv.Close()

	res, err := NewFetch(FetchOptions{URL: srv.URL, Method: http.MethodGet, Timeout: 5, MaxBodySize: 10}).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 10 {
		t.Fatalf("expected truncated length 10, got %d", len(res))
	}
}

func TestDo_Retry_On5xx(t *testing.T) {
	count := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		if count <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("err"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	res, err := NewFetch(FetchOptions{URL: srv.URL, Method: http.MethodGet, Timeout: 5, Retry: 2, RetryDelay: 0}).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(res) != "ok" {
		t.Fatalf("expected ok after retries, got %s", string(res))
	}
}
