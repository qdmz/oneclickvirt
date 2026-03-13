# 系统初始化脚本

$baseURL = "http://localhost:8890/api/v1/public/init"

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
        username = "oneclickvirt"
        password = "123456"
    }
} | ConvertTo-Json -Depth 10

Write-Host "开始系统初始化..." -ForegroundColor Green
Write-Host "请求数据: $payload" -ForegroundColor Yellow

try {
    $response = Invoke-RestMethod -Uri $baseURL -Method Post -Body $payload -ContentType "application/json" -TimeoutSec 30

    if ($response.code -eq 0 -or $response.code -eq 200) {
        Write-Host "✅ 系统初始化成功！" -ForegroundColor Green
        Write-Host "消息: $($response.data)" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "现在可以登录了：" -ForegroundColor Cyan
        Write-Host "管理员: http://localhost:8080/#/admin/login" -ForegroundColor White
        Write-Host "  用户名: admin" -ForegroundColor White
        Write-Host "  密码: admin123" -ForegroundColor White
        Write-Host ""
        Write-Host "普通用户: http://localhost:8080/#/login" -ForegroundColor White
        Write-Host "  用户名: user" -ForegroundColor White
        Write-Host "  密码: user123" -ForegroundColor White
    } else {
        Write-Host "❌ 初始化失败: $($response.msg)" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ 发生错误: $_" -ForegroundColor Red
    Write-Host "详细信息: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "按任意键退出..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
