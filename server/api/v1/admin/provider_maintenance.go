package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"oneclickvirt/constant"
	"oneclickvirt/global"
	providerModel "oneclickvirt/model/provider"
	"oneclickvirt/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// vncSession 缓存 Proxmox vncproxy 得到的临时 VNC 连接信息。
// Proxmox 的 VNC 以 unix socket 形式存在，必须通过 vncproxy API 获取一个
// 带有时效性的 TCP 端口与一次性密码；该会话需在“获取连接信息”与“建立 WebSocket”
// 两次请求之间保持一致，因此缓存在内存中。
type vncSession struct {
	Host      string
	Port      int
	Password  string
	Ticket    string
	User      string
	ExpiresAt time.Time
}

var (
	vncSessionMu sync.Mutex
	vncSessions  = make(map[uint]*vncSession)
)

const vncSessionTTL = 120 * time.Second

// buildInstanceVNCInfo 返回实例 WebVNC 是否可用；对于 Proxmox 虚拟机，
// 会同时返回本次 vncproxy 会话的 VNC 密码，供前端 noVNC 完成 RFB 认证。
func buildInstanceVNCInfo(instanceID uint, userID uint, admin bool) (gin.H, error) {
	host, port, password, err := resolveInstanceVNCProxyTarget(instanceID, userID, admin)
	if err != nil {
		return gin.H{"enabled": false, "reason": err.Error()}, nil
	}
	_ = host
	_ = port
	return gin.H{"enabled": true, "password": password}, nil
}

// resolveInstanceVNCProxyTarget 解析 VNC 代理目标。
// 对于 Proxmox 虚拟机，返回 127.0.0.1 + vncproxy 临时端口 + 一次性密码；
// 对于其他（或非 Proxmox）节点，保留原有 TCP 直连解析逻辑。
func resolveInstanceVNCProxyTarget(instanceID uint, userID uint, admin bool) (string, int, string, error) {
	var inst providerModel.Instance
	query := global.APP_DB.Select("id", "provider_id", "status", "instance_type", "provider_vm_id", "discovered_data", "user_id")
	if admin {
		query = query.Where("id = ?", instanceID)
	} else {
		query = query.Where("id = ? AND user_id = ?", instanceID, userID)
	}
	if err := query.First(&inst).Error; err != nil {
		return "", 0, "", err
	}
	if constant.IsBusyStatus(inst.Status) {
		return "", 0, "", fmt.Errorf("实例正在操作进行中（当前状态：%s），请等待当前任务完成", inst.Status)
	}
	if inst.Status != constant.InstanceStatusRunning {
		return "", 0, "", fmt.Errorf("实例未运行")
	}
	var p providerModel.Provider
	if err := global.APP_DB.Select("id", "type", "endpoint", "port_ip", "enable_vnc", "vnc_base_port", "vnc_host", "host_name", "token").
		First(&p, inst.ProviderID).Error; err != nil {
		return "", 0, "", err
	}
	if !p.EnableVNC {
		return "", 0, "", fmt.Errorf("节点未启用WebVNC")
	}

	// Proxmox 走 vncproxy API 会话：qemu 与 lxc 都支持（P1-5 补齐 lxc）。
	if isProxmoxConsoleProviderType(p.Type) {
		if _, resourceErr := proxmoxAPIResource(inst.InstanceType); resourceErr != nil {
			return "", 0, "", fmt.Errorf("当前实例类型不支持WebVNC")
		}
		sess, err := getOrCreateProxmoxVNCSession(&inst, &p)
		if err != nil {
			return "", 0, "", err
		}
		return sess.Host, sess.Port, sess.Password, nil
	}

	if inst.InstanceType != "vm" {
		return "", 0, "", fmt.Errorf("当前实例类型不支持WebVNC")
	}

	// 非 Proxmox 节点：保留原有 TCP 直连解析逻辑
	host := strings.TrimSpace(p.VNCHost)
	if host == "" {
		host = strings.TrimSpace(p.PortIP)
		if host == "" {
			host = strings.TrimSpace(p.Endpoint)
		}
	}
	host = strings.Trim(host, "[]")
	port := parseVNCDiscoveredPort(inst.DiscoveredData)
	if port == 0 {
		base := p.VNCBasePort
		if base == 0 {
			base = 5900
		}
		vmid, _ := strconv.Atoi(inst.ProviderVMID)
		if vmid > 0 {
			port = base + vmid
		} else {
			port = base
		}
	}
	if host == "" || port <= 0 || port > 65535 {
		return "", 0, "", fmt.Errorf("VNC目标不可用")
	}
	return host, port, "", nil
}

func parseVNCDiscoveredPort(raw string) int {
	if raw == "" {
		return 0
	}
	var obj map[string]interface{}
	if json.Unmarshal([]byte(raw), &obj) != nil {
		return 0
	}
	var parse func(value interface{}) int
	parse = func(value interface{}) int {
		switch x := value.(type) {
		case float64:
			return int(x)
		case json.Number:
			n, _ := strconv.Atoi(string(x))
			return n
		case int:
			return x
		case string:
			n, _ := strconv.Atoi(strings.TrimSpace(x))
			return n
		case map[string]interface{}:
			for _, key := range []string{"port", "vncPort", "vnc_port"} {
				if n := parse(x[key]); n > 0 {
					return n
				}
			}
		}
		return 0
	}
	for _, key := range []string{"vncPort", "vnc_port", "vnc"} {
		if n := parse(obj[key]); n > 0 {
			return n
		}
	}
	console, _ := obj["console"].(map[string]interface{})
	if console == nil {
		return 0
	}
	for _, key := range []string{"vncPort", "vnc_port", "vnc"} {
		if n := parse(console[key]); n > 0 {
			return n
		}
	}
	if strings.EqualFold(strings.TrimSpace(fmt.Sprint(console["protocol"])), consoleProtocolVNC) {
		return parse(console["port"])
	}
	return 0
}

// getOrCreateProxmoxVNCSession 调用 Proxmox API 的 vncproxy 获取临时 VNC 端口与密码。
// 结果按实例缓存 vncSessionTTL，确保“获取连接信息”与“建立 WebSocket 代理”两次
// 请求使用同一会话（同一端口与密码）。
// cachedVNCSession returns the live cached session for an instance, if any.
func cachedVNCSession(instanceID uint, now time.Time) (*vncSession, bool) {
	vncSessionMu.Lock()
	defer vncSessionMu.Unlock()
	s, ok := vncSessions[instanceID]
	if !ok || s == nil || !now.Before(s.ExpiresAt) {
		return nil, false
	}
	return s, true
}

func getOrCreateProxmoxVNCSession(inst *providerModel.Instance, p *providerModel.Provider) (*vncSession, error) {
	if s, ok := cachedVNCSession(inst.ID, time.Now()); ok {
		return s, nil
	}

	node := strings.TrimSpace(p.HostName)
	if node == "" {
		node = "pve"
	}
	host := proxmoxAPIHost(*p)
	if host == "" {
		return nil, fmt.Errorf("节点缺少 API 地址")
	}
	vmid := strings.TrimSpace(inst.ProviderVMID)
	if vmid == "" {
		return nil, fmt.Errorf("实例缺少 VMID")
	}
	if cleanProxmoxAPIToken(p.Token) == "" {
		return nil, fmt.Errorf("节点未配置 Proxmox API Token")
	}
	resource, err := proxmoxAPIResource(inst.InstanceType)
	if err != nil {
		return nil, err
	}

	result, err := requestProxmoxVNCProxy(context.Background(), nil, host, node, resource, p.Token, vmid, false)
	if err != nil {
		return nil, err
	}

	s := &vncSession{
		// The vncproxy TCP port listens on the PVE node itself, so the dial
		// target is the node's API host (the previous hardcoded 127.0.0.1
		// only worked when the controller shared the node).
		Host:      host,
		Port:      result.Port,
		Password:  result.PasswordOrTicket(),
		Ticket:    result.Ticket,
		User:      result.User,
		ExpiresAt: time.Now().Add(vncSessionTTL),
	}
	vncSessionMu.Lock()
	vncSessions[inst.ID] = s
	vncSessionMu.Unlock()
	return s, nil
}

// proxmoxVNCProxy 调用 Proxmox API 打开一个 VNC 代理会话，返回 TCP 端口与一次性 VNC 密码。
// 保留此签名供既有调用方使用；实现已收敛到 requestProxmoxVNCProxy（qemu、非 websocket 模式）。
func proxmoxVNCProxy(endpoint, node, apiToken, vmid string) (int, string, error) {
	result, err := requestProxmoxVNCProxy(context.Background(), nil, utils.ExtractHost(endpoint), node, "qemu", apiToken, vmid, false)
	if err != nil {
		return 0, "", err
	}
	return result.Port, result.PasswordOrTicket(), nil
}

var vncUpgrader = websocket.Upgrader{
	ReadBufferSize:  32768,
	WriteBufferSize: 32768,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		appConfig := global.GetAppConfig()
		return utils.OriginAllowedForRequest(r, origin, appConfig.System.FrontendURL, appConfig.Cors.Whitelist)
	},
}
func proxyVNCWebSocket(c *gin.Context, host string, port int) {
	ws, err := vncUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), 10*time.Second)
	if err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("VNC连接失败: "+err.Error()))
		return
	}
	defer conn.Close()
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer cancel()
		for {
			mt, msg, err := ws.ReadMessage()
			if err != nil {
				return
			}
			if mt == websocket.BinaryMessage || mt == websocket.TextMessage {
				_, _ = conn.Write(msg)
			}
		}
	}()
	go func() {
		defer wg.Done()
		defer cancel()
		buf := make([]byte, 32768)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					return
				}
			}
			if err != nil {
				return
			}
		}
	}()
	<-ctx.Done()
	if global.APP_LOG != nil {
		global.APP_LOG.Debug("WebVNC会话结束", zap.String("host", host), zap.Int("port", port))
	}
}
