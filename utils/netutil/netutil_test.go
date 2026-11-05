package netutil

import (
	"net"
	"testing"
)

func TestNetUtil_IPAndPort(t *testing.T) {
	t.Run("本机与端口", func(t *testing.T) {
		if !IsLocalIP("127.0.0.1") || !IsLocalIP("localhost") {
			t.Fatal("local")
		}
		if !IsValidPort(80) || IsValidPort(70000) {
			t.Fatal("port")
		}
	})
	t.Run("内网判断", func(t *testing.T) {
		// 对齐 Java TestNetUtil.testInnerIp
		if !IsInnerIP("192.168.1.200") {
			t.Fatal("192.168")
		}
		if !IsInnerIP("172.21.0.9") {
			t.Fatal("172.21")
		}
		if IsInnerIP("123.196.125.13") {
			t.Fatal("public")
		}
		if !IsInnerIP("10.154.197.234") {
			t.Fatal("10.x")
		}
		if !IsInnerIP("127.0.0.1") || !IsLocalIP("127.0.0.1") {
			t.Fatal("loopback")
		}
		if IsInnerIP("8.8.8.8") {
			t.Fatal("8.8.8.8")
		}
	})
	t.Run("IPv4校验", func(t *testing.T) {
		// 对齐 Java TestNetUtil.testIsValidIp
		if !IsValidIP4("192.168.1.255") || !IsValidIP4("172.21.0.9") || !IsValidIP4("123.196.125.13") {
			t.Fatal("valid ipv4")
		}
		if IsValidIP4("123.196.125.344") || IsValidIP4("https://baidu.com") {
			t.Fatal("invalid ipv4")
		}
		if IsValidIP4("1.2.3") || IsValidIP4("::1") {
			t.Fatal("reject non-v4")
		}
	})
	t.Run("IPv6校验", func(t *testing.T) {
		// 对齐 Java TestNetUtil.testIsValidIp
		if !IsValidIP6("::1") {
			t.Fatal("::1")
		}
		if !IsValidIP6("2408:4005:3e2:6d00:ca74:160b:7749:c6b8") {
			t.Fatal("full ipv6")
		}
		if IsValidIP6("2408:4005:3g2:6d00:ca74:160b:7749:c6b8") {
			t.Fatal("invalid hex")
		}
		if !IsValidIP6("2001:db8::/32") {
			t.Fatal("prefix")
		}
		if IsValidIP6("192.168.1.1") {
			t.Fatal("reject ipv4")
		}
	})
}

func TestNetUtil_LocalList(t *testing.T) {
	t.Run("本机地址列表可调用", func(t *testing.T) {
		// 对齐 Java testGetInnerIp / testGetLocalHostName / testAllInnerIp：可调用即可
		_ = LocalIPv4s()
		_ = LocalIPv6s()
		_ = GetAllInnerIP()
		_ = GetLocalHostName()
		_ = GetInnerIP()
		_ = ToIPList([]net.IP{net.ParseIP("127.0.0.1")})
	})
}
