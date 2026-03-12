# 功能修复检查清单

## 已修复的问题

### 1. ✅ 后台代理商设置菜单
**问题**: 管理员侧边栏看不到代理商管理菜单
**修复**: 在 `web/src/view/layout/components/Sidebar/index.vue` 的管理员路由中添加了：
- `/admin/agents` - 代理商管理
- `/admin/kyc` - 实名管理

**测试方法**:
1. 使用管理员账号登录
2. 查看侧边栏，应该能看到代理商管理和实名菜单
3. 点击菜单，应该能正常跳转到相应页面

### 2. ⚠️ 域名保存功能
**状态**: 代码检查正常，需要实际测试验证

**检查结果**:
- ✅ 前端 API 调用正确的 (`/v1/user/domains`)
- ✅ 后端路由配置正确
- ✅ 后端 handler 正确 (`CreateDomain`)
- ✅ 数据库模型正确 (`Domain` struct)
- ✅ 编译通过

**可能的问题原因**:
1. 数据库未初始化 `domain_configs` 表
2. 配额限制导致创建失败
3. 域名格式验证失败
4. 权限问题

**测试方法**:
1. 登录用户账号
2. 进入"域名管理"页面
3. 点击"添加域名"
4. 填写表单：
   - 域名: `test.example.com`
   - 实例ID: 选择一个存在的实例ID
   - 内部IP: 填写虚拟机内部IP (如 `10.0.0.10`)
   - 内部端口: 填写端口号 (如 `80`)
   - 协议: 选择 `http`
5. 点击"创建"
6. 观察是否成功，如果失败查看错误提示

## 数据库初始化检查

确保以下表已创建（GORM AutoMigrate 会自动创建）：
- `domains` - 域名绑定表
- `domain_configs` - 域名配置表（默认配置会在首次调用 `GetDomainConfig()` 时创建）
- `agents` - 代理商表
- `sub_user_relations` - 代理商子用户关系表
- `commissions` - 佣金记录表
- `kyc_records` - 实名认证记录表

## 编译验证

✅ 后端编译成功: `server/oneclickvirt.exe`
✅ 前端编译成功: `web/dist/` 目录

## 启动测试

### 后端启动
```bash
cd server
./oneclickvirt.exe
```
或使用 Docker:
```bash
docker-compose up -d
```

### 前端启动 (开发模式)
```bash
cd web
npm run dev
```

或使用已构建的前端 (Nginx):
访问 `http://your-server-ip` (需配置 Nginx)

## 功能列表

### 管理员功能
- [ ] 仪表盘 - `/admin/dashboard`
- [ ] 用户管理 - `/admin/users`
- [ ] 代理商管理 - `/admin/agents` ⭐ 新修复
- [ ] 域名配置 - `/admin/domain-config`
- [ ] 域名管理 - `/admin/domains`
- [ ] 实名管理 - `/admin/kyc` ⭐ 新修复
- [ ] 其他原有功能...

### 用户功能
- [ ] 仪表盘 - `/user/dashboard`
- [ ] 我的实例 - `/user/instances`
- [ ] 域名管理 - `/user/domains` ⭐ 需要验证
- [ ] 钱包 - `/user/wallet`
- [ ] 订单 - `/user/orders`
- [ ] 购买 - `/user/purchase`
- [ ] 实名认证 - `/user/kyc`
- [ ] 个人中心 - `/user/profile`
- [ ] 其他原有功能...

### 代理商功能
- [ ] 代理商仪表盘 - `/agent/dashboard`
- [ ] 子用户管理 - `/agent/sub-users`
- [ ] 佣金记录 - `/agent/commissions`
- [ ] 钱包管理 - `/agent/wallet`
- [ ] 代理商资料 - `/agent/profile`

## 如果域名保存失败

### 调试步骤

1. **查看后端日志**
   - 检查控制台输出，看是否有错误信息
   - 常见错误：
     - "域名格式无效" - 域名不符合规范
     - "该域名已被绑定" - 域名重复
     - "已达到域名绑定上限" - 配额限制
     - "不允许绑定此后缀的域名" - 后缀限制

2. **检查数据库**
   ```sql
   -- 查看域名配置是否存在
   SELECT * FROM domain_configs;

   -- 查看已绑定域名
   SELECT * FROM domains;
   ```

3. **检查配额**
   - 默认用户可以绑定 3 个域名
   - 管理员可以在 `/admin/domain-config` 中修改配额

4. **浏览器控制台**
   - 按 F12 打开开发者工具
   - 切换到 Console 标签
   - 查看是否有 JavaScript 错误
   - 切换到 Network 标签
   - 查看 API 请求的响应内容

## 配置说明

### config.yaml 关键配置

```yaml
# 实名认证配置
payment:
  enable-real-name: true        # 启用实名认证
  require-real-name: false      # 是否强制实名
  callback-url: "http://your-domain/api/callback"

# 邮件配置
email:
  smtp-host: "smtp.gmail.com"
  smtp-port: 587
  smtp-user: "your-email@gmail.com"
  smtp-pass: "your-password"

# 域名配置 (在数据库 domain_configs 表中)
# - MaxDomainsPerUser: 每用户最多域名数
# - MaxDomainsPerAgentUser: 代理商子用户最多域名数
# - AllowedSuffixes: 允许的域名后缀 (逗号分隔)
# - DNSType: DNS 类型 (dnsmasq/hosts)
# - DNSConfigPath: DNS配置文件路径
# - NginxConfigPath: Nginx配置路径
```

## 注意事项

1. **域名绑定需要服务器权限**
   - DNS 配置需要写入 `/etc/dnsmasq.d/` 或 `/etc/hosts`
   - Nginx 反代需要写入 `/etc/nginx/conf.d/`
   - 需要 systemctl reload 权限

2. **Windows 本地测试限制**
   - Windows 不支持 dnsmasq
   - Windows hosts 文件路径: `C:\Windows\System32\drivers\etc\hosts`
   - 建议在 Linux 服务器上部署测试

3. **代理商权限**
   - 用户可以在 `/agent/profile` 申请成为代理商
   - 管理员在 `/admin/agents` 中审核代理商申请
   - 只有审核通过的代理商才能访问代理商功能

## 下一步

1. 启动服务进行实际测试
2. 如果发现任何问题，记录错误信息
3. 根据错误信息进行针对性修复
4. 更新此检查清单

---

最后更新: 2026-03-12
