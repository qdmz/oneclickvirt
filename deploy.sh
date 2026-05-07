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

# SSH 执行远程命令
ssh_cmd() {
    sshpass -p "$SSH_PASS" ssh -o StrictHostKeyChecking=no -p $SSH_PORT $SSH_USER@$SERVER_IP "$1"
}

# SSH 执行远程命令(多行)
ssh_exec() {
    sshpass -p "$SSH_PASS" ssh -o StrictHostKeyChecking=no -p $SSH_PORT $SSH_USER@$SERVER_IP "$@"
}

# SCP 上传文件
scp_push() {
    sshpass -p "$SSH_PASS" scp -o StrictHostKeyChecking=no -P $SSH_PORT "$1" $SSH_USER@$SERVER_IP:"$2"
}

# 1. 安装基础依赖
install_deps() {
    echo "[2/8] 安装基础依赖..."
    ssh_exec << 'ENDSSH'
        apt-get update
        apt-get install -y curl wget git build-essential mysql-server nginx
ENDSSH
}

# 2. 安装 Go
install_go() {
    echo "[3/8] 安装 Go..."
    ssh_exec << 'ENDSSH'
        if ! command -v go &> /dev/null; then
            wget -q https://go.dev/dl/go1.22.0.linux-amd64.tar.gz -O /tmp/go.tar.gz
            rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tar.gz
            echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
            export PATH=$PATH:/usr/local/go/bin
        fi
        export PATH=$PATH:/usr/local/go/bin
        go version
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

# 4. 配置 MySQL
setup_mysql() {
    echo "[5/8] 配置 MySQL..."
    ssh_exec << 'ENDSSH'
        systemctl enable mysql
        systemctl start mysql
        
        mysql -e "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'oneclickvirt123';"
        mysql -e "FLUSH PRIVILEGES;"
        
        mysql -u root -poneclickvirt123 -e "CREATE DATABASE IF NOT EXISTS oneclickvirt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
        
        echo "MySQL 配置完成"
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

# 7. 编译后端
build_backend() {
    echo "[8/8] 编译并启动服务..."
    ssh_exec << 'ENDSSH'
        export PATH=$PATH:/usr/local/go/bin
        
        # 编译后端
        cd /opt/oneclickvirt/server
        go mod download
        go build -o oneclickvirt .
        
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
        
        # 启动后端服务
        nohup ./oneclickvirt > /var/log/oneclickvirt.log 2>&1 &
        echo $! > /var/run/oneclickvirt.pid
        
        # 编译前端
        cd /opt/oneclickvirt/front
        npm install
        npm run build
        
        # 使用 serve 托管前端
        nohup npx serve -s dist -l 30005 > /var/log/oneclickvirt-front.log 2>&1 &
        echo $! > /var/run/oneclickvirt-front.pid
        
        echo "服务启动完成"
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
    build_backend
    
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
