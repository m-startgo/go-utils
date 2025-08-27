package mhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestFetch_Get_ParamsAndHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("q") != "1" {
			t.Fatalf("expected query q=1, got %s", r.URL.Query().Get("q"))
		}
		if r.Header.Get("X-Test") != "ok" {
			t.Fatalf("expected header X-Test=ok, got %s", r.Header.Get("X-Test"))
		}
		w.Write([]byte("hello"))
	}))
	defer ts.Close()

	opts := FetchOptions{
		URL:     ts.URL,
		Params:  map[string]string{"q": "1"},
		Headers: map[string]string{"X-Test": "ok"},
		Timeout: 5,
	}

	b, err := NewFetch(opts).Get()
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if !bytes.Equal(b, []byte("hello")) {
		t.Fatalf("unexpected body: %s", string(b))
	}
}

func TestFetch_Post_DataBytes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		b, _ := io.ReadAll(r.Body)
		if string(b) != "payload-bytes" {
			t.Fatalf("unexpected body: %s", string(b))
		}
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	opts := FetchOptions{
		URL:     ts.URL,
		Data:    []byte("payload-bytes"),
		Headers: map[string]string{"X-Test": "post"},
		Timeout: 5,
	}

	b, err := NewFetch(opts).Post()
	if err != nil {
		t.Fatalf("Post error: %v", err)
	}
	if string(b) != "ok" {
		t.Fatalf("unexpected resp: %s", string(b))
	}
}

func TestFetch_Post_DataMap_JSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		b, _ := io.ReadAll(r.Body)
		var got map[string]any
		if err := json.Unmarshal(b, &got); err != nil {
			t.Fatalf("invalid json body: %v", err)
		}
		if got["a"] != "b" {
			t.Fatalf("unexpected json content: %v", got)
		}
		w.Write([]byte("okjson"))
	}))
	defer ts.Close()

	opts := FetchOptions{
		URL:     ts.URL,
		DataMap: map[string]any{"a": "b"},
		Timeout: 5,
	}

	b, err := NewFetch(opts).Post()
	if err != nil {
		t.Fatalf("Post(DataMap) error: %v", err)
	}
	if string(b) != "okjson" {
		t.Fatalf("unexpected resp: %s", string(b))
	}
}

func TestFetch_Retry_OnTimeout(t *testing.T) {
	var mu sync.Mutex
	count := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		c := count
		count++
		mu.Unlock()
		// 首次请求睡眠，导致超时；后续请求正常返回
		if c == 0 {
			time.Sleep(2 * time.Second)
			// 如果超时发生，客户端会已经返回；但 handler 仍会 finish
			return
		}
		w.Write([]byte("retry-ok"))
	}))
	defer ts.Close()

	opts := FetchOptions{
		URL:        ts.URL,
		Timeout:    1, // 1s 超时
		Retry:      1,
		RetryDelay: 1,
	}

	b, err := NewFetch(opts).Get()
	if err != nil {
		t.Fatalf("expected success after retry, got error: %v", err)
	}
	if string(b) != "retry-ok" {
		t.Fatalf("unexpected body after retry: %s", string(b))
	}
}

// go test -v -run Test_mo7
func Test_mo7(t *testing.T) {
	res, err := NewFetch(FetchOptions{
		URL:     "https://v1.hitokoto.cn",
		Timeout: 10,
	}).Get()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("res", string(res))
}
