# MySQL Database Setup Script for OneClickVirt

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "OneClickVirt MySQL Database Setup" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$mysqlPath = "C:\Program Files\MySQL\MySQL Server 9.6\bin\mysql.exe"

# Check if MySQL is installed
if (-not (Test-Path $mysqlPath)) {
    Write-Host "ERROR: MySQL not found at $mysqlPath" -ForegroundColor Red
    Write-Host "Please update the path in this script." -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host "MySQL found: $mysqlPath" -ForegroundColor Green
Write-Host ""

# Get root password
$rootPassword = Read-Host "Enter MySQL root password (press Enter if no password)" -AsSecureString
$rootPasswordPlain = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto([System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($rootPassword))

if ($rootPasswordPlain -eq "") {
    Write-Host "Using empty password for root..." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Creating database and user..." -ForegroundColor Yellow

# SQL commands
$sqlCommands = @"
CREATE DATABASE IF NOT EXISTS oneclickvirt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER IF NOT EXISTS 'oneclickvirt'@'%' IDENTIFIED BY '123456';
CREATE USER IF NOT EXISTS 'oneclickvirt'@'localhost' IDENTIFIED BY '123456';

GRANT ALL PRIVILEGES ON oneclickvirt.* TO 'oneclickvirt'@'%';
GRANT ALL PRIVILEGES ON oneclickvirt.* TO 'oneclickvirt'@'localhost';

FLUSH PRIVILEGES;
"@

try {
    if ($rootPasswordPlain -eq "") {
        $result = & $mysqlPath -u root -e $sqlCommands 2>&1
    } else {
        $result = & $mysqlPath -u root -p$rootPasswordPlain -e $sqlCommands 2>&1
    }

    if ($LASTEXITCODE -eq 0) {
        Write-Host "SUCCESS: Database and user created successfully!" -ForegroundColor Green
        Write-Host ""

        # Verify
        Write-Host "Verifying setup..." -ForegroundColor Yellow
        $verifySQL = "SELECT User, Host FROM mysql.user WHERE User = 'oneclickvirt';"
        if ($rootPasswordPlain -eq "") {
            & $mysqlPath -u root -e $verifySQL
        } else {
            & $mysqlPath -u root -p$rootPasswordPlain -e $verifySQL
        }

        Write-Host ""
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "Setup Complete!" -ForegroundColor Green
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "Database: oneclickvirt" -ForegroundColor White
        Write-Host "Username: oneclickvirt" -ForegroundColor White
        Write-Host "Password: 123456" -ForegroundColor White
        Write-Host ""

        Write-Host "Next steps:" -ForegroundColor Yellow
        Write-Host "1. Refresh the browser page" -ForegroundColor White
        Write-Host "2. Test connection with these settings:" -ForegroundColor White
        Write-Host "   Type: mysql" -ForegroundColor Gray
        Write-Host "   Host: localhost" -ForegroundColor Gray
        Write-Host "   Port: 3306" -ForegroundColor Gray
        Write-Host "   Database: oneclickvirt" -ForegroundColor Gray
        Write-Host "   Username: oneclickvirt" -ForegroundColor Gray
        Write-Host "   Password: 123456" -ForegroundColor Gray
        Write-Host ""
        Write-Host "3. Click 'Test Connection' button" -ForegroundColor White
    } else {
        Write-Host "ERROR: Failed to create database/user" -ForegroundColor Red
        Write-Host "Output: $result" -ForegroundColor Red
        Write-Host ""
        Write-Host "Possible reasons:" -ForegroundColor Yellow
        Write-Host "1. Incorrect root password" -ForegroundColor Yellow
        Write-Host "2. MySQL service not running" -ForegroundColor Yellow
        Write-Host "3. Root user lacks permissions" -ForegroundColor Yellow
    }
} catch {
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Read-Host "Press Enter to exit"
