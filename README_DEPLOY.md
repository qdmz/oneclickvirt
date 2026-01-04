# OneClickVirt 部署指南

## 📋 概述

OneClickVirt 是一个支持多虚拟化平台的云资源管理系统,现已新增以下功能:

✅ **站点配置管理** - 网站/图标/页眉页脚
✅ **产品套餐管理** - 灵活的产品配置
✅ **充值支付系统** - 支付宝/微信/余额支付
✅ **钱包系统** - 余额管理和交易记录
✅ **兑换码系统** - 灵活的兑换码生成
✅ **订单管理系统** - 完整的订单流程

---

## 🚀 快速部署 (Windows本地)

### 方法1: 自动化部署 (推荐)

#### 步骤1: 安装环境

双击运行 `install-env.ps1`

或在PowerShell中执行:

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
.\install-env.ps1
```

脚本会自动:
- 检查并安装Go (从阿里云镜像下载)
- 检查并安装Node.js
- 验证环境安装

#### 步骤2: 启动系统

双击运行 `quick-test.bat`

或在CMD中执行:

```cmd
cd c:\Users\Administrator\CodeBuddy\code1\oneclickvirt
quick-test.bat
```

脚本会自动:
- 编译后端
- 生成管理员密码
- 启动前后端服务

### 方法2: 手动部署

#### 1. 安装Go环境

```powershell
# 下载Go (使用阿里云镜像)
Invoke-WebRequest -Uri "https://mirrors.aliyun.com/golang/go1.22.1.windows-amd64.zip" -OutFile "go.zip"

# 解压到C:\Go
Expand-Archive -Path go.zip -DestinationPath C:\

# 添加到PATH (永久生效)
[Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path', 'User') + ';C:\Go\bin', 'User')

# 验证安装
go version
```

#### 2. 安装Node.js

访问 https://nodejs.org/ 下载安装LTS版本

#### 3. 启动后端

打开**终端1**:

```cmd
cd c:\Users\Administrator\CodeBuddy\code1\oneclickvirt\server

# 下载依赖
go mod download

# 编译
go build -o oneclickvirt.exe main.go

# 运行
oneclickvirt.exe
```

#### 4. 启动前端

打开**终端2**:

```cmd
cd c:\Users\Administrator\CodeBuddy\code1\oneclickvirt\web

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

---

## 🌐 云端部署 (轻量云/腾讯云/阿里云)

### 自动化部署脚本 (Ubuntu/Debian)

上传项目到服务器后执行:

```bash
chmod +x deploy-cloud.sh
./deploy-cloud.sh
```

脚本会自动:
- 安装Go、Node.js、MySQL、Nginx
- 配置系统服务
- 生成登录信息
- 配置反向代理

### 手动部署

#### 1. 安装Go (使用阿里云镜像)

```bash
# 下载Go
wget https://mirrors.aliyun.com/golang/go1.22.1.linux-amd64.tar.gz

# 解压
tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
source /etc/profile

# 验证
go version
```

#### 2. 安装Node.js

```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt install -y nodejs

# 验证
node --version
npm --version
```

#### 3. 安装MySQL

```bash
apt install -y mysql-server

# 启动服务
systemctl start mysql
systemctl enable mysql

# 配置MySQL
mysql -u root
```

#### 4. 配置项目

```bash
# 克隆或上传项目
cd /opt/oneclickvirt

# 配置后端
cd server
go mod download
go build -o oneclickvirt main.go

# 修改config.yaml配置MySQL连接信息

# 启动后端
./oneclickvirt
```

#### 5. 配置前端

```bash
cd /opt/oneclickvirt/web

# 安装依赖
npm install

# 构建生产版本
npm run build
```

#### 6. 配置Nginx

```bash
# 创建Nginx配置
cat > /etc/nginx/sites-available/oneclickvirt << 'EOF'
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /opt/oneclickvirt/web/dist;
        try_files $uri $uri/ /index.html;
    }

    # 后端API代理
    location /api {
        proxy_pass http://127.0.0.1:8888;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
EOF

# 启用配置
ln -s /etc/nginx/sites-available/oneclickvirt /etc/nginx/sites-enabled/
nginx -t
systemctl restart nginx
```

#### 7. 创建系统服务

```bash
cat > /etc/systemd/system/oneclickvirt.service << 'EOF'
[Unit]
Description=OneClickVirt Backend
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/oneclickvirt/server
ExecStart=/opt/oneclickvirt/server/oneclickvirt
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable oneclickvirt
systemctl start oneclickvirt
```

---

## 🔐 登录信息

### 首次启动

系统首次启动时会自动创建管理员账号,登录信息保存在 `login-info.txt`:

```
管理员账号:
  用户名: admin
  密码:   [随机生成的12位密码]

访问地址:
  前端: http://localhost:5173
  后端: http://localhost:8888
```

### 默认账号

如果自动创建失败,可使用默认账号:

```
用户名: admin
密码:   admin123456
```

⚠️ **重要**: 首次登录后请立即修改密码!

---

## 📱 访问地址

### 本地部署

- **前端**: http://localhost:5173
- **后端**: http://localhost:8888

### 云端部署

- **前端**: http://your-server-ip 或 http://your-domain.com
- **后端**: http://your-server-ip:8888 (通过Nginx代理可关闭8888端口)

---

## ✨ 新功能使用指南

### 管理员功能

#### 1. 站点配置

访问: `/admin/site-config`

功能:
- 网站名称、URL、图标配置
- 页眉、页脚自定义内容
- 联系信息(邮箱/电话)
- 公司信息和ICP备案号

#### 2. 产品管理

访问: `/admin/products`

功能:
- 创建/编辑/删除产品套餐
- 配置CPU/内存/磁盘/带宽/流量
- 设置价格和有效期
- 对应用户等级(1-5级)
- 启用/禁用产品

#### 3. 兑换码管理

访问: `/admin/redemption-codes`

功能:
- 单个创建或批量生成兑换码
- 三种类型: 余额兑换/等级兑换/产品兑换
- 设置使用次数和过期时间
- 查看详细使用记录

#### 4. 订单管理

访问: `/admin/orders`

功能:
- 查看所有订单
- 订单详情和状态
- 订单筛选和搜索

### 用户功能

#### 1. 钱包

访问: `/user/wallet`

功能:
- 查看实时余额
- 累计充值/消费统计
- 交易记录(分页/筛选)
- 充值功能

#### 2. 产品购买

访问: `/user/purchase`

功能:
- 浏览所有可用产品
- 查看详细配置信息
- 选择支付方式(支付宝/微信/余额)
- 立即购买

#### 3. 我的订单

访问: `/user/orders`

功能:
- 查看我的订单列表
- 订单详情
- 订单状态筛选
- 取消待支付订单

---

## ⚙️ 配置说明

### 后端配置: `server/config.yaml`

```yaml
system:
  addr: 8888              # 后端端口
  db-type: sqlite         # 数据库类型 (sqlite/mysql)
  env: development        # 环境 (development/production)
  frontend-url: "http://localhost:5173"  # 前端地址

quota:
  default-level: 1        # 默认用户等级

auth:
  enable-public-registration: false  # 是否允许公开注册

zap:
  level: info             # 日志级别
  log-in-console: true    # 是否输出到控制台
```

### 前端配置: `web/.env.development`

```env
VITE_API_URL=http://localhost:8888
```

### 生产环境: `web/.env.production`

```env
VITE_API_URL=https://your-domain.com
```

---

## 🔧 常见问题

### 1. Go命令找不到

**错误**: `'go' 不是内部或外部命令`

**解决**:
- 确认Go已安装
- 检查PATH环境变量是否包含 `C:\Go\bin`
- 重启终端窗口

### 2. MySQL连接失败

**错误**: `Error 2003: Can't connect to MySQL server`

**解决**:
- 使用SQLite测试 (已配置)
- 或正确配置MySQL连接信息

### 3. 端口被占用

**错误**: `bind: address already in use`

**解决**:
- 修改 `config.yaml` 中的端口号
- 或关闭占用端口的程序

### 4. 前端依赖安装失败

**错误**: `npm install` 报错

**解决**:
```cmd
cd web
rmdir /s /q node_modules
del package-lock.json
npm install
```

### 5. 后端编译失败

**错误**: `go build` 报错

**解决**:
```cmd
cd server
go mod tidy
go build -v -o oneclickvirt.exe main.go
```

---

## 📝 数据库说明

### SQLite (本地测试)

- 位置: `server/storage/oneclickvirt.db`
- 优点: 无需额外安装,开箱即用
- 适用: 开发测试环境

### MySQL (生产环境)

配置 `config.yaml`:

```yaml
mysql:
  path: 127.0.0.1
  port: "3306"
  username: root
  password: your-password
  db-name: oneclickvirt
```

---

## 🔐 安全建议

1. **修改默认密码**: 首次登录后立即修改admin密码
2. **HTTPS配置**: 生产环境使用HTTPS加密
3. **防火墙规则**: 只开放必要端口(80, 443)
4. **数据库安全**: 使用强密码,限制远程访问
5. **定期备份**: 定期备份数据库和配置文件

---

## 📞 技术支持

### 文档

- **功能说明**: `FEATURES.md`
- **部署指南**: `DEPLOYMENT.md`
- **实现总结**: `IMPLEMENTATION_SUMMARY.md`
- **快速启动**: `QUICKSTART.md`

### 日志

- 后端日志: `server/storage/logs/`
- 系统日志: `journalctl -u oneclickvirt -f` (Linux)

---

## 📊 部署检查清单

### 本地部署

- [ ] Go已安装并验证
- [ ] Node.js已安装并验证
- [ ] config.yaml已配置
- [ ] 后端编译成功
- [ ] 前端依赖安装成功
- [ ] 后端服务正常启动
- [ ] 前端服务正常启动
- [ ] 能够成功登录系统

### 云端部署

- [ ] 服务器系统已更新
- [ ] Go已安装
- [ ] Node.js已安装
- [ ] MySQL已安装并配置
- [ ] 项目已上传
- [ ] 后端编译成功
- [ ] 前端构建成功
- [ ] Nginx已配置
- [ ] 系统服务已创建
- [ ] 防火墙规则已配置
- [ ] 域名已解析(如有)
- [ ] HTTPS已配置(推荐)

---

## 🎉 部署成功后

部署成功后,您将拥有一个完整的云资源管理系统:

### 管理功能
✅ 用户和权限管理
✅ 实例和资源管理
✅ Provider配置
✅ 站点配置
✅ 产品套餐管理
✅ 兑换码管理
✅ 订单管理

### 用户功能
✅ 实例创建和管理
✅ 终端访问
✅ 钱包管理
✅ 充值和支付
✅ 产品购买
✅ 订单查询

---

**祝您使用愉快!** 🎊
