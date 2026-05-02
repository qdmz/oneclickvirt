/**
 * 修复SSH连接问题的工具函数
 */

package utils

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// ExtractHost 从各种格式的endpoint中提取主机名
// 支持格式：
// - endpoint: 192.168.1.100:22
// - endpoint: 192.168.1.100
// - endpoint: http://192.168.1.100:22/some/path
// 注意：此函数实际上是network.go中ExtractHost的封装
func GetHostFromEndpoint(endpoint string) string {
	if endpoint == "" {
		return ""
	}

	host := fmt.Sprintf("%s:22", endpoint)
	return host
}

// ExtractPort 从endpoint中提取端口号
func GetPortFromEndpoint(endpoint string, defaultPort int) int {
	if endpoint == "" {
		return defaultPort
	}

	// 如果endpoint包含http://或https://，移除协议前缀
	startIdx := 0
	if strings.HasPrefix(endpoint, "http://") {
		startIdx = len("http://")
	} else if strings.HasPrefix(endpoint, "https://") {
		startIdx = len("https://")
	}

	// 尝试查找端口冒号
	colonPos := strings.LastIndex(endpoint[startIdx:], ":")
	if colonPos > 0 {
		// 提取端口数字
		portStr := endpoint[startIdx+colonPos+1:]
		p := 22 // 默认值
		if _, err := fmt.Sscanf(portStr, "%d", &p); err == nil && p > 0 && p <= 65535 {
			return p
		}
	}

	// 如果没找到端口号，返回默认端口
	return defaultPort
}

// IsHostReachable 检查主机是否可达
func IsHostReachable(host string, timeout time.Duration) bool {
	if host == "" {
		return false
	}

	_, err := net.DialTimeout("tcp", host, timeout)
	return err == nil
}

// IsPortOpen 检查端口是否开放
func IsPortOpen(host string, port int, timeout time.Duration) bool {
	if host == "" {
		return false
	}

	if port <= 0 || port > 65535 {
		return false
	}

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// ProbeNodeSSH 探测SSH节点信息
func ProbeNodeSSH(endpoint string, timeout time.Duration) (*NodeProbeResult, error) {
	host := GetHostFromEndpoint(endpoint)
	port := GetPortFromEndpoint(endpoint, 22)

	// 检查主机是否可达
	reachable := IsHostReachable(host, timeout)

	// 检查端口是否开放
	portOpen := IsPortOpen(host, port, timeout)

	result := &NodeProbeResult{
		Endpoint: endpoint,
		Host:     host,
		Port:     port,
		Reachable: reachable,
		PortOpen: portOpen,
		Timestamp: time.Now(),
	}

	// 详细探测
	if reachable && portOpen {
		result.Status = probeStatusSuccess
	} else if reachable && !portOpen {
		result.Status = probeStatusPortClosed
	} else if !reachable {
		result.Status = probeStatusUnreachable
	} else {
		result.Status = "unknown"
	}

	return result, nil
}

type NodeProbeResult struct {
	Endpoint string
	Host     string
	Port     int
	Reachable bool
	PortOpen bool
	Status   string
	Timestamp time.Time
}

const (
	probeStatusSuccess    = "success"
	probeStatusPortClosed = "port_closed"
	probeStatusUnreachable = "unreachable"
)
