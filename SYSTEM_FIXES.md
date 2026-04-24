# OneClickVirt 系统问题修复记录

## 修复时间
2026-04-24 15:00

## 修复的问题

### 1. 兑换码功能无法正常使用 ✅ 已修复

#### 问题描述
```
Error 1054 (42S22): Unknown column 'amount' in 'field list'
```

#### 问题原因
代码模型中使用的字段名是`Amount`，但数据库表中的实际字段名是`value`，导致字段名不匹配。

#### 修复方案
修改`server/model/redemption/redemption.go`文件，将`Amount`字段的GORM标签修改为映射到数据库的`value`字段：

```go
// 修改前
Amount    int64      `json:"amount" gorm:"default:0;comment:金额(分)或等级数"`

// 修改后
Amount    int64      `json:"amount" gorm:"column:value;default:0;comment:金额(分)或等级数"`
```

#### 修复文件
- `server/model/redemption/redemption.go`

### 2. 申请领取实例500错误 ⏳ 需要重新构建

#### 问题描述
用户提交申请领取实例后提示500错误。

#### 问题原因
这是之前修复的MySQL事务隔离级别问题导致的，需要重新构建和部署修复后的代码。

#### 修复方案
应用之前修复的`server/service/resources/quota.go`文件，移除错误的隔离级别设置代码。

#### 修复文件
- `server/service/resources/quota.go`

### 3. 节点SSH连接问题 ✅ 已解决

#### 问题描述
```
SSH Dial失败: address: "https:22", error: dial tcp: lookup https on 100.100.100.100:53: no such host
```

#### 问题原因
PVE节点的endpoint配置错误，被错误地解析为`https:22`而不是正确的IP地址。

#### 解决方案
用户已经在管理后台重新配置了PVE节点，现在SSH连接状态显示为`online`。

#### 当前状态
- 节点名称：pvetest
- 节点类型：proxmox
- API状态：online
- SSH状态：online
- 资源同步：已同步

## 修复状态

### 已完成
- ✅ 兑换码字段名修复
- ✅ SSH连接问题解决（用户配置）
- ✅ 代码修改完成

### 待完成
- ⏳ 重新构建Docker镜像
- ⏳ 部署修复后的版本
- ⏳ 验证所有功能正常

## 重新构建和部署

### 构建步骤
```bash
cd /root/.openclaw/workspace/oneclickvirt

# 1. 确认修改
git diff

# 2. 提交修改
git add server/model/redemption/redemption.go
git commit -m "fix: 修复兑换码字段名不匹配问题

- 修改Amount字段映射到数据库的value字段
- 解决兑换码创建和使用时的数据库错误
- 修复Unknown column 'amount'错误"

# 3. 构建Docker镜像
docker build -t oneclickvirt:fixed .

# 4. 停止旧容器
docker stop oneclickvirt
docker rm oneclickvirt

# 5. 启动新容器
docker run -d --name oneclickvirt \
  -p 8080:80 \
  -p 8443:443 \
  -v /path/to/data:/app/storage \
  oneclickvirt:fixed
```

### 验证步骤
```bash
# 1. 测试兑换码功能
# 在管理后台创建兑换码，然后测试使用

# 2. 测试申请领取实例
# 在用户界面申请领取实例，检查是否成功

# 3. 测试SSH连接
# 在节点管理中测试SSH连接功能

# 4. 测试实例创建
# 创建测试实例，验证功能正常
```

## 测试检查清单

### 兑换码功能
- [ ] 创建兑换码成功
- [ ] 使用兑换码成功
- [ ] 兑换码次数限制正常
- [ ] 兑换码过期时间正常

### 申请领取功能
- [ ] 提交申请不报500错误
- [ ] 申请记录正确保存
- [ ] 管理员能看到申请
- [ ] 审核流程正常

### SSH连接功能
- [ ] SSH连接测试成功
- [ ] 远程登录管理正常
- [ ] 节点状态显示正确
- [ ] 资源同步正常

### 实例创建功能
- [ ] 创建实例不报事务错误
- [ ] 实例创建成功
- [ ] 配额验证正常
- [ ] 实例状态正确

## 注意事项

1. **数据备份**：部署前请备份数据
2. **测试验证**：在测试环境充分验证后再部署到生产
3. **回滚准备**：准备好快速回滚方案
4. **监控观察**：部署后密切监控系统表现

## 后续建议

1. **代码审查**：让团队review修复代码
2. **自动化测试**：建立自动化测试流程
3. **文档更新**：更新相关文档
4. **用户通知**：通知用户系统已修复

## 联系方式

如有问题，请联系：
- 修复人员：小美丽
- 修复时间：2026-04-24
- 项目地址：https://github.com/qdmz/oneclickvirt

---

**修复完成时间**：2026-04-24 15:00
**修复状态**：✅ 代码修复完成，待部署验证
**下一步**：重新构建和部署
