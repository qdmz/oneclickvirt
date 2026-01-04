# 快速部署和测试指南

## 一、环境要求

### 后端环境
- Go 1.24+ (从阿里云镜像获取: https://mirrors.aliyun.com/golang/)
- MySQL 5.7+
- 8888端口可用

### 前端环境
- Node.js 16+
- npm 或 yarn
- 80端口可用

---

## 二、安装Go环境

### Windows系统 (使用阿里云镜像)

#### 方法1: 下载安装包
1. 访问 https://mirrors.aliyun.com/golang/
2. 下载对应版本的 Windows 安装包
3. 解压到 `C:\Go`
4. 配置环境变量:
   ```
   GOPATH=C:\Go\projects
   GOROOT=C:\Go
   PATH=%PATH%;%GOROOT%\bin
   ```

#### 方法2: 使用chocolatey
```powershell
# 设置阿里云镜像
$env:GOPROXY = "https://mirrors.aliyun.com/goproxy/,direct"

# 安装go
choco install golang
```

#### 方法3: 使用winget
```powershell
winget install GoLang.Go --source winget
```

### 验证安装
```bash
go version
go env
```

---

## 三、数据库配置

### MySQL安装

#### Windows
```bash
# 下载MySQL 8.0
# https://dev.mysql.com/downloads/mysql/

# 或者使用MySQL Installer图形界面安装
```

### 创建数据库
```sql
CREATE DATABASE oneclickvirt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 配置数据库连接

编辑 `server/config.yaml`:
```yaml
mysql:
  db-name: oneclickvirt
  port: 3306
  username: root
  password: your_password
  auto-create: "true"
```

---

## 四、后端启动

### 安装依赖
```bash
cd server

# 设置Go代理(国内加速)
go env -w GOPROXY=https://goproxy.cn,direct

# 下载依赖
go mod download
```

### 启动服务
```bash
# 开发模式
go run main.go

# 或编译后运行
go build -o oneclivkvirt.exe main.go
./oneclivkvirt.exe
```

### 验证后端启动
访问: http://localhost:8888/health

应返回: `{"code":200,"message":"success"}`

---

## 五、前端启动

### 安装依赖
```bash
cd web

# 使用国内npm镜像加速
npm config set registry https://registry.npmmirror.com

# 安装依赖
npm install
```

### 启动开发服务器
```bash
npm run dev
```

### 访问前端
打开浏览器访问: http://localhost:5173

---

## 六、功能测试流程

### 第一步: 系统初始化

1. 访问 http://localhost:5173/init
2. 填写数据库连接信息
3. 创建管理员账号
4. 完成初始化

### 第二步: 管理员配置

1. 使用管理员账号登录
2. 进入"站点配置"页面
3. 配置网站信息:
   - 网站名称: OneClickVirt
   - 网站URL: http://localhost:5173
   - 网站图标: /favicon.ico
   - 页眉/页脚内容

4. 进入"产品管理"页面
5. 创建测试产品:
   - 名称: 基础套餐
   - 等级: 1
   - 价格: 10元
   - 有效期: 30天
   - CPU: 1核
   - 内存: 512MB
   - 磁盘: 10GB
   - 带宽: 100Mbps
   - 流量: 100GB
   - 最大实例: 1

6. 进入"兑换码管理"页面
7. 生成兑换码:
   - 数量: 5个
   - 类型: 余额
   - 金额: 10元
   - 使用次数: 1

### 第三步: 用户功能测试

1. 退出管理员账号
2. 注册新用户账号
3. 登录新用户账号

#### 测试兑换码功能
1. 进入"钱包"页面
2. 点击"使用兑换码"
3. 输入刚才生成的兑换码
4. 验证钱包余额增加10元

#### 测试充值功能
1. 在"钱包"页面
2. 选择充值金额: 50元
3. 选择支付方式: 支付宝
4. 点击"立即充值"
5. 查看生成的订单和二维码(模拟)

#### 测试产品购买
1. 进入"产品购买"页面
2. 查看"基础套餐"
3. 点击"立即购买"
4. 选择支付方式: 余额
5. 确认购买
6. 验证用户等级提升

#### 查看订单和交易记录
1. 进入"我的订单"页面
2. 查看所有订单记录
3. 点击详情查看订单信息
4. 进入"钱包"页面
5. 查看交易记录列表

### 第四步: 管理员订单管理

1. 使用管理员账号登录
2. 进入"订单管理"页面
3. 查看所有订单
4. 点击订单详情
5. 查看订单状态和支付信息

---

## 七、常见问题排查

### 问题1: 数据库连接失败
```
错误: Error 2003: Can't connect to MySQL server
解决:
1. 检查MySQL服务是否启动
2. 检查端口3306是否开放
3. 检查用户名密码是否正确
4. 检查数据库名称是否存在
```

### 问题2: 端口被占用
```
错误: bind: address already in use
解决:
# Windows查看端口占用
netstat -ano | findstr :8888
# 杀死进程
taskkill /PID <进程ID> /F
```

### 问题3: 前端无法访问后端API
```
错误: Network Error
解决:
1. 检查后端是否正常启动
2. 检查CORS配置
3. 检查API地址配置 (web/src/utils/request.js)
```

### 问题4: npm安装失败
```
解决:
# 使用淘宝镜像
npm config set registry https://registry.npmmirror.com

# 清理缓存
npm cache clean --force

# 删除node_modules重新安装
rm -rf node_modules
npm install
```

### 问题5: Go模块下载失败
```
解决:
# 使用国内Go代理
go env -w GOPROXY=https://goproxy.cn,direct

# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download
```

---

## 八、支付功能集成说明

### 当前状态
- ✅ 完整的订单系统
- ✅ 完整的钱包系统
- ✅ 完整的兑换码系统
- ✅ 支付接口框架
- ⚠️ 支付回调为模拟数据(需要真实SDK集成)

### 集成真实支付

#### 支付宝集成
1. 注册支付宝开放平台账号
2. 创建应用获取AppID
3. 配置公钥私钥
4. 下载Go SDK:
   ```bash
   go get github.com/smartwalle/alipay
   ```
5. 修改 `server/api/v1/payment/callback.go`
6. 替换模拟数据为真实SDK调用

#### 微信支付集成
1. 注册微信支付商户号
2. 获取商户号和API密钥
3. 下载Go SDK:
   ```bash
   go get github.com/wechatpay-apiv3/wechatpay-go
   ```
4. 修改 `server/api/v1/payment/callback.go`
5. 替换模拟数据为真实SDK调用

---

## 九、性能优化建议

### 后端优化
1. 启用数据库连接池
2. 使用Redis缓存热点数据
3. 添加API限流
4. 异步处理支付回调
5. 添加定时任务清理过期订单

### 前端优化
1. 路由懒加载
2. 图片懒加载
3. 组件按需引入
4. 启用Gzip压缩
5. CDN加速静态资源

---

## 十、生产环境部署

### 后端部署
```bash
# 编译生产版本
go build -ldflags="-s -w" -o oneclivkvirt main.go

# 使用systemd管理(Linux)
sudo nano /etc/systemd/system/oneclivkvirt.service

[Unit]
Description=OneClickVirt Server
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/oneclivkvirt
ExecStart=/var/www/oneclivkvirt/oneclivkvirt
Restart=on-failure

[Install]
WantedBy=multi-user.target

# 启动服务
sudo systemctl enable oneclivkvirt
sudo systemctl start oneclivkvirt
```

### 前端部署
```bash
# 构建生产版本
cd web
npm run build

# 使用nginx配置
sudo nano /etc/nginx/sites-available/oneclivkvirt

server {
    listen 80;
    server_name your-domain.com;

    root /var/www/oneclivkvirt/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8888;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# 启用配置
sudo ln -s /etc/nginx/sites-available/oneclivkvirt /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 十一、备份和恢复

### 数据库备份
```bash
# 备份
mysqldump -u root -p oneclickvirt > backup_$(date +%Y%m%d).sql

# 恢复
mysql -u root -p oneclickvirt < backup_20231201.sql
```

### 配置文件备份
```bash
# 备份配置
cp server/config.yaml server/config.yaml.backup

# 恢复配置
cp server/config.yaml.backup server/config.yaml
```

---

## 十二、日志查看

### 后端日志
```bash
# 实时查看日志
tail -f logs/server.log

# 查看错误日志
grep ERROR logs/server.log

# 查看支付相关日志
grep "payment\|order" logs/server.log
```

### 前端日志
打开浏览器开发者工具 Console 标签查看

---

## 十三、监控告警

### 推荐监控指标
1. API响应时间
2. 数据库连接数
3. 订单创建量
4. 支付成功率
5. 兑换码使用量
6. 钱包余额变化
7. 用户活跃度

### 告警配置
- 订单创建失败率 > 5%
- 支付回调失败
- 数据库连接超时
- API响应时间 > 3s

---

## 快速命令参考

### 启动服务
```bash
# 启动后端
cd server && go run main.go

# 启动前端(新终端)
cd web && npm run dev
```

### 重置数据库
```sql
DROP DATABASE oneclickvirt;
CREATE DATABASE oneclickvirt CHARACTER SET utf8mb4;
```

### 查看日志
```bash
# Windows
type logs\server.log

# Linux/Mac
tail -f logs/server.log
```

### 杀死进程
```bash
# Windows
tasklist | findstr oneclivkvirt
taskkill /PID <进程ID> /F

# Linux/Mac
ps aux | grep oneclivkvirt
kill -9 <进程ID>
```

---

## 技术支持

如有问题请查看:
1. 项目README.md
2. API文档: http://localhost:8888/swagger/index.html
3. 源码注释
4. GitHub Issues
