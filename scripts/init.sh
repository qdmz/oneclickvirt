#!/bin/bash

# OneClickVirt 数据库初始化脚本
# 用于手动初始化数据库或重置数据

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

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

# 检查是否在Docker容器内运行
check_docker() {
    if [ -f /.dockerenv ]; then
        log_info "检测到在Docker容器内运行"
        return 0
    else
        log_info "检测到在宿主机运行"
        return 1
    fi
}

# 检查MySQL是否就绪
wait_for_mysql() {
    local max_attempts=30
    local attempt=1
    
    log_info "等待MySQL服务启动..."
    
    while [ $attempt -le $max_attempts ]; do
        if mysql -uroot -e "SELECT 1" &>/dev/null; then
            log_info "MySQL服务已就绪"
            return 0
        fi
        log_info "等待MySQL启动... ($attempt/$max_attempts)"
        sleep 2
        attempt=$((attempt + 1))
    done
    
    log_error "MySQL服务启动超时"
    return 1
}

# 检查数据库是否已初始化
check_initialized() {
    local user_count=$(mysql -uroot oneclickvirt -N -e "SELECT COUNT(*) FROM users" 2>/dev/null || echo "0")
    if [ "$user_count" -gt 0 ]; then
        return 0
    fi
    return 1
}

# 执行初始化
do_init() {
    log_info "开始执行数据库初始化..."
    
    # 检查SQL文件是否存在
    if [ ! -f "$SCRIPT_DIR/init.sql" ]; then
        log_error "初始化SQL文件不存在: $SCRIPT_DIR/init.sql"
        exit 1
    fi
    
    # 执行SQL初始化
    mysql -uroot oneclickvirt < "$SCRIPT_DIR/init.sql"
    
    if [ $? -eq 0 ]; then
        log_info "数据库初始化成功"
        
        # 创建初始化标志文件
        mkdir -p /app/storage
        echo "System initialized at: $(date -Iseconds)" > /app/storage/.system_initialized
        log_info "已创建初始化标志文件"
        
        echo ""
        log_info "============================================"
        log_info "初始化完成！"
        log_info "默认管理员账户:"
        log_info "  用户名: admin"
        log_info "  密码: admin123456"
        log_info "============================================"
        echo ""
    else
        log_error "数据库初始化失败"
        exit 1
    fi
}

# 主函数
main() {
    log_info "OneClickVirt 数据库初始化工具"
    log_info "============================================"
    
    # 检查运行环境
    if check_docker; then
        # 在Docker容器内
        wait_for_mysql
        
        if check_initialized; then
            log_warn "数据库已初始化，跳过"
            log_info "如需重新初始化，请先清空数据库"
            exit 0
        fi
        
        do_init
    else
        # 在宿主机
        if [ -z "$1" ]; then
            echo "用法: $0 <容器名或容器ID>"
            echo "示例: $0 oneclickvirt"
            exit 1
        fi
        
        CONTAINER=$1
        
        # 检查容器是否运行
        if ! docker ps | grep -q "$CONTAINER"; then
            log_error "容器 $CONTAINER 未运行"
            exit 1
        fi
        
        log_info "在容器 $CONTAINER 中执行初始化..."
        
        # 复制初始化脚本到容器
        docker cp "$SCRIPT_DIR/init.sql" "$CONTAINER:/tmp/init.sql"
        docker cp "$SCRIPT_DIR/init.sh" "$CONTAINER:/tmp/init.sh"
        
        # 在容器内执行
        docker exec "$CONTAINER" bash /tmp/init.sh --docker-internal
        
        # 清理临时文件
        docker exec "$CONTAINER" rm -f /tmp/init.sql /tmp/init.sh
    fi
}

# 处理内部调用参数
if [ "$1" = "--docker-internal" ]; then
    wait_for_mysql
    if check_initialized; then
        log_warn "数据库已初始化，跳过"
        exit 0
    fi
    do_init
else
    main "$@"
fi
