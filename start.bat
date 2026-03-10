@echo off
chcp 65001 >nul
echo ========================================
echo    OneClickVirt 开发环境启动脚本
echo ========================================
echo.

echo [1/3] 检查环境变量...
set "PATH=%PATH%;C:\Program Files\Git\cmd;C:\Program Files\Go\bin;C:\Program Files\nodejs"

echo.
echo [2/3] 启动后端服务 (Go)...
echo 后端将在 http://localhost:8890 运行
start "OneClickVirt-Backend" cmd /k "cd /d %~dp0server && go run main.go"

echo.
echo [3/3] 启动前端服务 (Node.js)...
echo 前端将在 http: //localhost:5173 运行
start "OneClickVirt-Frontend" cmd /k "cd /d %~dp0web && npm run dev"

echo.
echo ========================================
echo    启动完成！
echo ========================================
echo.
echo 前端地址: http: //localhost:5173
echo 后端地址: http: //localhost:8890
echo.
echo 按任意键关闭此窗口（服务将继续运行）...
pause >nul
