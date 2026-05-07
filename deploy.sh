#!/bin/bash
# OneClickVirt 一键部署脚本（本地执行版）
# 使用方法: bash deploy.sh
#
# 部署配置
DB_PASS="oneclickvirt123"
BACKEND_PORT="30002"
FRONTEND_PORT="30005"
DOMAIN="tianchong.ypvps.com"

echo "=========================================="
echo "OneClickVirt 一键部署脚本"
echo "=========================================="

# 1. 创建必要目录
echo "[1/8] 创建必要目录..."
mkdir -p /opt/oneclickvirt/storage/logs
mkdir -p /opt/oneclickvirt/storage/uploads
mkdir -p /opt/oneclickvirt/storage/avatars
echo "目录创建完成"

# 2. 配置 MariaDB
echo "[2/8] 配置 MariaDB..."
systemctl enable mariadb
systemctl restart mariadb
sleep 3
mysql -e "SET PASSWORD FOR 'root'@'localhost' = PASSWORD('$DB_PASS'); FLUSH PRIVILEGES;"
mysql -u root -p"$DB_PASS" -e "CREATE DATABASE IF NOT EXISTS oneclickvirt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
echo "MariaDB 配置完成"

# 3. 拉取代码（如果目录不存在）
echo "[3/8] 拉取代码..."
if [ ! -d "/opt/oneclickvirt/.git" ]; then
    rm -rf /opt/oneclickvirt
    git clone https://github.com/qdmz/oneclickvirt.git /opt/oneclickvirt
fi
cd /opt/oneclickvirt
git pull origin main
echo "代码拉取完成"

# 4. 初始化数据库（使用 complete_init.sql - 支持重复运行）
echo "[4/8] 初始化数据库..."
cd /opt/oneclickvirt

# 创建数据库
mysql -u root -p"$DB_PASS" -e "CREATE DATABASE IF NOT EXISTS oneclickvirt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 导入完整的 init.sql
mysql -u root -p"$DB_PASS" oneclickvirt < /opt/oneclickvirt/scripts/init.sql 2>&1 || true

# 使用 complete_init.sql 导入默认数据（INSERT IGNORE 防止重复）
if [ -f "/opt/oneclickvirt/complete_init.sql" ]; then
    echo "导入默认数据..."
    mysql -u root -p"$DB_PASS" oneclickvirt < /opt/oneclickvirt/complete_init.sql 2>&1 || true
fi

# 运行 autosetup 修复脚本（如果存在）
if [ -f "/opt/oneclickvirt/scripts/autosetup.sql" ]; then
    echo "运行自动修复脚本..."
    mysql -u root -p"$DB_PASS" oneclickvirt < /opt/oneclickvirt/scripts/autosetup.sql 2>&1 || true
fi

echo "数据库初始化完成"

# 5. 编译后端
echo "[5/8] 编译后端..."
cd /opt/oneclickvirt/server

# 创建正确的配置文件格式 (匹配 config.Server 结构)
cat > config.yaml << CONFIG
system:
  addr: $BACKEND_PORT
  env: public
  db-type: mysql
  use-multipoint: false
  use-redis: false

jwt:
  expires-time: "7d"
  buffer-time: "1d"
  issuer: "oneclickvirt"

zap:
  level: "info"
  prefix: "[oneclickvirt]"
  format: "console"
  director: "logs"
  encode-level: "LowercaseColorLevelEncoder"
  stacktrace-key: "stacktrace"
  show-line: true
  log-in-console: true

mysql:
  path: "127.0.0.1"
  port: "3306"
  config: "charset=utf8mb4&parseTime=True&loc=Local"
  db-name: "oneclickvirt"
  username: "root"
  password: "$DB_PASS"
  prefix: ""
  singular: false
  engine: "InnoDB"
  max-idle-conns: 10
  max-open-conns: 100
  log-mode: "info"
  log-zap: false
  max-lifetime: 3600
  auto-create: true

auth:
  enable-email: false
  enable-telegram: false
  enable-qq: false
  enable-oauth2: false
  enable-public-registration: false

quota:
  default-level: 1
  instance-type-permissions:
    min-level-for-container: 1
    min-level-for-vm: 1
    min-level-for-delete-container: 1
    min-level-for-delete-vm: 1
    min-level-for-reset-container: 1
    min-level-for-reset-vm: 1
  level-limits:
    1:
      max-instances: 1
      max-resources:
        cpu: 1
        memory: 1025
        disk: 1
    2:
      max-instances: 3
      max-resources:
        cpu: 2
        memory: 1024
        disk: 20
    3:
      max-instances: 5
      max-resources:
        cpu: 4
        memory: 2048
        disk: 40
    4:
      max-instances: 10
      max-resources:
        cpu: 8
        memory: 4096
        disk: 80
    5:
      max-instances: 20
      max-resources:
        cpu: 16
        memory: 8192
        disk: 160

invite-code:
  enabled: false
  required: false

captcha:
  enabled: false
  width: 120
  height: 40
  length: 4
  expire-time: 5

cors:
  mode: "allow-all"
  whitelist:
    - "http://localhost:8080"
    - "http://127.0.0.1:8080"

redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0
CONFIG

# 编译
/usr/local/go/bin/go mod download
/usr/local/go/bin/go build -o oneclickvirt .

# 停止旧进程
pkill -f oneclickvirt 2>/dev/null || true

# 启动后端
nohup ./oneclickvirt > /var/log/oneclickvirt.log 2>&1 &
echo $! > /var/run/oneclickvirt.pid
sleep 3

# 检查后端端口
if netstat -tlnp 2>/dev/null | grep -q ":$BACKEND_PORT" || ss -tlnp 2>/dev/null | grep -q ":$BACKEND_PORT"; then
    echo "后端服务启动成功 (端口$BACKEND_PORT)"
else
    echo "后端服务启动失败，查看日志:"
    tail -30 /var/log/oneclickvirt.log
fi

# 6. 编译前端
echo "[6/8] 编译前端..."
cd /opt/oneclickvirt/web

# 安装pnpm并编译前端
npm install -g pnpm
pnpm install
pnpm build

# 停止旧前端进程
pkill -f "node.*proxy-server" 2>/dev/null || true
pkill -f "serve.*$FRONTEND_PORT" 2>/dev/null || true

# 使用代理服务器启动前端（支持API转发）
export FRONTEND_PORT=$FRONTEND_PORT
export BACKEND_PORT=$BACKEND_PORT
export BACKEND_HOST=127.0.0.1
nohup node proxy-server.js > /var/log/oneclickvirt-front.log 2>&1 &
echo $! > /var/run/oneclickvirt-front.pid
sleep 3

if netstat -tlnp 2>/dev/null | grep -q ":$FRONTEND_PORT" || ss -tlnp 2>/dev/null | grep -q ":$FRONTEND_PORT"; then
    echo "前端服务启动成功 (端口$FRONTEND_PORT)"
else
    echo "前端服务启动失败，查看日志:"
    tail -30 /var/log/oneclickvirt-front.log
fi

# 7. 启动 API 文档服务（可选，端口30003）
echo "[7/8] 启动Swagger文档服务..."
cd /opt/oneclickvirt/server
pkill -f "node.*swagger" 2>/dev/null || true
nohup npx swagger serve swagger/swagger.yaml -p 30003 --no-open > /var/log/oneclickvirt-swagger.log 2>&1 &
echo $! > /var/run/oneclickvirt-swagger.pid
sleep 2

# 8. 服务状态检查
echo "[8/8] 服务状态检查..."
echo ""
echo "=========================================="
echo "服务状态"
echo "=========================================="
ps aux | grep -E "oneclickvirt|proxy-server" | grep -v grep
netstat -tlnp 2>/dev/null | grep -E "$BACKEND_PORT|$FRONTEND_PORT|30003" || ss -tlnp | grep -E "$BACKEND_PORT|$FRONTEND_PORT|30003"

echo ""
echo "=========================================="
echo "部署完成!"
echo "=========================================="
echo "前端地址: http://$DOMAIN:$FRONTEND_PORT"
echo "后端地址: http://$DOMAIN:$BACKEND_PORT"
echo "API文档:  http://$DOMAIN:30003/docs"
echo ""
echo "管理员账号:"
echo "  用户名: admin"
echo "  密码: TestPass12!#"
echo "=========================================="
