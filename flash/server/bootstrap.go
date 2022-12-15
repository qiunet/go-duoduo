package server

import "flag"

// Hook 钩子参数
type Hook struct {
	ShutdownMsg, ReloadMsg, DeprecateMsg string
	Custom                               func(command string)
	Args                                 func(args *CommandArgs)
}

type CommandArgs struct {
	// 自定义的参数
	data map[string]any
}

type Bootstrap struct {
	configs []*BootstrapConfig
	args    *CommandArgs
	hook    *Hook
}

// NewBootstrap 新建一个 Bootstrap
func NewBootstrap(hook *Hook) *Bootstrap {
	b := new(Bootstrap)
	b.hook = hook
	parseArgs(b)

	return b
}

// 检查参数
func parseArgs(b *Bootstrap) {
	// 自定义的参数解析
	b.hook.Args(b.args)
	// 其它
	flag.Parse()
}

// Listener 监听
func (b *Bootstrap) Listener(config *BootstrapConfig) *Bootstrap {
	b.configs = append(b.configs, config)
	return b
}

// Start 启动服务
func (b *Bootstrap) Start() {

}
