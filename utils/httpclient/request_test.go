package httpclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestConfigureLongLived(t *testing.T) {
	EnableLongLivedConnections(16)
	clientMu.RLock()
	idle := idleConnTimeout
	perHost := maxIdleConnsPerHost
	clientMu.RUnlock()
	if idle != 0 {
		t.Fatalf("IdleConnTimeout=%v want 0", idle)
	}
	if perHost != 16 {
		t.Fatalf("MaxIdleConnsPerHost=%d want 16", perHost)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer ts.Close()

	body, err := Get(ts.URL).Executor()
	if err != nil || body != "ok" {
		t.Fatalf("body=%q err=%v", body, err)
	}
}

func TestGetExecutor(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method=%s", r.Method)
		}
		if r.URL.Query().Get("q") != "1" {
			t.Errorf("query=%s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte("hello"))
	}))
	defer ts.Close()

	body, err := Get(ts.URL).Params(map[string]string{"q": "1"}).Executor()
	if err != nil {
		t.Fatal(err)
	}
	if body != "hello" {
		t.Fatalf("body=%q", body)
	}
}

func TestPostJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("content-type=%s", ct)
		}
		data, _ := io.ReadAll(r.Body)
		var m map[string]any
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatal(err)
		}
		if m["name"] != "alice" {
			t.Fatalf("body=%s", data)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer ts.Close()

	body, err := Post(ts.URL).WithJSON(map[string]any{"name": "alice"}).Executor()
	if err != nil {
		t.Fatal(err)
	}
	if body != `{"ok":true}` {
		t.Fatalf("body=%q", body)
	}
}

func TestPostForm(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("a") != "1" || r.Form.Get("b") != "x y" {
			t.Fatalf("form=%v", r.Form)
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer ts.Close()

	got, err := Post(ts.URL).WithFormData(map[string]string{"a": "1", "b": "x y"}).Executor()
	if err != nil || got != "ok" {
		t.Fatalf("got=%q err=%v", got, err)
	}
}

func TestExecutorNonOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad"))
	}))
	defer ts.Close()

	_, err := Get(ts.URL).Executor()
	if err == nil {
		t.Fatal("expected error")
	}

	resp, err := Get(ts.URL).Do(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if resp.OK() || resp.Body() != "bad" {
		t.Fatalf("resp=%d %q", resp.StatusCode(), resp.Body())
	}
}

func TestAsync(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("async"))
	}))
	defer ts.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	Get(ts.URL).Async(CallbackFunc{
		OnSuccess: func(resp *Response) {
			defer wg.Done()
			if resp.Body() != "async" {
				t.Errorf("body=%q", resp.Body())
			}
		},
		OnFail: func(req *http.Request, err error) {
			defer wg.Done()
			t.Errorf("unexpected fail: %v", err)
		},
	})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("async timeout")
	}
}

func TestInvalidScheme(t *testing.T) {
	_, err := Get("ftp://example.com").Executor()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHeaderAndBytes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token") != "t1" {
			t.Errorf("token=%s", r.Header.Get("X-Token"))
		}
		data, _ := io.ReadAll(r.Body)
		_, _ = w.Write(data)
	}))
	defer ts.Close()

	got, err := Post(ts.URL).Header("X-Token", "t1").WithStringData("ping").ExecutorBytes()
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "ping" {
		t.Fatalf("got=%q", got)
	}
}
