package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// PostRequest POST 请求，对应 Java PostHttpRequest。
// 仅提供 body 相关能力，避免误用 GET 查询参数方法。
type PostRequest struct {
	*baseRequest
}

// Post 创建 POST 请求，对应 Java HttpRequest.post。
func Post(rawURL string) *PostRequest {
	return &PostRequest{baseRequest: newBaseRequest(http.MethodPost, rawURL)}
}

// Header 设置单个请求头。
func (r *PostRequest) Header(name, value string) *PostRequest {
	r.setHeader(name, value)
	return r
}

// Headers 批量设置请求头。
func (r *PostRequest) Headers(m map[string]string) *PostRequest {
	r.setHeaders(m)
	return r
}

// Charset 设置 Accept-Charset（默认 utf-8）。
func (r *PostRequest) Charset(charset string) *PostRequest {
	r.setCharset(charset)
	return r
}

// CloseConnectAlive 请求完成后关闭连接（Connection: close）。
func (r *PostRequest) CloseConnectAlive() *PostRequest {
	r.setCloseConnectAlive()
	return r
}

// SetReadTimeout 设置读超时（秒）。
func (r *PostRequest) SetReadTimeout(seconds int) *PostRequest {
	r.setReadTimeout(seconds)
	return r
}

// WithFormData 以 application/x-www-form-urlencoded 提交。对应 Java withFormData。
func (r *PostRequest) WithFormData(params map[string]string) *PostRequest {
	r.setHeader("Content-Type", "application/x-www-form-urlencoded")
	vals := make(url.Values, len(params))
	for k, v := range params {
		vals.Set(k, v)
	}
	r.body = []byte(vals.Encode())
	return r
}

// WithJSONString 以 application/json 提交 JSON 字符串。
func (r *PostRequest) WithJSONString(jsonStr string) *PostRequest {
	r.setHeader("Content-Type", "application/json")
	r.body = []byte(jsonStr)
	return r
}

// WithJSON 将对象序列化为 JSON 后提交。对应 Java withJsonData。
func (r *PostRequest) WithJSON(v any) *PostRequest {
	b, err := json.Marshal(v)
	if err != nil {
		r.buildErr = fmt.Errorf("utils/httpclient: json marshal: %w", err)
		return r
	}
	return r.WithJSONString(string(b))
}

// WithStringData 以 text/plain 提交字符串。对应 Java withStringData。
func (r *PostRequest) WithStringData(data string) *PostRequest {
	r.setHeader("Content-Type", "text/plain")
	r.body = []byte(data)
	return r
}

// WithBytes 以 multipart/form-data Content-Type 提交原始字节（与 Java withBytes 对齐）。
func (r *PostRequest) WithBytes(data []byte) *PostRequest {
	r.setHeader("Content-Type", "multipart/form-data")
	r.body = append([]byte(nil), data...)
	return r
}

// CustomBody 自定义 body，不改 Content-Type。对应 Java customBody。
func (r *PostRequest) CustomBody(data []byte) *PostRequest {
	r.body = append([]byte(nil), data...)
	return r
}
