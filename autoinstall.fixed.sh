#!/bin/bash

# OneClickVirt 自动安装脚本
# 启动服务并初始化配置

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 日志文件
LOG_FILE="/app/storage/logs/container.log"
mkdir -p "$(dirname "$LOG_FILE")"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" >> "$LOG_FILE"
}

# 等待MySQL就绪
wait_for_mysql() {
    local max_attempts=60
    local attempt=1

    log "等待MySQL服务启动..."

    while [ $attempt -le $max_attempts ]; do
        # 使用环境变量中的用户名和密码，如果没有配置则使用默认值
        local mysql_user="${DB_USER:-root}"
        local mysql_password="${DB_PASS:-}"
        local mysql_host="${DB_HOST:-oneclickvirt-mysql}"
        local mysql_port="${DB_PORT:-3306}"

        if mysql -h"$mysql_host" -P"$mysql_port" -u"$mysql_user" -p"$mysql_password" -e "SELECT 1" &>/dev/null; then
            log "MySQL服务已就绪"
            return 0
        fi
        log "等待MySQL启动... ($attempt/$max_attempts)"
        sleep 2
        attempt=$((attempt + 1))
    done

    log_error "MySQL服务启动超时"
    return 1
}

# 初始化数据库
init_db() {
    log "检查数据库是否存在..."
    local mysql_user="${DB_USER:-root}"
    local mysql_password="${DB_PASS:-}"
    local mysql_host="${DB_HOST:-oneclickvirt-mysql}"
    local mysql_port="${DB_PORT:-3306}"

    if mysql -h"$mysql_host" -P"$mysql_port" -u"$mysql_user" -p"$mysql_password" -e "USE oneclickvirt;" &>/dev/null; then
        log "数据库已存在"
        return 0
    fi

    log "创建数据库..."
    if mysql -h"$mysql_host" -P"$mysql_port" -u"$mysql_user" -p"$mysql_password" -e "CREATE DATABASE oneclickvirt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null; then
        log "数据库创建成功"
        return 0
    else
        log_error "数据库创建失败"
        return 1
    fi
}

# 启动应用程序
start_app() {
    log "启动OneClickVirt应用程序..."
    exec ./main
}

# 主函数
main() {
    log "开始启动OneClickVirt容器"
    log_info "OneClickVirt AutoInstall"
    log_info "============================================"

    # 初始化数据库
    init_db

    # 等待MySQL就绪
    if ! wait_for_mysql; then
        log_error "无法连接到MySQL"
        exit 1
    fi

    log "所有初始化完成，启动应用..."

    # 启动应用
    start_app
}

echo ""
log_info "OneClickVirt 自动安装脚本"
echo "============================================"
echo ""

main
