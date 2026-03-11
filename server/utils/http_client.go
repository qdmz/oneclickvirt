package utils

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"oneclickvirt/global"
)

var (
	defaultHTTPClient     *http.Client
	defaultHTTPClientOnce sync.Once
)

// GetDefaultHTTPClient returns the default HTTP client with connection pooling
func GetDefaultHTTPClient() *http.Client {
	defaultHTTPClientOnce.Do(func() {
		defaultHTTPClient = &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   10,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				ForceAttemptHTTP2:     true,
			},
		}
	})
	return defaultHTTPClient
}

var (
	sharedTransport       *http.Transport
	sharedTransportOnce   sync.Once
	insecureTransport     *http.Transport
	insecureTransportOnce sync.Once
)

func getSharedTransport() *http.Transport {
	sharedTransportOnce.Do(func() {
		sharedTransport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     true,
		}
	})
	return sharedTransport
}

// getInsecureTransport returns transport with InsecureSkipVerify for internal hypervisors only
func getInsecureTransport() *http.Transport {
	insecureTransportOnce.Do(func() {
		insecureTransport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				// Security: only for internal hypervisors (LXD/Incus) with self-signed certs.
				// Use TrustedFingerprint in Provider model for fingerprint verification when possible.
				InsecureSkipVerify: true,
			},
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     true,
		}
	})
	return insecureTransport
}

// GetHTTPClientWithTimeout creates HTTP client with custom timeout (reuses transport)
func GetHTTPClientWithTimeout(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: getSharedTransport(),
	}
}

// GetInsecureHTTPClient returns HTTP client that skips TLS verification
func GetInsecureHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: getInsecureTransport(),
	}
}

// CleanupHTTPTransports cleans up idle connections (call on app shutdown)
func CleanupHTTPTransports() {
	if sharedTransport != nil {
		sharedTransport.CloseIdleConnections()
		global.APP_LOG.Info("Cleaned up shared HTTP transport idle connections")
	}
	if insecureTransport != nil {
		insecureTransport.CloseIdleConnections()
		global.APP_LOG.Info("Cleaned up insecure HTTP transport idle connections")
	}
}
