# 安全审计报告 - OneClickVirt Server

**审计日期**: 2026-03-11  
**项目路径**: `server/`  
**审计工具**: gosec v2, govulncheck (编译错误无法运行), 手动代码审查  
**扫描文件数**: ~100+ Go 源文件

---

## 摘要

| 严重程度 | 数量 |
|---------|------|
| 🔴 Critical | 5 |
| 🟠 High | 9 |
| 🟡 Medium | 7 |
| 🔵 Low | 4 |

---

## 🔴 Critical

### C-01: 配置文件包含硬编码的敏感凭据

- **文件**: `config.yaml`
- **问题**: 配置文件中包含明文的邮箱密码、支付宝私钥、微信支付密钥、Telegram Bot Token、QQ App Key/ID、易支付密钥等敏感信息
- **影响**: 如果配置文件泄露（如代码仓库公开），所有第三方服务凭据将直接暴露
- **修复建议**: 所有密钥、密码、Token 应通过环境变量或密钥管理服务注入，不写入配置文件。在 `.gitignore` 中排除包含真实凭据的配置文件
- **修复代码示例**:
```yaml
# config.yaml - 仅保留占位符
auth:
  email-password: ${EMAIL_PASSWORD}
  telegram-bot-token: ${TELEGRAM_BOT_TOKEN}
payment:
  alipay-private-key: ${ALIPAY_PRIVATE_KEY}
  wechat-api-key: ${WECHAT_API_KEY}
  epay-key: ${EPAY_KEY}
```

### C-02: CORS 配置允许任意来源 + AllowCredentials: true

- **文件**: `router/setup.go:47-54`
- **问题**: `AllowOrigins: ["*"]` 与 `AllowCredentials: true` 同时使用。虽然浏览器通常会阻止这种组合，但这表明安全意识不足，且某些浏览器可能行为不一致
- **影响**: 潜在的跨域攻击风险，恶意网站可以携带用户Cookie发起请求
- **修复建议**: 配置明确的前端URL白名单
- **修复代码示例**:
```go
Router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{
        global.APP_CONFIG.System.FrontendURL, // 如 "https://heyun.ypvps.com"
    },
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Content-Type", "Authorization", "X-Token"},
    ExposeHeaders:    []string{"Content-Length", "Authorization", "X-New-Token"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))
```

### C-03: 支付回调未验证签名

- **文件**: `api/v1/payment/callback.go` (AlipayNotify, WechatNotify, EpayNotify)
- **问题**: 支付宝和微信支付回调中仅有 `// TODO: 验证签名` 的注释，未实现任何签名验证。任何人都可以伪造支付成功通知
- **影响**: 攻击者可以伪造支付回调，在不实际付款的情况下获得服务/余额。这是**直接的资金损失风险**
- **修复建议**: 实现完整的支付平台签名验证
- **修复代码示例**:
```go
func AlipayNotify(c *gin.Context) {
    // 1. 获取原始请求参数
    params := make(map[string]string)
    c.Request.ParseForm()
    for k, v := range c.Request.Form {
        params[k] = v[0]
    }
    
    // 2. 验证支付宝签名
    if !alipay.VerifySign(params, global.APP_CONFIG.Payment.AlipayPublicKey) {
        c.JSON(400, gin.H{"code": 400, "message": "签名验证失败"})
        return
    }
    
    // 3. 验证通知ID幂等性（防止重放攻击）
    notifyId := params["notify_id"]
    if isProcessed(notifyId) {
        c.String(200, "success")
        return
    }
    
    // ... 后续处理
}
```

### C-04: pprof 性能分析端点公开暴露

- **文件**: `main.go:5`
- **问题**: `_ "net/http/pprof"` 被导入，自动注册了 `/debug/pprof/*` 端点，这些端点可以暴露内存数据、goroutine 堆栈等敏感内部信息
- **影响**: 攻击者可以通过 pprof 端点获取服务器内存快照，可能包含密钥、用户数据等敏感信息
- **修复建议**: 在非 debug 环境下不导入 pprof，或将其限制在受保护的管理端口
- **修复代码示例**:
```go
// main.go - 条件性导入
package main

import (
    "fmt"
    "os"
)

func main() {
    // 仅在 debug 模式下启用 pprof
    if os.Getenv("ENABLE_PPROF") == "true" {
        import _ "net/http/pprof"
    }
    // ...
}
```
或者使用独立的管理端口：
```go
if global.APP_CONFIG.System.Env == "development" {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```

### C-05: SSH 连接使用 InsecureIgnoreHostKey

- **文件**: `utils/ssh.go` (dialSSH 函数, CreateSSHConnection, CreateSSHConnectionFromAddress)
- **问题**: 所有 SSH 连接均使用 `ssh.InsecureIgnoreHostKey()`，完全跳过主机密钥验证
- **影响**: 容易遭受中间人攻击（MITM）。攻击者可以冒充目标服务器，截获 SSH 连接中的所有数据和命令
- **修复建议**: 使用 known_hosts 机制或至少实现主机密钥的首次信任（TOFU）模式
- **修复代码示例**:
```go
// 使用 TOFU (Trust On First Use) 模式
func hostKeyCallback(host string, remote net.Addr, key ssh.PublicKey) error {
    hostKeyFile := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts_oneclickvirt")
    
    // 读取已知主机密钥
    known, err := os.ReadFile(hostKeyFile)
    if err == nil {
        // 验证主机密钥
        for _, line := range strings.Split(string(known), "\n") {
            if _, _, _, pubKey, _, err := ssh.ParseKnownHosts([]byte(line)); err == nil {
                if bytes.Equal(key.Marshal(), pubKey.Marshal()) {
                    return nil
                }
            }
        }
        return fmt.Errorf("HOST KEY MISMATCH for %s - possible MITM attack", host)
    }
    
    // 首次连接，保存主机密钥
    line := fmt.Sprintf("%s %s %s\n", host, key.Type(), base64.StdEncoding.EncodeToString(key.Marshal()))
    os.WriteFile(hostKeyFile, []byte(line), 0600)
    return nil
}

sshConfig := &ssh.ClientConfig{
    User:            config.Username,
    Auth:            authMethods,
    HostKeyCallback: hostKeyCallback,
    Timeout:         config.ConnectTimeout,
}
```

---

## 🟠 High

### H-01: 多处 TLS InsecureSkipVerify 设置为 true

- **文件**: 
  - `utils/http_client.go:82`
  - `provider/lxd/lxd.go:597`
  - `provider/incus/incus.go:585`
  - `provider/health/lxd.go:233`
  - `provider/health/incus.go:223`
  - `provider/health/base.go:50`
- **问题**: 6处 `TLSClientConfig.InsecureSkipVerify = true`，跳过 TLS 证书验证
- **影响**: 容易遭受 MITM 攻击，攻击者可以截获与 LXD/Incus Provider 之间的所有通信
- **修复建议**: 对于自签名证书场景，将 Provider 的 CA 证书配置到信任池中
- **修复代码示例**:
```go
// 方案1: 自定义 CA 证书池
caCert, err := os.ReadFile(config.CACertPath)
if err != nil {
    return nil, fmt.Errorf("failed to read CA cert: %w", err)
}
caCertPool := x509.NewCertPool()
caCertPool.AppendCertsFromPEM(caCert)

tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{cert},
    RootCAs:      caCertPool,
    // 不再设置 InsecureSkipVerify
}

// 方案2: 指纹验证（如果自签名证书固定不变）
tlsConfig := &tls.Config{
    InsecureSkipVerify: false,
    VerifyConnection: func(state tls.ConnectionState) error {
        expected := loadExpectedFingerprint(providerID)
        actual := sha256.Sum256(state.PeerCertificates[0].Raw)
        if !bytes.Equal(expected[:], actual[:]) {
            return fmt.Errorf("certificate fingerprint mismatch")
        }
        return nil
    },
}
```

### H-02: 使用弱随机数生成器

- **文件**:
  - `utils/instance.go:12` - `rand.Intn(65536)` 用于生成实例名
  - `service/auth/auth.go:954` - `mathRand` 作为 `crypto/rand` 的后备方案用于生成验证码
- **问题**: `math/rand` 不是密码学安全的随机数生成器，可被预测
- **影响**: 实例名可被预测，攻击者可能枚举用户实例。验证码后备方案使用不安全随机数
- **修复建议**: 全部替换为 `crypto/rand`
- **修复代码示例**:
```go
// utils/instance.go
func GenerateInstanceName(providerName string) string {
    n, err := rand.Int(rand.Reader, big.NewInt(65536))
    if err != nil {
        // 极端情况，使用更可靠的备用方案
        n = big.NewInt(time.Now().UnixNano() % 65536)
    }
    randomStr := fmt.Sprintf("%04x", n.Int64())
    // ...
}
```

### H-03: 开发环境下验证码可被绕过

- **文件**: `service/auth/auth.go` (verifyCaptcha, loginWithPassword, ForgotPassword)
- **问题**: 开发模式下，验证码输入 "test" 即可通过。如果生产环境误配置为 `development`，验证码保护将完全失效
- **影响**: 暴力破解攻击、自动化注册攻击
- **修复建议**: 即使在开发环境也应使用真实验证码，仅增加日志输出
- **修复代码示例**:
```go
func (s *AuthService) verifyCaptcha(captchaId, code string) error {
    if captchaId == "" || code == "" {
        return errors.New("验证码参数不完整")
    }
    // 删除开发模式绕过
    // if global.APP_CONFIG.System.Env == "development" && code == "test" {
    //     return nil
    // }
    match := global.APP_CAPTCHA_STORE.Verify(captchaId, code, true)
    if !match {
        return errors.New("验证码错误或已过期")
    }
    return nil
}
```

### H-04: 信任所有代理（SetTrustedProxies(nil)）

- **文件**: `router/setup.go:37`
- **问题**: `Router.SetTrustedProxies(nil)` 信任所有代理的 `X-Forwarded-*` 头，攻击者可以伪造客户端 IP 来绕过 IP 限制
- **影响**: IP 限制策略（如 `iplimit-count`）可被绕过
- **修复建议**: 配置明确的信任代理列表
- **修复代码示例**:
```go
// 仅信任已知代理（如 Cloudflare、Nginx 反代）
trustedProxies := []string{
    "127.0.0.1",
    "::1",
}
// 如果使用 Cloudflare，可以添加 Cloudflare IP 段
if global.APP_CONFIG.System.Env == "production" {
    trustedProxies = append(trustedProxies,
        "173.245.48.0/20", "103.21.244.0/22", // Cloudflare IP ranges
    )
}
Router.SetTrustedProxies(trustedProxies)
```

### H-05: 明文密码自动哈希升级逻辑存在安全隐患

- **文件**: `service/auth/auth.go:125-142` (loginWithPassword)
- **问题**: 密码验证失败时，自动将明文密码视为正确并升级为哈希值。这意味着如果数据库中存储的是明文密码，任何输入的密码都会被"验证通过"并更新到数据库
- **影响**: 如果存在明文密码的用户，第一次登录时任何人输入任何密码都能成功登录并覆盖原密码
- **修复建议**: 删除自动升级逻辑，通过单独的数据迁移脚本处理
- **修复代码示例**:
```go
// 验证密码 - 不做自动升级
if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
    global.APP_LOG.Debug("用户密码验证失败", zap.String("username", utils.SanitizeUserInput(req.Username)))
    return nil, "", common.NewError(common.CodeInvalidCredentials)
}
```

### H-06: 环境配置为 development

- **文件**: `config.yaml` (`system.env: "development"`)
- **问题**: 生产配置文件中 `system.env` 设置为 `development`，这可能导致调试信息泄露、验证码绕过等安全功能被禁用
- **修复建议**: 生产环境必须设置为 `production`
```yaml
system:
  env: production
```

### H-07: 使用 MD5 哈希（弱加密原语）

- **文件**:
  - `service/images/image.go:173` - 文件完整性校验使用 MD5
  - `service/images/download.go:120` - 文件名生成使用 MD5
  - `provider/proxmox/image.go:415` - 镜像校验使用 MD5
- **问题**: MD5 存在碰撞漏洞，不应用于安全相关场景
- **修复建议**: 替换为 SHA256
- **修复代码示例**:
```go
import "crypto/sha256"

hash := sha256.New()
if _, err := io.Copy(hash, file); err != nil {
    return "", err
}
checksum := fmt.Sprintf("%x", hash.Sum(nil))
```

### H-08: 缺少 CSRF 防护

- **文件**: 全局 - 所有路由
- **问题**: API 使用 JWT 认证（Bearer Token），不依赖 Cookie，因此传统 CSRF 风险较低。但 JWT 也通过 `X-New-Token` 响应头返回刷新令牌，且 WebSocket 连接通过 URL query 参数传递 token
- **影响**: 如果前端错误地将 JWT 存储在 Cookie 中，将面临 CSRF 攻击。WebSocket token 通过 URL 传递可能被日志/Referer 头泄露
- **修复建议**: 
  1. 确保前端始终将 JWT 存储在 `localStorage` 或 `sessionStorage` 中，不使用 Cookie
  2. WebSocket token 应通过首条消息传递而非 URL 参数
  3. 对状态修改请求（POST/PUT/DELETE）添加 CSRF token 验证作为纵深防御

### H-09: WebSocket 升级器允许任意来源

- **文件**: `api/v1/user/user_ssh.go:18-21`
- **问题**: WebSocket upgrader 的 `CheckOrigin` 函数始终返回 `true`，允许任何来源的 WebSocket 连接
- **影响**: 跨站 WebSocket 劫持（CSWSH）风险
- **修复建议**: 验证请求来源
- **修复代码示例**:
```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        if origin == "" {
            return false
        }
        allowedOrigins := []string{
            global.APP_CONFIG.System.FrontendURL,
        }
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                return true
            }
        }
        return false
    },
}
```

---

## 🟡 Medium

### M-01: JWT Token ID 生成可预测

- **文件**: `utils/auth.go:87-88`
- **问题**: `generateTokenID()` 使用 `time.Now().UnixNano() + os.Getpid()` 生成，可被预测
- **影响**: Token ID 是 JWT 黑名单的唯一标识，可预测意味着可能被绕过黑名单机制
- **修复建议**: 使用 UUID
```go
import "github.com/google/uuid"

func generateTokenID() string {
    return uuid.New().String()
}
```

### M-02: JWT signing-key 配置为空字符串

- **文件**: `config.yaml` (`jwt.signing-key: ""`)
- **问题**: 配置文件中 JWT 签名密钥为空。虽然代码中有数据库密钥管理系统的兜底逻辑，但如果数据库密钥也未初始化，JWT 将使用空字符串作为密钥
- **修复建议**: 启动时验证 JWT 密钥有效性，如果为空则拒绝启动
```go
func ValidateJWTConfig() error {
    key := GetJWTKey()
    if key == "" || len(key) < 32 {
        return fmt.Errorf("JWT signing key is not configured or too short (min 32 chars)")
    }
    return nil
}
```

### M-03: 密码重置链接使用 HTTP

- **文件**: `service/auth/auth.go:657`
- **问题**: `resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken)` 使用 HTTP 协议
- **影响**: Token 在 HTTP 连接中以明文传输，可被中间人截获
- **修复建议**: 使用 HTTPS
```go
resetURL := fmt.Sprintf("%s/reset-password?token=%s", global.APP_CONFIG.System.FrontendURL, resetToken)
```

### M-04: 路径遍历风险 - 配置文件写入

- **文件**: `service/system/init.go:283`
- **问题**: `os.WriteFile(backupPath, configData, 0644)` 中 `backupPath` 由 `configPath + ".backup"` 拼接，gosec 检测到潜在路径遍历
- **修复建议**: 验证路径在允许的目录内
```go
backupPath := configPath + ".backup"
absBackup, _ := filepath.Abs(backupPath)
if !strings.HasPrefix(absBackup, filepath.Dir(configPath)) {
    return fmt.Errorf("invalid backup path")
}
```

### M-05: 文件权限过于宽松

- **文件**: 多处（gosec G304/G301/G306 共 37 处）
- **问题**: 多处文件操作使用 `0644`（全局可读）权限，日志文件使用 `0666`
- **修复建议**: 敏感文件使用 `0600`，配置文件使用 `0640`
```go
// 替换所有敏感文件操作
os.WriteFile(path, data, 0600)  // 仅所有者可读写
os.MkdirAll(dir, 0750)          // 目录
os.WriteFile(path, data, 0640)  // 配置文件
```

### M-06: 整数溢出转换

- **文件**: 
  - `service/resources/monitoring.go:201` - `uint64 -> int64`
  - `api/v1/system/performance.go:588` - `uint64 -> int`
  - `api/v1/user/wallet.go:77` - `uint -> int`
  - `api/v1/user/order.go:104` - `uint -> int`
- **问题**: 大数值转换时可能发生溢出，导致逻辑错误
- **修复建议**: 使用安全转换或保持原类型

### M-07: 数据库密钥明文存储在 system_configs 表

- **文件**: `service/auth/jwt.go`
- **问题**: JWT 密钥存储在数据库 `system_configs` 表中，以 JSON 明文形式保存。如果数据库被 SQL 注入或备份泄露，密钥将暴露
- **修复建议**: 使用应用层加密存储敏感配置值，或使用专门的密钥管理服务（如 HashiCorp Vault）

---

## 🔵 Low

### L-01: JWT Token 通过 WebSocket URL 参数传递

- **文件**: `middleware/auth.go:46-47`
- **问题**: Token 也从 query 参数获取（`c.Query("token")`），可能被浏览器历史、代理日志、Referer 头等泄露
- **修复建议**: WebSocket 应通过首条消息或 Sec-WebSocket-Protocol 头传递 token

### L-02: SQL 查询总体安全（GORM 参数化）

- **文件**: 全局
- **评估**: 所有 GORM 查询均使用了参数化占位符（`?`），未发现原始 SQL 拼接（`getNextKeyVersion` 中的 `Raw()` 也使用了占位符）
- **结论**: ✅ SQL 注入风险低

### L-03: XSS 风险评估

- **文件**: 全局
- **评估**: API 返回 JSON 数据，不由服务端渲染 HTML。文件上传有安全扫描和 MIME 类型验证
- **结论**: ✅ XSS 风险低（依赖前端正确处理）

### L-04: 密码存储使用 bcrypt

- **文件**: `service/auth/auth.go`
- **评估**: 使用 `bcrypt.DefaultCost`（cost=10），符合当前安全标准
- **建议**: 考虑使用 `bcrypt.MinCost + 2`（cost=12）以应对未来算力增长
```go
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
```

---

## ✅ 做得好的方面

1. **JWT 密钥轮换机制**：实现了数据库级别的密钥版本管理和自动轮换
2. **权限校验完整性**：所有路由组都正确应用了 `RequireAuth` 中间件，用户路由 `AuthLevelUser`，管理员路由 `AuthLevelAdmin`
3. **文件上传安全**：有白名单验证（MIME + 扩展名）、安全扫描、UUID 文件名、路径遍历防护
4. **密码策略**：强制 8 位以上，要求大小写字母+数字+特殊字符，禁止弱密码
5. **认证中间件设计**：服务端验证用户状态和权限，不依赖 JWT claims 中的用户类型
6. **JWT 黑名单**：支持 Token 撤销（通过 JTI）
7. **支付回调幂等性**：通过订单状态检查防止重复处理

---

## govulncheck 结果

govulncheck 因编译错误无法运行，发现以下编译问题：
- `api/v1/admin/product.go:238` - 语法错误
- `api/v1/admin/product.go` - 未定义函数 `fillProductResourcesFromLevelLimit`
- `api/v1/user/order.go` - `gorm` 重复导入
- `scripts/` 目录多个文件存在 `main` 函数重复声明

**建议**: 修复编译错误后再运行 govulncheck 进行依赖漏洞扫描。

---

## gosec 扫描统计

| 规则ID | 描述 | 数量 |
|--------|------|------|
| G104 | 未检查的错误返回值 | 219 |
| G304 | 文件路径由用户输入构造 | 14 |
| G301 | 文件创建权限问题 | 13 |
| G401 | 使用弱加密原语 (MD5) | 12 |
| G306 | 文件写入权限问题 | 10 |
| G106 | 不安全的 chmod | 10 |
| G501 | 导入弱加密原语 (MD5) | 9 |
| G118 | context 取消函数未调用 | 9 |
| G402 | TLS InsecureSkipVerify | 6 |
| G115 | 整数溢出转换 | 4 |
| G101 | 硬编码凭据 | 4 |
| G404 | 弱随机数生成器 | 2 |

---

## 修复优先级建议

1. **立即修复**（P0）: C-01, C-02, C-03, C-04, H-05
2. **尽快修复**（P1）: C-05, H-01, H-02, H-06, H-08, H-09
3. **计划修复**（P2）: M-01 ~ M-07, H-03, H-04, H-07
4. **优化改进**（P3）: L-01, L-04, G104 未检查错误

---

*审计工具: [gosec](https://github.com/securego/gosec), [govulncheck](https://golang.org/x/vuln/cmd/govulncheck), 人工审查*
