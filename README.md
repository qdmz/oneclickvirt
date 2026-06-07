# OneClickVirt

[English](./README_EN.md) | 中文

## 项目简介

OneClickVirt 是一个现代化的虚拟服务管理平台，支持多种虚拟化技术（Docker、LXD、Incus、Proxmox VE），提供完整的产品管理、用户管理、代理商管理、域名绑定、订单管理、实名认证、工单系统和资源监控功能。

## 快速开始

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

首次部署时，访问网站会自动跳转到初始化页面，按照提示设置管理员账户和普通用户账户。

初始化完成后：
- **管理员账户**: `admin / Admin123!@#`
- **普通用户**: 初始化时自定义的用户名和密码

> 首次登录后请立即修改密码！

## 功能特性

### 管理员功能
- 站点配置管理
- 产品套餐管理（支持库存量和销售计数）
- 兑换码管理
- 订单管理
- 用户管理（批量操作：批量删除、批量修改等级、批量修改状态）
- 虚拟实例管理（Docker/LXD/Incus/Proxmox）
- 资源监控与性能分析
- 流量统计与限速管理
- 系统镜像管理
- 公告管理
- 邀请码管理
- 端口映射管理
- OAuth2 第三方登录配置（QQ/Telegram）
- 多种支付方式（支付宝/微信/余额/易支付/码支付）
- **代用户登录** — 管理员可直接以用户身份登录
- **实例转移归属** — 管理员可将实例在用户间转移
- **代理商系统管理** — 代理商审核、佣金调整、子用户管理
- **域名绑定管理** — DNS 内部解析配置、Nginx 反代、用户域名配额
- **实名认证管理** — 查看认证记录、手动审核
- **工单管理** — 查看所有工单、回复用户、更新工单状态（开启/处理中/已解决/已关闭）
- **智简魔方集成** — 支持智简魔方 API 管理，可通过智简魔方创建和管理实例

### 代理商功能
- 代理商申请与入驻
- 子用户管理（创建/删除/批量操作）
- 佣金记录与结算
- 代理商钱包与提现
- 代理商仪表盘（数据统计）
- 推广链接与邀请码

### 用户功能
- 虚拟实例管理（创建/启停/删除/控制台）
- 产品购买与订单管理
- 钱包管理与充值
- 流量监控
- SSH Web 终端连接
- 端口映射查看
- **邮件注册激活** — 注册后邮箱验证激活
- **密码找回/重置** — 通过邮件重置密码
- **域名绑定** — 绑定自定义域名到虚拟机内部 IP:端口
- **实名认证** — 支付宝实名认证（姓名+身份证号）
- **工单系统** — 创建工单（咨询/故障/功能建议/投诉）、回复工单、查看处理进度
- **实例申请** — 快速申请创建虚拟机或容器实例
- **深色/浅色主题切换**

## 项目结构

```
oneclickvirt/
├── server/                        # Go 后端
│   ├── api/v1/
│   │   ├── admin/                 # 管理员 API
│   │   ├── user/                  # 用户 API（含容器/虚拟机管理）
│   │   ├── agent/                 # 代理商 API
│   │   ├── auth/                  # 认证 API
│   │   ├── config/                # 配置任务 API
│   │   ├── oauth2/                # OAuth2 API
│   │   ├── payment/               # 支付回调 API
│   │   ├── provider/              # 节点 API
│   │   ├── public/                # 公开 API
│   │   ├── system/                # 系统 API
│   │   └── traffic/               # 流量 API
│   ├── config/                    # 配置管理
│   ├── constant/                  # 常量定义
│   ├── core/                      # 核心组件（日志/采样/编码）
│   ├── docs/                      # Swagger 文档
│   ├── global/                    # 全局变量
│   ├── initialize/                # 初始化（数据库/路由/配置）
│   ├── middleware/                # 中间件（认证/权限/代理权限/限流/日志）
│   ├── model/                     # 数据模型
│   │   ├── admin/                 # 管理员/任务/流量监控模型
│   │   ├── agent/                 # 代理商模型
│   │   ├── api/                   # 通用 API 响应模型
│   │   ├── auth/                  # 认证/角色/权限模型
│   │   ├── common/                # 通用请求/响应/错误模型
│   │   ├── config/                # 配置模型
│   │   ├── dashboard/             # 仪表盘模型
│   │   ├── domain/                # 域名绑定/域名配置模型
│   │   ├── image/                 # 镜像请求模型
│   │   ├── kyc/                   # 实名认证模型
│   │   ├── monitoring/            # 监控/性能/流量历史模型
│   │   ├── oauth2/                # OAuth2 提供商模型
│   │   ├── order/                 # 订单模型
│   │   ├── permission/            # 权限模型
│   │   ├── product/               # 产品模型（含库存字段）
│   │   ├── provider/              # 节点/实例/端口模型
│   │   ├── redemption/            # 兑换码模型
│   │   ├── resource/              # 资源请求模型
│   │   ├── site/                  # 站点配置模型
│   │   ├── system/                # 系统健康/JWT/存储模型
│   │   ├── ticket/                # 工单/工单回复模型
│   │   ├── user/                  # 用户/API Key 模型
│   │   └── wallet/                # 钱包模型
│   ├── service/                   # 业务逻辑
│   │   ├── admin/                 # 管理员服务（实例/邀请/节点/用户）
│   │   ├── agent/                 # 代理商服务
│   │   ├── auth/                  # 认证服务（JWT/权限/角色/黑名单）
│   │   ├── cache/                 # 用户缓存服务
│   │   ├── config/                # 配置任务服务
│   │   ├── database/              # 数据库服务
│   │   ├── domain/                # 域名服务（DNS/Nginx）
│   │   ├── email/                 # 邮件服务（SMTP）
│   │   ├── images/                # 镜像下载/健康检查服务
│   │   ├── interfaces/            # 接口定义
│   │   ├── kyc/                   # 实名认证服务（支付宝 API）
│   │   ├── lifecycle/             # 生命周期管理
│   │   ├── log/                   # 日志导出/轮转服务
│   │   ├── oauth2/                # OAuth2 服务
│   │   ├── pmacct/                # 流量采集/聚合/清理服务
│   │   ├── provider/              # 节点服务（API/证书/配置）
│   │   ├── resources/             # 资源服务（配额/监控/同步/统计）
│   │   ├── scheduler/             # 调度服务（健康检查/维护/监控）
│   │   ├── storage/               # 存储服务
│   │   ├── system/                # 系统服务（清理/初始化/JWT 密钥轮换）
│   │   ├── task/                  # 任务服务（实例操作/端口映射/重置）
│   │   ├── ticket/                # 工单服务
│   │   ├── traffic/               # 流量服务（聚合/限速/历史/同步）
│   │   └── user/                  # 用户服务（实例/通知/资料/节点/资源）
│   ├── provider/                  # 虚拟化提供商（统一接口）
│   ├── router/                    # 路由定义
│   ├── scripts/                   # 数据库脚本
│   ├── source/                    # 数据种子
│   ├── utils/                     # 工具函数
│   ├── config.yaml                # 配置文件
│   ├── go.mod                     # Go 模块依赖
│   └── main.go                    # 入口
├── web/                           # Vue 3 前端
│   ├── src/
│   │   ├── api/                   # API 封装
│   │   ├── assets/                # 静态资源（图片/样式变量）
│   │   ├── components/            # 公共组件（SSH 终端/流量图表/主题切换）
│   │   ├── composables/           # 组合式函数
│   │   ├── i18n/                  # 国际化（中/英）
│   │   │   └── locales/
│   │   │       ├── zh-CN/         # 中文语言包
│   │   │       └── en-US/         # 英文语言包
│   │   ├── pinia/                 # 状态管理
│   │   ├── router/                # 路由配置
│   │   ├── style/                 # 全局样式 + 深色/浅色主题系统
│   │   ├── utils/                 # 工具函数
│   │   ├── view/
│   │   │   ├── admin/             # 管理员页面
│   │   │   │   ├── agents/            # 代理商管理
│   │   │   │   ├── announcements/     # 公告管理
│   │   │   │   ├── config/            # 站点配置
│   │   │   │   ├── dashboard/         # 管理仪表盘
│   │   │   │   ├── domain-config/     # 域名系统配置
│   │   │   │   ├── domains/           # 域名管理
│   │   │   │   ├── instances/          # 实例管理
│   │   │   │   ├── invite-codes/      # 邀请码管理
│   │   │   │   ├── kyc/               # 实名认证管理
│   │   │   │   ├── oauth2/            # OAuth2 配置
│   │   │   │   ├── orders/            # 订单管理
│   │   │   │   ├── performance/       # 性能监控
│   │   │   │   ├── portmapping/       # 端口映射管理
│   │   │   │   ├── products/          # 产品管理
│   │   │   │   ├── providers/         # 节点管理
│   │   │   │   ├── redemption-codes/  # 兑换码管理
│   │   │   │   ├── site-config/       # 站点配置
│   │   │   │   ├── system-images/     # 系统镜像管理
│   │   │   │   ├── tasks/             # 任务管理
│   │   │   │   ├── tickets/           # 工单管理
│   │   │   │   ├── traffic/           # 流量管理
│   │   │   │   └── users/             # 用户管理
│   │   │   ├── agent/             # 代理商页面（仪表盘/子用户/佣金/钱包/资料）
│   │   │   ├── user/              # 用户页面
│   │   │   │   ├── api-keys/          # API Key 管理
│   │   │   │   ├── apply/             # 实例申请
│   │   │   │   ├── dashboard/         # 用户仪表盘
│   │   │   │   ├── domains/           # 域名管理
│   │   │   │   ├── instances/         # 实例管理
│   │   │   │   ├── kyc/               # 实名认证
│   │   │   │   ├── orders/            # 订单管理
│   │   │   │   ├── profile/           # 个人资料
│   │   │   │   ├── purchase/          # 产品购买
│   │   │   │   ├── tasks/             # 任务管理
│   │   │   │   ├── tickets/           # 工单系统
│   │   │   │   └── wallet/            # 钱包管理
│   │   │   └── layout/            # 布局组件（侧边栏/导航栏）
│   │   ├── App.vue
│   │   └── main.js
│   └── package.json
├── scripts/
│   ├── init.sql                   # 数据库初始化脚本
│   ├── autosetup.sql              # 自动建表脚本
│   ├── fix_database.sql           # 数据库修复脚本
│   ├── missing_tables.sql         # 补充缺失表脚本
│   └── init.sh                    # 初始化脚本
├── deploy/                        # 部署配置
│   └── my.cnf/                    # MySQL 配置
├── .github/workflows/             # CI/CD
│   ├── build.yml                  # 构建流程
│   └── build_docker.yml           # Docker 构建流程
├── docker-compose.yaml            # Docker Compose 编排
├── docker-compose.web.yml         # 前端 Docker Compose
├── Dockerfile                     # Docker 构建文件
├── autoinstall.sh                 # 一键安装脚本
├── deploy.sh                      # 部署脚本
├── SECURITY_AUDIT.md              # 安全审计报告
└── LICENSE                        # MIT 许可证
```

## 技术栈

### 后端
| 技术 | 用途 |
|------|------|
| Go 1.25+ | 后端语言 |
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

## 数据库表结构

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
| `tickets` | 工单表 |
| `ticket_replies` | 工单回复表 |

## v1.530 更新日志

### Bug 修复
- 修复产品创建时价格单位不匹配问题（前端元 -> 后端分）
- 修复产品编辑时布尔值与整型不匹配问题
- 修复购买产品时订单创建顺序错误（交易记录引用未创建的订单ID）
- 修复充值订单 ProductID 非空约束导致创建失败
- 修复支付回调中 PaymentTime 字段引用错误
- 修复 Order 模型冗余字段（统一使用 ExpiredAt/PaidAt）

### 新功能
- 新增容器/虚拟机专用管理接口（列表、创建、控制、删除）
- 新增实例操作日志查询接口
- 新增用户昵称快速更新接口
- 新增域名管理后端路由

### 优化
- 统一所有模型 JSON 命名规范（snake_case -> camelCase）
- 统一软删除字段标签（json:"-"）
- 补充关键模型索引和约束
- 清理项目无用文件（75+个临时文件）

### 文件变更
- 新增: server/api/v1/user/container_vm.go
- 修改: server/model/order/order.go, server/api/v1/user/order.go, server/api/v1/user/recharge.go
- 修改: server/api/v1/payment/callback.go
- 修改: server/model/monitoring/*.go, server/model/domain/domain.go, server/model/kyc/kyc.go
- 修改: server/router/admin.go, server/router/user.go
- 修改: web/src/api/api-key.ts, web/src/api/user.js

## v1.520 更新日志

### 新增功能
- 完整工单系统 -- 用户可创建/查看/回复工单，管理员可管理所有工单、回复、更新状态（开启/处理中/已解决/已关闭）
- 实例申请页面 -- `/user/apply` 路由，支持快速申请创建虚拟机或容器实例
- 密码重置页面 -- `/reset-password` 路由，支持通过邮件链接重置密码

### 修复问题
- 修复初始化时 `AutoMigrateTables()` 缺失 20+ 个表的问题（OAuth2/Product/Wallet/Order/Ticket 等）
- 修复 `SeedSystemImages()` 被重复调用的冗余问题
- 修复普通用户初始化状态从禁用改为默认启用
- 统一 `MinLevelForVM` 配置值为 3（config.go 与 seed.go 一致）
- 修复 `/user/apply` 路由未注册导致 404 的问题
- 添加 Vite `allowedHosts` 配置支持 `*.monkeycode-ai.online` 域名

## 安全措施

- JWT Token 认证 + Token 黑名单 + 密钥轮换
- bcrypt 密码哈希（cost=12）
- 基于角色的权限控制（RBAC）
- CORS 白名单配置
- 参数化 SQL 查询（防注入）
- 文件上传白名单 + 安全扫描
- SSH TOFU 主机密钥验证
- WebSocket Origin 验证
- pprof 仅开发环境暴露
- 敏感配置环境变量注入
- 实名认证身份证号加密存储 + SHA256 查重

完整安全审计报告见 [SECURITY_AUDIT.md](./SECURITY_AUDIT.md)

## 配置说明

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

## 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

[MIT License](./LICENSE)

## 联系方式

如有问题或建议，请创建 [Issue](https://github.com/qdmz/oneclickvirt/issues)。

---

**如果觉得有用，请给个 Star！**
