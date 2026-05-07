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

# 检查依赖
check_deps() {
    echo "[1/8] 检查依赖..."
    if ! command -v sshpass &> /dev/null; then
        echo "正在安装 sshpass..."
        apt-get update && apt-get install -y sshpass
    fi
}

# SSH 执行远程命令(多行)
ssh_exec() {
    sshpass -p "$SSH_PASS" ssh -o StrictHostKeyChecking=no -p $SSH_PORT $SSH_USER@$SERVER_IP "$@"
}

# 1. 安装基础依赖
install_deps() {
    echo "[2/8] 安装基础依赖..."
    ssh_exec << 'ENDSSH'
        export DEBIAN_FRONTEND=noninteractive
        apt-get update
        apt-get install -y curl wget git build-essential mariadb-server nginx
ENDSSH
}

# 2. 安装 Go
install_go() {
    echo "[3/8] 安装 Go..."
    ssh_exec << 'ENDSSH'
        export PATH=$PATH:/usr/local/go/bin
        if ! command -v go &> /dev/null; then
            wget -q https://go.dev/dl/go1.22.0.linux-amd64.tar.gz -O /tmp/go.tar.gz
            rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tar.gz
            echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
            export PATH=$PATH:/usr/local/go/bin
        fi
        /usr/local/go/bin/go version
ENDSSH
}

# 3. 安装 Node.js
install_node() {
    echo "[4/8] 安装 Node.js..."
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
    echo "[5/8] 配置 MariaDB..."
    ssh_exec << 'ENDSSH'
        export DEBIAN_FRONTEND=noninteractive
        systemctl enable mariadb
        systemctl start mariadb
        
        # 设置root密码
        mysql -e "ALTER USER 'root'@'localhost' IDENTIFIED BY 'oneclickvirt123';"
        mysql -e "FLUSH PRIVILEGES;"
        
        # 创建数据库
        mysql -u root -poneclickvirt123 -e "CREATE DATABASE IF NOT EXISTS oneclickvirt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
        
        echo "MariaDB 配置完成"
ENDSSH
}

# 5. 拉取代码
clone_code() {
    echo "[6/8] 拉取代码..."
    ssh_exec << ENDSSH
        mkdir -p $INSTALL_DIR
        cd $INSTALL_DIR
        if [ -d ".git" ]; then
            git pull origin main
        else
            git clone $REPO_URL $INSTALL_DIR
        fi
        echo "代码拉取完成"
ENDSSH
}

# 6. 初始化数据库
init_database() {
    echo "[7/8] 初始化数据库..."
    ssh_exec << ENDSSH
        cd $INSTALL_DIR
        mysql -u root -poneclickvirt123 oneclickvirt < scripts/init.sql
        echo "数据库初始化完成"
ENDSSH
}

# 7. 编译并启动服务
build_and_start() {
    echo "[8/8] 编译并启动服务..."
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
        [ -f /var/run/oneclickvirt-front.pid ] && kill $(cat /var/run/oneclickvirt-front.pid) 2>/dev/null || true
        
        # 启动后端服务
        nohup ./oneclickvirt > /var/log/oneclickvirt.log 2>&1 &
        echo $! > /var/run/oneclickvirt.pid
        sleep 2
        
        # 安装pnpm并编译前端
        cd /opt/oneclickvirt/front
        npm install -g pnpm
        pnpm install
        pnpm build
        
        # 使用npx serve托管前端
        nohup npx serve -s dist -l 30005 > /var/log/oneclickvirt-front.log 2>&1 &
        echo $! > /var/run/oneclickvirt-front.pid
        sleep 2
        
        echo ""
        echo "=========================================="
        echo "服务状态检查"
        echo "=========================================="
        ps aux | grep oneclickvirt | grep -v grep
        netstat -tlnp | grep -E '30002|30005' || ss -tlnp | grep -E '30002|30005'
        
        echo ""
        echo "服务启动完成!"
        echo "后端端口: 30002"
        echo "前端端口: 30005"
        echo "访问地址: http://tianchong.ypvps.com:30005"
ENDSSH
}

# 主流程
main() {
    check_deps
    install_deps
    install_go
    install_node
    setup_mysql
    clone_code
    init_database
    build_and_start
    
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
