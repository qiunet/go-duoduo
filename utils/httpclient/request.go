package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// baseRequest 请求公共实现（未导出），供 GetRequest / PostRequest 复用。
type baseRequest struct {
	method      string
	rawURL      string
	headers     http.Header
	query       url.Values
	body        []byte
	readTimeout time.Duration
	closeAlive  bool
	buildErr    error
}

// newBaseRequest 创建带默认头的公共请求。
func newBaseRequest(method, rawURL string) *baseRequest {
	h := make(http.Header)
	h.Set("Accept", "*/*")
	h.Set("Accept-Charset", "utf-8")
	return &baseRequest{
		method:      method,
		rawURL:      rawURL,
		headers:     h,
		query:       make(url.Values),
		readTimeout: time.Duration(defaultReadTimeoutSec) * time.Second,
	}
}

// setHeader 设置单个请求头。
func (r *baseRequest) setHeader(name, value string) {
	r.headers.Set(name, value)
}

// setHeaders 批量设置请求头。
func (r *baseRequest) setHeaders(m map[string]string) {
	for k, v := range m {
		r.headers.Set(k, v)
	}
}

// setCharset 设置 Accept-Charset。
func (r *baseRequest) setCharset(charset string) {
	r.headers.Set("Accept-Charset", charset)
}

// setCloseConnectAlive 标记 Connection: close。
func (r *baseRequest) setCloseConnectAlive() {
	r.closeAlive = true
	r.headers.Set("Connection", "close")
}

// setReadTimeout 设置读超时（秒）。
func (r *baseRequest) setReadTimeout(seconds int) {
	if seconds > 0 {
		r.readTimeout = time.Duration(seconds) * time.Second
	}
}

// Do 同步执行请求并返回完整响应（不强制状态码 200）。
func (r *baseRequest) Do(ctx context.Context) (*Response, error) {
	if r.buildErr != nil {
		return nil, r.buildErr
	}

	if _, ok := ctx.Deadline(); !ok && r.readTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.readTimeout)
		defer cancel()
	}

	req, u, err := r.buildHTTPRequest(ctx)
	if err != nil {
		return nil, err
	}

	httpResp, err := getClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	limited := io.LimitReader(httpResp.Body, maxBodySize()+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBodySize() {
		return nil, fmt.Errorf("utils/httpclient: response body exceeds max size %d", maxBodySize())
	}

	return &Response{
		statusCode: httpResp.StatusCode,
		header:     httpResp.Header.Clone(),
		body:       data,
		requestURL: u,
		proto:      httpResp.Proto,
	}, nil
}

// buildHTTPRequest 组装标准库 http.Request。
func (r *baseRequest) buildHTTPRequest(ctx context.Context) (*http.Request, *url.URL, error) {
	u, err := url.Parse(r.rawURL)
	if err != nil {
		return nil, nil, fmt.Errorf("utils/httpclient: invalid url: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, nil, fmt.Errorf("utils/httpclient: url [%s] not a valid http request", r.rawURL)
	}

	if len(r.query) > 0 {
		q := u.Query()
		for k, vs := range r.query {
			for _, v := range vs {
				q.Set(k, v)
			}
		}
		u.RawQuery = q.Encode()
	}

	var bodyReader io.Reader
	if len(r.body) > 0 {
		bodyReader = bytes.NewReader(r.body)
	}

	req, err := http.NewRequestWithContext(ctx, r.method, u.String(), bodyReader)
	if err != nil {
		return nil, nil, err
	}
	for k, vs := range r.headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	if r.closeAlive {
		req.Close = true
	}
	return req, u, nil
}

// Executor 同步执行并要求状态码 200，返回响应体字符串。
func (r *baseRequest) Executor() (string, error) {
	return r.ExecutorContext(context.Background())
}

// ExecutorContext 带 context 的同步执行，要求状态码 200，返回 body 字符串。
func (r *baseRequest) ExecutorContext(ctx context.Context) (string, error) {
	resp, err := r.Do(ctx)
	if err != nil {
		return "", fmt.Errorf("http client send request error: %w", err)
	}
	if !resp.OK() {
		return "", fmt.Errorf("request: %s fail, statusCode %d", resp.RequestURL(), resp.StatusCode())
	}
	return resp.Body(), nil
}

// ExecutorBytes 同步执行并要求状态码 200，返回响应体字节。
func (r *baseRequest) ExecutorBytes() ([]byte, error) {
	resp, err := r.Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("http client send request error: %w", err)
	}
	if !resp.OK() {
		return nil, fmt.Errorf("request: %s fail, statusCode %d", resp.RequestURL(), resp.StatusCode())
	}
	return resp.Bytes(), nil
}

// Async 异步执行请求；在新 goroutine 中回调。
func (r *baseRequest) Async(cb Callback) {
	go func() {
		resp, err := r.Do(context.Background())
		if err != nil {
			if cb != nil {
				cb.Fail(nil, err)
			}
			return
		}
		if cb != nil {
			cb.Success(resp)
		}
	}()
}
