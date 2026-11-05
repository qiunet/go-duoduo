package utils

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// FileChangeCallback 文件变动时的回调，对应 Java IFileChangeCallback。
type FileChangeCallback func(path string)

type fileChangeObserver struct {
	path         string
	length       int64
	lastModified time.Time
	callbacks    []FileChangeCallback
	mu           sync.Mutex
}

// checkAndNotify 检查文件是否变化，有变化则通知回调。
func (o *fileChangeObserver) checkAndNotify() {
	info, err := os.Stat(o.path)
	if err != nil {
		return
	}
	mod := info.ModTime()
	len_ := info.Size()
	if o.lastModified != mod || o.length != len_ {
		o.mu.Lock()
		cbs := append([]FileChangeCallback(nil), o.callbacks...)
		o.mu.Unlock()
		for _, cb := range cbs {
			cb(o.path)
		}
	}
	o.lastModified = mod
	o.length = len_
}

// addCallback 追加变动回调。
func (o *fileChangeObserver) addCallback(cb FileChangeCallback) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.callbacks = append(o.callbacks, cb)
}

type fileChangeMonitor struct {
	mu        sync.Mutex
	observers map[string]*fileChangeObserver
}

// addObserver 为已存在的普通文件注册观察者。
func (m *fileChangeMonitor) addObserver(path string, cb FileChangeCallback) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("utils: change listener requires an existing file: %s", path)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	obs, ok := m.observers[path]
	if !ok {
		obs = &fileChangeObserver{
			path:         path,
			length:       info.Size(),
			lastModified: info.ModTime(),
		}
		m.observers[path] = obs
	}
	obs.addCallback(cb)
	return nil
}

// runOnce 轮询检查所有观察者。
func (m *fileChangeMonitor) runOnce() {
	m.mu.Lock()
	list := make([]*fileChangeObserver, 0, len(m.observers))
	for _, o := range m.observers {
		list = append(list, o)
	}
	m.mu.Unlock()
	for _, o := range list {
		o.checkAndNotify()
	}
}

var (
	changeMonitor     = &fileChangeMonitor{observers: map[string]*fileChangeObserver{}}
	changeMonitorOnce sync.Once
)

// startChangeMonitor 启动与 Java 类似的定时轮询（约每 15 秒）。
func startChangeMonitor() {
	changeMonitorOnce.Do(func() {
		go func() {
			t := time.NewTicker(15 * time.Second)
			defer t.Stop()
			for range t.C {
				changeMonitor.runOnce()
			}
		}()
	})
}

// ChangeListener 监听文件变动；文件内容或长度变化时调用 callback。
// 对应 Java FileUtil.changeListener。仅支持已存在的普通文件。
func ChangeListener(path string, callback FileChangeCallback) error {
	if callback == nil {
		return fmt.Errorf("utils: change callback is nil")
	}
	startChangeMonitor()
	return changeMonitor.addObserver(path, callback)
}
