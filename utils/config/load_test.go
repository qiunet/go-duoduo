package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

type dbConfig struct {
	DSN     string `yaml:"dsn" env:"TEST_DB_DSN" default:"localhost"`
	MaxOpen int    `yaml:"maxOpen" default:"10"`
}

type appConfig struct {
	Port  int           `yaml:"port" env:"TEST_APP_PORT" default:"8080"`
	Host  string        `yaml:"host" env:"TEST_APP_HOST" default:"0.0.0.0"`
	Debug bool          `yaml:"debug" env:"TEST_APP_DEBUG" default:"false"`
	Wait  time.Duration `yaml:"wait" default:"1s"`
	DB    dbConfig      `yaml:"db"`
}

type flatPathConfig struct {
	C   string `yaml:"a.b.c" env:"TEST_ABC" default:"def"`
	Num int    `yaml:"a.b.n" default:"1"`
}

func TestLoadBytes_DefaultsOnly(t *testing.T) {
	t.Run("空YAML用默认值", func(t *testing.T) {
		var cfg appConfig
		if err := LoadBytes(nil, &cfg); err != nil {
			t.Fatal(err)
		}
		if cfg.Port != 8080 || cfg.Host != "0.0.0.0" || cfg.Debug || cfg.Wait != time.Second {
			t.Fatalf("%+v", cfg)
		}
		if cfg.DB.DSN != "localhost" || cfg.DB.MaxOpen != 10 {
			t.Fatalf("db=%+v", cfg.DB)
		}
	})
}

func TestLoadBytes_YAMLOverridesDefault(t *testing.T) {
	t.Run("YAML覆盖默认值", func(t *testing.T) {
		raw := []byte(`
port: 9090
host: 127.0.0.1
debug: true
wait: 2s
db:
  dsn: mysql://x
  maxOpen: 20
`)
		var cfg appConfig
		if err := LoadBytes(raw, &cfg); err != nil {
			t.Fatal(err)
		}
		if cfg.Port != 9090 || cfg.Host != "127.0.0.1" || !cfg.Debug || cfg.Wait != 2*time.Second {
			t.Fatalf("%+v", cfg)
		}
		if cfg.DB.DSN != "mysql://x" || cfg.DB.MaxOpen != 20 {
			t.Fatalf("db=%+v", cfg.DB)
		}
	})
}

func TestLoadBytes_EnvOverridesYAML(t *testing.T) {
	t.Run("ENV优先于YAML", func(t *testing.T) {
		t.Setenv("TEST_APP_PORT", "7000")
		t.Setenv("TEST_APP_HOST", "env-host")
		t.Setenv("TEST_DB_DSN", "env-dsn")
		raw := []byte(`
port: 9090
host: 127.0.0.1
db:
  dsn: file-dsn
`)
		var cfg appConfig
		if err := LoadBytes(raw, &cfg); err != nil {
			t.Fatal(err)
		}
		if cfg.Port != 7000 || cfg.Host != "env-host" || cfg.DB.DSN != "env-dsn" {
			t.Fatalf("%+v", cfg)
		}
		if cfg.DB.MaxOpen != 10 {
			t.Fatalf("maxOpen default lost: %d", cfg.DB.MaxOpen)
		}
	})
}

func TestLoadBytes_DottedPath(t *testing.T) {
	t.Run("点路径a.b.c绑定", func(t *testing.T) {
		raw := []byte(`
a:
  b:
    c: hello
    n: 42
`)
		var cfg flatPathConfig
		if err := LoadBytes(raw, &cfg); err != nil {
			t.Fatal(err)
		}
		if cfg.C != "hello" || cfg.Num != 42 {
			t.Fatalf("%+v", cfg)
		}
	})
	t.Run("点路径ENV覆盖", func(t *testing.T) {
		t.Setenv("TEST_ABC", "from-env")
		raw := []byte(`
a:
  b:
    c: hello
`)
		var cfg flatPathConfig
		if err := LoadBytes(raw, &cfg); err != nil {
			t.Fatal(err)
		}
		if cfg.C != "from-env" {
			t.Fatalf("%+v", cfg)
		}
	})
	t.Run("点路径缺省用default", func(t *testing.T) {
		var cfg flatPathConfig
		if err := LoadBytes([]byte("x: 1"), &cfg); err != nil {
			t.Fatal(err)
		}
		if cfg.C != "def" || cfg.Num != 1 {
			t.Fatalf("%+v", cfg)
		}
	})
}

func TestLoad_File(t *testing.T) {
	t.Run("从文件加载", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "app.yaml")
		content := "port: 3000\nhost: file-host\n"
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		var cfg appConfig
		if err := Load(path, &cfg); err != nil {
			t.Fatal(err)
		}
		if cfg.Port != 3000 || cfg.Host != "file-host" {
			t.Fatalf("%+v", cfg)
		}
	})
	t.Run("文件不存在报错", func(t *testing.T) {
		var cfg appConfig
		if err := Load(filepath.Join(t.TempDir(), "missing.yaml"), &cfg); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestLoadBytes_InvalidDst(t *testing.T) {
	t.Run("非法dst", func(t *testing.T) {
		if err := LoadBytes(nil, nil); err == nil {
			t.Fatal("nil")
		}
		var i int
		if err := LoadBytes(nil, &i); err == nil {
			t.Fatal("non-struct")
		}
		var cfg appConfig
		if err := LoadBytes(nil, cfg); err == nil {
			t.Fatal("non-pointer")
		}
	})
}

func TestLoadBytes_SkipYAMLDash(t *testing.T) {
	t.Run("yaml横杠跳过", func(t *testing.T) {
		type s struct {
			A string `yaml:"-" default:"keep"`
			B string `yaml:"b" default:"b0"`
		}
		var cfg s
		if err := LoadBytes([]byte("a: x\nb: y"), &cfg); err != nil {
			t.Fatal(err)
		}
		if cfg.A != "keep" || cfg.B != "y" {
			t.Fatalf("%+v", cfg)
		}
	})
}
