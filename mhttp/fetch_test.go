package mhttp_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/m-startgo/go-utils/mhttp"
)

func TestGet(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("id") != "12345" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing id"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}

	srv := httptest.NewServer(http.HandlerFunc(h))
	defer srv.Close()

	cfg := mhttp.RequestConfig{
		URL:     srv.URL,
		Params:  map[string]any{"id": 12345},
		Timeout: 5,
	}
	req, err := mhttp.Fetch(cfg)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	resp, err := req.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	var dst map[string]any
	if err := json.Unmarshal(resp.Body, &dst); err != nil {
		t.Fatalf("invalid body json: %v", err)
	}
	if _, ok := dst["ok"]; !ok {
		t.Fatalf("expected ok in response")
	}
}

func TestPostJSON(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if !strings.Contains(ct, "application/json") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid content-type"))
			return
		}
		body, _ := io.ReadAll(r.Body)
		var data map[string]any
		if err := json.Unmarshal(body, &data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid json"))
			return
		}
		if data["key"] != "value" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid payload"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("ok"))
	}

	srv := httptest.NewServer(http.HandlerFunc(h))
	defer srv.Close()

	cfg := mhttp.RequestConfig{
		URL:     srv.URL,
		Data:    map[string]any{"key": "value"},
		Header:  map[string]string{"User-Agent": "unit-test"},
		Timeout: 5,
	}

	req, err := mhttp.Fetch(cfg)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	resp, err := req.Post()
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
}

func TestPostForm(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if !strings.Contains(ct, "application/x-www-form-urlencoded") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid content-type"))
			return
		}
		body, _ := io.ReadAll(r.Body)
		vals, _ := url.ParseQuery(string(body))
		if vals.Get("a") != "1" || vals.Get("b") != "2" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid form"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}

	srv := httptest.NewServer(http.HandlerFunc(h))
	defer srv.Close()

	cfg := mhttp.RequestConfig{
		URL:     srv.URL,
		Data:    map[string]any{"a": 1, "b": 2},
		Header:  map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		Timeout: 5,
	}

	req, err := mhttp.Fetch(cfg)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	resp, err := req.Post()
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
}

func TestRetrySuccess(t *testing.T) {
	count := 0
	h := func(w http.ResponseWriter, r *http.Request) {
		if count == 0 {
			count++
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("fail"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}

	srv := httptest.NewServer(http.HandlerFunc(h))
	defer srv.Close()

	cfg := mhttp.RequestConfig{
		URL:        srv.URL,
		Timeout:    5,
		Retry:      1,
		RetryDelay: 1,
	}

	req, err := mhttp.Fetch(cfg)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	start := time.Now()
	resp, err := req.Get()
	delta := time.Since(start)
	if err != nil {
		t.Fatalf("expected success after retry, got err: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	// Ensure there was at least one retry delay (approx)
	if delta < time.Duration(cfg.RetryDelay)*time.Second {
		t.Fatalf("expected at least retry delay, took %v", delta)
	}
}

func TestTimeout(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		// sleep longer than client timeout
		time.Sleep(1500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(http.HandlerFunc(h))
	defer srv.Close()

	cfg := mhttp.RequestConfig{
		URL:     srv.URL,
		Timeout: 1, // seconds
		Retry:   0,
	}

	req, err := mhttp.Fetch(cfg)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	_, err = req.Get()
	if err == nil {
		t.Fatalf("expected timeout error")
	}
}
