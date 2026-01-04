@echo off
chcp 65001 >nul
echo ========================================
echo   安装前端依赖
echo ========================================
echo.

echo 当前目录: %CD%
echo.

echo [1/2] 配置npm镜像...
call "C:\Program Files\nodejs\npm.cmd" config set registry https://registry.npmmirror.com
echo [OK] 镜像配置完成
echo.

if exist "node_modules\" (
    echo [提示] 依赖已安装,跳过安装步骤
    goto :end
)

echo [2/2] 安装依赖...
echo 正在下载依赖包,这可能需要 5-10 分钟,请耐心等待...
echo.

call "C:\Program Files\nodejs\npm.cmd" install

if %errorlevel% neq 0 (
    echo.
    echo ========================================
    echo [错误] 依赖安装失败!
    echo ========================================
    echo.
    pause
    exit /b 1
)

echo.
echo ========================================
echo [OK] 依赖安装成功!
echo ========================================
echo.

:end
echo 现在可以启动前端了:
echo npm run dev
echo.
pause
