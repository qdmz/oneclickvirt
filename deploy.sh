#!/bin/bash
# OneClickVirt 一键部署脚本
# 使用方法: bash deploy.sh

set -e

# 配置
SERVER_IP="38.76.178.60"
SSH_PORT="30001"
SSH_USER="root"
SSH_PASS="9124a590f"
DOMAIN="tianchong.ypvps.com"
BACKEND_PORT="30002"
FRONTEND_PORT="30005"
DB_PASS="oneclickvirt123"
REPO_URL="https://github.com/qdmz/oneclickvirt.git"
INSTALL_DIR="/opt/oneclickvirt"

echo "=========================================="
echo "OneClickVirt 一键部署脚本"
echo "=========================================="

# SSH 执行远程命令(多行)
ssh_exec() {
    sshpass -p "$SSH_PASS" ssh -o StrictHostKeyChecking=no -p $SSH_PORT $SSH_USER@$SERVER_IP "$@"
}

# 1. 安装基础依赖
install_deps() {
    echo "[1/8] 安装基础依赖..."
    ssh_exec << 'ENDSSH'
        export DEBIAN_FRONTEND=noninteractive
        apt-get update
        apt-get install -y curl wget git build-essential mariadb-server nginx sshpass
ENDSSH
}

# 2. 安装 Go
install_go() {
    echo "[2/8] 安装 Go..."
    ssh_exec << 'ENDSSH'
        if ! command -v /usr/local/go/bin/go &> /dev/null; then
            wget -q https://go.dev/dl/go1.22.0.linux-amd64.tar.gz -O /tmp/go.tar.gz
            rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tar.gz
            echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
        fi
        /usr/local/go/bin/go version
ENDSSH
}

# 3. 安装 Node.js
install_node() {
    echo "[3/8] 安装 Node.js..."
    ssh_exec << 'ENDSSH'
        if ! command -v node &> /dev/null; then
            curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
            apt-get install -y nodejs
        fi
        node --version
        npm --version
ENDSSH
}

# 4. 配置 MariaDB
setup_mysql() {
    echo "[4/8] 配置 MariaDB..."
    ssh_exec << 'ENDSSH'
        export DEBIAN_FRONTEND=noninteractive
        systemctl enable mariadb
        systemctl restart mariadb
        
        # 等待MariaDB启动
        sleep 3
        
        # 安全初始化MariaDB，设置root密码
        mysql -e "SET PASSWORD FOR 'root'@'localhost' = PASSWORD('oneclickvirt123'); FLUSH PRIVILEGES;"
        
        # 创建数据库
        mysql -u root -poneclickvirt123 -e "CREATE DATABASE IF NOT EXISTS oneclickvirt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
        
        echo "MariaDB 配置完成"
ENDSSH
}

# 5. 拉取代码
clone_code() {
    echo "[5/8] 拉取代码..."
    ssh_exec << ENDSSH
        mkdir -p $INSTALL_DIR
        cd $INSTALL_DIR
        if [ -d ".git" ]; then
            git pull origin main
        else
            rm -rf $INSTALL_DIR
            git clone $REPO_URL $INSTALL_DIR
        fi
        echo "代码拉取完成"
ENDSSH
}

# 6. 初始化数据库
init_database() {
    echo "[6/8] 初始化数据库..."
    ssh_exec << 'ENDSSH'
        cd /opt/oneclickvirt
        
        # 先创建基础表结构（users表需要先存在）
        mysql -u root -poneclickvirt123 -e "
        CREATE TABLE IF NOT EXISTS users (
          id INT NOT NULL AUTO_INCREMENT,
          uuid VARCHAR(36) NOT NULL,
          username VARCHAR(64) NOT NULL UNIQUE,
          password VARCHAR(255) NOT NULL,
          nickname VARCHAR(64) DEFAULT NULL,
          email VARCHAR(128) DEFAULT NULL,
          phone VARCHAR(32) DEFAULT NULL,
          avatar VARCHAR(255) DEFAULT NULL,
          status INT DEFAULT 1,
          level INT DEFAULT 1,
          level_expire_at DATETIME DEFAULT NULL,
          user_type VARCHAR(20) DEFAULT 'user',
          balance DECIMAL(10,2) DEFAULT 0.00,
          total_spent DECIMAL(10,2) DEFAULT 0.00,
          total_orders INT DEFAULT 0,
          last_login_at DATETIME DEFAULT NULL,
          last_login_ip VARCHAR(64) DEFAULT NULL,
          created_at DATETIME(3) DEFAULT NULL,
          updated_at DATETIME(3) DEFAULT NULL,
          deleted_at DATETIME(3) DEFAULT NULL,
          PRIMARY KEY (id),
          UNIQUE KEY uuid (uuid),
          UNIQUE KEY username (username),
          KEY status (status),
          KEY user_type (user_type)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
        "
        
        # 执行完整初始化脚本
        mysql -u root -poneclickvirt123 oneclickvirt < scripts/init.sql
        
        echo "数据库初始化完成"
ENDSSH
}

# 7. 编译并启动服务
build_and_start() {
    echo "[7/8] 编译后端..."
    ssh_exec << 'ENDSSH'
        export PATH=$PATH:/usr/local/go/bin
        
        # 编译后端
        cd /opt/oneclickvirt/server
        /usr/local/go/bin/go mod download
        /usr/local/go/bin/go build -o oneclickvirt .
        
        # 创建配置文件
        cat > config.yaml << 'CONFIG'
app:
  port: 30002
  mode: debug

database:
  host: localhost
  port: 3306
  user: root
  password: oneclickvirt123
  name: oneclickvirt
  charset: utf8mb4

server:
  domain: tianchong.ypvps.com
  backend_port: 30002
  frontend_port: 30005
CONFIG
        
        # 停止旧进程
        [ -f /var/run/oneclickvirt.pid ] && kill $(cat /var/run/oneclickvirt.pid) 2>/dev/null || true
        pkill -f oneclickvirt 2>/dev/null || true
        
        # 启动后端服务
        nohup ./oneclickvirt > /var/log/oneclickvirt.log 2>&1 &
        echo $! > /var/run/oneclickvirt.pid
        sleep 3
        
        # 检查后端是否启动成功
        if ps aux | grep -v grep | grep oneclickvirt > /dev/null; then
            echo "后端服务启动成功"
        else
            echo "后端服务启动失败，查看日志:"
            cat /var/log/oneclickvirt.log
        fi
ENDSSH
}

# 8. 编译前端
build_frontend() {
    echo "[8/8] 编译前端..."
    ssh_exec << 'ENDSSH'
        # 安装pnpm并编译前端
        cd /opt/oneclickvirt/front
        npm install -g pnpm
        pnpm install
        pnpm build
        
        # 停止旧前端进程
        pkill -f "serve.*dist" 2>/dev/null || true
        
        # 使用npx serve托管前端
        nohup npx serve -s dist -l 30005 > /var/log/oneclickvirt-front.log 2>&1 &
        echo $! > /var/run/oneclickvirt-front.pid
        sleep 3
        
        # 检查前端是否启动成功
        if ps aux | grep -v grep | grep "serve.*30005" > /dev/null; then
            echo "前端服务启动成功"
        else
            echo "前端服务启动失败，查看日志:"
            cat /var/log/oneclickvirt-front.log
        fi
        
        echo ""
        echo "=========================================="
        echo "服务状态检查"
        echo "=========================================="
        ps aux | grep -E "oneclickvirt|serve.*30005" | grep -v grep
        netstat -tlnp 2>/dev/null | grep -E "30002|30005" || ss -tlnp | grep -E "30002|30005"
ENDSSH
}

# 主流程
main() {
    install_deps
    install_go
    install_node
    setup_mysql
    clone_code
    init_database
    build_and_start
    build_frontend
    
    echo ""
    echo "=========================================="
    echo "部署完成!"
    echo "=========================================="
    echo "前端地址: http://tianchong.ypvps.com:30005"
    echo "后端地址: http://tianchong.ypvps.com:30002"
    echo ""
    echo "管理员账号:"
    echo "  用户名: admin"
    echo "  密码: TestPass12!#"
    echo "=========================================="
}

main
