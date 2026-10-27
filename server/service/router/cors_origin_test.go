package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func corsTestContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "http://panel.example.com/api/v1/health", nil)
	c.Request.Host = "panel.example.com"
	return c
}

func TestCorsOriginAllowedLoopbackOnlyWhenAllowed(t *testing.T) {
	c := corsTestContext()
	allowed := map[string]struct{}{}
	if !corsOriginAllowed(c, "http://127.0.0.1:5173", "", allowed, true) {
		t.Fatal("loopback origin should be allowed when allowLoopback=true (non-production)")
	}
	if corsOriginAllowed(c, "http://127.0.0.1:5173", "", allowed, false) {
		t.Fatal("loopback origin must not be allowed when allowLoopback=false (production)")
	}
	if corsOriginAllowed(c, "https://evil.example.com", "", allowed, true) {
		t.Fatal("unrelated cross origin must not be allowed")
	}
}

func TestTrustedProxyCIDRsParse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	err := r.SetTrustedProxies([]string{
		"127.0.0.0/8", "::1/128",
		"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16",
		"fc00::/7",
	})
	if err != nil {
		t.Fatalf("trusted proxy CIDRs should parse: %v", err)
	}
}
