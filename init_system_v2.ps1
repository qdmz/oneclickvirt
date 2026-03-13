# 系统初始化脚本 - 使用 root 账户

$baseURL = "http://localhost:8890/api/v1/public/init"

# 先创建数据库
Write-Host "正在创建数据库 oneclickvirt..." -ForegroundColor Yellow
$mysqlPath = "C:\Program Files\MySQL\MySQL Server 9.6\bin\mysql.exe"

# 尝试连接并创建数据库（使用空密码的 root）
try {
    $createDB = @"
CREATE DATABASE IF NOT EXISTS oneclickvirt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
FLUSH PRIVILEGES;
"@

    $executeResult = & $mysqlPath -u root -e $createDB 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ 数据库创建成功" -ForegroundColor Green
    } else {
        Write-Host "⚠️ 使用 root 创建数据库可能需要密码，尝试直接初始化..." -ForegroundColor Yellow
    }
} catch {
    Write-Host "⚠️ 数据库创建跳过，继续初始化..." -ForegroundColor Yellow
}

# 使用 root 账户初始化
Write-Host "开始系统初始化（使用 root 账户）..." -ForegroundColor Green

$payload = @{
    admin = @{
        username = "admin"
        password = "admin123"
        email = "admin@example.com"
    }
    user = @{
        username = "user"
        password = "user123"
        email = "user@example.com"
    }
    database = @{
        type = "mysql"
        host = "localhost"
        port = "3306"
        database = "oneclickvirt"
        username = "root"
        password = ""
    }
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri $baseURL -Method Post -Body $payload -ContentType "application/json" -TimeoutSec 30

    if ($response.code -eq 0 -or $response.code -eq 200) {
        Write-Host "✅ 系统初始化成功！" -ForegroundColor Green
        Write-Host "消息: $($response.data)" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "现在可以登录了！" -ForegroundColor Cyan
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "👨‍💼 管理员登录:" -ForegroundColor White
        Write-Host "   地址: http://localhost:8080/#/admin/login" -ForegroundColor Gray
        Write-Host "   用户名: admin" -ForegroundColor Gray
        Write-Host "   密码: admin123" -ForegroundColor Gray
        Write-Host ""
        Write-Host "👤 普通用户登录:" -ForegroundColor White
        Write-Host "   地址: http://localhost:8080/#/login" -ForegroundColor Gray
        Write-Host "   用户名: user" -ForegroundColor Gray
        Write-Host "   密码: user123" -ForegroundColor Gray
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Cyan
    } else {
        Write-Host "❌ 初始化失败: $($response.msg)" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ 发生错误: $_" -ForegroundColor Red
    Write-Host "" -ForegroundColor Red
    Write-Host "如果是数据库连接错误，请检查：" -ForegroundColor Yellow
    Write-Host "1. MySQL 服务是否运行" -ForegroundColor Yellow
    Write-Host "2. root 用户密码是否正确" -ForegroundColor Yellow
    Write-Host "3. 可以手动访问初始化页面: http://localhost:8080/#/init" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "按任意键退出..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
