package httpclient

import (
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	defaultMaxReceivedContentLength = 1024 * 128
	defaultConnectTimeoutSec        = 10
	defaultReadTimeoutSec           = 10
	defaultKeepAlive                = 30 * time.Second
	defaultIdleConnTimeout          = 90 * time.Second
	defaultMaxIdleConns             = 100
	defaultMaxIdleConnsPerHost      = 32
)

// PoolConfig 连接池 / 长连接相关配置。
// 网关到少量逻辑服、希望尽量不释放 TCP 时，可将 IdleConnTimeout 设为 0（表示不因空闲关闭）。
type PoolConfig struct {
	// ConnectTimeout 建连超时；零值表示保持现有/默认。
	ConnectTimeout time.Duration
	// KeepAlive TCP 保活探测间隔；零值表示保持现有/默认。
	KeepAlive time.Duration
	// IdleConnTimeout 空闲连接多久后关闭；0 表示不因空闲关闭（长连接）。
	// 未设置时请用负数以外的哨兵：调用 Configure 时若字段为 0 且想表示「不关闭」，需显式传 0。
	// 为区分「未改」与「设为 0」，使用 IdleConnTimeoutSet。
	IdleConnTimeout time.Duration
	// IdleConnTimeoutSet 为 true 时才应用 IdleConnTimeout（含 0=永不因空闲关闭）。
	IdleConnTimeoutSet bool
	// MaxIdleConns 全局最大空闲连接数；零值表示不修改。
	MaxIdleConns int
	// MaxIdleConnsPerHost 每个 host 最大空闲连接数；零值表示不修改。
	// 对「就几台逻辑服」很关键，Go 默认仅 2。
	MaxIdleConnsPerHost int
}

var (
	clientMu                 sync.RWMutex
	maxReceivedContentLength = defaultMaxReceivedContentLength
	connectTimeout           = time.Duration(defaultConnectTimeoutSec) * time.Second
	tcpKeepAlive             = defaultKeepAlive
	idleConnTimeout          = defaultIdleConnTimeout
	maxIdleConns             = defaultMaxIdleConns
	maxIdleConnsPerHost      = defaultMaxIdleConnsPerHost
	sharedClient             *http.Client
)

// SetMaxReceivedContentLength 设置单次响应体最大字节数，对应 Java HttpRequest.setMaxReceivedContentLength。
func SetMaxReceivedContentLength(n int) {
	if n <= 0 {
		return
	}
	clientMu.Lock()
	maxReceivedContentLength = n
	clientMu.Unlock()
}

// SetConnectTimeout 设置连接超时（秒），对应 Java HttpRequest.setConnectTimeout。
func SetConnectTimeout(seconds int) {
	if seconds <= 0 {
		return
	}
	clientMu.Lock()
	connectTimeout = time.Duration(seconds) * time.Second
	sharedClient = nil
	clientMu.Unlock()
}

// Configure 配置连接池与 TCP 保活；修改后下次请求会重建共享 Client。
// 网关→少量逻辑服长连接示例：
//
//	httpclient.Configure(httpclient.PoolConfig{
//		IdleConnTimeoutSet:  true,
//		IdleConnTimeout:     0,              // 不因空闲关闭
//		KeepAlive:           30 * time.Second,
//		MaxIdleConnsPerHost: 64,
//	})
func Configure(cfg PoolConfig) {
	clientMu.Lock()
	defer clientMu.Unlock()

	if cfg.ConnectTimeout > 0 {
		connectTimeout = cfg.ConnectTimeout
	}
	if cfg.KeepAlive > 0 {
		tcpKeepAlive = cfg.KeepAlive
	}
	if cfg.IdleConnTimeoutSet {
		idleConnTimeout = cfg.IdleConnTimeout
	}
	if cfg.MaxIdleConns > 0 {
		maxIdleConns = cfg.MaxIdleConns
	}
	if cfg.MaxIdleConnsPerHost > 0 {
		maxIdleConnsPerHost = cfg.MaxIdleConnsPerHost
	}
	sharedClient = nil
}

// EnableLongLivedConnections 快捷开启「尽量保持 TCP」：空闲不关连接，并提高每 host 空闲连接数。
func EnableLongLivedConnections(maxPerHost int) {
	if maxPerHost <= 0 {
		maxPerHost = defaultMaxIdleConnsPerHost
	}
	Configure(PoolConfig{
		IdleConnTimeoutSet:  true,
		IdleConnTimeout:     0,
		MaxIdleConnsPerHost: maxPerHost,
		MaxIdleConns:        maxPerHost * 4,
	})
}

// getClient 返回共享的 http.Client（带连接池）。
func getClient() *http.Client {
	clientMu.RLock()
	c := sharedClient
	clientMu.RUnlock()
	if c != nil {
		return c
	}

	clientMu.Lock()
	defer clientMu.Unlock()
	if sharedClient != nil {
		return sharedClient
	}

	timeout := connectTimeout
	keepAlive := tcpKeepAlive
	idle := idleConnTimeout
	maxIdle := maxIdleConns
	maxPerHost := maxIdleConnsPerHost

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: keepAlive,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          maxIdle,
		MaxIdleConnsPerHost:   maxPerHost,
		IdleConnTimeout:       idle, // 0 表示不因空闲关闭
		TLSHandshakeTimeout:   timeout,
		ExpectContinueTimeout: 1 * time.Second,
	}
	sharedClient = &http.Client{
		Transport: transport,
		Timeout:   0, // 单次超时由 Request 的 context / readTimeout 控制
	}
	return sharedClient
}

// maxBodySize 读取当前响应体上限。
func maxBodySize() int64 {
	clientMu.RLock()
	defer clientMu.RUnlock()
	return int64(maxReceivedContentLength)
}
