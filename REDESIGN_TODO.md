# ONECLICKVIRT_REDESIGN — 改造项目

## ✅ 已完成
- [x] 移除用户等级限制（order.go）
- [x] 支付成功后自动创建 VM（callback.go）
- [x] 产品模型扩展（VMType, AutoProvision 字段）

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
- [ ] 确保 create_vm 任务能被正确执行
- [ ] 任务执行逻辑完善

### 4. 测试
- [ ] 测试支付回调流程
- [ ] 测试 VM 自动创建
