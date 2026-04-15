package health

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"go.uber.org/zap"
)

type ZJMFHealthChecker struct {
	*BaseHealthChecker
}

func NewZJMFHealthChecker(config HealthConfig, logger *zap.Logger) *ZJMFHealthChecker {
	return &ZJMFHealthChecker{
		BaseHealthChecker: NewBaseHealthChecker(config, logger),
	}
}

func (z *ZJMFHealthChecker) CheckHealth(ctx context.Context) (*HealthResult, error) {
	if !z.config.APIEnabled {
		return &HealthResult{
			Status:    HealthStatusUnknown,
			Timestamp: time.Now(),
			Errors:    []string{"API is not enabled"},
		}, nil
	}

	checks := []func(context.Context) CheckResult{
		z.createCheckFunc(CheckTypeAPI, z.checkAPI),
	}

	result := z.executeChecks(ctx, checks)
	return result, nil
}

func (z *ZJMFHealthChecker) checkAPI(ctx context.Context) error {
	if z.config.APIKey == "" || z.config.APISecret == "" {
		return fmt.Errorf("ZJMF API key or secret is not configured")
	}

	apiURL := z.config.Host
	if !hasScheme(apiURL) {
		apiURL = "http://" + apiURL
	}

	url := fmt.Sprintf("%s/api.php", apiURL)

	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	params := map[string]string{
		"action":    "Common/ping",
		"client_id": z.config.APIKey,
		"time":      timestamp,
	}

	signature := z.generateSignature(params)

	requestParams := map[string]interface{}{
		"action":    "Common/ping",
		"client_id": z.config.APIKey,
		"time":      timestamp,
		"sign":      signature,
	}

	jsonData, err := json.Marshal(requestParams)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := z.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var result ZJMFAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "success" {
		return fmt.Errorf("API error: %s", result.Msg)
	}

	return nil
}

func (z *ZJMFHealthChecker) generateSignature(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, params[k]))
	}
	sort.Strings(pairs)

	signStr := ""
	for _, pair := range pairs {
		if signStr != "" {
			signStr += "&"
		}
		signStr += pair
	}
	signStr += "&secret_key=" + z.config.APISecret

	hash := md5.Sum([]byte(signStr))
	return fmt.Sprintf("%x", hash)
}

type ZJMFAPIResponse struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func hasScheme(url string) bool {
	return len(url) > 7 && (url[:7] == "http://" || url[:8] == "https://")
}
