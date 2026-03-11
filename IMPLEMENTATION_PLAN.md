# OneClickVirt 功能扩展实施方案

## 项目现状分析

### 技术栈
- **后端**: Go 1.24 + Gin + GORM + MySQL/SQLite
- **前端**: Vue 3 + Element Plus + Vite + Pinia + vue-i18n + ECharts
- **部署**: Docker + docker-compose
- **已有认证**: JWT + 邮件SMTP + OAuth2(QQ/Telegram) + 验证码

### 已有基础
- ✅ 用户CRUD、批量操作（批量删/批量改等级/批量改状态）
- ✅ 邮件SMTP配置（config.yaml已配置）
- ✅ VerifyCode/PasswordReset 模型已存在
- ✅ 端口映射管理（手动添加端口段）
- ✅ 产品管理（CRUD、启用/禁用）
- ✅ 用户等级体系（5级配额限制）
- ✅ 忘记密码页面 (`forgot-password/index.vue`)
- ✅ 重置密码页面 (`reset-password/index.vue`)
- ✅ 角色系统 (admin/user)

### 需要新增/扩展的功能
1. 代理商系统（全新）
2. 邮件注册激活（需完善）
3. 找回密码邮件发送（需完善）
4. 产品库存量（新增字段）
5. 域名绑定系统（全新）
6. 代码安全审计
7. 前端UI美化（大改）

---

## 实施计划

### Phase 1: 代理商系统 (Agent System)

#### 1.1 数据模型
**新增 `server/model/agent/agent.go`**
```go
type Agent struct {
    ID           uint
    UserID       uint           // 关联用户
    Code         string         // 代理商唯一代码（用于注册推广链接）
    Name         string         // 代理商名称/公司名
    ContactName  string
    ContactEmail string
    ContactPhone string
    CommissionRate float64      // 佣金比例 0-100%
    MaxSubUsers  int            // 最大子用户数
    MaxDomainsPerUser int       // 每个子用户可绑域名数
    Status       int            // 0=待审核 1=正常 2=禁用
    Balance      int64          // 代理商余额（分）
    CreatedAt/UpdatedAt
}
```

**新增 `server/model/agent/sub_user.go`**
```go
type SubUserRelation struct {
    ID        uint
    AgentID   uint           // 代理商ID
    UserID    uint           // 子用户ID
    CreatedAt time.Time
}
```

#### 1.2 后端路由 & API
**新增 `server/router/agent.go`** + `server/api/v1/agent/` + `server/service/agent/`

**代理商端路由：**
- `GET /v1/agent/profile` - 代理商信息
- `PUT /v1/agent/profile` - 更新资料
- `GET /v1/agent/sub-users` - 子用户列表
- `POST /v1/agent/sub-users` - 创建子用户（可选）
- `PUT /v1/agent/sub-users/:id` - 修改子用户
- `DELETE /v1/agent/sub-users/:id` - 删除子用户
- `PUT /v1/agent/sub-users/:id/status` - 启用/禁用子用户
- `PUT /v1/agent/sub-users/batch-status` - 批量操作
- `POST /v1/agent/sub-users/batch-delete` - 批量删除
- `GET /v1/agent/statistics` - 统计（子用户数、收益等）
- `GET /v1/agent/commissions` - 佣金记录
- `GET /v1/agent/wallet` - 代理商钱包

**管理员路由（代理商管理）：**
- `GET /v1/admin/agents` - 代理商列表
- `POST /v1/admin/agents` - 创建代理商
- `PUT /v1/admin/agents/:id` - 更新代理商
- `DELETE /v1/admin/agents/:id` - 删除代理商
- `PUT /v1/admin/agents/:id/status` - 审核代理商
- `PUT /v1/admin/agents/:id/commission` - 调整佣金比例
- `GET /v1/admin/agents/:id/sub-users` - 查看代理商子用户

#### 1.3 中间件
**新增 `server/middleware/agent.go`**
- `RequireAgent()` - 要求代理商权限
- `CheckAgentSubUser()` - 管理员可管理代理商的子用户

---

### Phase 2: 邮件注册激活 & 找回密码

#### 2.1 注册激活
**修改 `server/model/user/user.go`**
- User 模型新增 `EmailVerified bool` 字段

**新增/修改邮件服务 `server/service/email/email.go`**
```go
type EmailService struct {
    smtpHost string
    smtpPort int
    username string
    password string
}

func (s *EmailService) SendVerificationCode(target, code string) error  // 注册验证码
func (s *EmailService) SendPasswordReset(target, token string) error    // 密码重置链接
func (s *EmailService) SendActivationEmail(target, token string) error  // 激活链接
func (s *EmailService) SendWelcomeEmail(target string) error            // 欢迎邮件
```

**修改注册流程：**
1. 用户注册 → 发送激活邮件（含激活token）
2. 用户点击邮件链接 → `GET /api/v1/auth/verify-email?token=xxx`
3. 验证通过 → `EmailVerified = true`，允许登录
4. 未验证 → 登录时提示需要验证邮箱

**修改 `config.yaml`：**
```yaml
auth:
    enable-email-verification: true    # 新增
    email-activation-expire: 24h       # 激活链接过期时间
```

#### 2.2 找回密码
**修改 `server/router/public.go`：**
- `POST /api/v1/auth/forgot-password` - 发送重置密码邮件
- `POST /api/v1/auth/reset-password` - 使用token重置密码
- `GET /api/v1/auth/verify-reset-token` - 验证重置token有效性

---

### Phase 3: 产品库存管理

#### 3.1 数据模型修改
**修改 `server/model/product/product.go`：**
```go
type Product struct {
    // ... 原有字段
    Stock        int       `json:"stock" gorm:"default:-1;comment:库存量(-1表示无限)"`
    SoldCount    int       `json:"soldCount" gorm:"default:0;comment:已售数量"`
    IsEnabled    int       `json:"isEnabled" gorm:"column:is_enabled;default:1"`
}
```

#### 3.2 业务逻辑
- 购买时检查库存：`Stock == -1 || Stock > 0`
- 购买成功后：`Stock--`, `SoldCount++`
- 库存为0时自动禁用产品（可选配置）
- 管理员可手动调整库存

#### 3.3 API
- `PUT /v1/admin/products/:id/stock` - 调整库存
- 产品列表返回库存信息
- 前端显示库存状态（充足/紧张/售罄）

---

### Phase 4: 域名绑定系统

#### 4.1 数据模型
**新增 `server/model/domain/domain.go`**
```go
type Domain struct {
    ID          uint
    UserID      uint           // 所属用户
    InstanceID  uint           // 绑定的实例ID
    Domain      string         // 域名
    Protocol    string         // http/https/tcp/udp
    InternalIP  string         // 虚拟机内部IP
    InternalPort int           // 虚拟机内部端口
    ExternalPort int           // 外部映射端口
    SSL         bool           // 是否启用SSL
    Status      int            // 0=待配置 1=正常 2=异常
    AgentID     *uint          // 代理商ID（如果是代理商绑定的）
    ExpiresAt   *time.Time     // 到期时间
    CreatedAt/UpdatedAt
}

type DomainConfig struct {
    ID                   uint
    MaxDomainsPerUser    int       // 每用户最大域名数
    DefaultTTL           int       // 默认DNS TTL
    AutoSSL              bool      // 自动SSL
    AllowedDomainSuffixes string   // 允许绑定的域名后缀（逗号分隔，空=不限）
}
```

#### 4.2 后端服务
**新增 `server/service/domain/`**
```go
// 域名解析管理（操作服务器上的DNS配置）
func ConfigureInternalDNS(domain, internalIP string, port int) error
func RemoveInternalDNS(domain string) error
func GenerateDNSConfig() error  // 重新生成整个DNS配置

// Nginx/反代配置
func ConfigureReverseProxy(domain, internalIP string, internalPort, externalPort int) error
func RemoveReverseProxy(domain string) error
```

#### 4.3 DNS内部解析实现
- 使用 dnsmasq 或 CoreDNS 做内部DNS
- 自动生成配置文件并 reload
- 或者直接操作 hosts 文件（简单方案）

#### 4.4 API
**用户端：**
- `GET /v1/user/domains` - 我的域名列表
- `POST /v1/user/domains` - 绑定域名
- `PUT /v1/user/domains/:id` - 修改域名绑定
- `DELETE /v1/user/domains/:id` - 解绑域名
- `GET /v1/user/domains/available` - 可绑域名配额

**代理商端：**
- `GET /v1/agent/domains` - 名下用户域名
- `PUT /v1/agent/domains/batch` - 批量管理

**管理员端：**
- `GET /v1/admin/domains` - 所有域名
- `GET /v1/admin/domain-config` - 域名配置
- `PUT /v1/admin/domain-config` - 更新配置
- `DELETE /v1/admin/domains/:id` - 删除域名
- `POST /v1/admin/domains/sync-dns` - 手动同步DNS

---

### Phase 5: 代码安全审计

#### 5.1 Go安全扫描
- 使用 `gosec` 扫描后端代码
- 使用 `govulncheck` 检查依赖漏洞
- SQL注入检查（GORM参数化查询审查）
- XSS/CSRF防护检查
- JWT安全配置检查
- 密码存储检查（bcrypt审查）
- 文件上传安全检查
- SSH连接安全检查
- API权限校验完整性

#### 5.2 前端安全
- 依赖漏洞扫描（npm audit）
- XSS防护审查
- 敏感信息泄露检查（API keys等）
- CORS配置审查

#### 5.3 配置安全
- 检查默认密码/密钥
- 检查敏感信息硬编码
- 检查 DEBUG 模式
- 检查 CORS 配置

---

### Phase 6: 前端UI美化

#### 6.1 整体设计方向
- **风格**: 现代深色/浅色双主题
- **配色方案**:
  - 主色: `#6366F1` (Indigo) / `#8B5CF6` (Violet)
  - 辅色: `#06B6D4` (Cyan) / `#10B981` (Emerald)
  - 背景: `#0F172A` → `#1E293B` (Slate 渐变)
  - 文字: `#F8FAFC` (Light) / `#94A3B8` (Muted)
- **组件库**: 升级 Element Plus 主题 + 自定义组件
- **布局**: 优化侧边栏、仪表盘卡片、表格样式

#### 6.2 具体改动
1. **全局主题变量** - 重写 CSS 变量，支持深色/浅色切换
2. **登录/注册页** - 全新设计，背景动画，玻璃态卡片
3. **仪表盘** - 重构卡片样式，渐变色，动画效果
4. **侧边栏** - 收缩动画，图标优化，渐变高亮
5. **表格** - 行悬浮效果，斑马纹优化，状态标签美化
6. **按钮/表单** - 圆角优化，过渡动画，Loading效果
7. **新增页面模板选择器** - 管理员可切换2-3套预设主题

#### 6.3 新增代理商端页面
- `web/src/view/agent/dashboard/index.vue` - 代理商仪表盘
- `web/src/view/agent/sub-users/index.vue` - 子用户管理
- `web/src/view/agent/domains/index.vue` - 域名管理
- `web/src/view/agent/wallet/index.vue` - 代理商钱包
- `web/src/view/agent/profile/index.vue` - 代理商资料

#### 6.4 新增域名管理页面
- `web/src/view/user/domains/index.vue` - 用户域名管理
- `web/src/view/admin/domains/index.vue` - 管理员域名管理
- `web/src/view/admin/domain-config/index.vue` - 域名配置

---

## 开发顺序建议

1. **Phase 5 (安全审计)** — 先审计现有代码，发现问题
2. **Phase 3 (产品库存)** — 最简单，先热身
3. **Phase 2 (邮件验证)** — 完善现有功能
4. **Phase 1 (代理商系统)** — 核心新功能，工作量最大
5. **Phase 4 (域名绑定)** — 依赖代理商系统的域名配额
6. **Phase 6 (前端美化)** — 最后做，所有功能确定后再美化

---

## 环境搭建

### Windows 本地开发环境
```bash
# 后端
cd server
go mod download
# 修改 config.yaml 中的数据库连接
go run main.go

# 前端
cd web
npm install
npm run dev
```

### Docker 部署
```bash
docker-compose -f docker-compose.yaml up -d
```

---

## 风险与注意事项

1. **代理商系统** - 涉及佣金结算，需要严谨的幂等性设计
2. **域名绑定** - 需要服务器 root 权限操作 DNS/Nginx
3. **邮件发送** - 需要确保 SMTP 配置正确，部分邮箱有频率限制
4. **安全审计** - 发现的漏洞需优先修复
5. **前端美化** - 不能影响现有功能，需逐步替换
