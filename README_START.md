# 🚀 OneClickVirt 本地部署完成

## ✅ 已完成的工作

### 1. Go环境配置
- ✅ 已配置Go国内镜像: `https://goproxy.cn,direct`
- ✅ 后端编译成功: `server/oneclickvirt.exe`

### 2. 菜单配置
- ✅ 用户钱包 (`/user/wallet`)
- ✅ 我的订单 (`/user/orders`)
- ✅ 产品购买 (`/user/purchase`)
- ✅ 管理员订单管理 (`/admin/orders`)

### 3. 翻译配置
- ✅ 中英文翻译已配置

## 📋 启动步骤

### 方式1: 一键启动(推荐)

```cmd
cd c:\Users\Administrator\CodeBuddy\code1\oneclickvirt
start.bat
```

### 方式2: 分步启动

#### 步骤1: 安装前端依赖(首次运行)

```cmd
cd c:\Users\Administrator\CodeBuddy\code1\oneclickvirt
install-web.bat
```

#### 步骤2: 启动服务

```cmd
start.bat
```

### 方式3: 手动启动

#### 启动后端(窗口1):
```cmd
cd c:\Users\Administrator\CodeBuddy\code1\oneclickvirt\server
oneclickvirt.exe
```

#### 启动前端(窗口2):
```cmd
cd c:\Users\Administrator\CodeBuddy\code1\oneclickvirt\web
npm run dev
```

## 🔐 登录信息

```
用户名: admin
密码:   admin123456

前端地址: http://localhost:5173
后端地址: http://localhost:8888
```

*登录信息保存在: `login-info.txt`*

## 📱 新功能测试清单

### 登录后检查菜单:

#### 普通用户视图:
- [ ] 仪表盘
- [ ] 我的实例
- [ ] 申请领取
- [ ] 任务列表
- [ ] **钱包** ⭐ 新增
- [ ] **我的订单** ⭐ 新增
- [ ] **产品购买** ⭐ 新增
- [ ] 个人中心

#### 管理员视图:
- [ ] 仪表盘
- [ ] 用户管理
- [ ] 邀请码管理
- [ ] 节点管理
- [ ] 任务管理
- [ ] 实例管理
- [ ] **订单管理** ⭐ 新增
- [ ] 流量管理
- [ ] 端口管理
- [ ] 系统镜像
- [ ] 公告管理
- [ ] OAuth2
- [ ] 系统配置
- [ ] 性能监控

## ⚠️ 注意事项

1. **首次启动需要安装前端依赖**
   - 运行 `install-web.bat` 或 `start.bat`
   - 安装过程可能需要几分钟
   - 之后可以直接运行 `start.bat`

2. **浏览器缓存**
   - 如果菜单看不到,清除浏览器缓存
   - 或使用无痕/隐私模式访问
   - 或按 Ctrl+F5 强制刷新

3. **服务端口**
   - 前端: 5173
   - 后端: 8888
   - 确保端口没有被占用

4. **关闭服务**
   - 直接关闭两个黑色的命令窗口即可

## 📁 相关文件

### 启动脚本:
- `start.bat` - 一键启动服务(自动检测依赖)
- `install-web.bat` - 安装前端依赖

### 配置文件:
- `login-info.txt` - 登录信息

### 源代码:
- `server/oneclickvirt.exe` - 后端可执行文件
- `web/` - 前端源码

## 🎯 快速开始

1. **启动服务**:
   ```cmd
   cd c:\Users\Administrator\CodeBuddy\code1\oneclickvirt
   start.bat
   ```

2. **等待服务启动**:
   - 会出现两个黑色命令窗口
   - 等待几秒钟让服务完全启动

3. **访问系统**:
   - 打开浏览器: http://localhost:5173
   - 使用管理员账号登录:
     - 用户名: `admin`
     - 密码: `admin123456`

4. **测试新功能**:
   - 查看侧边栏菜单
   - 点击"钱包"、"我的订单"、"产品购买"
   - 切换到管理员视图查看"订单管理"

## 🐛 故障排除

### 问题1: 前端依赖安装失败

**解决**:
```cmd
cd web
npm config set registry https://registry.npmmirror.com
npm install
```

### 问题2: 菜单看不到新功能

**解决**:
1. 清除浏览器缓存
2. 使用无痕模式访问
3. 按 Ctrl+F5 强制刷新
4. 检查前端服务是否正常运行

### 问题3: 端口被占用

**错误**: `bind: address already in use`

**解决**:
- 关闭占用8888或5173端口的程序
- 或修改配置文件中的端口号

## 📞 技术支持

如果遇到问题,请查看:
- 登录信息: `login-info.txt`
- 快速启动指南: `QUICKSTART.md`
- 部署指南: `LOCAL_DEPLOY_GUIDE.md`
- 功能说明: `FEATURES.md`

## ✨ 功能特性

- ✅ 用户钱包系统
- ✅ 订单管理系统
- ✅ 产品购买功能
- ✅ 管理员订单管理
- ✅ 支持中英文切换
- ✅ 响应式设计

---

**祝您使用愉快!** 🎉
