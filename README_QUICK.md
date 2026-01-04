# 🚀 快速部署指南

## 问题诊断: http://localhost:5173 打不开

**原因**: Go和Node.js环境未安装,服务未启动

---

## 方案1: 本地快速启动 (推荐)

### 步骤1: 安装环境 (必需)

#### 安装Go (使用阿里云镜像)

打开PowerShell (管理员)执行:

```powershell
# 下载Go
Invoke-WebRequest -Uri "https://mirrors.aliyun.com/golang/go1.22.1.windows-amd64.zip" -OutFile "$env:TEMP\go.zip"

# 解压到C:\Go
Expand-Archive -Path "$env:TEMP\go.zip" -DestinationPath "C:\"

# 添加到系统PATH
[Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path', 'User') + ';C:\Go\bin', 'User')

# 验证
go version
```

#### 安装Node.js

访问 https://nodejs.org/ 下载安装LTS版本

---

### 步骤2: 一键启动

双击运行: `auto-install-and-start.bat`

或手动执行:

```cmd
cd c:\Users\Administrator\CodeBuddy\code1\oneclickvirt
auto-install-and-start.bat
```

脚本会自动:
- 检查Go和Node环境
- 编译后端
- 准备前端
- 启动服务
- 打开浏览器

---

### 步骤3: 访问系统

浏览器打开: **http://localhost:5173**

登录:
```
用户名: admin
密码:   admin123456
```

---

## 方案2: 云端部署 (轻量云)

### 上传项目到服务器

#### 方法1: 使用scp上传

```bash
# 在本地执行
scp -r c:/Users/Administrator/CodeBuddy/code1/oneclickvirt root@your-server-ip:/opt/
```

#### 方法2: 使用Git

```bash
# 在服务器上执行
cd /opt
git clone your-repo-url oneclickvirt
cd oneclickvirt
```

#### 方法3: 手动上传

1. 将整个 `oneclickvirt` 文件夹压缩
2. 通过SFTP或云控制台上传到 `/opt/oneclickvirt`
3. 解压: `tar -xzf oneclickvirt.tar.gz -C /opt/`

---

### 运行部署脚本

在服务器上执行:

```bash
cd /opt/oneclickvirt
chmod +x deploy-to-cloud.sh
sudo ./deploy-to-cloud.sh
```

脚本会自动:
- 安装Go (从阿里云镜像)
- 安装Node.js
- 安装Nginx
- 配置项目
- 启动服务
- 配置防火墙

---

### 访问云端系统

```
前端: http://your-server-ip
后端: http://your-server-ip:8888

登录账号:
  用户名: admin
  密码:   admin123456
```

---

## 🔧 故障排除

### 本地启动问题

#### 问题: 端口被占用

```cmd
# 查看占用端口的进程
netstat -ano | findstr "5173"
netstat -ano | findstr "8888"

# 结束进程
taskkill /F /PID <进程ID>
```

#### 问题: Go命令找不到

重新打开CMD/PowerShell窗口,让PATH生效

#### 问题: npm安装失败

```cmd
cd web
rmdir /s /q node_modules
del package-lock.json
npm install
```

---

### 云端部署问题

#### 问题: 端口无法访问

1. 检查服务状态:
```bash
systemctl status oneclickvirt
systemctl status nginx
```

2. 检查防火墙:
```bash
# Ubuntu/Debian
sudo ufw status

# CentOS/RHEL
sudo firewall-cmd --list-all
```

3. **重要**: 在云服务器控制台开放端口:
   - 端口: 80, 443
   - 协议: TCP

#### 问题: 服务启动失败

查看日志:
```bash
journalctl -u oneclickvirt -n 50
```

---

## 📱 测试清单

### 本地测试

- [ ] Go环境已安装 (`go version`)
- [ ] Node环境已安装 (`node --version`)
- [ ] 后端启动成功 (看到 "Listening on :8888")
- [ ] 前端启动成功 (看到 "Local: http://localhost:5173")
- [ ] 能访问 http://localhost:5173
- [ ] 能成功登录

### 云端测试

- [ ] 项目已上传到服务器
- [ ] 部署脚本执行成功
- [ ] 后端服务运行中 (`systemctl is-active oneclickvirt`)
- [ ] Nginx服务运行中 (`systemctl is-active nginx`)
- [ ] 防火墙已配置或安全组已开放
- [ ] 能访问 http://your-server-ip
- [ ] 能成功登录

---

## 🎯 功能测试

登录后测试新功能:

### 管理员功能
- [ ] 站点配置 (http://localhost:5173/admin/site-config)
- [ ] 产品管理 (http://localhost:5173/admin/products)
- [ ] 兑换码管理 (http://localhost:5173/admin/redemption-codes)
- [ ] 订单管理 (http://localhost:5173/admin/orders)

### 用户功能
- [ ] 钱包 (http://localhost:5173/user/wallet)
- [ ] 产品购买 (http://localhost:5173/user/purchase)
- [ ] 我的订单 (http://localhost:5173/user/orders)

---

## 📞 获取帮助

- 本地启动失败: 查看 `auto-install-and-start.bat` 错误信息
- 云端部署失败: 查看脚本输出和日志
- 系统登录问题: 检查后端服务是否启动

---

## 💡 提示

1. **首次启动较慢**: 编译和安装依赖需要时间
2. **不要关闭窗口**: 后端和前端服务需要保持运行
3. **日志查看**: 后端日志在 `server/storage/logs/`
4. **数据存储**: SQLite数据库在 `server/storage/oneclickvirt.db`

---

**现在开始安装环境并启动系统!** 🚀
