# OneClickVirt

[English](./README_EN.md) | 中文

## 📖 项目简介

OneClickVirt 是一个现代化的虚拟服务管理平台，支持多种虚拟化技术（Docker、LXD、Incus、Proxmox VE），提供完整的产品管理、用户管理、代理商管理、域名绑定、订单管理、实名认证和资源监控功能。

## 🚀 快速开始

### 本地开发

#### 前端开发
```bash
cd web
npm install
npm run dev
```
前端将运行在 `http://localhost:8080`

#### 后端开发
```bash
cd server
go mod download
go run main.go
```
后端将运行在 `http://localhost:8890`

### Docker 部署（推荐）

```bash
docker-compose -f docker-compose.yaml up -d
```

#### 数据初始化

首次部署时，访问网站会自动跳转到初始化页面，按照提示设置管理员账户。

也可以手动初始化：
```bash
docker exec -i oneclickvirt mysql -uroot oneclickvirt < scripts/init.sql
```

默认管理员账户：
- **用户名**: `admin`
- **密码**: `admin123456`

> ⚠️ 首次登录后请立即修改密码！

## ✨ 功能特性

### 管理员功能
- ✅ 站点配置管理
- ✅ 产品套餐管理（支持库存量和销售计数）
- ✅ 兑换码管理
- ✅ 订单管理
- ✅ 用户管理（批量操作：批量删除、批量修改等级、批量修改状态）
- ✅ 虚拟实例管理（Docker/LXD/Incus/Proxmox）
- ✅ 资源监控与性能分析
- ✅ 流量统计与限速管理
- ✅ 系统镜像管理
- ✅ 公告管理
- ✅ 邀请码管理
- ✅ 端口映射管理
- ✅ OAuth2 第三方登录配置（QQ/Telegram）
- ✅ 多种支付方式（支付宝/微信/余额/易支付/码支付）
- ✅ **代用户登录** — 管理员可直接以用户身份登录
- ✅ **实例转移归属** — 管理员可将实例在用户间转移
- ✅ **代理商系统管理** — 代理商审核、佣金调整、子用户管理
- ✅ **域名绑定管理** — DNS 内部解析配置、Nginx 反代、用户域名配额
- ✅ **实名认证管理** — 查看认证记录、手动审核
- ✅ **智简魔方集成** — 支持智简魔方 API 管理，可通过智简魔方创建和管理实例

### 代理商功能
- ✅ 代理商申请与入驻
- ✅ 子用户管理（创建/删除/批量操作）
- ✅ 佣金记录与结算
- ✅ 代理商钱包与提现
- ✅ 代理商仪表盘（数据统计）
- ✅ 推广链接与邀请码

### 用户功能
- ✅ 虚拟实例管理（创建/启停/删除/控制台）
- ✅ 产品购买与订单管理
- ✅ 钱包管理与充值
- ✅ 流量监控
- ✅ SSH Web 终端连接
- ✅ 端口映射查看
- ✅ **邮件注册激活** — 注册后邮箱验证激活
- ✅ **密码找回** — 通过邮件重置密码
- ✅ **域名绑定** — 绑定自定义域名到虚拟机内部 IP:端口
- ✅ **实名认证** — 支付宝实名认证（姓名+身份证号）
- ✅ **深色/浅色主题切换**

## 📁 项目结构

```
oneclickvirt/
├── server/                    # Go 后端
│   ├── api/v1/
│   │   ├── admin/             # 管理员 API
│   │   ├── user/              # 用户 API
│   │   ├── agent/             # 代理商 API
│   │   ├── payment/           # 支付回调 API
│   │   ├── public/            # 公开 API
│   │   └── system/            # 系统 API
│   ├── config/                # 配置管理
│   ├── middleware/            # 中间件（认证/权限/代理权限）
│   ├── model/                 # 数据模型
│   │   ├── user/              # 用户模型
│   │   ├── product/           # 产品模型（含库存字段）
│   │   ├── agent/             # 代理商/子用户/佣金模型
│   │   ├── domain/            # 域名绑定/域名配置模型
│   │   ├── kyc/               # 实名认证模型
│   │   ├── auth/              # 认证/角色模型
│   │   ├── order/             # 订单模型
│   │   ├── wallet/            # 钱包模型
│   │   └── provider/          # 节点/实例/端口模型
│   ├── service/               # 业务逻辑
│   │   ├── auth/              # 认证服务
│   │   ├── agent/             # 代理商服务
│   │   ├── domain/            # 域名服务（DNS/Nginx）
│   │   ├── email/             # 邮件服务（SMTP）
│   │   └── kyc/               # 实名认证服务（支付宝 API）
│   ├── provider/              # 虚拟化提供商
│   │   ├── docker/            # Docker 支持
│   │   ├── lxd/               # LXD 支持
│   │   ├── incus/             # Incus 支持
│   │   ├── proxmox/           # Proxmox VE 支持
│   │   └── zjmf/              # 智简魔方支持
│   ├── router/                # 路由定义
│   ├── initialize/            # 初始化（数据库/路由）
│   ├── utils/                 # 工具函数
│   ├── config.yaml            # 配置文件
│   └── main.go               # 入口
├── web/                       # Vue 3 前端
│   ├── src/
│   │   ├── api/               # API 封装
│   │   ├── view/
│   │   │   ├── admin/         # 管理员页面
│   │   │   │   ├── agents/        # 代理商管理
│   │   │   │   ├── domains/       # 域名管理
│   │   │   │   ├── kyc/           # 实名认证管理
│   │   │   │   ├── products/      # 产品管理
│   │   │   │   ├── users/         # 用户管理
│   │   │   │   └── ...
│   │   │   ├── agent/         # 代理商页面（仪表盘/子用户/佣金/钱包/资料）
│   │   │   └── user/          # 用户页面（实例/订单/钱包/域名/实名认证）
│   │   ├── components/        # 公共组件
│   │   ├── style/             # 全局样式 + 深色/浅色主题系统
│   │   ├── pinia/             # 状态管理
│   │   ├── router/            # 路由配置
│   │   └── i18n/              # 国际化（中/英）
│   └── package.json
├── scripts/
│   ├── init.sql               # 数据库初始化脚本
│   └── init.sh                # 初始化脚本
├── docker-compose.yaml        # Docker Compose 编排
├── Dockerfile                 # Docker 构建文件
└── SECURITY_AUDIT.md          # 安全审计报告
```

## 🔧 技术栈

### 后端
| 技术 | 用途 |
|------|------|
| Go 1.24+ | 后端语言 |
| Gin | Web 框架 |
| GORM | ORM 框架 |
| MySQL / SQLite | 数据库 |
| JWT | 身份认证 |
| net/smtp | 邮件服务 |
| crypto/rsa | 支付宝签名（RSA2） |
| WebSocket | SSH 终端 |

### 前端
| 技术 | 用途 |
|------|------|
| Vue 3 | 前端框架 |
| Element Plus | UI 组件库 |
| Vite | 构建工具 |
| Pinia | 状态管理 |
| vue-i18n | 国际化 |
| ECharts | 图表库 |
| xterm.js | Web SSH 终端 |
| SCSS + CSS Variables | 主题系统 |

## 📊 数据库表结构

### 核心表（自动创建）
| 表名 | 说明 |
|------|------|
| `users` | 用户表（含 email_verified, real_name_verified） |
| `roles` | 角色表 |
| `user_roles` | 用户角色关联 |
| `products` | 产品表（含 stock, sold_count） |
| `product_purchases` | 产品购买记录 |
| `orders` | 订单表 |
| `payment_records` | 支付记录 |
| `instances` | 虚拟实例 |
| `providers` | 节点/提供商 |
| `ports` | 端口映射 |
| `user_wallets` | 用户钱包 |
| `wallet_transactions` | 钱包交易记录 |

### 新增功能表（自动创建）
| 表名 | 说明 |
|------|------|
| `agents` | 代理商表 |
| `sub_user_relations` | 代理商-子用户关联 |
| `commissions` | 佣金记录 |
| `domains` | 域名绑定 |
| `domain_configs` | 域名系统配置 |
| `kyc_records` | 实名认证记录 |

## 🔒 安全措施

- ✅ JWT Token 认证 + Token 黑名单 + 密钥轮换
- ✅ bcrypt 密码哈希（cost=12）
- ✅ 基于角色的权限控制（RBAC）
- ✅ CORS 白名单配置
- ✅ 参数化 SQL 查询（防注入）
- ✅ 文件上传白名单 + 安全扫描
- ✅ SSH TOFU 主机密钥验证
- ✅ WebSocket Origin 验证
- ✅ pprof 仅开发环境暴露
- ✅ 敏感配置环境变量注入
- ✅ 实名认证身份证号加密存储 + SHA256 查重

完整安全审计报告见 [SECURITY_AUDIT.md](./SECURITY_AUDIT.md)

## ⚙️ 配置说明

### 环境变量（敏感信息）

```bash
# 邮件 SMTP
export EMAIL_PASSWORD="your_smtp_password"

# 支付宝
export ALIPRAY_APP_ID="your_app_id"
export ALIPAY_PRIVATE_KEY="your_private_key"
export ALIPAY_PUBLIC_KEY="alipay_public_key"

# 微信支付
export WECHAT_API_KEY="your_api_key"
export WECHAT_API_V3_KEY="your_v3_key"
export WECHAT_APP_ID="your_app_id"

# 第三方支付
export EPAY_KEY="your_epay_key"
export MAPAY_KEY="your_mapay_key"

# OAuth2
export TELEGRAM_BOT_TOKEN="your_bot_token"
export QQ_APP_ID="your_qq_app_id"
export QQ_APP_KEY="your_qq_app_key"
```

### config.yaml 关键配置

```yaml
system:
  env: production                    # 生产环境设为 production
  frontend-url: "https://your-domain.com"

auth:
  enable-email: true
  enable-email-verification: true    # 邮箱注册激活
  email-activation-expire-hours: 24
  enable-public-registration: true

payment:
  enable-alipay: true
  enable-wechat: true
  enable-epay: true
  enable-mapay: true
  enable-real-name: false            # 实名认证开关
  require-real-name: false           # 强制实名
  real-name-callback-url: "https://your-domain.com/api/v1/kyc/callback"
```

## 🤝 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

[MIT License](./LICENSE)

## 📞 联系方式

如有问题或建议，请创建 [Issue](https://github.com/qdmz/oneclickvirt/issues)。

---

**⭐ 如果觉得有用，请给个 Star！**
