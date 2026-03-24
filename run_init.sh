#!/bin/bash

export MYSQL_DATABASE=${MYSQL_DATABASE:-oneclickvirt}

echo "开始初始化数据库..."

# 检查 MySQL 是否运行
for i in {1..30}; do
    if mysql -h localhost -u root -e "SELECT 1" >/dev/null 2>&1; then
        echo "MySQL 已启动"
        break
    fi
    echo "等待 MySQL 启动... ($i/30)"
    if [ $i -eq 30 ]; then
        echo "MySQL 启动失败，无法执行初始化"
        exit 1
    fi
    sleep 1
done

# 创建数据库（如果不存在）
mysql -h localhost -u root -e "CREATE DATABASE IF NOT EXISTS \`${MYSQL_DATABASE}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
echo "数据库 ${MYSQL_DATABASE} 已准备就绪"

# 执行初始化脚本
echo "执行数据库初始化脚本..."
mysql -h localhost -u root ${MYSQL_DATABASE} < /root/oneclickvirt-new/complete_init.sql

if [ $? -eq 0 ]; then
    echo "数据库初始化成功！"
    echo "默认管理员账号: admin"
    echo "默认密码: password"
    echo "默认普通用户账号: user"
    echo "默认密码: password"
    echo "OAuth2 提供商已配置：GitHub, Google"
    echo "请在管理后台配置 OAuth2 提供商的客户端 ID 和密钥"
else
    echo "数据库初始化失败，请检查错误信息"
    exit 1
fi

echo "初始化完成！"
echo "系统已准备就绪，可以开始使用了。"
