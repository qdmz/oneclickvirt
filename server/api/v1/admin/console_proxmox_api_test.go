package admin

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	providerModel "oneclickvirt/model/provider"

	"github.com/gorilla/websocket"
)

// fakePVE serves the handful of PVE API endpoints used by the API console
// path and records how it was called.
type fakePVE struct {
	server       *httptest.Server
	vncPath      string
	vncAuth      string
	vncFormWS    string
	clusterCalls int
}

func newFakePVE(t *testing.T) *fakePVE {
	t.Helper()
	f := &fakePVE{}
	mux := http.NewServeMux()
	mux.HandleFunc("/api2/json/nodes/pve/qemu/100/vncproxy", func(w http.ResponseWriter, r *http.Request) {
		f.vncPath = r.URL.Path
		f.vncAuth = r.Header.Get("Authorization")
		_ = r.ParseForm()
		f.vncFormWS = r.Form.Get("websocket")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":{"ticket":"PVEVNC:fake-ticket","port":5901,"user":"root@pam","cert":"x"}}`)
	})
	mux.HandleFunc("/api2/json/nodes/pve/lxc/101/vncproxy", func(w http.ResponseWriter, r *http.Request) {
		f.vncPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":{"password":"ct-pass","port":"5902"}}`)
	})
	mux.HandleFunc("/api2/json/cluster/resources", func(w http.ResponseWriter, r *http.Request) {
		f.clusterCalls++
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[{"vmid":100,"name":"vm100","type":"qemu","node":"pve","status":"running"}]}`)
	})
	f.server = httptest.NewTLSServer(mux)
	t.Cleanup(f.server.Close)
	return f
}

func (f *fakePVE) host() string { return strings.TrimPrefix(f.server.URL, "https://") }

// useFakePVEPort points the PVE API client helpers at the fake server port.
func useFakePVEPort(t *testing.T, f *fakePVE) {
	t.Helper()
	u, err := url.Parse(f.server.URL)
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		t.Fatal(err)
	}
	original := proxmoxAPIPort
	proxmoxAPIPort = port
	t.Cleanup(func() { proxmoxAPIPort = original })
}

func TestProxmoxVNCProxyViaFakeAPI(t *testing.T) {
	f := newFakePVE(t)
	useFakePVEPort(t, f)
	ctx := context.Background()

	result, err := requestProxmoxVNCProxy(ctx, f.server.Client(), f.host(), "pve", "qemu", "root@pam!tok=secret\r\n", "100", true)
	if err != nil {
		t.Fatalf("requestProxmoxVNCProxy() error: %v", err)
	}
	if f.vncPath != "/api2/json/nodes/pve/qemu/100/vncproxy" {
		t.Fatalf("unexpected vncproxy path: %s", f.vncPath)
	}
	if f.vncAuth != "PVEAPIToken=root@pam!tok=secret" {
		t.Fatalf("unexpected Authorization header: %q", f.vncAuth)
	}
	if f.vncFormWS != "1" {
		t.Fatalf("websocket form field = %q, want 1", f.vncFormWS)
	}
	if result.Port != 5901 || result.Ticket != "PVEVNC:fake-ticket" || result.PasswordOrTicket() != "PVEVNC:fake-ticket" {
		t.Fatalf("unexpected result: %+v", result)
	}

	// LXC guests use the lxc resource path and may return a password instead
	// of a ticket, with the port encoded as a string.
	result, err = requestProxmoxVNCProxy(ctx, f.server.Client(), f.host(), "pve", "lxc", "root@pam!tok=secret", "101", false)
	if err != nil {
		t.Fatalf("lxc requestProxmoxVNCProxy() error: %v", err)
	}
	if f.vncPath != "/api2/json/nodes/pve/lxc/101/vncproxy" {
		t.Fatalf("unexpected lxc path: %s", f.vncPath)
	}
	if result.Port != 5902 || result.PasswordOrTicket() == "" {
		t.Fatalf("unexpected lxc result: %+v", result)
	}
}

func TestParseProxmoxVNCAPIResponseShapes(t *testing.T) {
	if _, err := parseProxmoxVNCAPIResponse([]byte(`{"data":{"port":70000,"ticket":"t"}}`)); err == nil {
		t.Fatal("out-of-range port unexpectedly accepted")
	}
	if _, err := parseProxmoxVNCAPIResponse([]byte(`{"data":{"port":5900}}`)); err == nil {
		t.Fatal("missing credential unexpectedly accepted")
	}
	got, err := parseProxmoxVNCAPIResponse([]byte(`{"data":{"port":"5903","password":"p","user":"root@pam"}}`))
	if err != nil || got.Port != 5903 || got.User != "root@pam" || got.PasswordOrTicket() == "" {
		t.Fatalf("string port shape = %+v, %v", got, err)
	}
}

func TestProxmoxVNCWebSocketURLAssembly(t *testing.T) {
	got := proxmoxVNCWebSocketURL("pve.example.com", "pve", "qemu", "100", 5901, "PVEVNC:a+b/c=")
	for _, want := range []string{
		"wss://pve.example.com:8006/api2/json/nodes/pve/qemu/100/vncwebsocket?",
		"port=5901",
		"vncticket=PVEVNC%3Aa%2Bb%2Fc%3D",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("vncwebsocket URL %q missing %q", got, want)
		}
	}
	if resource, err := proxmoxAPIResource("container"); err != nil || resource != "lxc" {
		t.Fatalf("proxmoxAPIResource(container) = %q, %v", resource, err)
	}
	if _, err := proxmoxAPIResource("sparc"); err == nil {
		t.Fatal("unknown runtime type unexpectedly accepted")
	}
}

func TestOpenProxmoxVNCConnViaAPIFake(t *testing.T) {
	f := newFakePVE(t)
	useFakePVEPort(t, f)

	var dialedURL, dialedToken string
	originalDial := proxmoxVNCWebSocketDial
	proxmoxVNCWebSocketDial = func(wsURL, apiToken string) (net.Conn, error) {
		dialedURL, dialedToken = wsURL, apiToken
		client, server := net.Pipe()
		go func() { _ = server.Close() }()
		return client, nil
	}
	defer func() { proxmoxVNCWebSocketDial = originalDial }()

	target := consoleTarget{
		protocol:    consoleProtocolVNC,
		proxmoxVNC:  true,
		transport:   "api",
		runtimeID:   "100",
		runtimeNode: "pve",
		runtimeType: "qemu",
		provider: providerModel.Provider{
			ConnectionType: "api",
			Endpoint:       "https://" + f.host(),
			Token:          "root@pam!tok=secret",
			HostName:       "pve",
		},
	}
	// The fake client must trust the httptest TLS cert; route the package's
	// default client through the test server client for this call.
	conn, cleanup, credential, err := openProxmoxVNCConnViaAPIWithClient(target, f.server.Client())
	if err != nil {
		t.Fatalf("openProxmoxVNCConnViaAPIWithClient() error: %v", err)
	}
	defer cleanup()
	if conn == nil {
		t.Fatal("expected a connection")
	}
	if credential != "PVEVNC:fake-ticket" {
		t.Fatalf("credential = %q, want ticket", credential)
	}
	if !strings.Contains(dialedURL, "/vncwebsocket?") || !strings.Contains(dialedURL, "vncticket=PVEVNC%3Afake-ticket") {
		t.Fatalf("unexpected dialed URL: %s", dialedURL)
	}
	if dialedToken != "root@pam!tok=secret" {
		t.Fatalf("unexpected dial token: %q", dialedToken)
	}
}

func TestLookupProxmoxConsoleRuntimeViaAPIFake(t *testing.T) {
	f := newFakePVE(t)
	useFakePVEPort(t, f)
	inst := providerModel.Instance{ID: 7, Name: "vm100", ProviderVMID: "100"}
	p := providerModel.Provider{
		ConnectionType: "api",
		Endpoint:       "https://" + f.host(),
		Token:          "root@pam!tok=secret",
		HostName:       "pve",
	}
	runtime, err := lookupProxmoxConsoleRuntimeViaAPIWithClient(inst, p, f.server.Client())
	if err != nil {
		t.Fatalf("lookupProxmoxConsoleRuntimeViaAPIWithClient() error: %v", err)
	}
	if runtime.ID != "100" || runtime.Type != "qemu" || runtime.Node != "pve" || runtime.Status != "running" {
		t.Fatalf("unexpected runtime: %+v", runtime)
	}
}

func TestProxmoxAPITransportSelection(t *testing.T) {
	withToken := providerModel.Provider{ConnectionType: "api", Endpoint: "https://pve.example.com", Token: "root@pam!tok=secret"}
	if transport, reason := proxmoxConsoleVNCTransport(withToken); transport != "api" || reason != "" {
		t.Fatalf("api transport = %q, %q", transport, reason)
	}
	noToken := providerModel.Provider{ConnectionType: "api", Endpoint: "https://pve.example.com"}
	if _, reason := proxmoxConsoleVNCTransport(noToken); !strings.Contains(reason, "API Token") {
		t.Fatalf("missing-token reason = %q", reason)
	}
	implicit := providerModel.Provider{Endpoint: "https://pve.example.com", Token: "root@pam!tok=secret"}
	if transport, _ := proxmoxConsoleVNCTransport(implicit); transport != "api" {
		t.Fatalf("token-only provider transport = %q, want api", transport)
	}
}

func TestCachedVNCSessionExpiry(t *testing.T) {
	const instanceID = uint(987654)
	vncSessionMu.Lock()
	vncSessions[instanceID] = &vncSession{Host: "pve", Port: 5901, ExpiresAt: time.Now().Add(-time.Second)}
	vncSessionMu.Unlock()
	if _, ok := cachedVNCSession(instanceID, time.Now()); ok {
		t.Fatal("expired session unexpectedly returned")
	}
	vncSessionMu.Lock()
	vncSessions[instanceID] = &vncSession{Host: "pve", Port: 5901, ExpiresAt: time.Now().Add(time.Minute)}
	vncSessionMu.Unlock()
	s, ok := cachedVNCSession(instanceID, time.Now())
	if !ok || s.Port != 5901 {
		t.Fatalf("live session = %+v, %v", s, ok)
	}
	vncSessionMu.Lock()
	delete(vncSessions, instanceID)
	vncSessionMu.Unlock()
}

func TestWSNetConnEcho(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer ws.Close()
		for {
			mt, msg, err := ws.ReadMessage()
			if err != nil {
				return
			}
			if err := ws.WriteMessage(mt, msg); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial echo server: %v", err)
	}
	conn := &wsNetConn{ws: ws}
	defer conn.Close()

	payload := []byte("RFB 003.008\npartial-frame")
	if _, err := conn.Write(payload); err != nil {
		t.Fatalf("write: %v", err)
	}
	// Read back in two chunks to exercise the buffered-reader path.
	buf := make([]byte, 5)
	if _, err := io.ReadFull(conn, buf); err != nil {
		t.Fatalf("read chunk 1: %v", err)
	}
	rest := make([]byte, len(payload)-5)
	if _, err := io.ReadFull(conn, rest); err != nil {
		t.Fatalf("read chunk 2: %v", err)
	}
	if string(buf)+string(rest) != string(payload) {
		t.Fatalf("echo mismatch: %q + %q", buf, rest)
	}
}
