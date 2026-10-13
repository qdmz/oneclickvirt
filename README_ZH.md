# OneHost 云计算平台

> 基于 [OneClickVirt](https://github.com/oneclickvirt/oneclickvirt) 的定制发行版，面向 **IDC / 云主机售卖** 场景。
> 在原有统一虚拟化管理能力之上，新增 **易支付、代金券、余额钱包、商城售卖** 等商业化闭环能力，
> 并针对低配 VPS（1 核 1G 级别）做了完整的部署优化。

[![Build and Release](https://github.com/qdmz/onehost/actions/workflows/build.yml/badge.svg)](https://github.com/qdmz/onehost/actions/workflows/build.yml)

[![Build and Push Docker Images](https://github.com/qdmz/onehost/actions/workflows/build_docker.yml/badge.svg)](https://github.com/qdmz/onehost/actions/workflows/build_docker.yml)

[![Integration Tests](https://github.com/qdmz/onehost/actions/workflows/integration-tests.yml/badge.svg)](https://github.com/qdmz/onehost/actions/workflows/integration-tests.yml)

## 项目简介

OneHost 是一套开箱即用的云计算售卖与管理平台。后端统一调度 LXD、Incus、Docker、Podman、Containerd、
Proxmox VE、QEMU/KVM、KubeVirt 等虚拟化后端；前端同时提供**用户自助控制台**与**功能完整的管理后台**。

它解决的是「把虚拟化资源当成商品卖出去」的完整链路：

```
产品定价 → 用户下单 → 在线支付 → 自动开通 → 自助运维 → 到期续费 → 流量/快照/重装 → 工单支持
```

**与上游 OneClickVirt 的主要差异**：上游是通用的虚拟化资源管理平台，本 fork 补齐了商业化售卖闭环
（产品、库存、限购、订单、支付、代金券、余额、续费），并针对国内低配服务器与国内支付网关做了适配。
详见 [本 fork 定制功能](#本-fork-定制功能)。

前端控制台基于 Vue 3、Vite 和 Element Plus 构建，已针对桌面、平板、安卓尺寸和 iOS 尺寸视口进行响应式布局检查，
界面支持中英文双语切换。v0.23 起采用靛蓝→天蓝设计系统，亮色 / 暗色双主题可一键切换（跟随系统并记住选择），
首页、登录注册、管理员与用户仪表盘均已重做，图表随主题自动换色。

## 在线演示

- 线上站点：<https://webhost.ypvps.com>（当前运行 v0.25.1，Docker Compose 部署于 1 核 1G 级 Debian 12 VPS）
- 首页、商城与注册页无需登录即可浏览；完整截图见 [演示截图](#演示截图)。

## 语言

- **简体中文**：当前文档（`README.md` 与 `README_ZH.md` 内容一致）
- **English**：暂未同步，可参考[上游 OneClickVirt 英文文档](https://github.com/oneclickvirt/oneclickvirt#readme)

平台界面本身支持中英文双语切换（基于 vue-i18n），可在后台「系统配置」中调整默认语言。

## 功能总览

平台分为**用户端**与**管理后台**两套界面，共 40+ 功能模块、1700+ 个 API 接口、62 张业务数据表。

| 能力域 | 用户端 | 管理后台 |
|--------|--------|----------|
| 实例（云主机/容器） | 创建、开关机、重装、VNC、续费、快照、端口映射 | 全局实例管理、批量操作、机型权限、配额 |
| 资源节点 | 查看可用节点与规格 | 接入/配置多虚拟化节点、IPv4 池、系统镜像 |
| 交易 | 商城选购、下单、支付、续费 | 产品定价、库存、订单管理、代金券、兑换码 |
| 资金 | 余额、充值、代金券兑换、流水 | 易支付配置、余额调整、充值订单 |
| 用户 | 注册登录、实名认证、邀请码 | 用户管理、用户组、KYC 审核、权限 |
| 网络 | 端口映射、域名解析 | 端口池、防火墙规则、域名配置 |
| 运维 | 任务中心、流量查看 | 流量监控、性能监控、日志审计、快照计划 |
| 运营 | 签到、工单、公告浏览 | 公告发布、工单处理、邀新、站点配置 |

## 用户端功能

### 实例管理
- **开通实例**：选择节点、系统镜像、机型规格后一键开通，任务进度实时可查
- **生命周期操作**：开机、关机、重启、强制停止、重装系统、重置密码
- **远程接入**：Web VNC 控制台（noVNC）、WebSSH 终端（xterm.js）
- **快照**：创建/恢复/删除快照，界面显示允许数量配额，超限自动禁用新建
- **续费**：实例到期前自助续费（本 fork 新增入口）
- **端口映射**：自助添加/删除端口映射规则
- **共享链接**：生成临时分享链接，便于协作或技术支持
- **流量与监控**：实时流量曲线、历史流量、资源占用（CPU/内存/磁盘）

### 商城与订单
- 按机型、节点、价格浏览在售产品，列表直接显示**库存余量**
- 下单购买，受后台设置的**每人限购数量**约束
- 订单详情完整展示规格、周期、金额、支付状态与开通结果
- 支持 0 元免费产品（后台可配置）

### 钱包与支付（本 fork 新增）
- 账户余额查看与充值，对接易支付网关（支付宝/微信）
- **代金券兑换**：输入券码充值到余额
- 余额流水明细，每笔变动可追溯
- 充值订单记录与状态查询

### 账号与运营
- 注册/登录（支持邮箱验证码、OAuth2 第三方登录）
- **实名认证**（KYC）：提交证件资料，等待后台审核
- **每日签到**：领取签到奖励
- **工单系统**：提交问题、上传截图、跟踪处理进度
- **邀请码**：使用邀请码注册，查看邀请记录
- **API 令牌**：生成令牌供程序化调用开放接口
- **个人中心**：修改资料、修改密码、绑定邮箱

## 管理后台功能

| 模块 | 说明 |
|------|------|
| 仪表盘 | 平台总览：实例数、用户数、收入、节点健康度 |
| 用户管理 | 用户列表、启停、余额调整、用户组、权限分配 |
| 用户组 | 分组管理与差异化配额、机型权限 |
| 实例管理 | 全局实例查看与管控、批量操作、强制删除 |
| 节点管理 | 接入 LXD/Incus/Docker/PVE/QEMU 等节点，配置 IPv4 池 |
| 系统镜像 | 镜像的拉取、同步、删除与可见性控制 |
| 产品管理 | 产品定价、库存、限购、机型绑定、上下架 |
| 订单管理 | 产品订单与充值订单查询、状态处理 |
| 代金券 | 批量生成、作废、删除、统计（本 fork 新增） |
| 兑换码 | 兑换码生成与核销记录 |
| 邀请码 | 邀请码管理与邀请关系追踪 |
| 易支付配置 | 商户号、密钥、网关地址配置（本 fork 新增） |
| 站点配置 | 站名、Logo、Favicon、首页栏目开关、自定义页头页脚 |
| 站点链接 | 导航链接、友情链接管理 |
| 公告管理 | 公告发布、置顶、定时上下线 |
| 工单管理 | 工单分派、回复、分类、关闭 |
| 实名审核 | KYC 资料审核、通过/驳回 |
| 快照管理 | 快照查看、批量操作、定时计划任务 |
| 端口映射 | 端口池管理、映射规则审核 |
| 防火墙规则 | 黑名单/白名单规则下发 |
| 域名管理 | 域名解析配置与模板 |
| 流量管理 | 流量统计、异常流量预警、限流策略 |
| 性能监控 | 节点与实例性能指标采集与展示 |
| 日志审计 | 操作日志、登录日志、安全事件 |
| 任务中心 | 异步任务队列、失败重试、执行日志 |
| OAuth2 | 第三方登录 provider 配置 |
| API 令牌 | 后台令牌管理与权限范围 |
| 系统配置 | SMTP 邮件、验证码、注册开关、语言、安全策略 |

## 本 fork 定制功能

以下是 OneHost 相对上游 OneClickVirt **新增或增强**的能力，也是本项目的核心价值所在。

### 1. 易支付集成（新增）
内置国内主流第四方支付网关「易支付」，打通在线充值链路。

- 后台「易支付配置」页填写商户号（PID）、商户密钥、网关地址即可启用
- 支持 MD5 签名，签名算法已针对网关特性做过适配（小写签名 + 密钥自检按钮）
- 用户充值生成支付订单 → 跳转网关 → 异步回调 → 自动入账
- 相关代码：`server/api/v1/admin/yipay.go`、`server/service/product/yipay.go`

### 2. 代金券系统（新增）
上游原本没有此能力，本 fork 实现了完整的代金券体系。

- 后台批量生成券码（指定面额、数量、有效期），支持作废、删除、统计
- 用户在钱包页输入券码兑换，余额实时到账
- 兑换走数据库事务 + 行锁，杜绝并发重复兑换
- 相关代码：`server/api/v1/admin/voucher.go`、`server/model/product/voucher.go`

### 3. 余额钱包（新增）
与代金券配套的资金账户体系。

- 用户余额账户，支持充值、消费、退款
- 每笔变动写入 `user_balance_logs` 流水，全程可追溯
- 后台可直接调整指定用户余额（带流水记录）
- 余额扣减走 gorm 事务 + 行锁，避免并发超扣

### 4. 商城产品增强
在原产品模块上补齐售卖所需的细节：

| 增强点 | 说明 |
|--------|------|
| 支持 0 元定价 | 修复 Element Plus 把数字 0 当空值拦截的问题，允许免费产品 |
| 每人限购数量 | 新增 `maxPerUser` 字段（0 = 不限），防止恶意囤货 |
| 库存显示 | 商城列表与详情页展示实时库存 |
| 订单详情字段补齐 | 完整展示规格、周期、金额、支付与开通状态 |

### 5. 站点配置动态化
让平台外观完全可在后台配置，无需改代码：

- **站名 / Logo / Favicon 动态化**：浏览器标签页标题与图标随后台配置实时变化
- **首页栏目开关**：平台介绍、赞助方、推荐产品等区块可独立开关
- **自定义页头页脚**：支持注入自定义 HTML（如备案号、统计代码）
- **站点链接管理**：导航栏与页脚链接可配置

### 6. 公告系统改进
- 公告内容区分 `content`（纯文本）与 `contentHtml`（富文本），避免前端重复包裹 `<p>` 标签
- 支持置顶、定时上下线

### 7. 快照配额管控
- 快照页面显示"已用/允许"数量
- 达到配额或配额为 0 时，自动禁用新建与上传按钮，避免用户无效操作

### 8. 实例详情增强
- 补齐 **VNC 连接菜单**与**续费按钮**入口
- 节点未启用 IPv6 时，概览自动隐藏 IPv6 信息，避免误导用户

### 9. 国际化修复
- 修复工单分类、产品管理、商城页面等处未汉化或翻译丢失的问题
- 管理后台与用户端均支持中英文切换

### 10. 低配 VPS 部署优化（v0.18）
针对 1 核 1G ~ 2 核 2G 的入门级服务器做了专项优化，详见 [低配 VPS 部署优化](#低配-vps-部署优化)。

## 技术栈

| 层次 | 技术 |
|------|------|
| 后端 | Go 1.25、Gin、GORM、JWT、Zap、WebSocket、cron |
| 前端 | Vue 3、Vite、Element Plus、Pinia、Vue Router、vue-i18n |
| 可视化 | ECharts（流量/性能图表）、xterm.js（WebSSH）、noVNC（VNC） |
| 数据库 | MySQL 5.7+ / MariaDB 10.11+（可选 Redis 缓存） |
| 网关 | Nginx（反向代理、静态资源、WebSocket 透传、限流） |
| 部署 | Docker / Docker Compose，支持裸机与二进制部署 |

## 详细说明

[www.spiritlhl.net](https://www.spiritlhl.net/guide/oneclickvirt/oneclickvirt_precheck.html)

## 集成测试报告

自动化集成测试报告地址: [oneclickvirt.github.io/oneclickvirt](https://oneclickvirt.github.io/oneclickvirt/)

报告支持中英双语显示、亮色/暗色主题切换、Git ref/SHA/run 元数据和失败用例服务端日志展开，覆盖 200+ API 接口的功能测试、权限测试、边界测试和安全测试。详见 [`action_tests/`](action_tests/) 目录。

## 支持的虚拟化平台

| 类型标识 | 平台 | 实例类型 | 仓库地址 |
|---------|------|---------|---------|
| `lxd` | LXD | container, vm | [oneclickvirt/lxd](https://github.com/oneclickvirt/lxd) |
| `incus` | Incus | container, vm | [oneclickvirt/incus](https://github.com/oneclickvirt/incus) |
| `docker` | Docker | container | [oneclickvirt/docker](https://github.com/oneclickvirt/docker) |
| `podman` | Podman | container | [oneclickvirt/podman](https://github.com/oneclickvirt/podman) |
| `containerd` | Containerd (nerdctl) | container | [oneclickvirt/containerd](https://github.com/oneclickvirt/containerd) |
| `proxmox` | Proxmox VE | container, vm | [oneclickvirt/pve](https://github.com/oneclickvirt/pve) |
| `qemu` | QEMU | vm | [oneclickvirt/qemu](https://github.com/oneclickvirt/qemu) |
| `kubevirt` | KubeVirt | vm | [oneclickvirt/kubevirt](https://github.com/oneclickvirt/kubevirt) |

后端还包含面向本地或桌面虚拟化实验场景的适配器，例如 `orbstack`、`multipass`、`vagrant`、`virtualbox` 和 `vmware`。实现细节和支持范围见 [`server/provider/README.md`](server/provider/README.md)。

## 快速部署

尽量不要自行编译，推荐使用二进制文件分离部署或直接docker拉取镜像部署

### 方式零：使用 1Panel 第三方应用商店

[okxlin/appstore](https://github.com/okxlin/appstore) 已收录 OneClickVirt。已安装 1Panel 的用户，可以按该仓库说明添加或同步本地应用商店，然后在本地应用列表中选择 `oneclickvirt` 部署。

### 方式一：使用预构建镜像

使用已构建好的多架构镜像，会自动根据当前系统架构下载对应版本。

**镜像标签说明：**

| 镜像标签 | 说明 | 适用场景 |
|---------|------|---------|
| `oneclickvirt/oneclickvirt:latest` | 一体化版本（内置数据库）最新版 | 快速部署 |
| `oneclickvirt/oneclickvirt:20260925` | 一体化版本特定日期版本 | 需要固定版本 |
| `oneclickvirt/oneclickvirt:no-db` | 独立数据库版本最新版 | 不内置数据库 |
| `oneclickvirt/oneclickvirt:no-db-20260925` | 独立数据库版本特定日期 | 不内置数据库 |

所有镜像均支持 `linux/amd64` 和 `linux/arm64` 架构。

<details>
<summary>展开查看一体化版本（内置数据库）</summary>

**基础使用（不配置域名）：**

```bash
docker run -d \
  --name oneclickvirt \
  -p 80:80 \
  -p 10000-10099:10000-10099 \
  -e OCV_CONTROLLER_PORT_RANGE_START=10000 \
  -e OCV_CONTROLLER_PORT_RANGE_END=10099 \
  -v oneclickvirt-data:/var/lib/mysql \
  -v oneclickvirt-storage:/app/storage \
  --restart unless-stopped \
  oneclickvirt/oneclickvirt:latest
```

**配置域名访问：**

如果你需要配置域名，需要设置 `FRONTEND_URL` 环境变量：

```bash
docker run -d \
  --name oneclickvirt \
  -p 80:80 \
  -p 10000-10099:10000-10099 \
  -e OCV_CONTROLLER_PORT_RANGE_START=10000 \
  -e OCV_CONTROLLER_PORT_RANGE_END=10099 \
  -e FRONTEND_URL="https://your-domain.com" \
  -v oneclickvirt-data:/var/lib/mysql \
  -v oneclickvirt-storage:/app/storage \
  --restart unless-stopped \
  oneclickvirt/oneclickvirt:latest
```

或者使用 GitHub Container Registry：

```bash
docker run -d \
  --name oneclickvirt \
  -p 80:80 \
  -p 10000-10099:10000-10099 \
  -e OCV_CONTROLLER_PORT_RANGE_START=10000 \
  -e OCV_CONTROLLER_PORT_RANGE_END=10099 \
  -e FRONTEND_URL="https://your-domain.com" \
  -v oneclickvirt-data:/var/lib/mysql \
  -v oneclickvirt-storage:/app/storage \
  --restart unless-stopped \
  ghcr.io/oneclickvirt/oneclickvirt:latest
```

</details>

<details>
<summary>展开查看独立数据库版本</summary>

使用外部数据库，镜像更小，启动更快：

```bash
docker run -d \
  --name oneclickvirt \
  -p 80:80 \
  -p 10000-10099:10000-10099 \
  -e OCV_CONTROLLER_PORT_RANGE_START=10000 \
  -e OCV_CONTROLLER_PORT_RANGE_END=10099 \
  -e FRONTEND_URL="https://your-domain.com" \
  -e DB_HOST="your-mysql-host" \
  -e DB_PORT="3306" \
  -e DB_NAME="oneclickvirt" \
  -e DB_USER="root" \
  -e DB_PASSWORD="your-password" \
  -v oneclickvirt-storage:/app/storage \
  --restart unless-stopped \
  oneclickvirt/oneclickvirt:no-db
```

**环境变量说明：**
- `FRONTEND_URL`: 前端访问地址（必填，支持 http/https）
- `DB_HOST`: 数据库主机地址
- `DB_PORT`: 数据库端口（默认 3306）
- `DB_NAME`: 数据库名称
- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `OCV_CONTROLLER_PORT_RANGE_START` / `OCV_CONTROLLER_PORT_RANGE_END`：Agent 控制端内穿使用的 TCP 端口范围。宿主机与容器必须发布完全相同的范围，文档默认值为 `10000-10099`。仅在预期并发内穿数量确实需要时扩大，因为 Docker 会把发布范围展开为逐端口绑定。

控制端内穿监听实际运行在面板容器内部。Docker 的 `-p` 范围必须与两个 `OCV_CONTROLLER_PORT_RANGE_*` 环境变量完全一致；配置范围外的端口会被拒绝，不再出现数据库显示 active、但公网不可达的情况。源码或裸机部署不设置这两个变量时，继续保留原有的 `10000-65535` 分配范围。

`no-db` 镜像会将运行时配置保存到 `oneclickvirt-storage` 卷内的 `/app/storage/config.yaml`。更新镜像或重建容器时必须继续挂载同一个存储卷；初始化页面写入的数据库配置和系统级配置会随该卷保留。非空的 `DB_*` 环境变量优先于配置文件，因此重建时也可继续传入同一组数据库环境变量。显式挂载 `/app/config.yaml` 的部署仍会优先使用该文件。

</details>

> **说明**：`FRONTEND_URL` 用于配置前端访问地址，影响 CORS、OAuth2 回调等功能。系统会自动检测 HTTP/HTTPS 协议并调整相应配置，协议头可以是http或https。

### 方式二：使用 Docker Compose

<details>
<summary>展开查看 Docker Compose 部署</summary>

使用 Docker Compose 可以一键部署完整的开发环境，采用**分容器部署**架构，包括独立的前端容器、后端容器和数据库容器：

```bash
git clone https://github.com/oneclickvirt/oneclickvirt.git
cd oneclickvirt
cat > .env << 'EOF'
MYSQL_ROOT_PASSWORD=change-this-root-password
MYSQL_PASSWORD=change-this-app-password
EOF
docker-compose up -d --build || docker compose up -d --build
```

**默认配置说明：**

- 前端服务：`http://localhost:8888`
- 后端 API：通过前端代理访问
- MariaDB 数据库：端口 3306，数据库名 `oneclickvirt`
- 数据库凭据：来自 `.env` 的 `MYSQL_ROOT_PASSWORD` 和 `MYSQL_PASSWORD`
- 数据持久化：
  - 数据库数据：Docker volume `mysql_data`
  - 应用存储：`./data/app/`

**初始化配置：**

首次访问时会进入初始化界面，数据库配置请填写：
- 数据库地址：`mysql`（容器名称，不是 127.0.0.1）
- 数据库端口：`3306`
- 数据库名称：`oneclickvirt`
- 数据库用户：`oneclickvirt`
- 数据库密码：使用 `.env` 中的 `MYSQL_PASSWORD`

**自定义端口（可选）：**

如果需要修改前端访问端口，编辑 `docker-compose.yaml` 文件中的 ports 配置：

```yaml
services:
  web:
    ports:
      - "你的端口:80"  # 例如 "80:80" 或 "8080:80"
```

Agent 控制端内穿范围通过 `.env` 同时设置 `OCV_CONTROLLER_PORT_RANGE_START` 和 `OCV_CONTROLLER_PORT_RANGE_END`；Compose 会自动在 `api` 服务发布相同的 TCP 范围。范围应按实际容量设置，不要无条件发布全部端口。

**停止服务：**

```bash
docker-compose down
```

**查看日志：**

```bash
docker-compose logs -f
```

**清理数据：**

```bash
docker-compose down
rm -rf ./data
```

</details>

### 方式三：裸机全依赖安装

<details>
<summary>展开查看全量安装脚本</summary>

`scripts/install_full.sh` 会在一个流程中安装数据库、反向代理、TLS 配置、前端、后端和系统服务，支持 MySQL 兼容本地数据库（MySQL 或 MariaDB）以及 Caddy/Nginx/OpenResty。

安装器会自动识别常见 Linux 与类 Unix 目标，包括 Debian/Ubuntu、RHEL/CentOS/Rocky/Alma/Fedora/Amazon Linux、openSUSE/SLES、Arch/Manjaro、Alpine 和 BSD 包管理器；同时识别 systemd、OpenRC、rc.d/service 和无 init 环境。在原生 MySQL 包不可用或不稳定的发行版上，安装器会自动回退到 MariaDB 作为 MySQL 兼容后端；如需禁用该行为可使用 `--no-db-fallback`。BSD 安装需要存在对应 OS/架构的 release 二进制，否则请使用 Docker/Linux 或从源码构建服务端。

域名输入会自动识别协议前缀：输入 `https://panel.example.com` 自动启用 TLS，输入 `http://panel.example.com` 自动关闭 TLS，无前缀则交互询问。

```bash
curl -fsSL https://raw.githubusercontent.com/oneclickvirt/oneclickvirt/main/scripts/install_full.sh -o install_full.sh
bash install_full.sh
```

非交互部署示例：

```bash
# HTTPS 自动启用 TLS
bash install_full.sh \
  --non-interactive \
  --domain https://panel.example.com \
  --email admin@example.com \
  --db-type mariadb \
  --proxy caddy

# HTTP 纯端口模式，不启用 TLS
bash install_full.sh \
  --non-interactive \
  --domain http://192.168.1.100 \
  --proxy caddy
```

常用自动化参数：

```bash
bash install_full.sh --version v1.2.3 --db-wait-timeout 300
bash install_full.sh --db-type mysql --no-db-fallback
```

安装脚本默认要求至少 20 GB 可用磁盘和 4 GB 内存。生成的数据库密码会在安装摘要中输出，请在关闭终端前保存。

</details>

### 主控面板版本管理

登录超级管理员后，可在主控页面页脚的“升级管理”中查看当前版本、可用 Release、远程回退版本、本地备份和对应的人工命令。

- 面板仅会在 Linux、root 用户、受控 systemd 服务、主控二进制和待更新 Web 目录均位于受控安装根目录且无符号链接时执行自动升级、回退或重启。切换前会创建最多五份本地备份，失败时尝试恢复主控和 Web 资产；数据库迁移不会自动逆向回退。
- Release 必须包含 `SHA256SUMS` 资产。面板先下载校验清单，再校验所选 Linux 主控包和（部署配置为管理受控静态目录时的）`web-dist.zip`，随后才会解包或修改本机文件。历史 Release 没有该清单时会显示为不可自动应用，仍可使用原脚本手动处理。
- Docker、Docker Compose 和源码部署不会由面板写入运行环境，只会显示命令。Compose 当前项目使用 `docker compose up -d --build --force-recreate api web` 重建 API 与 Web，保留 `mysql_data` 命名卷；Docker 自定义端口、域名、环境变量或 `no-db` 部署应先导出并复用原始容器参数。
- 已配置的反向代理只有在 `ONECLICKVIRT_PROXY_SERVICES` 明确列出时才会在受控 systemd 重启后执行 `reload`。Release 元数据和资产下载会依次使用 GitHub、`ONECLICKVIRT_UPDATE_API_ENDPOINTS` 和现有 CDN/API 反代端点；所有面板下载地址均要求 HTTPS。

常用覆盖项（建议写入受控服务的环境文件后重启服务）：

| 变量 | 用途 |
| --- | --- |
| `ONECLICKVIRT_UPDATE_ENABLED=false` | 完全禁用面板升级、回退和重启按钮。 |
| `ONECLICKVIRT_UPDATE_MODE` | 强制部署模式：`systemd`、`docker`、`compose`、`source`、`embedded`、`unknown` 或 `disabled`。Compose 容器无法可靠自识别时可显式设为 `compose`。 |
| `ONECLICKVIRT_UPDATE_PROXY` | 逗号分隔的 HTTPS 发布资产/CDN 反代前缀。 |
| `ONECLICKVIRT_UPDATE_API_ENDPOINTS` | 逗号分隔的 HTTPS GitHub API 或 API 反代根地址。 |
| `ONECLICKVIRT_UPDATE_REPO` | Release 仓库，默认 `oneclickvirt/oneclickvirt`。 |
| `ONECLICKVIRT_UPDATE_FLAVOR` | `standalone` 或 `allinone`；通常由安装器记录的 `SERVER_ASSET` 自动识别。 |
| `ONECLICKVIRT_UPDATE_WEB` | 显式启用或禁用受控静态 Web 目录的替换。`install_full.sh` 会设为 `true`，因为其反向代理会提供 `web-dist.zip` 静态文件。 |
| `ONECLICKVIRT_INSTALL_ROOT`、`ONECLICKVIRT_SERVER_BIN`、`ONECLICKVIRT_WEB_DIR` | 受控安装根目录、主控二进制和受管理的静态 Web 目录。自动模式要求相应更新目标在安装根目录内。 |
| `ONECLICKVIRT_SERVICE_NAME`、`ONECLICKVIRT_SERVICE_FILE` | 受控 systemd 服务及其 unit 文件。 |
| `ONECLICKVIRT_PROXY_SERVICES` | 逗号分隔的 Nginx/OpenResty/Caddy 等 systemd 服务名，升级或重启后执行 `reload`。 |
| `ONECLICKVIRT_UPDATE_SCRIPT` | 页脚“命令”页使用的原安装脚本路径；未找到时显示官方下载后执行的命令。 |
| `ONECLICKVIRT_UPDATE_HEALTH_PORT` | 重启后的本机 `/api/v1/health` 检查端口，默认读取主控配置或使用 `8888`。 |
| `ONECLICKVIRT_UPDATE_ALLOW_UNVERIFIED=true` | 仅用于受控恢复场景，允许缺少 `SHA256SUMS` 的旧 Release；生产环境不建议设置。 |

### 方式四：自己编译打包

<details>
<summary>展开查看编译步骤</summary>

如果需要修改源码或自定义构建：

**一体化版本（内置数据库）：**

```bash
git clone https://github.com/oneclickvirt/oneclickvirt.git
cd oneclickvirt
docker build -t oneclickvirt .
docker run -d \
  --name oneclickvirt \
  -p 80:80 \
  -p 10000-10099:10000-10099 \
  -e OCV_CONTROLLER_PORT_RANGE_START=10000 \
  -e OCV_CONTROLLER_PORT_RANGE_END=10099 \
  -v oneclickvirt-data:/var/lib/mysql \
  -v oneclickvirt-storage:/app/storage \
  --restart unless-stopped \
  oneclickvirt
```

Docker 构建会自动内嵌 `scripts/install_agent.sh`。如果你还希望控制端镜像直接提供本地 Agent 发布包，而不是在下载时 302 跳转到 GitHub Releases，请在执行 `docker build` 前把下面这些文件放到 `server/assets/agent/`：

```text
install_agent.sh
oneclickvirt-agent-linux-amd64.tar.gz
oneclickvirt-agent-linux-arm64.tar.gz
```

**独立数据库版本：**

```bash
git clone https://github.com/oneclickvirt/oneclickvirt.git
cd oneclickvirt
docker build -f Dockerfile.no-db -t oneclickvirt:no-db .
docker run -d \
  --name oneclickvirt \
  -p 80:80 \
  -p 10000-10099:10000-10099 \
  -e OCV_CONTROLLER_PORT_RANGE_START=10000 \
  -e OCV_CONTROLLER_PORT_RANGE_END=10099 \
  -e FRONTEND_URL="https://your-domain.com" \
  -e DB_HOST="your-mysql-host" \
  -e DB_PORT="3306" \
  -e DB_NAME="oneclickvirt" \
  -e DB_USER="root" \
  -e DB_PASSWORD="your-password" \
  -v oneclickvirt-storage:/app/storage \
  --restart unless-stopped \
  oneclickvirt:no-db
```

更新或重建 `no-db` 容器时继续挂载同一个 `oneclickvirt-storage` 卷，运行时配置位于卷内的 `/app/storage/config.yaml`，无需重新初始化数据库。

直接执行 Go 源码编译时也是同样逻辑：`server/assets/agent/` 里的本地 Agent 资源是可选的，缺失时会回退到官方 GitHub 安装脚本和 Release 包，不会因此导致控制端构建失败。

</details>

### Proxmox VE 集成检查

`proxmoxve` 集成任务使用 [`oneclickvirt/pve`](https://github.com/oneclickvirt/pve) 的安装脚本，在可销毁的 Worker 上执行分离式安装阶段，等待内核重启以及可能由 `ifupdown2` 引导服务触发的第二次重启，然后校验 PVE 运行时、网桥、NAT 状态和主控侧 Provider 链路。只有在分离任务无法继续观测时才会将失联归类为基础设施跳过，不会把它当作模块测试通过。PVE 仓库本身还会执行脚本语法、网络回归和 ShellCheck 检查。

### 方式五：手动开发部署

<details>
<summary>展开查看开发部署步骤</summary>

#### 环境要求

* Go 1.25.0
* Node.js 22+
* MySQL 5.7+
* npm 或 yarn

#### 环境部署

1. 构建前端
```bash
cd web
npm i
npm run serve
```

2. 构建后端
```bash
cd server
go mod tidy
go run main.go
```

3. 开发模式下不需要反代后端，vite已自带后端代理请求。

5. 在mysql中创建一个空的数据库```oneclickvirt```，记录对应的账户和密码。

6. 访问前端地址，自动跳转到初始化界面，填写数据库信息和相关信息，点击初始化。

7. 完成初始化后会自动跳转到首页，可以开始开发测试了。

#### 本地开发

* 前端：[http://localhost:8080](http://localhost:8080)
* 后端 API：[http://localhost:8888](http://localhost:8888)
* API 文档：[http://localhost:8888/swagger/index.html](http://localhost:8888/swagger/index.html)

</details>

## 初始账户

首次初始化时会根据初始化表单创建管理员账户。快捷填充会每次生成随机强密码，请在提交表单前保存生成的密码。

## 配置文件

主要配置文件位于 `server/config.yaml`

> **注意**：该文件会被运行时写入数据库明文密码，**请勿提交到版本库**。
> `.gitignore` 已忽略 `.env`、`*.log`、`data/`，提交前建议自检：`git diff --cached | grep -iE 'password|secret|"'`

## 低配 VPS 部署优化

在 **1 核 1G / 2 核 2G** 级别的入门级服务器上，默认的 MariaDB 配置会直接触发 OOM
（进程被系统杀掉，容器退出码 `137`）。v0.18 已内置低配优化配置：

| 参数 | 原默认值 | 低配优化值 | 原因 |
|------|---------|-----------|------|
| `innodb_buffer_pool_size` | 256M | **64M** | 1G 内存机器上 256M 缓冲池过大 |
| `max_connections` | 500 | **80** | 每连接都有内存开销 |
| `performance_schema` | 开启 | **关闭** | 可省约 100MB 内存 |
| `bind-address` | 127.0.0.1 | **注释掉** | 否则 api 容器无法通过 `mysql:3306` 连接 |
| `innodb_redo_log_capacity` | 256M | **移除** | MySQL 8.0.30+ 专有参数，MariaDB 不识别会导致启动失败 |

配置文件位于 [`deploy/my.cnf`](deploy/my.cnf)。实测效果：**MariaDB 内存占用由 400MB+ 降至 12.5MB**，
整机（2 核 960Mi）跑完 web + api + mysql 三容器后仍有 350MB 以上可用内存。

其他低配部署建议：

- **构建必须串行**：先 `docker compose build api`（2 核机器上约 2.5 小时），再 `docker compose build web`（约 25 分钟）。并行构建必然 OOM。
- **提前准备 swap**：建议 2G 以上，`go build -a` 编译期会大量占用。
- **磁盘预留 10G 以上**：构建中间层与镜像较占空间。

## HTTPS 配置

v0.18 起的 nginx 配置已原生支持 HTTPS：**80 端口仅保留 `/.well-known/acme-challenge/` 供证书续期，
其余请求 301 跳转 HTTPS**；443 端口承载全部业务（API、WebSocket、静态资源、前端路由），并启用 HTTP/2 与 HSTS。

配置步骤（以 Docker Compose 部署为例）：

```bash
# 1. 安装 certbot
apt-get update && apt-get install -y certbot

# 2. 申请证书（webroot 模式，无需停机）
certbot certonly --webroot -w /opt/onehost/data/certbot-webroot -d 你的域名 \
  --non-interactive --agree-tos -m 你的邮箱

# 3. 重启 web 容器使证书生效
cd /opt/onehost && docker compose restart web

# 4. 验证续期（certbot 自带 systemd timer，到期前自动续期）
certbot renew --dry-run
```

**关键设计**：`deploy/default.conf` 与 `/etc/letsencrypt` 都是通过 volume 挂载进容器的，
所以修改 nginx 配置或证书续期后，**只需 `docker compose restart web`，无需重建镜像**（重建要 25 分钟）。

证书续期后会自动触发 `/etc/letsencrypt/renewal-hooks/deploy/reload-web.sh` 重新加载 nginx。

## 版本历史

| 版本 | 主要变更 |
|------|---------|
| **v0.25.3** | 文档：README 截图改为 jsDelivr CDN 加载，修复国内无法显示 |
| v0.25.2 | 文档：README 按线上站点 webhost.ypvps.com 重写演示截图（首页亮/暗/移动端、登录、管理后台、用户控制台），新增在线演示与版本历史 |
| v0.25.1 | 首页平台轮播 / 赞助方未配置图标时回退内置图标 |
| v0.25 | 已初始化数据库启动时执行"只增不改"的增量结构同步，修复升级缺列；版本号与发布号一致 |
| v0.24 | 首页落地页重做、ECharts 亮/暗主题、管理员与用户仪表盘重做、404 页面 |
| v0.23 | 全新靛蓝→天蓝设计系统（亮/暗双主题）、侧边栏分组、分屏式登录注册、找回密码页 |
| v0.18 | 低配 VPS 部署优化（MariaDB 内存 400MB+ → 12.5MB）；HTTPS 支持（80 跳 443、ACME 续期路径）；修复 `.dockerignore` 排除整个 `deploy/` 导致前端镜像构建失败；nginx 配置改为 volume 挂载（改配置免重建）；加固 `.gitignore` 防止数据库密码误提交 |
| v0.17 | 首页推荐产品卡片点击跳转到具体产品详情页 |
| v0.07 | 易支付签名小写 + 密钥自检按钮；默认节点 / OS 选项恢复；`selectProvider` 默认节点优先兜底；站点链接类型映射修复 |
| v0.06 | Owner 15 项需求收尾：续费入口、快照配额、订单详情、库存显示、站名动态化、代金券与余额体系 |
| v0.05 | 修复首页栏目开关失效（gorm 模型列名与迁移列名不匹配，导致开关存不进去） |
| v0.04 | 修复站点配置保存 500（JSON 键名未映射为真实数据库列名） |
| v0.03 | 首页动态化、站点配置 `Save()` 全列覆盖修复、SMTP 配置唯一索引缺陷修复 |

> v0.08 ~ v0.16 的变更详见 [提交历史](https://github.com/qdmz/onehost/commits/main)。

## 常见问题

### 初始化页「测试数据库连接」提示"系统内部错误"

数据库主机必须填**容器名 `mysql`**，不能填 `127.0.0.1` 或 `localhost`。
数据库运行在独立的 mysql 容器中，api 容器内的 `127.0.0.1` 并没有数据库服务，两者靠 compose 内部网络通信。

### 初始化完成后状态一直是 `starting`

正常现象。初始化后配置管理器需要重载，执行一次即可变为 `ready`：

```bash
docker compose restart api
```

### 前端构建报错 `deploy/default.conf: not found`

原 `.dockerignore` 中的 `deploy/*` 规则排除了整个 deploy 目录，导致构建上下文拿不到该文件。
**v0.18 已修复**（放行 `default.conf` 与 `nginx.dockerfile`）。若使用旧版本，需手动在 `.dockerignore` 末尾追加：

```
!deploy/default.conf
!deploy/nginx.dockerfile
```

### 修改了 nginx 配置但不生效

配置是 volume 挂载的，执行 `docker compose restart web` 即可，**不要重建镜像**。

### MariaDB 容器启动后很快退出（退出码 137）

这是 OOM（内存不足）被系统杀掉的典型表现。请确认使用的是 v0.18 的 `deploy/my.cnf`，
并参考 [低配 VPS 部署优化](#低配-vps-部署优化) 调低内存参数、增加 swap。

### 易支付提示"MD5 签名校验失败"

签名算法本身已验证无误，绝大多数情况是**后台填写的商户密钥与支付网关侧不一致**。
请从易支付商户后台重新复制密钥填入「易支付配置」页，页面提供密钥自检按钮可辅助排查。

## 赞助方

感谢以下团体或个人赞助 OneClickVirt 项目：

[![Docker Sponsored OSS](https://img.shields.io/badge/Docker-Sponsored%20OSS-2496ED?logo=docker&logoColor=white)](https://hub.docker.com/r/oneclickvirt/oneclickvirt)

<p>
  <a href="https://dartnode.com?aff=bonus">
    <img src="https://cdn.jsdelivr.net/gh/qdmz/onehost@main/web/src/assets/images/dartnode.png" alt="DartNode" height="44">
  </a>
  &nbsp;&nbsp;
  <a href="https://console.zmto.com/?affid=1524">
    <img src="https://console.zmto.com/templates/2019/dist/images/logo_dark.svg" alt="zmto" height="44">
  </a>
  &nbsp;&nbsp;
  <a href="https://community.ibm.com/zsystems/form/l1cc-oss-vm-request/">
    <img src="https://cdn.jsdelivr.net/gh/qdmz/onehost@main/web/src/assets/images/ibm-linuxone.png" alt="IBM LinuxONE OSS Community Cloud" height="44">
  </a>
  &nbsp;&nbsp;
  <a href="https://fossvps.org/">
    <img src="https://lowendspirit.com/uploads/userpics/793/nHSR7IOVIBO84.png" alt="fossvps" height="44">
  </a>
  &nbsp;&nbsp;
  <a href="https://linux.do/">
    <img src="https://cdn3.ldstatic.com/original/4X/d/1/4/d146c68151340881c884d95e0da4acdf369258c6.png" alt="Linux DO" height="44">
  </a>
</p>

## LICENSE

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Foneclickvirt%2Foneclickvirt.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Foneclickvirt%2Foneclickvirt?ref=badge_large&issueType=license)

## 演示截图

以下截图于 2026-10-13 从线上站点 [webhost.ypvps.com](https://webhost.ypvps.com)（v0.25.1）实时截取，覆盖首页（亮色 / 暗色 / 移动端）、登录页、管理后台和用户控制台。节点 API 地址已打码。图片经 jsDelivr CDN 加载，国内网络也可正常显示。

### 首页

| 亮色 | 暗色 |
|------|------|
| ![首页 亮色](https://cdn.jsdelivr.net/gh/qdmz/onehost@main/.back/home-light.jpg) | ![首页 暗色](https://cdn.jsdelivr.net/gh/qdmz/onehost@main/.back/home-dark.jpg) |

**移动端首页**

<img src="https://cdn.jsdelivr.net/gh/qdmz/onehost@main/.back/home-mobile.jpg" alt="移动端首页" width="320">

**登录页**

![登录页](https://cdn.jsdelivr.net/gh/qdmz/onehost@main/.back/login.jpg)

### 管理后台

**管理员仪表盘**（CPU / 内存 / 磁盘进度环、实例分布、资源分配）

![管理员仪表盘](https://cdn.jsdelivr.net/gh/qdmz/onehost@main/.back/admin-dashboard.jpg)

**管理员仪表盘 · 暗色模式**

![管理员仪表盘 暗色](https://cdn.jsdelivr.net/gh/qdmz/onehost@main/.back/admin-dashboard-dark.jpg)

**产品管理**

![产品管理](https://cdn.jsdelivr.net/gh/qdmz/onehost@main/.back/admin-products.jpg)

**节点管理**

![节点管理](https://cdn.jsdelivr.net/gh/qdmz/onehost@main/.back/admin-providers.jpg)

**站点配置**

![站点配置](https://cdn.jsdelivr.net/gh/qdmz/onehost@main/.back/admin-site-config.jpg)

### 用户控制台

**用户仪表盘**（快捷入口与六项资源配额）

![用户仪表盘](https://cdn.jsdelivr.net/gh/qdmz/onehost@main/.back/user-dashboard.jpg)

**产品商城**

![产品商城](https://cdn.jsdelivr.net/gh/qdmz/onehost@main/.back/user-store.jpg)

**我的钱包**（余额、代金券兑换、余额流水）

![我的钱包](https://cdn.jsdelivr.net/gh/qdmz/onehost@main/.back/user-wallet.jpg)
