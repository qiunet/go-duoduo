package httpclient

import "net/http"

// Callback 异步请求回调，对应 Java IHttpCallBack。
type Callback interface {
	// Success 请求成功（已拿到响应，不保证状态码一定为 200）。
	Success(resp *Response)
	// Fail 请求失败（网络错误、超时、构建请求失败等）。
	Fail(req *http.Request, err error)
}

// CallbackFunc 用函数适配 Callback；Fail 可选。
type CallbackFunc struct {
	OnSuccess func(resp *Response)
	OnFail    func(req *http.Request, err error)
}

// Success 调用成功回调。
func (c CallbackFunc) Success(resp *Response) {
	if c.OnSuccess != nil {
		c.OnSuccess(resp)
	}
}

// Fail 调用失败回调；未设置时忽略。
func (c CallbackFunc) Fail(req *http.Request, err error) {
	if c.OnFail != nil {
		c.OnFail(req, err)
	}
}
