package httpclient

import (
	"net/http"
	"net/url"
)

// Response 表示一次 HTTP 响应，对应 Java HttpResponse。
type Response struct {
	statusCode int
	header     http.Header
	body       []byte
	requestURL *url.URL
	proto      string
}

// StatusCode 返回 HTTP 状态码。
func (r *Response) StatusCode() int {
	return r.statusCode
}

// OK 判断状态码是否为 200。
func (r *Response) OK() bool {
	return r.statusCode == http.StatusOK
}

// Header 返回响应头（只读使用，勿修改）。
func (r *Response) Header() http.Header {
	return r.header
}

// Bytes 返回响应体原始字节。
func (r *Response) Bytes() []byte {
	return r.body
}

// Body 以 UTF-8 返回响应体字符串。
func (r *Response) Body() string {
	return string(r.body)
}

// RequestURL 返回请求 URL。
func (r *Response) RequestURL() *url.URL {
	return r.requestURL
}

// Proto 返回协议版本，例如 HTTP/1.1。
func (r *Response) Proto() string {
	return r.proto
}
