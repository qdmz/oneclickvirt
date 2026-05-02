#!/bin/bash

# OneClickVirt 容器入口脚本
# 根据 CONTAINER_NAME 环境变量决定启动方式

set -e

_db_wait() {
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "SELECT 1" 2>/dev/null; then
            return 0
        fi
        sleep 2
        attempt=$((attempt + 1))
    done
    return 1
}

# 当前容器名称
CONTAINER_NAME="${CONTAINER_NAME:-${HOSTNAME}}"

if [ "$CONTAINER_NAME" = "oneclickvirt-init" ]; then
    # Init 服务：初始化数据库
    echo "[INFO] OneClickVirt 初始化服务"
    
    # 等待 MySQL 就绪
    if ! _db_wait; then
        echo "[ERROR] 无法连接到 MySQL"
        exit 1
    fi
    
    # 初始化数据库
    if ! mysql -h"$DB_HOST" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < /app/complete_init.sql 2>/dev/null; then
        echo "[ERROR] 数据库初始化失败"
        exit 1
    fi
    
    echo "[INFO] 数据库初始化完成"
    exit 0
    
elif [ "$CONTAINER_NAME" = "oneclickvirt-api" ]; then
    # API 服务：启动 Backend
    echo "[INFO] OneClickVirt API 服务启动"

    # 等待 MySQL 就绪
    if ! _db_wait; then
        echo "[ERROR] 无法连接到 MySQL"
        exit 1
    fi
    
    # 启动 API
    exec ./main
    
elif [ "$CONTAINER_NAME" = "oneclickvirt-web" ]; then
    # Web 服务：启动 Nginx（静态文件 + 代理到 API）
    echo "[INFO] OneClickVirt Web 服务启动"
    
    # 启动 Nginx
    exec nginx -g "daemon off;"
    
else
    echo "[ERROR] 未知容器: $CONTAINER_NAME"
    echo "当前容器名: $CONTAINER_NAME"
    exit 1
fi
