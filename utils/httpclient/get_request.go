package httpclient

import "net/http"

// GetRequest GET 请求，对应 Java GetHttpRequest。
// 仅提供查询参数相关能力，避免误用 POST body 方法。
type GetRequest struct {
	*baseRequest
}

// Get 创建 GET 请求，对应 Java HttpRequest.get。
func Get(rawURL string) *GetRequest {
	return &GetRequest{baseRequest: newBaseRequest(http.MethodGet, rawURL)}
}

// Header 设置单个请求头。
func (r *GetRequest) Header(name, value string) *GetRequest {
	r.setHeader(name, value)
	return r
}

// Headers 批量设置请求头。
func (r *GetRequest) Headers(m map[string]string) *GetRequest {
	r.setHeaders(m)
	return r
}

// Charset 设置 Accept-Charset（默认 utf-8）。
func (r *GetRequest) Charset(charset string) *GetRequest {
	r.setCharset(charset)
	return r
}

// CloseConnectAlive 请求完成后关闭连接（Connection: close）。
func (r *GetRequest) CloseConnectAlive() *GetRequest {
	r.setCloseConnectAlive()
	return r
}

// SetReadTimeout 设置读超时（秒）。
func (r *GetRequest) SetReadTimeout(seconds int) *GetRequest {
	r.setReadTimeout(seconds)
	return r
}

// Params 设置查询参数（拼到 URL 后）。对应 Java GetHttpRequest.params。
func (r *GetRequest) Params(params map[string]string) *GetRequest {
	for k, v := range params {
		r.query.Set(k, v)
	}
	return r
}
