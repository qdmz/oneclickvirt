# Login Issue Resolution - 网络连接失败

## 问题原因
前端代理配置正确，但可能有缓存问题。

## 解决方案

### 方法 1：清除浏览器缓存并硬刷新（最快）

1. 按 `Ctrl + Shift + Delete` 打开清除浏览器数据
2. 选择"缓存的图片和文件"
3. 时间范围选"全部"
4. 点击"清除数据"

然后硬刷新页面：
- Windows: `Ctrl + Shift + R`
- Mac: `Cmd + Shift + R`

### 方法 2：使用隐身/无痕模式
- 按 `Ctrl + Shift + N` (Chrome) 或 `Ctrl + Shift + P` (Firefox)
- 访问 http://localhost:8080
- 尝试登录

### 方法 3：检查前端服务
前端应该运行在:
- 本地: http://localhost:8080
- 网络: http://192.168.11.3:8080

如果前端服务停止了，需要重启:
```powershell
cd C:\Users\admin\.openclaw-autoclaw\workspace\oneclickvirt\web
npm run dev
```

## 备用方案：直接使用 localhost

建议使用本地地址访问，避免网络问题:
```
http://localhost:8080/#/admin/login
```

## 验证 API 正常
测试 API:
- 访问 http://localhost:8890/api/v1/public/init/check
- 应该返回: `{"code":0,"data":{"message":"数据库已初始化","needInit":false},"msg":"success"}`

## 登录信息

管理员登录:
- URL: http://localhost:8080/#/admin/login
- Username: admin
- Password: admin123

用户登录:
- URL: http://localhost:8080/#/login
- Username: user
- Password: user123

## 登录后检查清单
- [ ] 代理商管理菜单
- [ ] 实名管理菜单
- [ ] 域名配置菜单
- [ ] 域名管理菜单
- [ ] 主题切换功能
- [ ] 深色/浅色主题
