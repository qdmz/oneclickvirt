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

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "OneClickVirt System Initialization" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Initializing system with MySQL..." -ForegroundColor Yellow

try {
    $response = Invoke-RestMethod -Uri $baseURL -Method Post -Body $payload -ContentType "application/json" -TimeoutSec 45

    if ($response.code -eq 0 -or $response.code -eq 200) {
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Green
        Write-Host "SUCCESS! System initialized!" -ForegroundColor Green
        Write-Host "========================================" -ForegroundColor Green
        Write-Host ""
        Write-Host "Dashboard: http://localhost:8080" -ForegroundColor White
        Write-Host ""
        Write-Host "Admin Login:" -ForegroundColor Cyan
        Write-Host "  URL: http://localhost:8080/#/admin/login" -ForegroundColor Gray
        Write-Host "  Username: admin" -ForegroundColor Gray
        Write-Host "  Password: admin123" -ForegroundColor Gray
        Write-Host ""
        Write-Host "User Login:" -ForegroundColor Cyan
        Write-Host "  URL: http://localhost:8080/#/login" -ForegroundColor Gray
        Write-Host "  Username: user" -ForegroundColor Gray
        Write-Host "  Password: user123" -ForegroundColor Gray
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "Test Features:" -ForegroundColor Cyan
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "Admin菜单中应该看到:" -ForegroundColor White
        Write-Host "  - 代理商管理" -ForegroundColor Gray
        Write-Host "  - 实名管理" -ForegroundColor Gray
        Write-Host "  - 域名配置" -ForegroundColor Gray
        Write-Host "  - 域名管理" -ForegroundColor Gray
        Write-Host ""
        Write-Host "User菜单中应该看到:" -ForegroundColor White
        Write-Host "  - 域名管理" -ForegroundColor Gray
        Write-Host "  - 实名认证" -ForegroundColor Gray
        Write-Host "  - 钱包" -ForegroundColor Gray
        Write-Host ""
        Write-Host "点击右上角按钮切换深色主题!" -ForegroundColor Yellow
        Write-Host "========================================" -ForegroundColor Cyan
    } else {
        Write-Host ""
        Write-Host "FAILED: $($response.msg)" -ForegroundColor Red
        Write-Host "Code: $($response.code)" -ForegroundColor Red
        Write-Host ""
        Write-Host "Please refresh browser and try manual init:" -ForegroundColor Yellow
        Write-Host "http://localhost:8080/#/init" -ForegroundColor White
    }
} catch {
    Write-Host ""
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
    Write-Host "Status Code: $($_.Exception.Response.StatusCode)" -ForegroundColor Yellow

    if ($_.Exception.Response.StatusCode -eq 403) {
        Write-Host ""
        Write-Host "403 Error - Permission denied" -ForegroundColor Yellow
        Write-Host "This is a permission/security issue in the backend" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "Workaround: Use manual initialization in browser" -ForegroundColor Cyan
        Write-Host "http://localhost:8080/#/init" -ForegroundColor White
        Write-Host ""
        Write-Host "Database config:" -ForegroundColor White
        Write-Host "  Type: mysql" -ForegroundColor Gray
        Write-Host "  Host: localhost" -ForegroundColor Gray
        Write-Host "  Port: 3306" -ForegroundColor Gray
        Write-Host "  Database: oneclickvirt" -ForegroundColor Gray
        Write-Host "  Username: oneclickvirt" -ForegroundColor Gray
        Write-Host "  Password: 123456" -ForegroundColor Gray
    } else {
        Write-Host ""
        Write-Host "Please check:" -ForegroundColor Yellow
        Write-Host "1. Backend service is running (port 8890)" -ForegroundColor Yellow
        Write-Host "2. MySQL service is running (port 3306)" -ForegroundColor Yellow
        Write-Host "3. Database oneclickvirt exists" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "Press any key to exit..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
