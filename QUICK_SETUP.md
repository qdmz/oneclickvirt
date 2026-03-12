# 快速测试 setup

## 方案 A: 使用 Docker MySQL（推荐，最简单）

```powershell
# 1. 安装 Docker Desktop（如果没有）
# https://www.docker.com/products/docker-desktop/

# 2. 启动 MySQL 容器
docker run -d --name oneclickvirt-mysql ^
  -e MYSQL_ROOT_PASSWORD=123456 ^
  -e MYSQL_DATABASE=oneclickvirt ^
  -p 3306:3306 ^
  mysql:8.0

# 3. 等待 MySQL 启动（约 10 秒）
# 可以检查日志：docker logs -f oneclickvirt-mysql

# 4. 在浏览器访问 http://localhost:8080
# 填写数据库配置：
# - 类型: mysql
# - 主机: localhost
# - 端口: 3306
# - 数据库: oneclickvirt
# - 用户名: root
# - 密码: 123456

# 5. 填写管理员和默认用户信息，然后初始化
```

## 方案 B: 使用 SQLite（无需外部数据库）

修改 `server/config.yaml`：

```yaml
mysql:
  db-type: sqlite  # 改为 sqlite
  path: ./data/oneclickvirt.db  # SQLite 文件路径
```

然后重启后端服务。

## 方案 C: 使用已有的 MySQL 服务

如果本地已有 MySQL：

```powershell
# 1. 创建数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS oneclickvirt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 2. 在浏览器 http://localhost:8080 填写配置
```

## 初始化完成后

系统会自动创建：
- 管理员账号
- 默认用户账号
- 所有数据表
- 系统配置

### 测试账号

**管理员账号**（初始化时填写）：
- 路径: http://localhost:8080/#/admin/login
- 可以访问所有新增功能：
  - 代理商管理
  - 实名管理
  - 域名配置
  - 域名管理

**用户账号**（初始化时填写）：
- 路径: http://localhost:8080/#/login
- 可以访问：
  - 域名管理
  - 钱包
  - 订单
  - 实名认证
  - 其他用户功能

## 检查侧边栏菜单

登录后，检查以下菜单项是否显示：

### 管理员侧边栏
- [ ] 代理商管理 (`/admin/agents`)
- [ ] 实名管理 (`/admin/kyc`)
- [ ] 域名配置 (`/admin/domain-config`)
- [ ] 域名管理 (`/admin/domains`)

### 用户侧边栏
- [ ] 域名管理 (`/user/domains`)
- [ ] 实名认证 (`/user/kyc`)
- [ ] 钱包 (`/user/wallet`)
- [ ] 订单 (`/user/orders`)

## 前端主题测试

查看深色/浅色主题切换：
1. 登录后
2. 点击右上角主题切换按钮
3. 界面切换主题

## 测试域名功能（需要 Linux 服务器）

域名绑定功能需要在有 dnsmasq/nginx 的 Linux 服务器上测试：
- Windows 本地测试会显示功能界面但无法实际配置 DNS
- 建议在虚拟机或云服务器上完整测试
