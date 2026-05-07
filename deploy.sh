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
echo "[1/7] 创建必要目录..."
mkdir -p /opt/oneclickvirt/storage/logs
mkdir -p /opt/oneclickvirt/storage/uploads
mkdir -p /opt/oneclickvirt/storage/avatars
echo "目录创建完成"

# 2. 配置 MariaDB
echo "[2/7] 配置 MariaDB..."
systemctl enable mariadb
systemctl restart mariadb
sleep 3
mysql -e "SET PASSWORD FOR 'root'@'localhost' = PASSWORD('$DB_PASS'); FLUSH PRIVILEGES;"
mysql -u root -p"$DB_PASS" -e "CREATE DATABASE IF NOT EXISTS oneclickvirt CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
echo "MariaDB 配置完成"

# 3. 拉取代码（如果目录不存在）
echo "[3/7] 拉取代码..."
if [ ! -d "/opt/oneclickvirt/.git" ]; then
    rm -rf /opt/oneclickvirt
    git clone https://github.com/qdmz/oneclickvirt.git /opt/oneclickvirt
fi
cd /opt/oneclickvirt
git pull origin main
echo "代码拉取完成"

# 4. 初始化数据库
echo "[4/7] 初始化数据库..."
cd /opt/oneclickvirt

# 先创建users表（依赖其他表）
mysql -u root -p"$DB_PASS" oneclickvirt << 'EOF'
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
  KEY status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO users (uuid, username, password, nickname, email, phone, status, level, level_expire_at, user_type, created_at, updated_at) VALUES
('admin-001', 'admin', 'TestPass12!#', '管理员', 'admin@example.com', '13800138000', 1, 5, '2099-12-31 23:59:59', 'admin', NOW(), NOW());

CREATE TABLE IF NOT EXISTS roles (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL,
  description TEXT,
  status INT DEFAULT 1,
  created_at DATETIME(3) DEFAULT NULL,
  updated_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO roles (name, code, description, status, created_at, updated_at) VALUES
('admin', 'admin', '系统管理员角色', 1, NOW(), NOW()),
('user', 'user', '普通用户角色', 1, NOW(), NOW());

CREATE TABLE IF NOT EXISTS user_roles (
  id INT NOT NULL AUTO_INCREMENT,
  user_id INT NOT NULL,
  role_id INT NOT NULL,
  created_at DATETIME(3) DEFAULT NULL,
  updated_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY user_id (user_id),
  KEY role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO user_roles (user_id, role_id, created_at, updated_at) VALUES (1, 1, NOW(), NOW());

CREATE TABLE IF NOT EXISTS products (
  id INT NOT NULL AUTO_INCREMENT,
  uuid VARCHAR(36) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  price DECIMAL(10,2) DEFAULT 0.00,
  cpu_limit INT DEFAULT 1,
  memory_limit INT DEFAULT 512,
  disk_limit INT DEFAULT 10240,
  bandwidth_limit INT DEFAULT 100,
  instance_limit INT DEFAULT 1,
  features TEXT,
  status INT DEFAULT 1,
  sort_order INT DEFAULT 0,
  is_enabled INT DEFAULT 1,
  cpu INT DEFAULT 1,
  memory INT DEFAULT 512,
  disk INT DEFAULT 10240,
  bandwidth INT DEFAULT 100,
  traffic INT DEFAULT 0,
  period INT DEFAULT 30,
  allow_repeat INT DEFAULT 1,
  created_at DATETIME(3) DEFAULT NULL,
  updated_at DATETIME(3) DEFAULT NULL,
  deleted_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO products (uuid, name, description, price, cpu_limit, memory_limit, disk_limit, bandwidth_limit, instance_limit, status, sort_order, is_enabled, cpu, memory, disk, bandwidth, period, allow_repeat, created_at, updated_at) VALUES
(UUID(), '入门套餐', '适合个人用户', 0, 1, 512, 10240, 100, 1, 1, 1, 1, 1, 512, 10240, 100, 30, 1, NOW(), NOW()),
(UUID(), '标准套餐', '适合小型团队', 990, 2, 1024, 20480, 200, 3, 1, 2, 1, 2, 1024, 20480, 200, 30, 1, NOW(), NOW());

CREATE TABLE IF NOT EXISTS system_configs (
  id INT NOT NULL AUTO_INCREMENT,
  \`key\` VARCHAR(128) NOT NULL,
  value TEXT,
  description TEXT,
  created_at DATETIME(3) DEFAULT NULL,
  updated_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY \`key\` (\`key\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO system_configs (\`key\`, value, description, created_at, updated_at) VALUES
('site_name', 'OneClickVirt', '网站名称', NOW(), NOW()),
('enable_registration', 'true', '是否开启注册', NOW(), NOW()),
('default_user_level', '1', '默认用户等级', NOW(), NOW());

CREATE TABLE IF NOT EXISTS site_configs (
  id INT NOT NULL AUTO_INCREMENT,
  \`key\` VARCHAR(128) NOT NULL,
  value TEXT,
  type VARCHAR(32) DEFAULT 'string',
  \`group\` VARCHAR(32) DEFAULT 'basic',
  description TEXT,
  created_at DATETIME(3) DEFAULT NULL,
  updated_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY \`key\` (\`key\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO site_configs (\`key\`, value, type, \`group\`, description, created_at, updated_at) VALUES
('site_name', 'OneClickVirt', 'string', 'basic', '网站名称', NOW(), NOW());

CREATE TABLE IF NOT EXISTS announcements (
  id INT NOT NULL AUTO_INCREMENT,
  uuid VARCHAR(36) NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT,
  content_html TEXT,
  type VARCHAR(32) DEFAULT NULL,
  status INT DEFAULT 1,
  priority INT DEFAULT 0,
  is_sticky INT DEFAULT 0,
  start_at DATETIME(3) DEFAULT NULL,
  end_at DATETIME(3) DEFAULT NULL,
  sort_order INT DEFAULT 0,
  created_at DATETIME(3) DEFAULT NULL,
  updated_at DATETIME(3) DEFAULT NULL,
  deleted_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uuid (uuid),
  KEY type (type),
  KEY status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO announcements (uuid, title, content, content_html, type, status, priority, is_sticky, created_at, updated_at) VALUES
(UUID(), '欢迎使用', '欢迎使用虚拟化管理平台', '<p>欢迎使用</p>', 'homepage', 1, 10, 1, NOW(), NOW());
EOF
echo "数据库初始化完成"

# 5. 编译后端
echo "[5/7] 编译后端..."
cd /opt/oneclickvirt/server

# 创建配置文件
cat > config.yaml << CONFIG
app:
  port: $BACKEND_PORT
  mode: debug

database:
  host: localhost
  port: 3306
  user: root
  password: $DB_PASS
  name: oneclickvirt
  charset: utf8mb4

server:
  domain: $DOMAIN
  backend_port: $BACKEND_PORT
  frontend_port: $FRONTEND_PORT
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

if netstat -tlnp 2>/dev/null | grep -q $BACKEND_PORT || ss -tlnp 2>/dev/null | grep -q $BACKEND_PORT; then
    echo "后端服务启动成功 (端口$BACKEND_PORT)"
else
    echo "后端服务启动失败，查看日志:"
    tail -20 /var/log/oneclickvirt.log
fi

# 6. 编译前端
echo "[6/7] 编译前端..."
cd /opt/oneclickvirt/web

# 安装pnpm并编译前端
npm install -g pnpm
pnpm install
pnpm build

# 停止旧前端进程
pkill -f "serve.*$FRONTEND_PORT" 2>/dev/null || true

# 使用npx serve托管前端
nohup npx serve -s dist -l $FRONTEND_PORT > /var/log/oneclickvirt-front.log 2>&1 &
echo $! > /var/run/oneclickvirt-front.pid
sleep 3

if netstat -tlnp 2>/dev/null | grep -q $FRONTEND_PORT || ss -tlnp 2>/dev/null | grep -q $FRONTEND_PORT; then
    echo "前端服务启动成功 (端口$FRONTEND_PORT)"
else
    echo "前端服务启动失败，查看日志:"
    tail -20 /var/log/oneclickvirt-front.log
fi

# 7. 服务状态检查
echo "[7/7] 服务状态检查..."
echo ""
echo "=========================================="
echo "服务状态"
echo "=========================================="
ps aux | grep -E "oneclickvirt|serve.*$FRONTEND_PORT" | grep -v grep
netstat -tlnp 2>/dev/null | grep -E "$BACKEND_PORT|$FRONTEND_PORT" || ss -tlnp | grep -E "$BACKEND_PORT|$FRONTEND_PORT"

echo ""
echo "=========================================="
echo "部署完成!"
echo "=========================================="
echo "前端地址: http://$DOMAIN:$FRONTEND_PORT"
echo "后端地址: http://$DOMAIN:$BACKEND_PORT"
echo ""
echo "管理员账号:"
echo "  用户名: admin"
echo "  密码: TestPass12!#"
echo "=========================================="
