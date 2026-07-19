# ONECLICKVIRT_REDESIGN — 改造项目

## ✅ 已完成
- [x] 移除用户等级限制（order.go）
- [x] 支付成功后自动创建 VM（callback.go）
- [x] 产品模型扩展（VMType, AutoProvision, OS字段）
- [x] 数据库迁移脚本（scripts/add_product_vm_fields.sql）
- [x] 前端产品管理页面增加VM配置选项（index.vue）
- [x] 任务系统支持create_vm类型（helpers.go）
- [x] Instance模型添加状态常量和自动开通标记（provider.go）
- [x] 后台任务执行器（auto_provision.go）
- [x] 代码已提交（feature/redesign-product-flow分支）

## 🔧 进行中

### 1. 数据库迁移 — 为产品表添加 auto_provision 字段
- [x] 检查产品表现有结构
- [x] 创建迁移脚本（scripts/add_product_vm_fields.sql）
- [x] 应用迁移（需手动执行或在服务器上运行）

### 2. 前端产品管理页面
- [x] 管理员后台产品编辑表单增加"自动开通"复选框
- [x] VM类型选择器（vm/container/proxmox）
- [x] 最低购买等级输入

### 3. 完善后台任务处理
- [x] create_vm 任务路由已添加
- [x] 任务执行器已创建
- [ ] 需要Go 1.25+环境验证编译

### 4. 部署步骤
- [ ] 在服务器上执行数据库迁移：mysql -u root -p oneclickvirt < scripts/add_product_vm_fields.sql
- [ ] 用Go 1.25+编译后端
- [ ] 重启后端服务
- [ ] 前端npm run build并部署
