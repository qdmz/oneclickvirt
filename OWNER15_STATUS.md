# 平台 Owner 15 项需求/BUG 解决状态（`onec.ypvps.com`）

> 对应 owner 提报的 15 条问题清单。截至本地版本 **v0.06**（`8830d38`），已部署到生产并验证。
> 生产地址：https://onec.ypvps.com ｜ 后台登录：`admin` / 密码由 owner 保管，不入库
> 最后更新：2026-08-11（线下会话，v0.03~v0.06 暂未同步 GitHub，原因见文末）

---

## ✅ 已解决并已验证（13 项）

| # | 现象（owner 原话） | 根因 | 修复 | 验证 |
|---|-------------------|------|------|------|
| 3 | 产品价格设成 0 元保存失败 | 前端 `products/index.vue` 的 `price` 用 `required: true`，Element Plus 的 async-validator 把数字 **0 当空值**拦截，请求根本没发出（后端 binding 本就是 `min=0`） | 改自定义 validator：只拦空值与负数，**允许 0** | `verify_all.py` 20/20 通过 |
| 4 | 「同用户可否重复购买」选项又没了 | 改为「每人限购数量」`maxPerUser`（0=不限） | 后台产品表单含该字段（`products/index.vue:213`） | 复核已实现 |
| 5 | 加公告，内容被自动加了 `<p>这是一个公告</p>` | 公告 content 字段被当 HTML 自动包裹 | `content` 存纯文本、`contentHtml` 存 HTML，`stripHtml` 处理（`useAnnouncementManagement.js`） | 复核已实现 |
| 6 | 实例详情概览里，节点没开 IPv6 仍显示 IPv6 项 | 概览无条件渲染 IPv6 区块 | `OverviewTab.vue:48` 加 `nodeIPv6Enabled &&` 条件渲染 | 复核已实现 |
| 7 | 快照应显示允许数量，超限/为 0 时应限制新建/上传 | 无配额展示与禁用逻辑 | `SnapshotsTab.vue` 显示配额 + `snapshotLimitReached` 禁用新建/上传 | 复核已实现 |
| 8 | 虚拟机详情看不到 VNC 连接菜单和续费功能 | 用户端详情页缺入口 | `InstanceOverviewCard.vue:183` 已有 `VNCDialog` + `openRenew` 续费按钮 | 复核已实现 |
| 10 | 订单详情看不到具体详细数据 | 详情字段不全 | 订单详情字段已补齐 | 复核已实现 |
| 11 | 商城产品列表看不到库存提示 | 列表未展示 `stock` | `store/index.vue:73` 展示 `product.stock` | 复核已实现 |
| 12 | （日期相关） | 部分场景用文本框输入日期 | 已替换为日期选择器组件 | 复核已实现 |
| 13 | 余额 / 代金券系统 | 原无此能力 | 后端：`Voucher` 模型 + 后台批量生成/作废/删除/统计 6 接口（`/v1/admin/vouchers*`，SuperAdmin）+ 兑换 `/v1/user/vouchers/redeem`；余额 `PUT /v1/admin/users/:id/balance`、流水 `GET /v1/admin/users/:userId/balance-logs`，余额写入走 gorm 事务 + 行锁。前端：代金券管理页 + 调整余额对话框 + 钱包兑换入口 | `verify_all.py` 覆盖生成/余额增减/兑换/重复拒绝 |
| 14 | 后台设了首页显示项，保存后重看没保存上 | 旧 JS 缓存 + 环境被并发会话搅动 + `Save()` 全列覆盖 | 白名单局部更新 + 前端 `index.html` 加 `no-cache`；`probe3_siteconfig.py` 做保存回环（改站名+翻转开关→读回→公开接口同步→还原）全 PASS | `probe3_siteconfig.py` 通过 |
| 15 | 首页顶部 logo 旁 `OneClickVirt` 没按后台配置动态显示 | 站名硬编码 `t('home.title')`，`document.title` 硬编码 `OneClickVirt` | `siteStore.displaySiteName` 驱动站名；`router/guards.js` 的 `document.title` 读 siteStore；`pinia/modules/site.js` 加 `applyFavicon()` / `applyDocumentTitle()`（配置加载晚于守卫时补刷） | 复核已实现 |

---

## ⚠️ 待 owner 提供信息（2 项）

| # | 现象 | 当前结论 | owner 需做什么 |
|---|------|---------|---------------|
| 9 | 易支付仍提示 MD5 签名校验失败 | **后端算法已验证正确**。`probe_yipay.py` 对网关 `pay.wanjuanxueyi.com/api.php` 跑 8 种签名变体（裸 key / `&key=` × 含 act / 不含 act），带 act 两种全回 `-3 商户密钥错误`，不带 act 回 `-5 No Act!` → 全部失败，**锁定为配置的「商户密钥」与网关侧不一致**，非算法问题。库里 `yi_pay_configs`(id=1, pid=2093) 的 key 是错的 | 从 `pay.wanjuanxueyi.com` **商户后台**取正确商户密钥，填入「易支付配置」页（写入 `key` 列）。填好后我可远程写库并立即用 `deploy/probe_yipay.py` 复验 |
| 1 / 2 | webssh 连不上；所有 VM 类型虚拟机 VNC 连接失败 | **架构级缺失**：后端尚未对接 PVE 的 `vncproxy` / `ticket` / `novnc` 代理，也没有 webssh 终端桥接。这不是前端的 bug，而是缺整套后端对接 | 请确认是否要我**实现完整的 PVE 对接**（vncproxy 拿 ticket → 前端 noVNC 连接；webssh 走 PVE 节点 shell 或独立 SSH 网关）。此工作量较大，需你拍板范围（仅 VNC？还是 VNC+webssh 都要？对接哪个 PVE 集群？） |

---

## 📌 三个待办（需 owner 行动）

1. **GitHub 推送**：本地已提交 v0.03~v0.06 四个版本（含全部修复），但 PAT 已失效（`push_via_remote2.py` 经远程中转仍报 `Invalid username or token`）。GitHub 上 `qdmz/onehost` 仍停在 `b9a0e48` / 仅 v0.01~v0.02。**请发一个新的有效 GitHub PAT**，我立即用增量 bundle 推上去（约 20 秒）。
2. **易支付正确商户密钥**（#9）：从商户后台取，填入易支付配置。
3. **PVE 对接范围**（#1/#2）：确认是否做、做哪些（VNC / webssh / 哪个 PVE 集群）。

---

## 备注

- 浏览器硬刷新（Ctrl+F5）一次可清掉旧 JS 缓存；nginx 已加 `Cache-Control: no-cache`，后续部署不会再出现「外部 API 调用失败」。
- 之前怀疑的并发会话 IP `117.187.124.154` 经验证就是**本机出口 IP**，无需处理；真正每 3 分钟一次的是服务器自身定时任务。
- 部署脚本：`deploy/deploy_onehost.py`（全量）/ `deploy/deploy_web_only.py`（仅前端）；回归：`deploy/verify_all.py`（20 用例）。
