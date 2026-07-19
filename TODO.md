# ONECLICKVIRT_REDESIGN — 重构 oneclickvirt 项目功能

## 项目信息
- 仓库：https://github.com/oneclickvirt/oneclickvirt.git
- 路径：/run/csi/mount-root/nas/4079184d856ecc166ed19d4887083405/workspaces/default/coding_projects/oneclickvirt
- 后端：Go 1.25+ (Gin + GORM + MySQL)
- 前端：Vue 3 (Vite + Element Plus)

## 用户需求

### 1. 易支付接口用户充值功能
- [ ] 实现易支付充值接口（epay/yipay）
- [ ] 支持充值订单生成、支付、回调处理
- [ ] 充值成功后自动充入用户钱包

### 2. 虚拟机产品管理
- [ ] 产品模型重构（产品 = 虚拟机规格 + 价格）
- [ ] 管理员后台产品 CRUD
- [ ] 产品包含：名称、CPU、内存、磁盘、带宽、价格、周期、库存

### 3. 去掉用户等级资源限制
- [ ] 移除用户等级等级相关的资源限制逻辑
- [ ] 实例创建不再基于等级配额
- [ ] 改为基于产品配置

### 4. 申请领取实例改为在线订购产品
- [ ] 移除旧的申请/领取流程
- [ ] 改为产品列表 + 在线订购 + 支付流程
- [ ] 支付成功后自动创建虚拟机

### 5. 支付成功后自动开通虚拟机
- [ ] 支付回调自动触发实例创建
- [ ] 自动分配节点、镜像
- [ ] 自动设置密码/密钥

### 6. noVNC 在线管理
- [ ] 集成 noVNC WebSocket 代理
- [ ] 实例详情页添加 VNC 控制台入口
- [ ] 支持 Proxmox VE 的 VNC 接口

### 7. 管理员后台手动开通虚拟机
- [ ] 管理员可选择用户 + 产品 + 节点 手动创建实例
- [ ] 支持批量创建

### 8. 工单系统
- [ ] 已有工单模型和 API，需确认功能完整性
- [ ] 用户可创建/回复工单
- [ ] 管理员可管理/回复/关闭工单

## 技术栈
- 后端：Go 1.25, Gin, GORM, MySQL
- 前端：Vue 3, Vite, Element Plus, Pinia

## 关键文件
- server/model/product/product.go - 产品模型
- server/model/user/user.go - 用户模型（含等级/配额字段）
- server/api/v1/user/order.go - 订单 API
- server/api/v1/user/recharge.go - 充值 API
- server/api/v1/user/container_vm.go - 容器/虚拟机 API
- server/router/user.go - 用户路由
- server/model/provider/provider.go - 节点/实例模型
- web/src/view/user/apply/ - 申请实例页面
- web/src/view/user/purchase/ - 购买产品页面
- web/src/view/admin/products/ - 管理员产品管理
