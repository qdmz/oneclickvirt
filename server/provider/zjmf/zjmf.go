package zjmf

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"oneclickvirt/global"
	"oneclickvirt/provider"
	"oneclickvirt/provider/health"
	"oneclickvirt/utils"

	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type ZJMFProvider struct {
	config        provider.NodeConfig
	apiKey        string
	secretKey     string
	apiURL        string
	httpClient    *http.Client
	connected     bool
	healthChecker health.HealthChecker
	version       string
	mu            sync.RWMutex
}

func NewZJMFProvider() provider.Provider {
	return &ZJMFProvider{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (z *ZJMFProvider) GetType() string {
	return "zjmf"
}

func (z *ZJMFProvider) GetName() string {
	return z.config.Name
}

func (z *ZJMFProvider) GetSupportedInstanceTypes() []string {
	return []string{"container", "vm"}
}

func (z *ZJMFProvider) Connect(ctx context.Context, config provider.NodeConfig) error {
	z.config = config

	z.apiKey = config.APIKey
	z.secretKey = config.APISecret

	if config.APIURL != "" {
		z.apiURL = config.APIURL
	} else {
		z.apiURL = config.Host
	}

	if !strings.HasPrefix(z.apiURL, "http") {
		z.apiURL = "http://" + z.apiURL
	}

	global.APP_LOG.Info("ZJMF provider开始连接",
		zap.String("host", utils.TruncateString(config.Host, 32)),
		zap.String("apiURL", z.apiURL))

	healthConfig := health.HealthConfig{
		Host:       config.Host,
		Port:       config.Port,
		APIEnabled: true,
		APIURL:     z.apiURL,
		APIKey:     z.apiKey,
		APISecret:  z.secretKey,
		SSHEnabled: false,
		Timeout:    30 * time.Second,
	}

	zapLogger, _ := zap.NewProduction()
	z.healthChecker = health.NewZJMFHealthChecker(healthConfig, zapLogger)

	z.connected = true

	global.APP_LOG.Info("ZJMF provider连接成功",
		zap.String("host", utils.TruncateString(config.Host, 32)),
		zap.String("apiURL", z.apiURL))

	return nil
}

func (z *ZJMFProvider) Disconnect(ctx context.Context) error {
	z.connected = false
	return nil
}

func (z *ZJMFProvider) IsConnected() bool {
	return z.connected
}

func (z *ZJMFProvider) HealthCheck(ctx context.Context) (*health.HealthResult, error) {
	if z.healthChecker == nil {
		return nil, fmt.Errorf("health checker not initialized")
	}
	return z.healthChecker.CheckHealth(ctx)
}

func (z *ZJMFProvider) GetHealthChecker() health.HealthChecker {
	return z.healthChecker
}

func (z *ZJMFProvider) GetVersion() string {
	z.mu.RLock()
	defer z.mu.RUnlock()
	return z.version
}

func (z *ZJMFProvider) SetInstancePassword(ctx context.Context, instanceID, password string) error {
	result, err := z.doRequest("Host/reset_password", map[string]interface{}{
		"host_id":      instanceID,
		"new_password": password,
	})
	if err != nil {
		return err
	}
	if result.Status != "success" {
		return fmt.Errorf("%s", result.Msg)
	}
	return nil
}

func (z *ZJMFProvider) ResetInstancePassword(ctx context.Context, instanceID string) (string, error) {
	password := utils.GeneratePassword(16)
	err := z.SetInstancePassword(ctx, instanceID, password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (z *ZJMFProvider) ExecuteSSHCommand(ctx context.Context, command string) (string, error) {
	return "", fmt.Errorf("ZJMF provider does not support direct SSH command execution")
}

func (z *ZJMFProvider) generateSignature(params map[string]string) string {
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

	signStr := strings.Join(pairs, "&") + "&secret_key=" + z.secretKey
	hash := md5.Sum([]byte(signStr))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}

type ZJMFAPIResponse struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"
}

func (z *ZJMFProvider) doRequest(action string, params map[string]interface{}) (*ZJMFAPIResponse, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	requestParams := map[string]string{
		"action":    action,
		"client_id": z.apiKey,
		"time":      timestamp,
	}

	for k, v := range params {
		requestParams[k] = fmt.Sprintf("%v", v)
	}

	requestParams["sign"] = z.generateSignature(requestParams)

	jsonData, err := json.Marshal(requestParams)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api.php", z.apiURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := z.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result ZJMFAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

func (z *ZJMFProvider) ListInstances(ctx context.Context) ([]provider.Instance, error) {
	result, err := z.doRequest("Host/get_list", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("%s", result.Msg)
	}

	var instances []provider.Instance
	data, ok := result.Data.([]interface{})
	if !ok {
		return instances, nil
	}

	for _, item := range data {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		status := cast.ToString(itemMap["status"])
		instanceStatus := "unknown"
		switch status {
		case "Active", "Running":
			instanceStatus = "running"
		case "Suspended":
			instanceStatus = "suspended"
		case "Terminated", "Deleted":
			instanceStatus = "deleted"
		default:
			instanceStatus = "stopped"
		}

		instance := provider.Instance{
			ID:       cast.ToString(itemMap["id"]),
			Name:     cast.ToString(itemMap["name"]),
			Status:   instanceStatus,
			Image:    cast.ToString(itemMap["os"]),
			IP:       cast.ToString(itemMap["ipaddress"]),
			CPU:      cast.ToString(itemMap["cpu"]),
			Memory:   cast.ToString(itemMap["memory"]),
			Disk:     cast.ToString(itemMap["disk"]),
			Created:  cast.ToTime(itemMap["create_time"]),
			Metadata: make(map[string]string),
		}

		if instance.Metadata == nil {
			instance.Metadata = make(map[string]string)
		}
		instance.Metadata["product_id"] = cast.ToString(itemMap["product_id"])
		instance.Metadata["product_name"] = cast.ToString(itemMap["product_name"])
		instance.Metadata["dc_id"] = cast.ToString(itemMap["dc_id"])
		instance.Metadata["dc_name"] = cast.ToString(itemMap["dc_name"])

		instances = append(instances, instance)
	}

	return instances, nil
}

func (z *ZJMFProvider) CreateInstance(ctx context.Context, config provider.InstanceConfig) error {
	return fmt.Errorf("use CreateInstanceWithProgress for ZJMF provider")
}

func (z *ZJMFProvider) CreateInstanceWithProgress(ctx context.Context, config provider.InstanceConfig, progressCallback provider.ProgressCallback) error {
	progressCallback(10, "准备创建实例...")

	hostname := config.Name
	if hostname == "" {
		hostname = fmt.Sprintf("instance-%d", time.Now().Unix())
	}

	billingCycle := "monthly"
	if config.BillingCycle != "" {
		billingCycle = config.BillingCycle
	}

	params := map[string]interface{}{
		"product_id":   config.ProductID,
		"hostname":     hostname,
		"username":     config.Username,
		"password":     config.Password,
		"billing_cycle": billingCycle,
	}

	if config.Image != "" {
		params["os_id"] = config.Image
	}

	progressCallback(30, "正在提交订单...")

	result, err := z.doRequest("Host/create", params)
	if err != nil {
		progressCallback(0, fmt.Sprintf("创建失败: %s", err.Error()))
		return err
	}

	if result.Status != "success" {
		progressCallback(0, fmt.Sprintf("创建失败: %s", result.Msg))
		return fmt.Errorf("%s", result.Msg)
	}

	progressCallback(80, "实例创建成功，等待开通...")

	progressCallback(100, "创建完成")

	return nil
}

func (z *ZJMFProvider) StartInstance(ctx context.Context, id string) error {
	result, err := z.doRequest("Host/boot", map[string]interface{}{
		"host_id": id,
	})
	if err != nil {
		return err
	}
	if result.Status != "success" {
		return fmt.Errorf("%s", result.Msg)
	}
	return nil
}

func (z *ZJMFProvider) StopInstance(ctx context.Context, id string) error {
	result, err := z.doRequest("Host/shutdown", map[string]interface{}{
		"host_id": id,
	})
	if err != nil {
		return err
	}
	if result.Status != "success" {
		return fmt.Errorf("%s", result.Msg)
	}
	return nil
}

func (z *ZJMFProvider) RestartInstance(ctx context.Context, id string) error {
	result, err := z.doRequest("Host/reboot", map[string]interface{}{
		"host_id": id,
	})
	if err != nil {
		return err
	}
	if result.Status != "success" {
		return fmt.Errorf("%s", result.Msg)
	}
	return nil
}

func (z *ZJMFProvider) DeleteInstance(ctx context.Context, id string) error {
	result, err := z.doRequest("Host/terminated", map[string]interface{}{
		"host_id": id,
	})
	if err != nil {
		return err
	}
	if result.Status != "success" {
		return fmt.Errorf("%s", result.Msg)
	}
	return nil
}

func (z *ZJMFProvider) GetInstance(ctx context.Context, id string) (*provider.Instance, error) {
	result, err := z.doRequest("Host/get_details", map[string]interface{}{
		"host_id": id,
	})
	if err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("%s", result.Msg)
	}

	data, ok := result.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	status := cast.ToString(data["status"])
	instanceStatus := "unknown"
	switch status {
	case "Active", "Running":
		instanceStatus = "running"
	case "Suspended":
		instanceStatus = "suspended"
	case "Terminated", "Deleted":
		instanceStatus = "deleted"
	default:
		instanceStatus = "stopped"
	}

	instance := &provider.Instance{
		ID:       cast.ToString(data["id"]),
		Name:     cast.ToString(data["name"]),
		Status:   instanceStatus,
		Image:    cast.ToString(data["os"]),
		IP:       cast.ToString(data["ipaddress"]),
		CPU:      cast.ToString(data["cpu"]),
		Memory:   cast.ToString(data["memory"]),
		Disk:     cast.ToString(data["disk"]),
		Metadata: make(map[string]string),
	}

	if instance.Metadata == nil {
		instance.Metadata = make(map[string]string)
	}
	instance.Metadata["product_id"] = cast.ToString(data["product_id"])
	instance.Metadata["product_name"] = cast.ToString(data["product_name"])

	return instance, nil
}

func (z *ZJMFProvider) ListImages(ctx context.Context) ([]provider.Image, error) {
	result, err := z.doRequest("Product/get_operate_os_list", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("%s", result.Msg)
	}

	var images []provider.Image
	data, ok := result.Data.([]interface{})
	if !ok {
		return images, nil
	}

	for _, item := range data {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		image := provider.Image{
			ID:      cast.ToString(itemMap["id"]),
			Name:    cast.ToString(itemMap["name"]),
			Type:    cast.ToString(itemMap["type"]),
			OS:      cast.ToString(itemMap["os"]),
			Arch:    cast.ToString(itemMap["arch"]),
			URL:     cast.ToString(itemMap["url"]),
			License: cast.ToString(itemMap["license"]),
		}

		images = append(images, image)
	}

	return images, nil
}

func (z *ZJMFProvider) PullImage(ctx context.Context, image string) error {
	return fmt.Errorf("ZJMF provider does not support image pull")
}

func (z *ZJMFProvider) DeleteImage(ctx context.Context, id string) error {
	return fmt.Errorf("ZJMF provider does not support image deletion")
}

func (z *ZJMFProvider) GetProducts() ([]map[string]interface{}, error) {
	result, err := z.doRequest("Product/products_list", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("%s", result.Msg)
	}

	data, ok := result.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	var products []map[string]interface{}
	for _, item := range data {
		if itemMap, ok := item.(map[string]interface{}); ok {
			products = append(products, itemMap)
		}
	}

	return products, nil
}

func (z *ZJMFProvider) GetLocations() ([]map[string]interface{}, error) {
	result, err := z.doRequest("Product/get_datacenter_list", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("%s", result.Msg)
	}

	data, ok := result.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	var locations []map[string]interface{}
	for _, item := range data {
		if itemMap, ok := item.(map[string]interface{}); ok {
			locations = append(locations, itemMap)
		}
	}

	return locations, nil
}

func (z *ZJMFProvider) GetVNCInfo(hostID string) (map[string]interface{}, error) {
	result, err := z.doRequest("Host/get_vnc", map[string]interface{}{
		"host_id": hostID,
	})
	if err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("%s", result.Msg)
	}

	data, ok := result.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data")
	}

	return data, nil
}

func (z *ZJMFProvider) ReinstallSystem(hostID, osID, rootPassword string) error {
	result, err := z.doRequest("Host/reinstall", map[string]interface{}{
		"host_id":       hostID,
		"os_id":        osID,
		"root_password": rootPassword,
	})
	if err != nil {
		return err
	}
	if result.Status != "success" {
		return fmt.Errorf("%s", result.Msg)
	}
	return nil
}

func (z *ZJMFProvider) TestConnection() error {
	_, err := z.doRequest("Common/ping", map[string]interface{}{})
	return err
}

func init() {
	provider.RegisterProvider("zjmf", NewZJMFProvider)
}
