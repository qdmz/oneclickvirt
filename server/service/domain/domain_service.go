package domain

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"oneclickvirt/global"
	domainModel "oneclickvirt/model/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DomainService struct {
	db *gorm.DB
}

func NewDomainService(db *gorm.DB) *DomainService {
	return &DomainService{db: db}
}

// CreateDomain 创建域名绑定
func (s *DomainService) CreateDomain(userID, instanceID uint, domainName, internalIP string, internalPort int, protocol string, agentID *uint) (*domainModel.Domain, error) {
	// 域名格式验证
	if !isValidDomain(domainName) {
		return nil, fmt.Errorf("域名格式无效")
	}

	// 唯一性检查
	var count int64
	s.db.Model(&domainModel.Domain{}).Where("domain = ?", domainName).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("该域名已被绑定")
	}

	// 配额检查
	maxDomains := s.getMaxDomains(userID, agentID)
	var userCount int64
	s.db.Model(&domainModel.Domain{}).Where("user_id = ?", userID).Count(&userCount)
	if int(userCount) >= maxDomains {
		return nil, fmt.Errorf("已达到域名绑定上限(%d)", maxDomains)
	}

	// 后缀检查
	config, _ := s.GetDomainConfig()
	if config != nil && config.AllowedSuffixes != "" {
		if !isAllowedSuffix(domainName, config.AllowedSuffixes) {
			return nil, fmt.Errorf("不允许绑定此后缀的域名")
		}
	}

	d := &domainModel.Domain{
		UserID:       userID,
		InstanceID:   instanceID,
		Domain:       domainName,
		Protocol:     protocol,
		InternalIP:   internalIP,
		InternalPort: internalPort,
		Status:       0,
		AgentID:      agentID,
	}

	if err := s.db.Create(d).Error; err != nil {
		return nil, fmt.Errorf("创建域名绑定失败: %v", err)
	}

	// 尝试配置DNS
	if err := s.configureDNS(d.Domain, d.InternalIP); err != nil {
		global.APP_LOG.Warn("DNS配置失败", zap.String("domain", d.Domain), zap.Error(err))
		d.Status = 2
		s.db.Save(d)
	} else {
		d.Status = 1
		s.db.Save(d)
	}

	return d, nil
}

// GetDomains 获取用户域名列表
func (s *DomainService) GetDomains(userID uint, page, pageSize int) ([]domainModel.Domain, int64, error) {
	var domains []domainModel.Domain
	var total int64

	db := s.db.Model(&domainModel.Domain{}).Where("user_id = ?", userID)
	db.Count(&total)

	offset := (page - 1) * pageSize
	if err := db.Order("id desc").Offset(offset).Limit(pageSize).Find(&domains).Error; err != nil {
		return nil, 0, err
	}
	return domains, total, nil
}

// UpdateDomain 更新域名绑定
func (s *DomainService) UpdateDomain(domainID, userID uint, updates map[string]interface{}) error {
	var d domainModel.Domain
	if err := s.db.Where("id = ? AND user_id = ?", domainID, userID).First(&d).Error; err != nil {
		return fmt.Errorf("域名记录不存在")
	}

	if domainName, ok := updates["domain"]; ok {
		if newDomain, ok := domainName.(string); ok && newDomain != d.Domain {
			if !isValidDomain(newDomain) {
				return fmt.Errorf("域名格式无效")
			}
			var count int64
			s.db.Model(&domainModel.Domain{}).Where("domain = ? AND id != ?", newDomain, domainID).Count(&count)
			if count > 0 {
				return fmt.Errorf("该域名已被绑定")
			}
			// 更新DNS
			s.removeDNS(d.Domain)
			updates["domain"] = newDomain
		}
	}

	if err := s.db.Model(&d).Updates(updates).Error; err != nil {
		return err
	}

	// 重新配置DNS
	if err := s.configureDNS(d.Domain, d.InternalIP); err != nil {
		global.APP_LOG.Warn("DNS更新失败", zap.Error(err))
	}

	return nil
}

// DeleteDomain 删除域名绑定
func (s *DomainService) DeleteDomain(domainID, userID uint) error {
	var d domainModel.Domain
	if err := s.db.Where("id = ? AND user_id = ?", domainID, userID).First(&d).Error; err != nil {
		return fmt.Errorf("域名记录不存在")
	}

	s.removeDNS(d.Domain)
	s.removeNginx(d.Domain)
	return s.db.Delete(&d).Error
}

// GetUserDomainCount 获取用户域名数量
func (s *DomainService) GetUserDomainCount(userID uint) (int64, error) {
	var count int64
	err := s.db.Model(&domainModel.Domain{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// GetAvailableQuota 获取可用配额
func (s *DomainService) GetAvailableQuota(userID uint, agentID *uint) (used int64, max int, err error) {
	max = s.getMaxDomains(userID, agentID)
	used, err = s.GetUserDomainCount(userID)
	return
}

// AdminGetDomains 管理员获取所有域名
func (s *DomainService) AdminGetDomains(page, pageSize int, filters map[string]interface{}) ([]domainModel.Domain, int64, error) {
	var domains []domainModel.Domain
	var total int64

	db := s.db.Model(&domainModel.Domain{})
	if v, ok := filters["userId"]; ok {
		db = db.Where("user_id = ?", v)
	}
	if v, ok := filters["agentId"]; ok {
		db = db.Where("agent_id = ?", v)
	}
	if v, ok := filters["domain"]; ok {
		db = db.Where("domain LIKE ?", "%"+v.(string)+"%")
	}
	if v, ok := filters["status"]; ok {
		db = db.Where("status = ?", v)
	}

	db.Count(&total)
	offset := (page - 1) * pageSize
	if err := db.Order("id desc").Offset(offset).Limit(pageSize).Find(&domains).Error; err != nil {
		return nil, 0, err
	}
	return domains, total, nil
}

// AdminDeleteDomain 管理员删除域名
func (s *DomainService) AdminDeleteDomain(domainID uint) error {
	var d domainModel.Domain
	if err := s.db.First(&d, domainID).Error; err != nil {
		return fmt.Errorf("域名记录不存在")
	}
	s.removeDNS(d.Domain)
	s.removeNginx(d.Domain)
	return s.db.Delete(&d).Error
}

// GetDomainConfig 获取域名配置
func (s *DomainService) GetDomainConfig() (*domainModel.DomainConfig, error) {
	var config domainModel.DomainConfig
	if err := s.db.First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建默认配置
			config = domainModel.DomainConfig{
				MaxDomainsPerUser:     3,
				MaxDomainsPerAgentUser: 5,
				DefaultTTL:            300,
				DNSType:               "dnsmasq",
				DNSConfigPath:         "/etc/dnsmasq.d/oneclickvirt-hosts.conf",
				NginxConfigPath:       "/etc/nginx/conf.d/oneclickvirt-domains",
			}
			s.db.Create(&config)
			return &config, nil
		}
		return nil, err
	}
	return &config, nil
}

// UpdateDomainConfig 更新域名配置
func (s *DomainService) UpdateDomainConfig(config *domainModel.DomainConfig) error {
	return s.db.Save(config).Error
}

// SyncAllDNS 全量DNS同步
func (s *DomainService) SyncAllDNS() error {
	return s.rebuildAllDNS()
}

// configureDNS 配置DNS解析
func (s *DomainService) configureDNS(domainName, internalIP string) error {
	config, err := s.GetDomainConfig()
	if err != nil {
		return fmt.Errorf("获取域名配置失败: %v", err)
	}

	switch config.DNSType {
	case "hosts":
		return s.configureHosts(domainName, internalIP)
	default: // dnsmasq
		return s.configureDnsmasq(domainName, internalIP, config.DNSConfigPath)
	}
}

// configureDnsmasq 配置dnsmasq
func (s *DomainService) configureDnsmasq(domainName, internalIP, configPath string) error {
	// 确保目录存在
	dir := configPath[:strings.LastIndex(configPath, "/")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 读取现有配置
	existing := ""
	if data, err := os.ReadFile(configPath); err == nil {
		existing = string(data)
	}

	// 检查是否已存在
	entry := fmt.Sprintf("address=/%s/%s", domainName, internalIP)
	if strings.Contains(existing, entry) {
		return nil // 已存在
	}

	// 追加新记录
	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开配置文件失败: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(entry + "\n"); err != nil {
		return fmt.Errorf("写入配置失败: %v", err)
	}

	// 重载dnsmasq
	if err := exec.Command("systemctl", "reload", "dnsmasq").Run(); err != nil {
		global.APP_LOG.Warn("重载dnsmasq失败，尝试restart", zap.Error(err))
		exec.Command("systemctl", "restart", "dnsmasq").Run()
	}

	return nil
}

// configureHosts 配置/etc/hosts
func (s *DomainService) configureHosts(domainName, internalIP string) error {
	data, err := os.ReadFile("/etc/hosts")
	if err != nil {
		return fmt.Errorf("读取hosts文件失败: %v", err)
	}

	content := string(data)
	entry := fmt.Sprintf("%s %s", internalIP, domainName)
	if strings.Contains(content, domainName) {
		return nil
	}

	f, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开hosts文件失败: %v", err)
	}
	defer f.Close()

	f.WriteString(entry + "\n")
	return nil
}

// removeDNS 移除DNS记录
func (s *DomainService) removeDNS(domainName string) error {
	config, err := s.GetDomainConfig()
	if err != nil {
		return err
	}

	switch config.DNSType {
	case "hosts":
		return s.removeHostsEntry(domainName)
	default:
		return s.removeDnsmasqEntry(domainName, config.DNSConfigPath)
	}
}

// removeDnsmasqEntry 移除dnsmasq记录
func (s *DomainService) removeDnsmasqEntry(domainName, configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	for _, line := range lines {
		if !strings.Contains(line, "/"+domainName+"/") {
			newLines = append(newLines, line)
		}
	}

	if err := os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return err
	}

	exec.Command("systemctl", "reload", "dnsmasq").Run()
	return nil
}

// removeHostsEntry 移除hosts记录
func (s *DomainService) removeHostsEntry(domainName string) error {
	data, err := os.ReadFile("/etc/hosts")
	if err != nil {
		return nil
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	for _, line := range lines {
		if !strings.Contains(line, domainName) || strings.HasPrefix(line, "#") {
			newLines = append(newLines, line)
		}
	}

	return os.WriteFile("/etc/hosts", []byte(strings.Join(newLines, "\n")), 0644)
}

// rebuildAllDNS 重建整个DNS配置
func (s *DomainService) rebuildAllDNS() error {
	config, err := s.GetDomainConfig()
	if err != nil {
		return err
	}

	var domains []domainModel.Domain
	s.db.Where("status = ?", 1).Find(&domains)

	switch config.DNSType {
	case "hosts":
		return s.rebuildHosts(domains)
	default:
		return s.rebuildDnsmasq(domains, config.DNSConfigPath)
	}
}

func (s *DomainService) rebuildDnsmasq(domains []domainModel.Domain, configPath string) error {
	dir := configPath[:strings.LastIndex(configPath, "/")]
	os.MkdirAll(dir, 0755)

	var content string
	for _, d := range domains {
		content += fmt.Sprintf("address=/%s/%s\n", d.Domain, d.InternalIP)
	}

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return err
	}

	exec.Command("systemctl", "reload", "dnsmasq").Run()
	return nil
}

func (s *DomainService) rebuildHosts(domains []domainModel.Domain) error {
	data, _ := os.ReadFile("/etc/hosts")
	lines := strings.Split(string(data), "\n")
	var newLines []string
	for _, line := range lines {
		keep := true
		for _, d := range domains {
			if strings.Contains(line, d.Domain) {
				keep = false
				break
			}
		}
		if keep {
			newLines = append(newLines, line)
		}
	}

	for _, d := range domains {
		newLines = append(newLines, fmt.Sprintf("%s %s", d.InternalIP, d.Domain))
	}

	return os.WriteFile("/etc/hosts", []byte(strings.Join(newLines, "\n")), 0644)
}

// configureNginx 配置Nginx反代
func (s *DomainService) configureNginx(domainName, internalIP string, internalPort, externalPort int, protocol string, ssl bool) error {
	config, _ := s.GetDomainConfig()
	if config == nil {
		return nil
	}
	configPath := config.NginxConfigPath
	dir := configPath[:strings.LastIndex(configPath, "/")]
	os.MkdirAll(dir, 0755)

	if externalPort == 0 {
		externalPort = 80
	}

	listenDirective := fmt.Sprintf("listen %d;", externalPort)
	if ssl {
		listenDirective = fmt.Sprintf("listen %d ssl;", externalPort)
	}

	nginxConf := fmt.Sprintf(`server {
    %s
    server_name %s;
    location / {
        proxy_pass %s://%s:%d;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
`, listenDirective, domainName, protocol, internalIP, internalPort)

	filePath := fmt.Sprintf("%s/%s.conf", dir, domainName)
	if err := os.WriteFile(filePath, []byte(nginxConf), 0644); err != nil {
		return err
	}

	exec.Command("nginx", "-t").Run()
	exec.Command("systemctl", "reload", "nginx").Run()
	return nil
}

// removeNginx 移除Nginx配置
func (s *DomainService) removeNginx(domainName string) error {
	config, _ := s.GetDomainConfig()
	if config == nil {
		return nil
	}
	dir := config.NginxConfigPath[:strings.LastIndex(config.NginxConfigPath, "/")]
	filePath := fmt.Sprintf("%s/%s.conf", dir, domainName)
	os.Remove(filePath)
	exec.Command("nginx", "-t").Run()
	exec.Command("systemctl", "reload", "nginx").Run()
	return nil
}

// getMaxDomains 获取用户最大域名数
func (s *DomainService) getMaxDomains(userID uint, agentID *uint) int {
	config, _ := s.GetDomainConfig()
	if config == nil {
		return 3
	}
	if agentID != nil && *agentID > 0 {
		return config.MaxDomainsPerAgentUser
	}
	return config.MaxDomainsPerUser
}

func isValidDomain(domain string) bool {
	// 简单域名验证
	re := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	return re.MatchString(domain)
}

func isAllowedSuffix(domain, allowedSuffixes string) bool {
	suffixes := strings.Split(allowedSuffixes, ",")
	for _, s := range suffixes {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if strings.HasSuffix(domain, s) {
			return true
		}
	}
	return false
}
