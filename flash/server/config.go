package server

type BootstrapConfig struct {
	httpConfig       *HttpBootstrapConfig
	kcpConfig        *KcpBootstrapConfig
	tcpConfig        *TcpBootstrapConfig
	protocolHeader   *ProtocolHeader
	maxReceiveLength int
	serverName       string
	port             int
}

func NewBootstrapConfig(serverName string, port int) *BootstrapConfig {
	return &BootstrapConfig{
		serverName: serverName,
		port:       port,
	}
}

func (conf *BootstrapConfig) SetKcpBootstrapConfig(config *KcpBootstrapConfig) *BootstrapConfig {
	conf.kcpConfig = config
	return conf
}

func (conf *BootstrapConfig) SetTcpBootstrapConfig(config *TcpBootstrapConfig) *BootstrapConfig {
	conf.tcpConfig = config
	return conf
}

func (conf *BootstrapConfig) SetHttpBootstrapConfig(config *HttpBootstrapConfig) *BootstrapConfig {
	conf.httpConfig = config
	return conf
}

func (conf *BootstrapConfig) SetProtocolHeader(header *ProtocolHeader) *BootstrapConfig {
	conf.protocolHeader = header
	return conf
}

// DEFAULT_KCP_PARAM 默认的参数
var DEFAULT_KCP_PARAM = KcpStartupParam{512, 512, 20, true, 512, 2, true}

// KcpStartupParam kcp 启动的一些额外参数
type KcpStartupParam struct {
	sndWin     int
	rcvWin     int
	interval   int
	noDelay    bool
	mtu        int
	fastResend int
	noCwnd     bool
}

type KcpBootstrapConfig struct {
	param KcpStartupParam
	// 端口个数. 监听的udp需要多个端口
	portCount int
	// 是否跟主端口偏移
	portOffset int
}

type TcpBootstrapConfig struct {
	// kcp是依赖于Tcp协商的
	kcpConfig KcpBootstrapConfig
}

// HttpBootstrapConfig HTTP的启动参数
type HttpBootstrapConfig struct {
	// 游戏请求uri地址.
	gamePath string
	// 升级ws的uri地址.
	wsPath string
}
