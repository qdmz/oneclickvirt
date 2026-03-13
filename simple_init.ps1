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

Write-Host "Initializing OneClickVirt system..." -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Uri $baseURL -Method Post -Body $payload -ContentType "application/json" -TimeoutSec 45

    if ($response.code -eq 0 -or $response.code -eq 200) {
        Write-Host "SUCCESS!" -ForegroundColor Green
        Write-Host "Admin: http://localhost:8080/#/admin/login (admin/admin123)" -ForegroundColor White
        Write-Host "User: http://localhost:8080/#/login (user/user123)" -ForegroundColor White
    } else {
        Write-Host "FAILED: $($response.msg)" -ForegroundColor Red
        Write-Host "Use browser: http://localhost:8080/#/init" -ForegroundColor Yellow
    }
} catch {
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Use browser: http://localhost:8080/#/init" -ForegroundColor Yellow
}

Start-Sleep -Seconds 2
