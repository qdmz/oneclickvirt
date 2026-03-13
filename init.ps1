# System Initialization Script

$baseURL = "http://localhost:8890/api/v1/public/init"

Write-Host "Starting system initialization (using root account)..." -ForegroundColor Green

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

Write-Host "Request data: $payload" -ForegroundColor Yellow

try {
    $response = Invoke-RestMethod -Uri $baseURL -Method Post -Body $payload -ContentType "application/json" -TimeoutSec 30

    if ($response.code -eq 0 -or $response.code -eq 200) {
        Write-Host "SUCCESS: System initialized!" -ForegroundColor Green
        Write-Host "Message: $($response.data)" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "You can now login!" -ForegroundColor Cyan
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "Admin Login:" -ForegroundColor White
        Write-Host "   URL: http://localhost:8080/#/admin/login" -ForegroundColor Gray
        Write-Host "   Username: admin" -ForegroundColor Gray
        Write-Host "   Password: admin123" -ForegroundColor Gray
        Write-Host ""
        Write-Host "User Login:" -ForegroundColor White
        Write-Host "   URL: http://localhost:8080/#/login" -ForegroundColor Gray
        Write-Host "   Username: user" -ForegroundColor Gray
        Write-Host "   Password: user123" -ForegroundColor Gray
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Cyan
    } else {
        Write-Host "FAILED: $($response.msg)" -ForegroundColor Red
    }
} catch {
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
    Write-Host "If database connection fails, please check:" -ForegroundColor Yellow
    Write-Host "1. MySQL service is running" -ForegroundColor Yellow
    Write-Host "2. Try manual initialization: http://localhost:8080/#/init" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Press any key to exit..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
