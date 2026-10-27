package admin

// console_proxmox_api.go wires the PVE HTTP-API VNC path (see
// proxmox_vnc_api.go) into the provider-neutral console system. It is used
// when a Proxmox provider's connection type is "api": the controller asks
// the PVE API for a websocket vncproxy ticket, dials vncwebsocket, and
// adapts that WebSocket to the net.Conn the existing RFB machinery expects.
// The ticket never leaves server memory.

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	providerModel "oneclickvirt/model/provider"
	"oneclickvirt/utils"

	"github.com/gorilla/websocket"
)

// proxmoxAPIHost resolves the PVE API host from a provider record.
func proxmoxAPIHost(p providerModel.Provider) string {
	if host := utils.ExtractHost(p.Endpoint); host != "" {
		return host
	}
	return utils.ExtractHost(p.PortIP)
}

// proxmoxVNCWebSocketDial is indirected for tests.
var proxmoxVNCWebSocketDial = func(wsURL, apiToken string) (net.Conn, error) {
	dialer := &websocket.Dialer{
		HandshakeTimeout: proxmoxVNCSetupTimeout,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true}, // #nosec G402 -- PVE 节点自签名证书
	}
	header := http.Header{}
	if token := cleanProxmoxAPIToken(apiToken); token != "" {
		header.Set("Authorization", "PVEAPIToken="+token)
	}
	ws, _, err := dialer.Dial(wsURL, header)
	if err != nil {
		return nil, err
	}
	return &wsNetConn{ws: ws}, nil
}

// openProxmoxVNCConnViaAPI implements the "api" transport for
// openProxmoxVNCConn: vncproxy ticket request + vncwebsocket dial.
func openProxmoxVNCConnViaAPI(target consoleTarget) (net.Conn, func(), string, error) {
	return openProxmoxVNCConnViaAPIWithClient(target, nil)
}

func openProxmoxVNCConnViaAPIWithClient(target consoleTarget, client *http.Client) (net.Conn, func(), string, error) {
	p := target.provider
	host := proxmoxAPIHost(p)
	if host == "" {
		return nil, nil, "", fmt.Errorf("PVE API 控制台缺少节点 API 地址")
	}
	resource, err := proxmoxAPIResource(target.runtimeType)
	if err != nil {
		return nil, nil, "", err
	}
	node := strings.TrimSpace(target.runtimeNode)
	if node == "" {
		node = strings.TrimSpace(p.HostName)
	}
	if node == "" {
		node = "pve"
	}
	ctx, cancel := context.WithTimeout(context.Background(), proxmoxVNCSetupTimeout)
	defer cancel()
	result, err := requestProxmoxVNCProxy(ctx, client, host, node, resource, p.Token, target.runtimeID, true)
	if err != nil {
		return nil, nil, "", err
	}
	wsURL := proxmoxVNCWebSocketURL(host, node, resource, target.runtimeID, result.Port, result.Ticket)
	if wsURL == "" {
		return nil, nil, "", fmt.Errorf("PVE VNC WebSocket 地址拼装失败")
	}
	conn, err := proxmoxVNCWebSocketDial(wsURL, p.Token)
	if err != nil {
		return nil, nil, "", fmt.Errorf("连接 PVE VNC WebSocket 失败: %w", err)
	}
	return conn, func() { _ = conn.Close() }, result.PasswordOrTicket(), nil
}

// lookupProxmoxConsoleRuntimeViaAPI resolves the live qemu/lxc resource for
// an instance through the PVE HTTP API, for providers without a usable
// shell executor. It reuses the same resource parsing/matching as the
// pvesh-based lookup so both paths agree on identity.
func lookupProxmoxConsoleRuntimeViaAPI(inst providerModel.Instance, p providerModel.Provider) (proxmoxConsoleRuntime, error) {
	return lookupProxmoxConsoleRuntimeViaAPIWithClient(inst, p, nil)
}

func lookupProxmoxConsoleRuntimeViaAPIWithClient(inst providerModel.Instance, p providerModel.Provider, client *http.Client) (proxmoxConsoleRuntime, error) {
	host := proxmoxAPIHost(p)
	if host == "" {
		return proxmoxConsoleRuntime{}, fmt.Errorf("PVE API 控制台缺少节点 API 地址")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	body, err := proxmoxAPIGet(ctx, client, host, p.Token, "/api2/json/cluster/resources?type=vm")
	if err == nil {
		var resources []proxmoxConsoleResource
		if resources, err = parseProxmoxConsoleResources(string(body)); err == nil {
			if runtime, findErr := findProxmoxConsoleRuntime(resources, inst, p.HostName); findErr == nil {
				return runtime, nil
			} else {
				err = findErr
			}
		}
	}
	// Cluster inventory can be restricted on delegated API tokens; fall back
	// to an exact per-node status lookup with the stored numeric VMID.
	if direct, directErr := lookupProxmoxConsoleRuntimeViaAPIByID(ctx, client, host, inst, p); directErr == nil {
		return direct, nil
	} else {
		return proxmoxConsoleRuntime{}, fmt.Errorf("读取 PVE 实例运行态失败(API): %v；单资源回退失败: %v", err, directErr)
	}
}

func lookupProxmoxConsoleRuntimeViaAPIByID(ctx context.Context, client *http.Client, host string, inst providerModel.Instance, p providerModel.Provider) (proxmoxConsoleRuntime, error) {
	runtimeID := strings.TrimSpace(inst.ProviderVMID)
	if !isProxmoxConsoleRuntimeID(runtimeID) {
		return proxmoxConsoleRuntime{}, fmt.Errorf("实例未保存有效 VMID/CTID，无法在无清单权限时安全恢复")
	}
	nodesBody, err := proxmoxAPIGet(ctx, client, host, p.Token, "/api2/json/nodes")
	if err != nil {
		return proxmoxConsoleRuntime{}, err
	}
	var nodesEnvelope struct {
		Data []struct {
			Node string `json:"node"`
		} `json:"data"`
	}
	if err := json.Unmarshal(nodesBody, &nodesEnvelope); err != nil {
		return proxmoxConsoleRuntime{}, fmt.Errorf("解析 PVE 节点清单失败: %w", err)
	}
	nodes := make([]string, 0, len(nodesEnvelope.Data)+1)
	if configured := strings.TrimSpace(p.HostName); configured != "" {
		nodes = append(nodes, configured)
	}
	for _, n := range nodesEnvelope.Data {
		if node := strings.TrimSpace(n.Node); node != "" {
			duplicate := false
			for _, existing := range nodes {
				if existing == node {
					duplicate = true
					break
				}
			}
			if !duplicate {
				nodes = append(nodes, node)
			}
		}
	}
	matches := make([]proxmoxConsoleRuntime, 0, 1)
	for _, node := range nodes {
		for _, resource := range []string{"qemu", "lxc"} {
			path := fmt.Sprintf("/api2/json/nodes/%s/%s/%s/status/current", node, resource, runtimeID)
			statusBody, statusErr := proxmoxAPIGet(ctx, client, host, p.Token, path)
			if statusErr != nil {
				continue
			}
			var statusEnvelope struct {
				Data struct {
					Status string `json:"status"`
				} `json:"data"`
			}
			if json.Unmarshal(statusBody, &statusEnvelope) != nil || statusEnvelope.Data.Status == "" {
				continue
			}
			matches = append(matches, proxmoxConsoleRuntime{ID: runtimeID, Type: resource, Node: node, Status: statusEnvelope.Data.Status})
		}
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) > 1 {
		return proxmoxConsoleRuntime{}, fmt.Errorf("PVE VMID/CTID %q 在 API 单资源探测中匹配多个对象，拒绝猜测", runtimeID)
	}
	return proxmoxConsoleRuntime{}, fmt.Errorf("PVE 中未找到 VMID/CTID %q 的 qemu 或 lxc 运行时对象", runtimeID)
}

// wsNetConn adapts a gorilla WebSocket carrying a raw RFB byte stream to
// net.Conn. Reads consume message payloads in order; writes send one binary
// message per call, which matches how the RFB proxy pumps bytes.
type wsNetConn struct {
	ws     *websocket.Conn
	reader io.Reader
	mu     sync.Mutex
}

func (c *wsNetConn) Read(b []byte) (int, error) {
	for {
		if c.reader != nil {
			n, err := c.reader.Read(b)
			if err == io.EOF {
				c.reader = nil
				if n > 0 {
					return n, nil
				}
				continue
			}
			return n, err
		}
		_, reader, err := c.ws.NextReader()
		if err != nil {
			return 0, err
		}
		c.reader = reader
	}
}

func (c *wsNetConn) Write(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.ws.WriteMessage(websocket.BinaryMessage, b); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *wsNetConn) Close() error                       { return c.ws.Close() }
func (c *wsNetConn) LocalAddr() net.Addr                { return c.ws.LocalAddr() }
func (c *wsNetConn) RemoteAddr() net.Addr               { return c.ws.RemoteAddr() }
func (c *wsNetConn) SetDeadline(t time.Time) error      { _ = c.ws.SetReadDeadline(t); return c.ws.SetWriteDeadline(t) }
func (c *wsNetConn) SetReadDeadline(t time.Time) error  { return c.ws.SetReadDeadline(t) }
func (c *wsNetConn) SetWriteDeadline(t time.Time) error { return c.ws.SetWriteDeadline(t) }
