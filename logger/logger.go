package logger

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type LoggerType string

const (
	LoggerTypeInfo    LoggerType = "info"
	LoggerTypeError   LoggerType = "error"
	LoggerTypeWarning LoggerType = "warning"
	LoggerTypeDebug   LoggerType = "debug"
	LoggerTypeTrace   LoggerType = "trace"
)

type LoggerName string

const (
	Duoduo LoggerName = "duoduo"
	Test   LoggerName = "test"
	SQL    LoggerName = "sql"
	Redis  LoggerName = "redis"
)

var (
	base  *logrus.Logger
	mu    sync.Mutex
	cache = map[string]*logrus.Entry{}
)

func init() {
	base = logrus.New()
	base.SetOutput(os.Stdout)
	base.SetLevel(logrus.InfoLevel)
	base.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// GetLogger 按名字获取 logger，用法类似 Java：
//
//	log := logger.GetLogger("MyService")
//	log.Info("hello")
//	log.Error("oops")
func GetLogger(name string) *logrus.Entry {
	mu.Lock()
	defer mu.Unlock()
	if e, ok := cache[name]; ok {
		return e
	}
	e := base.WithField("logger", name)
	cache[name] = e
	return e
}
