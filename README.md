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

```
前端地址: http://localhost:8080
后端API:  http://localhost:8890

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
- ✅ **第三方支付配置** - 支持易支付和码支付接口配置

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

## 📞 联系方式

如有问题或建议，请创建Issue或联系项目维护者。

**祝您使用愉快!** 🎉
