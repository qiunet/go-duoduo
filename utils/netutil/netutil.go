package netutil

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"utils/commonutil"
	"utils/strutil"
)

const (
	// PortRangeMax 最大端口。对应 Java NetUtil.PORT_RANGE_MAX。
	PortRangeMax = 0xFFFF
	// PortRangeMin 常用最小非特权端口。对应 Java NetUtil.PORT_RANGE_MIN。
	PortRangeMin = 1024
)

// IsLocalIP 是否本机回环地址。对应 Java NetUtil.isLocalIp。
func IsLocalIP(ip string) bool {
	return commonutil.ExistInSlice(ip, "localhost", "0:0:0:0:0:0:0:1", "127.0.0.1", "::1")
}

// IsInnerIP 是否内网 / 回环 IP。对应 Java NetUtil.isInnerIp。
func IsInnerIP(ip string) bool {
	addr := net.ParseIP(ip)
	if addr == nil {
		return false
	}
	return addr.IsLoopback() || addr.IsPrivate() || addr.IsLinkLocalUnicast()
}

// IsValidPort 端口是否在 0~65535。对应 Java NetUtil.isValidPort。
func IsValidPort(port int) bool {
	return port >= 0 && port <= PortRangeMax
}

// IsValidIP4 是否合法 IPv4 字符串。使用 net/netip 解析（Go 语义，非严格对齐 Java 手写校验）。
func IsValidIP4(host string) bool {
	addr, err := netip.ParseAddr(host)
	return err == nil && addr.Is4()
}

// IsValidIP6 是否合法 IPv6（支持 zone，以及 /前缀）。使用 net/netip 解析。
func IsValidIP6(inet6Address string) bool {
	if p, err := netip.ParsePrefix(inet6Address); err == nil {
		return p.Addr().Is6()
	}
	addr, err := netip.ParseAddr(inet6Address)
	return err == nil && addr.Is6()
}

// GetLocalHostName 本机主机名。对应 Java NetUtil.getLocalHostName。
func GetLocalHostName() string {
	host, err := os.Hostname()
	if err != nil {
		return ""
	}
	return host
}

// LocalAddressList 获取本机地址列表（可过滤）。对应 Java NetUtil.localAddressList。
func LocalAddressList(filter func(net.IP) bool) []net.IP {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	out := make([]net.IP, 0)
	seen := map[string]struct{}{}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil {
				continue
			}
			if filter != nil && !filter(ip) {
				continue
			}
			key := ip.String()
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, ip)
		}
	}
	return out
}

// ToIPList 地址转字符串列表。对应 Java NetUtil.toIpList。
func ToIPList(addrs []net.IP) []string {
	out := make([]string, 0, len(addrs))
	for _, a := range addrs {
		out = append(out, a.String())
	}
	return out
}

// LocalIPv4s 本机 IPv4 列表。对应 Java NetUtil.localIpv4s。
func LocalIPv4s() []string {
	return ToIPList(LocalAddressList(func(ip net.IP) bool {
		return ip.To4() != nil
	}))
}

// LocalIPv6s 本机 IPv6 列表。对应 Java NetUtil.localIpv6s。
func LocalIPv6s() []string {
	return ToIPList(LocalAddressList(func(ip net.IP) bool {
		return ip.To4() == nil && ip.To16() != nil
	}))
}

// GetAllInnerIP 本机所有 IP。对应 Java NetUtil.getAllInnerIp。
func GetAllInnerIP() []string {
	return ToIPList(LocalAddressList(nil))
}

// GetInnerIP 得到一个非回环内网 IPv4。对应 Java NetUtil.getInnerIp。
func GetInnerIP() string {
	for _, ip := range LocalIPv4s() {
		if !IsLocalIP(ip) {
			return ip
		}
	}
	return ""
}

// GetPublicIP6 得到一个可用公网 IPv6。对应 Java NetUtil.getPublicIp6。
func GetPublicIP6() (string, error) {
	for _, s := range LocalIPv6s() {
		lower := strings.ToLower(s)
		if strings.HasPrefix(lower, "0:0:0:0:0:0:0:1") ||
			strings.HasPrefix(lower, "::1") ||
			strings.HasPrefix(lower, "fe80") ||
			strings.HasPrefix(lower, "fc00") ||
			strings.HasPrefix(lower, "fec0") {
			continue
		}
		ip := strutil.SplitString(s, "%")[0]
		return ip, nil
	}
	return "", fmt.Errorf("utils: not support ipv6")
}

// GetPublicIP4 获取公网 IPv4：先本地非内网，再并发请求公共查询服务。
// 对应 Java NetUtil.getPublicIp4。
func GetPublicIP4() (string, error) {
	addrs := LocalAddressList(func(ip net.IP) bool {
		v4 := ip.To4()
		return v4 != nil && !IsInnerIP(v4.String())
	})
	if len(addrs) > 0 {
		return addrs[0].String(), nil
	}

	services := []string{
		"https://checkip.amazonaws.com",
		"https://ipv4.icanhazip.com",
		"https://icanhazip.com",
		"https://ipinfo.io/ip",
	}
	type result struct {
		ip  string
		err error
	}
	ch := make(chan result, len(services))
	var wg sync.WaitGroup
	client := &http.Client{Timeout: 3 * time.Second}
	for _, u := range services {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := client.Get(url)
			if err != nil {
				ch <- result{err: err}
				return
			}
			defer resp.Body.Close()
			line, err := bufio.NewReader(resp.Body).ReadString('\n')
			if err != nil && len(strings.TrimSpace(line)) == 0 {
				ch <- result{err: err}
				return
			}
			ip := strings.TrimSpace(line)
			if !IsValidIP4(ip) {
				ch <- result{err: fmt.Errorf("ip invalid: %s", ip)}
				return
			}
			ch <- result{ip: ip}
		}(u)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for r := range ch {
		if r.err == nil && r.ip != "" {
			return r.ip, nil
		}
	}
	return "", fmt.Errorf("utils: no public ip get")
}

// SocketSendData 使用 TCP 发送数据。对应 Java NetUtil.socketSendData。
func SocketSendData(host string, port int, data []byte) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(data)
	return err
}

// UDPSendData 使用 UDP 发送数据。对应 Java NetUtil.udpSendData。
func UDPSendData(host string, port int, data []byte) error {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(data)
	return err
}
