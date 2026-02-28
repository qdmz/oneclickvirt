# OneClickVirt

## 📖 项目简介

OneClickVirt是一个现代化的虚拟服务管理平台，支持多种虚拟化技术，提供完整的产品管理、用户管理、订单管理和资源监控功能。

## 🚀 快速开始

### 本地开发

#### 前端开发
```bash
cd web
npm install
npm run dev
```

前端将运行在 `http://localhost:8080`

#### 后端开发
```bash
cd server
go mod download
go run main.go
```

后端将运行在 `http://localhost:8890`

### Docker部署

#### 使用docker-compose (推荐)

```bash
docker-compose -f docker-compose.yaml up -d
```

或

```bash
docker-compose up -d
```

#### 手动构建和运行

1. 构建镜像
```bash
docker build -t oneclickvirt .
```

2. 运行容器

**基本运行（HTTP）**
```bash
docker run -d --name oneclickvirt -p 80:80 -e FRONTEND_URL="http://your-domain.com" oneclickvirt
```

**使用HTTPS（自动生成自签名证书）**
```bash
docker run -d --name oneclickvirt -p 80:80 -p 443:443 -e FRONTEND_URL="https://your-domain.com" oneclickvirt
```

**使用自定义SSL证书**
```bash
docker run -d --name oneclickvirt \
  -p 80:80 -p 443:443 \
  -e FRONTEND_URL="https://your-domain.com" \
  -e SSL_CERT_PATH="/certs/cert.pem" \
  -e SSL_KEY_PATH="/certs/key.pem" \
  -v /path/to/your/certs:/certs \
  oneclickvirt
```

**数据持久化运行（推荐用于生产环境）**
```bash
# 创建数据目录
mkdir -p docker-data/mysql docker-data/storage docker-data/ssl

# 运行容器（带数据卷挂载）
docker run -d --name oneclickvirt \
  -p 80:80 -p 443:443 \
  -v $(pwd)/docker-data/mysql:/var/lib/mysql \
  -v $(pwd)/docker-data/storage:/app/storage \
  -v $(pwd)/docker-data/ssl:/etc/nginx/ssl \
  -v $(pwd)/docker-data/config.yaml:/app/config.yaml \
  oneclickvirt
```

**环境变量说明**
- `FRONTEND_URL`: 前端访问URL（如 `https://your-domain.com` 或 `http://your-domain.com`）
- `SSL_CERT_PATH`: SSL证书文件路径（可选）
- `SSL_KEY_PATH`: SSL私钥文件路径（可选）

**数据持久化说明**
- `docker-data/mysql`: MySQL数据库文件（用户数据、订单、配置等）
- `docker-data/storage`: 应用存储（日志、上传文件、缓存等）
- `docker-data/ssl`: SSL证书文件
- `docker-data/config.yaml`: 应用配置文件

⚠️ **重要提示**: 使用数据卷挂载后，即使删除容器，数据也会保留在宿主机上。重新创建容器时只需使用相同的 `-v` 参数即可恢复数据。

#### 手动构建

1. 构建前端
```bash
cd web
npm install
npm run build
```

2. 构建后端
```bash
cd server
go build -o oneclickvirt main.go
```

3. 运行服务
```bash
./server/oneclickvirt
```

## 🔐 登录信息

**本地开发环境:**
```
前端地址: http://localhost:8080
后端API:  http://localhost:8890

管理员账号:
  用户名: admin
  密码:   admin123456
```

**Docker部署环境:**
```
前端地址: http://your-domain.com (或 https://your-domain.com)
后端API:  自动配置，无需手动访问

管理员账号:
  用户名: admin
  密码:   admin123456
```

⚠️ 首次登录后建议立即修改密码!

## 🌐 部署选项

### 本地部署
- 适合开发和测试环境
- 快速启动，便于调试

### Docker部署
- 适合生产环境
- 包含完整的服务栈
- 便于扩展和维护

### 云服务器部署
- 适合线上生产环境
- 支持高可用配置

## ✨ 功能特性

### 管理员功能
- ✅ 站点配置管理
- ✅ 产品套餐管理
- ✅ 兑换码管理
- ✅ 订单管理
- ✅ 用户管理
- ✅ 资源监控
- ✅ 流量统计
- ✅ **代用户登录** - 管理员可直接以用户身份登录系统
- ✅ **实例转移归属** - 管理员可将实例从一个用户转移到另一个用户
- ✅ **第三方支付配置** - 支持易支付和码支付接口配置（详见下方配置说明）

### 用户功能
- ✅ 虚拟实例管理
- ✅ 产品购买
- ✅ 钱包管理
- ✅ 订单管理
- ✅ 流量监控
- ✅ **多种支付方式** - 支持支付宝、微信支付、余额支付、易支付、码支付

## 📁 项目结构

```
├── deploy/           # 部署相关配置
│   ├── default.conf     # Nginx配置
│   ├── my.cnf           # MySQL配置
│   ├── nginx.dockerfile # Nginx Dockerfile
│   └── server.dockerfile # 后端服务Dockerfile
├── server/           # 后端代码
│   ├── api/             # API路由
│   ├── config/          # 配置管理
│   ├── model/           # 数据模型
│   ├── provider/        # 虚拟化提供商
│   └── main.go          # 主入口
├── web/              # 前端代码
│   ├── src/             # 源码
│   ├── Dockerfile       # 前端Dockerfile
│   └── package.json     # 依赖管理
├── docker-compose.yaml # Docker Compose配置
└── README.md          # 项目说明
```

## 🔧 技术栈

### 后端
- Go 1.24
- Gin Web框架
- GORM ORM框架
- MySQL / SQLite
- Redis (可选)

### 前端
- Vue 3
- Element Plus UI框架
- Vite构建工具

### 新增技术组件
- **JWT Token认证** - 支持管理员代用户登录的权限降级
- **MD5签名验证** - 第三方支付回调安全验证
- **异步支付处理** - 支持多平台支付回调处理

## 📝 开发规范

### 后端
- 遵循Go语言最佳实践
- 分层架构设计
- RESTful API设计
- 详细的日志记录

### 前端
- 组件化开发
- TypeScript支持
- 响应式设计
- 现代化UI风格

## 🔒 安全措施

- JWT身份认证
- 密码哈希存储
- 权限控制
- 防止SQL注入
- 防止XSS攻击
- **支付签名验证** - MD5签名确保支付回调安全性
- **权限降级机制** - 管理员代用户登录时权限安全控制
- **操作审计日志** - 记录敏感操作便于追踪

## 🤝 贡献指南

欢迎提交Issue和Pull Request!

1. Fork项目
2. 创建特性分支
3. 提交修改
4. 推送分支
5. 创建Pull Request

## 📄 许可证

MIT License

## 💳 支付配置说明

### 易支付 (Epay) 配置

在后台管理 → 系统设置 → 支付配置中配置以下参数：

| 参数 | 说明 | 示例 |
|------|------|------|
| `商户ID (PID)` | 易支付平台分配的商户ID | `1234` |
| `商户密钥 (Key)` | 易支付平台分配的商户密钥 | `your_secret_key_here` |
| `API接口地址` | 易支付平台的API地址 | `https://pay.example.com/` |
| **回调URL** | 支付成功后的异步通知地址 | `https://your-domain.com/api/v1/payment/epay/notify` |
| **返回URL** | 支付成功后跳转的页面地址 | `https://your-domain.com/user/wallet` |

⚠️ **重要提示**:
- 回调URL必须是外网可访问的HTTPS地址
- 回调URL末尾**不要**加空格或其他字符
- 返回URL通常是用户钱包页面或订单页面
- 确保易支付平台的回调白名单中已添加您的服务器IP

### 码支付 (Mapay) 配置

| 参数 | 说明 | 示例 |
|------|------|------|
| `商户ID` | 码支付平台分配的商户ID | `5678` |
| `商户密钥` | 码支付平台分配的商户密钥 | `your_mapay_key_here` |
| `API接口地址` | 码支付平台的API地址 | `https://mapay.example.com/` |
| **回调URL** | 支付成功后的异步通知地址 | `https://your-domain.com/api/v1/payment/mapay/notify` |
| **返回URL** | 支付成功后跳转的页面地址 | `https://your-domain.com/user/wallet` |

### 支付配置检查清单

- [ ] 已启用对应的支付方式（易支付/码支付）
- [ ] 已正确填写商户ID和密钥
- [ ] API接口地址以 `/` 结尾
- [ ] 回调URL格式正确且可访问
- [ ] 返回URL指向正确的页面
- [ ] 支付平台已配置回调白名单

## 📞 联系方式

如有问题或建议，请创建Issue或联系项目维护者。

**祝您使用愉快!** 🎉
