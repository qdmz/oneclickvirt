package admin

// proxmox_vnc_api.go provides the HTTP-API flavour of the Proxmox VNC flow
// (P1-5). Providers that are managed purely through the PVE API (API token,
// no SSH/Agent/local executor on the node) cannot use the pvesh-based path
// in console_proxmox_vnc.go. The helpers here talk to the PVE HTTP API
// directly:
//
//	POST /api2/json/nodes/{node}/{qemu|lxc}/{vmid}/vncproxy  (websocket=1)
//	→ {data:{ticket, port, user, cert}} and then
//	GET  /api2/json/nodes/{node}/{qemu|lxc}/{vmid}/vncwebsocket?port=&vncticket=
//
// The ticket doubles as the RFB password; it is short-lived, stays in server
// memory only, and is never sent to the browser (the controller terminates
// the RFB authentication itself, see proxyProxmoxVNCWebSocket).

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"oneclickvirt/utils"
)

// proxmoxAPIPort is the PVE API port; a variable so tests can point the
// client at a fake PVE server.
var proxmoxAPIPort = 8006

// proxmoxVNCAPIResult is the parsed payload of a PVE vncproxy call.
type proxmoxVNCAPIResult struct {
	Port     int
	Ticket   string
	Password string
	User     string
}

// PasswordOrTicket returns the credential used for the RFB handshake: the
// classic vncproxy password, or the websocket-mode ticket.
func (r proxmoxVNCAPIResult) PasswordOrTicket() string {
	if r.Password != "" {
		return r.Password
	}
	return r.Ticket
}

// cleanProxmoxAPIToken strips whitespace/control characters that can sneak
// into the persisted provider token (mirrors provider/proxmox behaviour).
func cleanProxmoxAPIToken(token string) string {
	token = strings.TrimSpace(token)
	token = strings.ReplaceAll(token, "\n", "")
	token = strings.ReplaceAll(token, "\r", "")
	return token
}

// proxmoxAPIResource maps a runtime/instance type to the PVE API resource
// collection name ("qemu" or "lxc").
func proxmoxAPIResource(runtimeType string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(runtimeType)) {
	case "qemu", "vm":
		return "qemu", nil
	case "lxc", "container":
		return "lxc", nil
	case "":
		return "qemu", nil
	default:
		return "", fmt.Errorf("不支持的 PVE 运行时类型: %s", runtimeType)
	}
}

// proxmoxVNCProxyAPIURL builds the vncproxy endpoint URL for a guest.
func proxmoxVNCProxyAPIURL(host, node, resource, vmid string) string {
	path := fmt.Sprintf("/api2/json/nodes/%s/%s/%s/vncproxy",
		url.PathEscape(strings.TrimSpace(node)), resource, url.PathEscape(strings.TrimSpace(vmid)))
	return utils.BuildEndpointURL("https", host, proxmoxAPIPort, path)
}

// proxmoxVNCWebSocketURL assembles the PVE vncwebsocket URL for a ticket
// issued by a websocket=1 vncproxy call.
func proxmoxVNCWebSocketURL(host, node, resource, vmid string, port int, ticket string) string {
	path := fmt.Sprintf("/api2/json/nodes/%s/%s/%s/vncwebsocket",
		url.PathEscape(strings.TrimSpace(node)), resource, url.PathEscape(strings.TrimSpace(vmid)))
	base := utils.BuildEndpointURL("wss", host, proxmoxAPIPort, path)
	if base == "" {
		return ""
	}
	query := url.Values{}
	query.Set("port", strconv.Itoa(port))
	query.Set("vncticket", ticket)
	return base + "?" + query.Encode()
}

// parseProxmoxVNCAPIResponse parses a vncproxy API response body. PVE returns
// port as a JSON number in websocket mode and as a string in some wrappers,
// so both shapes are accepted (same convention as the pvesh parser).
func parseProxmoxVNCAPIResponse(body []byte) (proxmoxVNCAPIResult, error) {
	var envelope struct {
		Data struct {
			Port     json.RawMessage `json:"port"`
			Ticket   string          `json:"ticket"`
			Password string          `json:"password"`
			User     string          `json:"user"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return proxmoxVNCAPIResult{}, fmt.Errorf("解析 Proxmox VNC 代理响应失败: %w", err)
	}
	port, err := parseProxmoxVNCProxyPort(envelope.Data.Port)
	if err != nil || port <= 0 || port > 65535 {
		return proxmoxVNCAPIResult{}, fmt.Errorf("Proxmox VNC 代理未返回有效端口")
	}
	vncPass := strings.TrimSpace(envelope.Data.Password)
	result := proxmoxVNCAPIResult{
		Port:   port,
		Ticket: strings.TrimSpace(envelope.Data.Ticket),
		User:   strings.TrimSpace(envelope.Data.User),
	}
	result.Password = vncPass
	if result.PasswordOrTicket() == "" {
		return proxmoxVNCAPIResult{}, fmt.Errorf("Proxmox VNC 代理未返回认证信息")
	}
	return result, nil
}

// defaultProxmoxAPIClient mirrors the TLS posture of the existing PVE API
// callers (self-signed node certificates are the norm for PVE).
func defaultProxmoxAPIClient() *http.Client {
	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402 -- PVE 节点自签名证书，与既有 PVE API 调用一致
		},
	}
}

// requestProxmoxVNCProxy issues the vncproxy call against the PVE HTTP API
// using the provider's API token. websocket=true asks PVE for a ticket that
// can be used with the vncwebsocket endpoint.
func requestProxmoxVNCProxy(ctx context.Context, client *http.Client, host, node, resource, apiToken, vmid string, websocket bool) (proxmoxVNCAPIResult, error) {
	apiToken = cleanProxmoxAPIToken(apiToken)
	if apiToken == "" {
		return proxmoxVNCAPIResult{}, fmt.Errorf("节点未配置 Proxmox API Token")
	}
	apiURL := proxmoxVNCProxyAPIURL(host, node, resource, vmid)
	if apiURL == "" {
		return proxmoxVNCAPIResult{}, fmt.Errorf("节点缺少有效的 API 地址")
	}
	form := url.Values{}
	if websocket {
		form.Set("websocket", "1")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return proxmoxVNCAPIResult{}, err
	}
	req.Header.Set("Authorization", "PVEAPIToken="+apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if client == nil {
		client = defaultProxmoxAPIClient()
	}
	resp, err := client.Do(req)
	if err != nil {
		return proxmoxVNCAPIResult{}, fmt.Errorf("调用 Proxmox VNC 代理失败: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return proxmoxVNCAPIResult{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return proxmoxVNCAPIResult{}, fmt.Errorf("Proxmox VNC 代理请求失败: HTTP %d", resp.StatusCode)
	}
	return parseProxmoxVNCAPIResponse(body)
}

// proxmoxAPIGet performs an authenticated GET against the PVE HTTP API and
// returns the raw body (used by the API runtime lookup).
func proxmoxAPIGet(ctx context.Context, client *http.Client, host, apiToken, path string) ([]byte, error) {
	apiToken = cleanProxmoxAPIToken(apiToken)
	if apiToken == "" {
		return nil, fmt.Errorf("节点未配置 Proxmox API Token")
	}
	apiURL := utils.BuildEndpointURL("https", host, proxmoxAPIPort, path)
	if apiURL == "" {
		return nil, fmt.Errorf("节点缺少有效的 API 地址")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "PVEAPIToken="+apiToken)
	if client == nil {
		client = defaultProxmoxAPIClient()
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用 Proxmox API 失败: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Proxmox API 请求失败: HTTP %d (%s)", resp.StatusCode, path)
	}
	return body, nil
}
