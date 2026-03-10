@echo off
chcp 65001 >nul
echo ========================================
echo    OneClickVirt 停止服务脚本
echo ========================================
echo.

echo 正在停止 OneClickVirt-Backend 进程...
taskkill /FI "WINDOWTITLE eq OneClickVirt-Backend*" /F 2>nul

echo 正在停止 OneClickVirt-Frontend 进程...
taskkill /FI "WINDOWTITLE eq OneClickVirt-Frontend*" /F 2>nul

echo.
echo ========================================
echo    服务已停止
echo ========================================
echo.
pause
