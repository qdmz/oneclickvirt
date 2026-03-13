<# Disable Captcha and Restart Backend #>

# MySQL commands as string with escaped backticks
$sql = @"
USE oneclickvirt;
SELECT id, \`key\`, value FROM system_configs WHERE \`key\`='enable_captcha';
UPDATE system_configs SET value='false' WHERE \`key\`='enable_captcha';
SELECT id, \`key\`, value FROM system_configs WHERE \`key\`='enable_captcha';
"@

Write-Host "Disabling captcha in database..." -ForegroundColor Cyan

# Save SQL to temp file
$tempFile = "$env:TEMP\disable_captcha_temp.sql"
$sql | Out-File -FilePath $tempFile -Encoding utf8

# Execute with MySQL
$mysqlPath = "C:\Program Files\MySQL\MySQL Server 9.6\bin\mysql.exe"
& $mysqlPath -u oneclickvirt -p123456 -D oneclickvirt < $tempFile 2>&1

Write-Host ""
Write-Host "Restarting backend..." -ForegroundColor Yellow

# Stop backend
Stop-Process -Name oneclickvirt -Force -ErrorAction SilentlyContinue
Start-Sleep -Seconds 3

# Start backend
$serverDir = "C:\Users\admin\.openclaw-autoclaw\workspace\oneclickvirt\server"
Set-Location $serverDir
Start-Process -FilePath ".\oneclickvirt.exe" -RedirectStandardOutput "server.log" -RedirectStandardError "server_error.log" -WindowStyle Normal

Write-Host "Waiting for backend to start..." -ForegroundColor Cyan
Start-Sleep -Seconds 5

# Check status
Write-Host ""
Write-Host "Testing backend..." -ForegroundColor Cyan
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8890/api/v1/public/init/check" -UseBasicParsing -TimeoutSec 5
    Write-Host "Backend is running!" -ForegroundColor Green
    Write-Host "Response: $($response.Content)" -ForegroundColor White
} catch {
    Write-Host "Backend test failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Done! Now try to login:" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Admin: http://localhost:8080/#/admin/login" -ForegroundColor White
Write-Host "User:  http://localhost:8080/#/login" -ForegroundColor White
Write-Host ""
Write-Host "Credentials:" -ForegroundColor Yellow
Write-Host "  Username: admin  Password: admin123" -ForegroundColor White
Write-Host "  Username: user   Password: user123" -ForegroundColor White
Write-Host "========================================" -ForegroundColor Cyan

# Cleanup
Remove-Item $tempFile -ErrorAction SilentlyContinue
