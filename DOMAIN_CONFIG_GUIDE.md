# 域名配置修改指南

## 需要修改的配置项

### 1. 后端配置文件 (config.yaml)

#### 修改 frontend-url
```yaml
# 修改前
system:
    frontend-url: "http://localhost:8080"

# 修改后
system:
    frontend-url: "https://oneclickvirt.ypvps.com"
```

#### 检查支付回调地址
```yaml
# 易支付回调地址（已经是正确的）
payment:
    epay-notify-url: https://oneclickvirt.ypvps.com/api/v1/payment/epay/notify
    epay-return-url: https://oneclickvirt.ypvps.com/user/wallet

# 马支付回调地址（可能需要修改）
payment:
    mapay-notify-url: https://heyun.ypvps.com/api/v1/payment/mapay/notify
    mapay-return-url: https://heyun.ypvps.com/user/wallet

# 实名认证回调地址（已经是正确的）
payment:
    real-name-callback-url: https://oneclickvirt.ypvps.com/api/v1/kyc/callback
```

### 2. 后端代码修改

#### 修改 yipay.go 中的硬编码 URL

**文件位置**: `/root/.openclaw/workspace/oneclickvirt/server/api/v1/payment/yipay.go`

**需要修改的地方**:

```go
// 修改前
req := &CreateOrderRequest{
    OutTradeNo: fmt.Sprintf("TEST%d", time.Now().Unix()),
    Subject:    "测试商品",
    TotalFee:   0.01,
    NotifyURL:  "http://localhost:8890/api/v1/payment/yipay/notify",
    ReturnURL:  "http://localhost:8890/api/v1/payment/yipay/return",
    Param:      "test",
}

// 修改后
req := &CreateOrderRequest{
    OutTradeNo: fmt.Sprintf("TEST%d", time.Now().Unix()),
    Subject:    "测试商品",
    TotalFee:   0.01,
    NotifyURL:  global.APP_CONFIG.System.FrontendURL + "/api/v1/payment/yipay/notify",
    ReturnURL:  global.APP_CONFIG.System.FrontendURL + "/api/v1/payment/yipay/return",
    Param:      "test",
}
```

同样需要修改 CreateYipayOrder 函数中的默认回调地址：

```go
// 修改前
if req.NotifyURL == "" {
    req.NotifyURL = "http://localhost:8890/api/v1/payment/yipay/notify"
}
if req.ReturnURL == "" {
    req.ReturnURL = "http://localhost:8890/api/v1/payment/yipay/return"
}

// 修改后
if req.NotifyURL == "" {
    req.NotifyURL = global.APP_CONFIG.System.FrontendURL + "/api/v1/payment/yipay/notify"
}
if req.ReturnURL == "" {
    req.ReturnURL = global.APP_CONFIG.System.FrontendURL + "/api/v1/payment/yipay/return"
}
```

### 3. 前端配置

前端代码已经使用了相对路径 `/api`，不需要修改。前端会通过 nginx 代理到后端 API。

### 4. Nginx 配置

确保 nginx 配置正确代理 API 请求：

```nginx
server {
    listen 80;
    server_name oneclickvirt.ypvps.com;

    # 前端静态文件
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://oneclickvirt-api:8890/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 5. CORS 配置

如果需要跨域访问，需要修改 CORS 配置：

```yaml
# config.yaml
cors:
    mode: "allow-all"
    whitelist:
        - "https://oneclickvirt.ypvps.com"
        - "http://oneclickvirt.ypvps.com"
```

## 修改步骤

1. 修改 config.yaml 中的 frontend-url
2. 修改 yipay.go 中的硬编码 URL
3. 重新构建后端镜像
4. 重启后端容器
5. 测试域名访问

## 测试清单

- [ ] 前端页面可以正常访问
- [ ] API 请求正常工作
- [ ] 支付回调地址正确
- [ ] 邮件发送功能正常
- [ ] 实名认证回调正常

## 注意事项

1. 确保域名 DNS 解析正确
2. 确保防火墙允许 80 和 443 端口访问
3. 如果使用 HTTPS，需要配置 SSL 证书
4. 修改配置后需要重启服务才能生效
