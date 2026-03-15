#!/bin/bash
set -e
echo "Starting OneClickVirt..."

export MYSQL_DATABASE=${MYSQL_DATABASE:-oneclickvirt}

# Configure SSL certificates if provided
if [ ! -z "$SSL_CERT_PATH" ] && [ ! -z "$SSL_KEY_PATH" ]; then
    echo "Configuring SSL certificates..."
    if [ -f "$SSL_CERT_PATH" ] && [ -f "$SSL_KEY_PATH" ]; then
        cp "$SSL_CERT_PATH" /etc/nginx/ssl/cert.pem
        cp "$SSL_KEY_PATH" /etc/nginx/ssl/key.pem
        chmod 644 /etc/nginx/ssl/cert.pem
        chmod 600 /etc/nginx/ssl/key.pem
        echo "SSL certificates configured successfully"
    else
        echo "Warning: SSL certificate files not found, using self-signed certificates"
        openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/nginx/ssl/key.pem -out /etc/nginx/ssl/cert.pem -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
    fi
fi

# Use config from persistence directory if exists
if [ -f /app/config/config.yaml ]; then
    echo "Using persistent config from /app/config/config.yaml"
    cp /app/config/config.yaml /app/config.yaml
fi

# Update config.yaml with FRONTEND_URL if provided
if [ ! -z "$FRONTEND_URL" ]; then
    echo "Configuring frontend-url: $FRONTEND_URL"
    sed -i "s|frontend-url:.*|frontend-url: \"$FRONTEND_URL\"|g" /app/config.yaml
    
    # Extract domain from FRONTEND_URL
    DOMAIN=$(echo "$FRONTEND_URL" | sed -e "s|^[^/]*//||" -e "s|/.*$||")
    echo "Extracted domain: $DOMAIN"
    
    # Update nginx server_name for both HTTP and HTTPS servers
    sed -i "s/server_name localhost;/server_name $DOMAIN;/g" /etc/nginx/nginx.conf
    
    # Detect if URL is HTTPS and update nginx config accordingly
    if echo "$FRONTEND_URL" | grep -q "^https://"; then
        echo "Detected HTTPS frontend, SSL will be enabled"
        # Use Certbot to obtain free SSL certificate if not provided
        if [ ! -f "/etc/nginx/ssl/cert.pem" ] || [ ! -f "/etc/nginx/ssl/key.pem" ] || [ ! -s "/etc/nginx/ssl/cert.pem" ] || [ ! -s "/etc/nginx/ssl/key.pem" ]; then
            echo "Obtaining free SSL certificate for $DOMAIN using Certbot..."
            # Create self-signed certificate for initial setup
            echo "Creating self-signed SSL certificate for $DOMAIN"
            openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/nginx/ssl/key.pem -out /etc/nginx/ssl/cert.pem -subj "/C=US/ST=State/L=City/O=Organization/CN=$DOMAIN"
            chmod 644 /etc/nginx/ssl/cert.pem
            chmod 600 /etc/nginx/ssl/key.pem
            echo "Self-signed SSL certificate created successfully"
            
            # Add HTTPS server configuration to nginx.conf
            # First, remove the closing } of the http block
            sed -i '$d' /etc/nginx/nginx.conf
            
            # Now add the HTTPS server configuration
            cat >> /etc/nginx/nginx.conf << 'EOF'

    # HTTPS server
    server {
        listen 443 ssl;
        server_name DOMAIN_PLACEHOLDER;
        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;
        root /var/www/html;
        index index.html;
        client_max_body_size 10M;
        
        location /api/ {
            proxy_pass http://127.0.0.1:8890;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header REMOTE-HOST $remote_addr;
            proxy_set_header X-Forwarded-Proto https;
            proxy_set_header X-Forwarded-Port 443;
            
            # WebSocket support
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection $connection_upgrade;
            proxy_http_version 1.1;
            
            # SSL settings
            proxy_ssl_server_name off;
            proxy_ssl_name $proxy_host;
            
            # Timeout settings for SSH connections
            proxy_connect_timeout 60s;
            proxy_send_timeout 600s;
            proxy_read_timeout 600s;
            
            # Disable buffering for real-time data
            proxy_buffering off;
            add_header X-Cache $upstream_cache_status;
            add_header Cache-Control no-cache;
        }
        
        location /swagger/ {
            proxy_pass http://127.0.0.1:8890;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }
        
        # WebSocket endpoints for SSH connections
        location /v1/ {
            proxy_pass http://127.0.0.1:8890;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header REMOTE-HOST $remote_addr;
            proxy_set_header X-Forwarded-Proto https;
            proxy_set_header X-Forwarded-Port 443;
            
            # WebSocket support
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection $connection_upgrade;
            proxy_http_version 1.1;
            
            # Timeout settings for SSH connections
            proxy_connect_timeout 60s;
            proxy_send_timeout 600s;
            proxy_read_timeout 600s;
            
            # Disable buffering for real-time data
            proxy_buffering off;
            add_header X-Cache $upstream_cache_status;
            add_header Cache-Control no-cache;
        }
        
        location / {
            try_files $uri $uri/ /index.html;
        }
    }
}
EOF
            
            # Replace the domain placeholder
            sed -i "s/DOMAIN_PLACEHOLDER/$DOMAIN/g" /etc/nginx/nginx.conf
            
            # Create cron job for automatic certificate renewal attempt
            echo "0 0 * * * certbot --nginx --non-interactive --agree-tos --email admin@$DOMAIN --domains $DOMAIN && nginx -s reload || true" > /etc/cron.d/certbot-renewal
            chmod 644 /etc/cron.d/certbot-renewal
            echo "Created cron job for automatic certificate renewal"
        fi
    else
        echo "Detected HTTP frontend, using default nginx config"
    fi
fi

# Detect architecture and set database type
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    DB_TYPE="mysql"
    DB_DAEMON="mysqld"
else
    DB_TYPE="mariadb"
    DB_DAEMON="mariadbd"
fi
echo "Detected architecture: $ARCH, using database: $DB_TYPE"

chown -R mysql:mysql /var/lib/mysql /var/run/mysqld /var/log/mysql
chmod 755 /var/run/mysqld

# Check if database needs initialization
INIT_NEEDED=false
# Create database initialization flag file path (different from business init)
DB_INIT_FLAG="/var/lib/mysql/.mysql_initialized"

# Check if there are existing database files (persistent data)
EXISTING_DATA=false
if [ -d "/var/lib/mysql/mysql" ] && [ "$(ls -A /var/lib/mysql 2>/dev/null | wc -l)" -gt 0 ]; then
    EXISTING_DATA=true
    echo "Existing database data found - using persistent data..."
else
    echo "No existing database data found - checking initialization status..."
    # Check various conditions for initialization
    if [ ! -f "$DB_INIT_FLAG" ]; then
        echo "Database initialization flag not found - database needs initialization"
        INIT_NEEDED=true
    elif [ ! -d "/var/lib/mysql/mysql" ]; then
        echo "Database system directory not found - reinitializing database..."
        INIT_NEEDED=true
    elif [ "$(ls -A /var/lib/mysql 2>/dev/null | wc -l)" -eq 0 ]; then
        echo "Database directory is empty - reinitializing database..."
        INIT_NEEDED=true
    else
        echo "Database already initialized (flag exists and data present), skipping initialization..."
    fi
fi

# Always check and import default users if needed
CHECK_USERS=true

if [ "$INIT_NEEDED" = "true" ]; then
    # Stop any running database processes
    pkill -f "$DB_DAEMON" || true
    sleep 2
    # Remove old/corrupted data only when needed
    rm -rf /var/lib/mysql/*
    # Initialize database based on type
    if [ "$DB_TYPE" = "mysql" ]; then
        mysqld --initialize-insecure --user=mysql --datadir=/var/lib/mysql --skip-name-resolve
    else
        mariadb-install-db --user=mysql --datadir=/var/lib/mysql --skip-name-resolve
    fi
    if [ $? -ne 0 ]; then
        echo "$DB_TYPE initialization failed"
        exit 1
    fi
fi

# Configure database users and permissions only if initialization was needed
if [ "$INIT_NEEDED" = "true" ]; then
    echo "Configuring $DB_TYPE users and permissions..."
    pkill -f "$DB_DAEMON" || true
    sleep 2
    
    # Start temporary database server for configuration
    echo "Starting temporary $DB_TYPE server for configuration..."
    $DB_DAEMON --user=mysql --skip-networking --skip-grant-tables --socket=/var/run/mysqld/mysqld.sock --pid-file=/var/run/mysqld/mysqld.pid --log-error=/var/log/mysql/error.log &
    mysql_pid=$!
    
    for i in {1..30}; do
        if mysql --socket=/var/run/mysqld/mysqld.sock -e "SELECT 1" >/dev/null 2>&1; then
            echo "$DB_TYPE started successfully"
            break
        fi
        echo "Waiting for $DB_TYPE to start... ($i/30)"
        if [ $i -eq 30 ]; then
            echo "$DB_TYPE failed to start"
            kill $mysql_pid 2>/dev/null || true
            exit 1
        fi
        sleep 1
    done
    
    echo "Configuring $DB_TYPE users and database..."
    if [ "$DB_TYPE" = "mysql" ]; then
        mysql --socket=/var/run/mysqld/mysqld.sock <<SQLEND
FLUSH PRIVILEGES;
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '';
CREATE USER IF NOT EXISTS 'root'@'127.0.0.1' IDENTIFIED WITH mysql_native_password BY '';
CREATE USER IF NOT EXISTS 'root'@'%' IDENTIFIED WITH mysql_native_password BY '';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'root'@'127.0.0.1' WITH GRANT OPTION;
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;
CREATE DATABASE IF NOT EXISTS \`${MYSQL_DATABASE}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
FLUSH PRIVILEGES;
SQLEND
    else
        mysql --socket=/var/run/mysqld/mysqld.sock <<SQLEND
FLUSH PRIVILEGES;
SET PASSWORD FOR 'root'@'localhost' = PASSWORD('');
GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION;
CREATE USER IF NOT EXISTS 'root'@'127.0.0.1' IDENTIFIED BY '';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'127.0.0.1' WITH GRANT OPTION;
CREATE USER IF NOT EXISTS 'root'@'%' IDENTIFIED BY '';
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;
CREATE DATABASE IF NOT EXISTS \`${MYSQL_DATABASE}\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
FLUSH PRIVILEGES;
SQLEND
    fi
    
    # Import default admin and user data if no users exist
    echo "Checking if users exist..."
    # First, check if users table exists
    TABLE_EXISTS=$(mysql --socket=/var/run/mysqld/mysqld.sock -e "USE oneclickvirt; SHOW TABLES LIKE 'users';" 2>/dev/null | tail -n 1 || echo "")
    if [ -z "$TABLE_EXISTS" ]; then
        echo "Users table does not exist, creating and importing default data..."
        mysql --socket=/var/run/mysqld/mysqld.sock <<SQLEND
USE oneclickvirt;
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    level INT NOT NULL DEFAULT 1,
    status INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
INSERT INTO users (id, username, email, password, level, status, created_at, updated_at) VALUES
    (1, "admin", "admin@example.com", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", 5, 1, NOW(), NOW()),
    (2, "user", "user@example.com", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", 1, 1, NOW(), NOW());

-- Import system image default data
CREATE TABLE IF NOT EXISTS announcements (
  id bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at datetime(3) NULL DEFAULT NULL,
  updated_at datetime(3) NULL DEFAULT NULL,
  deleted_at datetime(3) NULL DEFAULT NULL,
  title varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  content longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  content_html longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  type varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'homepage',
  priority bigint NULL DEFAULT 0,
  status bigint NULL DEFAULT 1,
  is_sticky tinyint(1) NULL DEFAULT 0,
  start_time datetime(3) NULL DEFAULT NULL,
  end_time datetime(3) NULL DEFAULT NULL,
  created_by bigint UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (id) USING BTREE,
  INDEX idx_announcements_deleted_at(deleted_at ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

INSERT INTO announcements VALUES (1, '2025-12-30 15:38:08.631', '2025-12-30 15:38:08.631', NULL, '欢迎使用虚拟化管理平台', '欢迎使用虚拟化管理平台，支持Docker、LXD、Incus、Proxmox VE等多种虚拟化技术。本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。', '<p>欢迎使用虚拟化管理平台，支持<strong>Docker</strong>、<strong>LXD</strong>、<strong>Incus</strong>、<strong>Proxmox VE</strong>等多种虚拟化技术。</p><p>本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。</p>', 'homepage', 10, 1, 1, NULL, NULL, NULL);
INSERT INTO announcements VALUES (2, '2025-12-30 15:38:08.633', '2025-12-30 15:38:08.633', NULL, '系统维护通知', '为了提供更好的服务质量，我们会定期进行系统维护。维护期间可能会影响部分功能的使用，请您谅解。', '<p>为了提供更好的服务质量，我们会定期进行系统维护。</p>', 'topbar', 5, 1, 0, NULL, NULL, NULL);
INSERT INTO announcements VALUES (3, '2025-12-30 15:38:08.644', '2025-12-30 15:38:08.644', NULL, '新手使用指南', '如果您是第一次使用本平台，建议先阅读使用文档。您可以在右上角的帮助菜单中找到详细的操作指南。', '<p>如果您是第一次使用本平台，建议先阅读使用文档。</p>', 'homepage', 8, 1, 0, NULL, NULL, NULL);

CREATE TABLE IF NOT EXISTS invite_codes (
  id bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  code varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  creator_id bigint UNSIGNED NOT NULL,
  creator_name varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  description varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  max_uses bigint NOT NULL DEFAULT 1,
  used_count bigint NOT NULL DEFAULT 0,
  expires_at datetime(3) NULL DEFAULT NULL,
  status bigint NOT NULL DEFAULT 1,
  created_at datetime(3) NULL DEFAULT NULL,
  updated_at datetime(3) NULL DEFAULT NULL,
  deleted_at datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (id) USING BTREE,
  UNIQUE INDEX idx_invite_codes_code(code ASC) USING BTREE,
  INDEX idx_invite_codes_creator_id(creator_id ASC) USING BTREE,
  INDEX idx_invite_codes_expires_at(expires_at ASC) USING BTREE,
  INDEX idx_invite_codes_status(status ASC) USING BTREE,
  INDEX idx_invite_codes_deleted_at(deleted_at ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

INSERT INTO invite_codes VALUES (1, 'SC0Q19BW', 1, '', '', 1, 0, NULL, 1, '2025-12-31 10:59:55.167', '2025-12-31 10:59:55.167', NULL);

CREATE TABLE IF NOT EXISTS jwt_secrets (
  id bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  secret_key varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'JWT签名密钥',
  created_at datetime(3) NULL DEFAULT NULL,
  updated_at datetime(3) NULL DEFAULT NULL,
  deleted_at datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (id) USING BTREE,
  UNIQUE INDEX idx_jwt_secrets_secret_key(secret_key ASC) USING BTREE,
  INDEX idx_jwt_secrets_deleted_at(deleted_at ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

INSERT INTO jwt_secrets VALUES (1, 'b64dca17bf31d0e725285cccf00a6911a43b0e2c8d8d26ed458cdbf16e6a14b5', '2025-12-30 16:00:51.689', '2025-12-30 16:00:51.689', NULL);
SQLEND
        echo "Default admin, user and system image data imported successfully"
    else
        # Table exists, check if it has data
        USER_COUNT=$(mysql --socket=/var/run/mysqld/mysqld.sock -e "USE oneclickvirt; SELECT COUNT(*) FROM users;" 2>/dev/null | tail -n 1 || echo 0)
        if [ "$USER_COUNT" -eq 0 ]; then
            echo "Users table exists but is empty, importing default admin and user data..."
            mysql --socket=/var/run/mysqld/mysqld.sock <<SQLEND
USE oneclickvirt;
INSERT INTO users (id, username, email, password, level, status, created_at, updated_at) VALUES
    (1, "admin", "admin@example.com", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", 5, 1, NOW(), NOW()),
    (2, "user", "user@example.com", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", 1, 1, NOW(), NOW());

-- Import system image default data if not exists
INSERT IGNORE INTO announcements (id, created_at, updated_at, deleted_at, title, content, content_html, type, priority, status, is_sticky, start_time, end_time, created_by) VALUES
(1, '2025-12-30 15:38:08.631', '2025-12-30 15:38:08.631', NULL, '欢迎使用虚拟化管理平台', '欢迎使用虚拟化管理平台，支持Docker、LXD、Incus、Proxmox VE等多种虚拟化技术。本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。', '<p>欢迎使用虚拟化管理平台，支持<strong>Docker</strong>、<strong>LXD</strong>、<strong>Incus</strong>、<strong>Proxmox VE</strong>等多种虚拟化技术。</p><p>本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。</p>', 'homepage', 10, 1, 1, NULL, NULL, NULL),
(2, '2025-12-30 15:38:08.633', '2025-12-30 15:38:08.633', NULL, '系统维护通知', '为了提供更好的服务质量，我们会定期进行系统维护。维护期间可能会影响部分功能的使用，请您谅解。', '<p>为了提供更好的服务质量，我们会定期进行系统维护。</p>', 'topbar', 5, 1, 0, NULL, NULL, NULL),
(3, '2025-12-30 15:38:08.644', '2025-12-30 15:38:08.644', NULL, '新手使用指南', '如果您是第一次使用本平台，建议先阅读使用文档。您可以在右上角的帮助菜单中找到详细的操作指南。', '<p>如果您是第一次使用本平台，建议先阅读使用文档。</p>', 'homepage', 8, 1, 0, NULL, NULL, NULL);

INSERT IGNORE INTO invite_codes (id, code, creator_id, creator_name, description, max_uses, used_count, expires_at, status, created_at, updated_at, deleted_at) VALUES
(1, 'SC0Q19BW', 1, '', '', 1, 0, NULL, 1, '2025-12-31 10:59:55.167', '2025-12-31 10:59:55.167', NULL);

INSERT IGNORE INTO jwt_secrets (id, secret_key, created_at, updated_at, deleted_at) VALUES
(1, 'b64dca17bf31d0e725285cccf00a6911a43b0e2c8d8d26ed458cdbf16e6a14b5', '2025-12-30 16:00:51.689', '2025-12-30 16:00:51.689', NULL);

-- Create products table if not exists
CREATE TABLE IF NOT EXISTS products (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品名称',
  description text COLLATE utf8mb4_unicode_ci COMMENT '产品描述',
  level bigint NOT NULL COMMENT '产品等级',
  price bigint NOT NULL COMMENT '价格(分)',
  period bigint NOT NULL COMMENT '周期(月), 0为永久',
  cpu bigint NOT NULL COMMENT 'CPU核心数',
  memory bigint NOT NULL COMMENT '内存(MB)',
  disk bigint NOT NULL COMMENT '磁盘(MB)',
  bandwidth bigint NOT NULL COMMENT '带宽(Mbps)',
  traffic bigint NOT NULL COMMENT '流量限制(MB)',
  max_instances bigint NOT NULL COMMENT '最大实例数',
  is_enabled bigint DEFAULT '1' COMMENT '是否启用(1:启用, 0:禁用)',
  sort_order bigint DEFAULT '0' COMMENT '排序',
  features text COLLATE utf8mb4_unicode_ci COMMENT '特性(JSON格式)',
  allow_repeat bigint DEFAULT '1' COMMENT '是否允许重复购买(1:允许, 0:不允许)',
  stock bigint DEFAULT '-1' COMMENT '库存(-1为无限)',
  sold_count bigint DEFAULT '0' COMMENT '已售数量',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create roles table if not exists
CREATE TABLE IF NOT EXISTS roles (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  code varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  description text COLLATE utf8mb4_unicode_ci,
  status bigint DEFAULT '1',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_roles_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create user_roles table if not exists
CREATE TABLE IF NOT EXISTS user_roles (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id bigint unsigned NOT NULL,
  role_id bigint unsigned NOT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_user_roles_user_role (user_id,role_id),
  KEY idx_user_roles_role_id (role_id),
  CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create system_configs table if not exists
CREATE TABLE IF NOT EXISTS system_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text COLLATE utf8mb4_unicode_ci NOT NULL,
  description varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_system_configs_key (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create site_configs table if not exists
CREATE TABLE IF NOT EXISTS site_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `group` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  description varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_site_configs_key (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create domain_configs table if not exists
CREATE TABLE IF NOT EXISTS domain_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  max_domains_per_user bigint DEFAULT '3',
  max_domains_per_agent_user bigint DEFAULT '5',
  default_ttl bigint DEFAULT '300',
  auto_ssl bigint DEFAULT '0',
  allowed_suffixes text COLLATE utf8mb4_unicode_ci,
  dns_type varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT 'dnsmasq',
  dns_config_path varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '/etc/dnsmasq.d/oneclickvirt-hosts.conf',
  nginx_config_path varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '/etc/nginx/conf.d/oneclickvirt-domains',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert roles data
INSERT IGNORE INTO roles (name, code, description, status, created_at, updated_at) VALUES
('admin', 'admin', '系统管理员角色', 1, NOW(), NOW()),
('user', 'user', '普通用户角色', 1, NOW(), NOW());

-- Update users table to include necessary fields
ALTER TABLE users ADD COLUMN IF NOT EXISTS uuid varchar(36) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS nickname varchar(50) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone varchar(20) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS level_expire_at datetime DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS user_type varchar(20) DEFAULT 'user';
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified tinyint(1) DEFAULT '0';
ALTER TABLE users ADD COLUMN IF NOT EXISTS real_name_verified tinyint(1) DEFAULT '0';
ALTER TABLE users ADD UNIQUE KEY IF NOT EXISTS idx_users_uuid (uuid);

-- Update existing users with missing fields
UPDATE users SET uuid = IFNULL(uuid, CONCAT('user-', id)), nickname = IFNULL(nickname, username), user_type = IFNULL(user_type, 'user') WHERE id IN (1, 2);

-- Insert user_roles data
INSERT IGNORE INTO user_roles (user_id, role_id, created_at, updated_at) VALUES
(1, 1, NOW(), NOW()),
(2, 2, NOW(), NOW());

-- Insert system_configs data
INSERT IGNORE INTO system_configs (`key`, `value`, description, created_at, updated_at) VALUES
('site_name', 'OneClickVirt', '网站名称', NOW(), NOW()),
('site_description', '虚拟化管理平台', '网站描述', NOW(), NOW()),
('site_keywords', '虚拟化,Docker,LXD,Incus,Proxmox', '网站关键词', NOW(), NOW()),
('enable_registration', 'true', '是否开启注册', NOW(), NOW()),
('enable_email_verify', 'false', '是否开启邮箱验证', NOW(), NOW()),
('default_user_level', '1', '默认用户等级', NOW(), NOW()),
('max_instances_per_user', '10', '每个用户最大实例数', NOW(), NOW()),
('default_instance_expiry_days', '30', '默认实例过期天数', NOW(), NOW()),
('enable_email_verification', 'false', '是否开启邮箱验证（注册后需验证邮箱）', NOW(), NOW()),
('email_activation_expire_hours', '24', '邮箱激活链接过期时间（小时）', NOW(), NOW()),
('enable_real_name', 'false', '是否开启实名认证', NOW(), NOW()),
('require_real_name', 'false', '是否强制实名认证后才能使用服务', NOW(), NOW()),
('enable_agent', 'true', '是否开启代理商功能', NOW(), NOW());

-- Insert site_configs data
INSERT IGNORE INTO site_configs (`key`, `value`, `type`, `group`, description, created_at, updated_at) VALUES
('site_name', 'OneClickVirt', 'string', 'basic', '网站名称', NOW(), NOW()),
('site_icon_url', '/favicon.ico', 'string', 'basic', '网站图标URL', NOW(), NOW()),
('site_logo_url', '/logo.png', 'string', 'basic', '网站Logo URL', NOW(), NOW()),
('footer_text', '© 2025 OneClickVirt. All rights reserved.', 'string', 'basic', '页脚文字', NOW(), NOW()),
('icp_number', '', 'string', 'basic', 'ICP备案号', NOW(), NOW()),
('police_number', '', 'string', 'basic', '公安备案号', NOW(), NOW());

-- Insert domain_configs data
INSERT IGNORE INTO domain_configs (max_domains_per_user, max_domains_per_agent_user, default_ttl, auto_ssl, allowed_suffixes, dns_type, dns_config_path, nginx_config_path, created_at, updated_at) VALUES
(3, 5, 300, 0, '', 'dnsmasq', '/etc/dnsmasq.d/oneclickvirt-hosts.conf', '/etc/nginx/conf.d/oneclickvirt-domains', NOW(), NOW());

-- Update products table to match init.sql structure
ALTER TABLE products ADD COLUMN IF NOT EXISTS billing_cycle varchar(20) DEFAULT 'monthly';
ALTER TABLE products ADD COLUMN IF NOT EXISTS cpu_limit bigint DEFAULT cpu;
ALTER TABLE products ADD COLUMN IF NOT EXISTS memory_limit bigint DEFAULT memory;
ALTER TABLE products ADD COLUMN IF NOT EXISTS disk_limit bigint DEFAULT disk;
ALTER TABLE products ADD COLUMN IF NOT EXISTS bandwidth_limit bigint DEFAULT bandwidth;
ALTER TABLE products ADD COLUMN IF NOT EXISTS instance_limit bigint DEFAULT max_instances;
ALTER TABLE products ADD COLUMN IF NOT EXISTS status bigint DEFAULT is_enabled;

-- Insert products data
INSERT IGNORE INTO products (id, name, description, price, billing_cycle, cpu_limit, memory_limit, disk_limit, bandwidth_limit, instance_limit, features, status, sort_order, created_at, updated_at) VALUES
(1, '入门套餐', '适合个人用户的基础套餐，包含基本的虚拟化功能', 0.00, 'monthly', 1, 512, 10240, 100, 1, '{"cpu": "1核", "memory": "512MB", "disk": "10GB", "bandwidth": "100Mbps", "instances": "1个实例"}', 1, 1, NOW(), NOW()),
(2, '标准套餐', '适合小型团队的标准套餐，包含更多资源', 9.90, 'monthly', 2, 1024, 20480, 200, 3, '{"cpu": "2核", "memory": "1GB", "disk": "20GB", "bandwidth": "200Mbps", "instances": "3个实例"}', 1, 2, NOW(), NOW()),
(3, '专业套餐', '适合中型团队的专业套餐，包含完整功能', 29.90, 'monthly', 4, 2048, 40960, 500, 5, '{"cpu": "4核", "memory": "2GB", "disk": "40GB", "bandwidth": "500Mbps", "instances": "5个实例"}', 1, 3, NOW(), NOW()),
(4, '企业套餐', '适合大型团队的企业套餐，包含无限资源', 99.90, 'monthly', 8, 4096, 102400, 1000, 10, '{"cpu": "8核", "memory": "4GB", "disk": "100GB", "bandwidth": "1000Mbps", "instances": "10个实例"}', 1, 4, NOW(), NOW());
SQLEND
            echo "Default admin, user and system image data imported successfully"
        else
            echo "Users already exist, checking system image data..."
            # Check if system image data exists, import if not
            ANNOUNCEMENT_COUNT=$(mysql --socket=/var/run/mysqld/mysqld.sock -e "USE oneclickvirt; SELECT COUNT(*) FROM announcements;" 2>/dev/null | tail -n 1 || echo 0)
            if [ "$ANNOUNCEMENT_COUNT" -eq 0 ]; then
                echo "Importing system image default data..."
                mysql --socket=/var/run/mysqld/mysqld.sock <<SQLEND
USE oneclickvirt;
INSERT INTO announcements (id, created_at, updated_at, deleted_at, title, content, content_html, type, priority, status, is_sticky, start_time, end_time, created_by) VALUES
(1, '2025-12-30 15:38:08.631', '2025-12-30 15:38:08.631', NULL, '欢迎使用虚拟化管理平台', '欢迎使用虚拟化管理平台，支持Docker、LXD、Incus、Proxmox VE等多种虚拟化技术。本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。', '<p>欢迎使用虚拟化管理平台，支持<strong>Docker</strong>、<strong>LXD</strong>、<strong>Incus</strong>、<strong>Proxmox VE</strong>等多种虚拟化技术。</p><p>本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。</p>', 'homepage', 10, 1, 1, NULL, NULL, NULL),
(2, '2025-12-30 15:38:08.633', '2025-12-30 15:38:08.633', NULL, '系统维护通知', '为了提供更好的服务质量，我们会定期进行系统维护。维护期间可能会影响部分功能的使用，请您谅解。', '<p>为了提供更好的服务质量，我们会定期进行系统维护。</p>', 'topbar', 5, 1, 0, NULL, NULL, NULL),
(3, '2025-12-30 15:38:08.644', '2025-12-30 15:38:08.644', NULL, '新手使用指南', '如果您是第一次使用本平台，建议先阅读使用文档。您可以在右上角的帮助菜单中找到详细的操作指南。', '<p>如果您是第一次使用本平台，建议先阅读使用文档。</p>', 'homepage', 8, 1, 0, NULL, NULL, NULL);

INSERT INTO invite_codes (id, code, creator_id, creator_name, description, max_uses, used_count, expires_at, status, created_at, updated_at, deleted_at) VALUES
(1, 'SC0Q19BW', 1, '', '', 1, 0, NULL, 1, '2025-12-31 10:59:55.167', '2025-12-31 10:59:55.167', NULL);

INSERT INTO jwt_secrets (id, secret_key, created_at, updated_at, deleted_at) VALUES
(1, 'b64dca17bf31d0e725285cccf00a6911a43b0e2c8d8d26ed458cdbf16e6a14b5', '2025-12-30 16:00:51.689', '2025-12-30 16:00:51.689', NULL);

-- Create products table if not exists
CREATE TABLE IF NOT EXISTS products (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品名称',
  description text COLLATE utf8mb4_unicode_ci COMMENT '产品描述',
  level bigint NOT NULL COMMENT '产品等级',
  price bigint NOT NULL COMMENT '价格(分)',
  period bigint NOT NULL COMMENT '周期(月), 0为永久',
  cpu bigint NOT NULL COMMENT 'CPU核心数',
  memory bigint NOT NULL COMMENT '内存(MB)',
  disk bigint NOT NULL COMMENT '磁盘(MB)',
  bandwidth bigint NOT NULL COMMENT '带宽(Mbps)',
  traffic bigint NOT NULL COMMENT '流量限制(MB)',
  max_instances bigint NOT NULL COMMENT '最大实例数',
  is_enabled bigint DEFAULT '1' COMMENT '是否启用(1:启用, 0:禁用)',
  sort_order bigint DEFAULT '0' COMMENT '排序',
  features text COLLATE utf8mb4_unicode_ci COMMENT '特性(JSON格式)',
  allow_repeat bigint DEFAULT '1' COMMENT '是否允许重复购买(1:允许, 0:不允许)',
  stock bigint DEFAULT '-1' COMMENT '库存(-1为无限)',
  sold_count bigint DEFAULT '0' COMMENT '已售数量',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create roles table if not exists
CREATE TABLE IF NOT EXISTS roles (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  code varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  description text COLLATE utf8mb4_unicode_ci,
  status bigint DEFAULT '1',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_roles_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create user_roles table if not exists
CREATE TABLE IF NOT EXISTS user_roles (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id bigint unsigned NOT NULL,
  role_id bigint unsigned NOT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_user_roles_user_role (user_id,role_id),
  KEY idx_user_roles_role_id (role_id),
  CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create system_configs table if not exists
CREATE TABLE IF NOT EXISTS system_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text COLLATE utf8mb4_unicode_ci NOT NULL,
  description varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_system_configs_key (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create site_configs table if not exists
CREATE TABLE IF NOT EXISTS site_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `group` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  description varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_site_configs_key (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create domain_configs table if not exists
CREATE TABLE IF NOT EXISTS domain_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  max_domains_per_user bigint DEFAULT '3',
  max_domains_per_agent_user bigint DEFAULT '5',
  default_ttl bigint DEFAULT '300',
  auto_ssl bigint DEFAULT '0',
  allowed_suffixes text COLLATE utf8mb4_unicode_ci,
  dns_type varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT 'dnsmasq',
  dns_config_path varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '/etc/dnsmasq.d/oneclickvirt-hosts.conf',
  nginx_config_path varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '/etc/nginx/conf.d/oneclickvirt-domains',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert roles data
INSERT IGNORE INTO roles (name, code, description, status, created_at, updated_at) VALUES
('admin', 'admin', '系统管理员角色', 1, NOW(), NOW()),
('user', 'user', '普通用户角色', 1, NOW(), NOW());

-- Update users table to include necessary fields
ALTER TABLE users ADD COLUMN IF NOT EXISTS uuid varchar(36) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS nickname varchar(50) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone varchar(20) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS level_expire_at datetime DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS user_type varchar(20) DEFAULT 'user';
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified tinyint(1) DEFAULT '0';
ALTER TABLE users ADD COLUMN IF NOT EXISTS real_name_verified tinyint(1) DEFAULT '0';
ALTER TABLE users ADD UNIQUE KEY IF NOT EXISTS idx_users_uuid (uuid);

-- Update existing users with missing fields
UPDATE users SET uuid = IFNULL(uuid, CONCAT('user-', id)), nickname = IFNULL(nickname, username), user_type = IFNULL(user_type, 'user') WHERE id IN (1, 2);

-- Insert user_roles data
INSERT IGNORE INTO user_roles (user_id, role_id, created_at, updated_at) VALUES
(1, 1, NOW(), NOW()),
(2, 2, NOW(), NOW());

-- Insert system_configs data
INSERT IGNORE INTO system_configs (`key`, `value`, description, created_at, updated_at) VALUES
('site_name', 'OneClickVirt', '网站名称', NOW(), NOW()),
('site_description', '虚拟化管理平台', '网站描述', NOW(), NOW()),
('site_keywords', '虚拟化,Docker,LXD,Incus,Proxmox', '网站关键词', NOW(), NOW()),
('enable_registration', 'true', '是否开启注册', NOW(), NOW()),
('enable_email_verify', 'false', '是否开启邮箱验证', NOW(), NOW()),
('default_user_level', '1', '默认用户等级', NOW(), NOW()),
('max_instances_per_user', '10', '每个用户最大实例数', NOW(), NOW()),
('default_instance_expiry_days', '30', '默认实例过期天数', NOW(), NOW()),
('enable_email_verification', 'false', '是否开启邮箱验证（注册后需验证邮箱）', NOW(), NOW()),
('email_activation_expire_hours', '24', '邮箱激活链接过期时间（小时）', NOW(), NOW()),
('enable_real_name', 'false', '是否开启实名认证', NOW(), NOW()),
('require_real_name', 'false', '是否强制实名认证后才能使用服务', NOW(), NOW()),
('enable_agent', 'true', '是否开启代理商功能', NOW(), NOW());

-- Insert site_configs data
INSERT IGNORE INTO site_configs (`key`, `value`, `type`, `group`, description, created_at, updated_at) VALUES
('site_name', 'OneClickVirt', 'string', 'basic', '网站名称', NOW(), NOW()),
('site_icon_url', '/favicon.ico', 'string', 'basic', '网站图标URL', NOW(), NOW()),
('site_logo_url', '/logo.png', 'string', 'basic', '网站Logo URL', NOW(), NOW()),
('footer_text', '© 2025 OneClickVirt. All rights reserved.', 'string', 'basic', '页脚文字', NOW(), NOW()),
('icp_number', '', 'string', 'basic', 'ICP备案号', NOW(), NOW()),
('police_number', '', 'string', 'basic', '公安备案号', NOW(), NOW());

-- Insert domain_configs data
INSERT IGNORE INTO domain_configs (max_domains_per_user, max_domains_per_agent_user, default_ttl, auto_ssl, allowed_suffixes, dns_type, dns_config_path, nginx_config_path, created_at, updated_at) VALUES
(3, 5, 300, 0, '', 'dnsmasq', '/etc/dnsmasq.d/oneclickvirt-hosts.conf', '/etc/nginx/conf.d/oneclickvirt-domains', NOW(), NOW());

-- Update products table to match init.sql structure
ALTER TABLE products ADD COLUMN IF NOT EXISTS billing_cycle varchar(20) DEFAULT 'monthly';
ALTER TABLE products ADD COLUMN IF NOT EXISTS cpu_limit bigint DEFAULT cpu;
ALTER TABLE products ADD COLUMN IF NOT EXISTS memory_limit bigint DEFAULT memory;
ALTER TABLE products ADD COLUMN IF NOT EXISTS disk_limit bigint DEFAULT disk;
ALTER TABLE products ADD COLUMN IF NOT EXISTS bandwidth_limit bigint DEFAULT bandwidth;
ALTER TABLE products ADD COLUMN IF NOT EXISTS instance_limit bigint DEFAULT max_instances;
ALTER TABLE products ADD COLUMN IF NOT EXISTS status bigint DEFAULT is_enabled;

-- Insert products data
INSERT IGNORE INTO products (id, name, description, price, billing_cycle, cpu_limit, memory_limit, disk_limit, bandwidth_limit, instance_limit, features, status, sort_order, created_at, updated_at) VALUES
(1, '入门套餐', '适合个人用户的基础套餐，包含基本的虚拟化功能', 0.00, 'monthly', 1, 512, 10240, 100, 1, '{"cpu": "1核", "memory": "512MB", "disk": "10GB", "bandwidth": "100Mbps", "instances": "1个实例"}', 1, 1, NOW(), NOW()),
(2, '标准套餐', '适合小型团队的标准套餐，包含更多资源', 9.90, 'monthly', 2, 1024, 20480, 200, 3, '{"cpu": "2核", "memory": "1GB", "disk": "20GB", "bandwidth": "200Mbps", "instances": "3个实例"}', 1, 2, NOW(), NOW()),
(3, '专业套餐', '适合中型团队的专业套餐，包含完整功能', 29.90, 'monthly', 4, 2048, 40960, 500, 5, '{"cpu": "4核", "memory": "2GB", "disk": "40GB", "bandwidth": "500Mbps", "instances": "5个实例"}', 1, 3, NOW(), NOW()),
(4, '企业套餐', '适合大型团队的企业套餐，包含无限资源', 99.90, 'monthly', 8, 4096, 102400, 1000, 10, '{"cpu": "8核", "memory": "4GB", "disk": "100GB", "bandwidth": "1000Mbps", "instances": "10个实例"}', 1, 4, NOW(), NOW());
SQLEND
                echo "System image default data imported successfully"
            else
                echo "System image data already exists, skipping import"
            fi
        fi
    fi
    
    kill $mysql_pid
    wait $mysql_pid 2>/dev/null || true
    echo "$DB_TYPE configuration completed."
    # Create database initialization flag to prevent re-initialization
    echo "$(date): Database initialized successfully with $DB_TYPE" > "$DB_INIT_FLAG"
    echo "Created database initialization flag at $DB_INIT_FLAG"
else
    echo "Database already configured, skipping user configuration..."
fi

# Create supervisor configuration dynamically
echo "Creating supervisor configuration for $DB_TYPE..."
cat > /etc/supervisor/conf.d/supervisord.conf <<SUPEREND
[unix_http_server]
file=/var/run/supervisor.sock
chmod=0700

[supervisorctl]
serverurl=unix:///var/run/supervisor.sock

[rpcinterface:supervisor]
supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface

[supervisord]
nodaemon=true
user=root

[program:mysql]
SUPEREND

if [ "$DB_TYPE" = "mysql" ]; then
    echo "command=/usr/sbin/mysqld --defaults-file=/etc/mysql/conf.d/custom.cnf --lc-messages=en_US" >> /etc/supervisor/conf.d/supervisord.conf
else
    echo "command=/usr/sbin/mariadbd --defaults-file=/etc/mysql/conf.d/custom.cnf" >> /etc/supervisor/conf.d/supervisord.conf
fi

cat >> /etc/supervisor/conf.d/supervisord.conf <<SUPEREND2
autostart=true
autorestart=true
user=mysql
priority=1
stdout_logfile=/var/log/supervisor/mysql.log
stderr_logfile=/var/log/supervisor/mysql_error.log
stdout_logfile_maxbytes=10MB
stderr_logfile_maxbytes=10MB
startsecs=10
startretries=3

[program:app]
command=/bin/bash -c "sleep 12 && /app/main"
directory=/app
autostart=true
autorestart=true
user=root
priority=2
environment=DB_HOST="127.0.0.1",DB_PORT="3306",DB_USER="root",DB_PASSWORD="",DB_NAME="oneclickvirt"
startsecs=1

[program:nginx]
command=/usr/sbin/nginx -g "daemon off;"
autostart=true
autorestart=true
user=root
priority=3
SUPEREND2

export DB_HOST="127.0.0.1"
export DB_PORT="3306"
export DB_NAME="$MYSQL_DATABASE"
export DB_USER="root"
export DB_PASSWORD=""

# Create a script to check and import users after services start
cat > /check_users.sh << 'EOF2'
#!/bin/bash

DB_HOST="localhost"

# Wait for MySQL to start
echo "Checking if users exist..."
for i in {1..60}; do
    if mysql -h "$DB_HOST" -u root -e "SELECT 1" >/dev/null 2>&1; then
        echo "MySQL started successfully"
        break
    fi
    echo "Waiting for MySQL to start... ($i/60)"
    if [ $i -eq 60 ]; then
        echo "MySQL failed to start"
        exit 1
    fi
    sleep 1
done

# First, check if users table exists
TABLE_EXISTS=$(mysql -h "$DB_HOST" -u root -e "USE oneclickvirt; SHOW TABLES LIKE 'users';" 2>/dev/null | tail -n 1 || echo "")
if [ -z "$TABLE_EXISTS" ]; then
    echo "Users table does not exist, skipping import"
else
    # Table exists, check if it has data
    USER_COUNT=$(mysql -h "$DB_HOST" -u root -e "USE oneclickvirt; SELECT COUNT(*) FROM users;" 2>/dev/null | tail -n 1 || echo 0)
    if [ "$USER_COUNT" -eq 0 ]; then
        echo "Users table exists but is empty, importing default admin and user data..."
        # Generate UUIDs for users
        ADMIN_UUID=$(cat /proc/sys/kernel/random/uuid)
        USER_UUID=$(cat /proc/sys/kernel/random/uuid)
        
        # Use simple password 'password' and let the system hash it later
        # For now, we'll use a placeholder and the system will handle the hashing
        mysql -h "$DB_HOST" -u root -e "USE oneclickvirt; INSERT INTO users (id, uuid, username, password, email, level, user_type, status, created_at, updated_at, max_instances, max_cpu, max_memory, max_disk) VALUES (1, '${ADMIN_UUID}', 'admin', 'password', 'admin@example.com', 5, 'admin', 1, NOW(), NOW(), 10, 8, 8192, 102400), (2, '${USER_UUID}', 'user', 'password', 'user@example.com', 1, 'user', 1, NOW(), NOW(), 1, 1, 512, 10240);"
        
        # Import system image default data
        echo "Importing system image default data..."
        mysql -h "$DB_HOST" -u root <<SQLEND
USE oneclickvirt;

-- Create tables if they don't exist
CREATE TABLE IF NOT EXISTS announcements (
  id bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at datetime(3) NULL DEFAULT NULL,
  updated_at datetime(3) NULL DEFAULT NULL,
  deleted_at datetime(3) NULL DEFAULT NULL,
  title varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  content longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  content_html longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  type varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'homepage',
  priority bigint NULL DEFAULT 0,
  status bigint NULL DEFAULT 1,
  is_sticky tinyint(1) NULL DEFAULT 0,
  start_time datetime(3) NULL DEFAULT NULL,
  end_time datetime(3) NULL DEFAULT NULL,
  created_by bigint UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (id) USING BTREE,
  INDEX idx_announcements_deleted_at(deleted_at ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

CREATE TABLE IF NOT EXISTS invite_codes (
  id bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  code varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  creator_id bigint UNSIGNED NOT NULL,
  creator_name varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  description varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  max_uses bigint NOT NULL DEFAULT 1,
  used_count bigint NOT NULL DEFAULT 0,
  expires_at datetime(3) NULL DEFAULT NULL,
  status bigint NOT NULL DEFAULT 1,
  created_at datetime(3) NULL DEFAULT NULL,
  updated_at datetime(3) NULL DEFAULT NULL,
  deleted_at datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (id) USING BTREE,
  UNIQUE INDEX idx_invite_codes_code(code ASC) USING BTREE,
  INDEX idx_invite_codes_creator_id(creator_id ASC) USING BTREE,
  INDEX idx_invite_codes_expires_at(expires_at ASC) USING BTREE,
  INDEX idx_invite_codes_status(status ASC) USING BTREE,
  INDEX idx_invite_codes_deleted_at(deleted_at ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

CREATE TABLE IF NOT EXISTS jwt_secrets (
  id bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  secret_key varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'JWT签名密钥',
  created_at datetime(3) NULL DEFAULT NULL,
  updated_at datetime(3) NULL DEFAULT NULL,
  deleted_at datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (id) USING BTREE,
  UNIQUE INDEX idx_jwt_secrets_secret_key(secret_key ASC) USING BTREE,
  INDEX idx_jwt_secrets_deleted_at(deleted_at ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- Insert data
INSERT IGNORE INTO announcements (id, created_at, updated_at, deleted_at, title, content, content_html, type, priority, status, is_sticky, start_time, end_time, created_by) VALUES
(1, '2025-12-30 15:38:08.631', '2025-12-30 15:38:08.631', NULL, '欢迎使用虚拟化管理平台', '欢迎使用虚拟化管理平台，支持Docker、LXD、Incus、Proxmox VE等多种虚拟化技术。本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。', '<p>欢迎使用虚拟化管理平台，支持<strong>Docker</strong>、<strong>LXD</strong>、<strong>Incus</strong>、<strong>Proxmox VE</strong>等多种虚拟化技术。</p><p>本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。</p>', 'homepage', 10, 1, 1, NULL, NULL, NULL),
(2, '2025-12-30 15:38:08.633', '2025-12-30 15:38:08.633', NULL, '系统维护通知', '为了提供更好的服务质量，我们会定期进行系统维护。维护期间可能会影响部分功能的使用，请您谅解。', '<p>为了提供更好的服务质量，我们会定期进行系统维护。</p>', 'topbar', 5, 1, 0, NULL, NULL, NULL),
(3, '2025-12-30 15:38:08.644', '2025-12-30 15:38:08.644', NULL, '新手使用指南', '如果您是第一次使用本平台，建议先阅读使用文档。您可以在右上角的帮助菜单中找到详细的操作指南。', '<p>如果您是第一次使用本平台，建议先阅读使用文档。</p>', 'homepage', 8, 1, 0, NULL, NULL, NULL);

INSERT IGNORE INTO invite_codes (id, code, creator_id, creator_name, description, max_uses, used_count, expires_at, status, created_at, updated_at, deleted_at) VALUES
(1, 'SC0Q19BW', 1, '', '', 1, 0, NULL, 1, '2025-12-31 10:59:55.167', '2025-12-31 10:59:55.167', NULL);

INSERT IGNORE INTO jwt_secrets (id, secret_key, created_at, updated_at, deleted_at) VALUES
(1, 'b64dca17bf31d0e725285cccf00a6911a43b0e2c8d8d26ed458cdbf16e6a14b5', '2025-12-30 16:00:51.689', '2025-12-30 16:00:51.689', NULL);

-- Create products table if not exists
CREATE TABLE IF NOT EXISTS products (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品名称',
  description text COLLATE utf8mb4_unicode_ci COMMENT '产品描述',
  level bigint NOT NULL COMMENT '产品等级',
  price bigint NOT NULL COMMENT '价格(分)',
  period bigint NOT NULL COMMENT '周期(月), 0为永久',
  cpu bigint NOT NULL COMMENT 'CPU核心数',
  memory bigint NOT NULL COMMENT '内存(MB)',
  disk bigint NOT NULL COMMENT '磁盘(MB)',
  bandwidth bigint NOT NULL COMMENT '带宽(Mbps)',
  traffic bigint NOT NULL COMMENT '流量限制(MB)',
  max_instances bigint NOT NULL COMMENT '最大实例数',
  is_enabled bigint DEFAULT '1' COMMENT '是否启用(1:启用, 0:禁用)',
  sort_order bigint DEFAULT '0' COMMENT '排序',
  features text COLLATE utf8mb4_unicode_ci COMMENT '特性(JSON格式)',
  allow_repeat bigint DEFAULT '1' COMMENT '是否允许重复购买(1:允许, 0:不允许)',
  stock bigint DEFAULT '-1' COMMENT '库存(-1为无限)',
  sold_count bigint DEFAULT '0' COMMENT '已售数量',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create roles table if not exists
CREATE TABLE IF NOT EXISTS roles (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  code varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  description text COLLATE utf8mb4_unicode_ci,
  status bigint DEFAULT '1',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_roles_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create user_roles table if not exists
CREATE TABLE IF NOT EXISTS user_roles (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id bigint unsigned NOT NULL,
  role_id bigint unsigned NOT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_user_roles_user_role (user_id,role_id),
  KEY idx_user_roles_role_id (role_id),
  CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create system_configs table if not exists
CREATE TABLE IF NOT EXISTS system_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text COLLATE utf8mb4_unicode_ci NOT NULL,
  description varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_system_configs_key (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create site_configs table if not exists
CREATE TABLE IF NOT EXISTS site_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `group` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  description varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_site_configs_key (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create domain_configs table if not exists
CREATE TABLE IF NOT EXISTS domain_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  max_domains_per_user bigint DEFAULT '3',
  max_domains_per_agent_user bigint DEFAULT '5',
  default_ttl bigint DEFAULT '300',
  auto_ssl bigint DEFAULT '0',
  allowed_suffixes text COLLATE utf8mb4_unicode_ci,
  dns_type varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT 'dnsmasq',
  dns_config_path varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '/etc/dnsmasq.d/oneclickvirt-hosts.conf',
  nginx_config_path varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '/etc/nginx/conf.d/oneclickvirt-domains',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert roles data
INSERT IGNORE INTO roles (name, code, description, status, created_at, updated_at) VALUES
('admin', 'admin', '系统管理员角色', 1, NOW(), NOW()),
('user', 'user', '普通用户角色', 1, NOW(), NOW());

-- Update users table to include necessary fields
ALTER TABLE users ADD COLUMN IF NOT EXISTS uuid varchar(36) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS nickname varchar(50) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone varchar(20) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS level_expire_at datetime DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS user_type varchar(20) DEFAULT 'user';
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified tinyint(1) DEFAULT '0';
ALTER TABLE users ADD COLUMN IF NOT EXISTS real_name_verified tinyint(1) DEFAULT '0';
ALTER TABLE users ADD UNIQUE KEY IF NOT EXISTS idx_users_uuid (uuid);

-- Update existing users with missing fields
UPDATE users SET uuid = IFNULL(uuid, CONCAT('user-', id)), nickname = IFNULL(nickname, username), user_type = IFNULL(user_type, 'user') WHERE id IN (1, 2);

-- Insert user_roles data
INSERT IGNORE INTO user_roles (user_id, role_id, created_at, updated_at) VALUES
(1, 1, NOW(), NOW()),
(2, 2, NOW(), NOW());

-- Insert system_configs data
INSERT IGNORE INTO system_configs (`key`, `value`, description, created_at, updated_at) VALUES
('site_name', 'OneClickVirt', '网站名称', NOW(), NOW()),
('site_description', '虚拟化管理平台', '网站描述', NOW(), NOW()),
('site_keywords', '虚拟化,Docker,LXD,Incus,Proxmox', '网站关键词', NOW(), NOW()),
('enable_registration', 'true', '是否开启注册', NOW(), NOW()),
('enable_email_verify', 'false', '是否开启邮箱验证', NOW(), NOW()),
('default_user_level', '1', '默认用户等级', NOW(), NOW()),
('max_instances_per_user', '10', '每个用户最大实例数', NOW(), NOW()),
('default_instance_expiry_days', '30', '默认实例过期天数', NOW(), NOW()),
('enable_email_verification', 'false', '是否开启邮箱验证（注册后需验证邮箱）', NOW(), NOW()),
('email_activation_expire_hours', '24', '邮箱激活链接过期时间（小时）', NOW(), NOW()),
('enable_real_name', 'false', '是否开启实名认证', NOW(), NOW()),
('require_real_name', 'false', '是否强制实名认证后才能使用服务', NOW(), NOW()),
('enable_agent', 'true', '是否开启代理商功能', NOW(), NOW());

-- Insert site_configs data
INSERT IGNORE INTO site_configs (`key`, `value`, `type`, `group`, description, created_at, updated_at) VALUES
('site_name', 'OneClickVirt', 'string', 'basic', '网站名称', NOW(), NOW()),
('site_icon_url', '/favicon.ico', 'string', 'basic', '网站图标URL', NOW(), NOW()),
('site_logo_url', '/logo.png', 'string', 'basic', '网站Logo URL', NOW(), NOW()),
('footer_text', '© 2025 OneClickVirt. All rights reserved.', 'string', 'basic', '页脚文字', NOW(), NOW()),
('icp_number', '', 'string', 'basic', 'ICP备案号', NOW(), NOW()),
('police_number', '', 'string', 'basic', '公安备案号', NOW(), NOW());

-- Insert domain_configs data
INSERT IGNORE INTO domain_configs (max_domains_per_user, max_domains_per_agent_user, default_ttl, auto_ssl, allowed_suffixes, dns_type, dns_config_path, nginx_config_path, created_at, updated_at) VALUES
(3, 5, 300, 0, '', 'dnsmasq', '/etc/dnsmasq.d/oneclickvirt-hosts.conf', '/etc/nginx/conf.d/oneclickvirt-domains', NOW(), NOW());

-- Update products table to match init.sql structure
ALTER TABLE products ADD COLUMN IF NOT EXISTS billing_cycle varchar(20) DEFAULT 'monthly';
ALTER TABLE products ADD COLUMN IF NOT EXISTS cpu_limit bigint DEFAULT cpu;
ALTER TABLE products ADD COLUMN IF NOT EXISTS memory_limit bigint DEFAULT memory;
ALTER TABLE products ADD COLUMN IF NOT EXISTS disk_limit bigint DEFAULT disk;
ALTER TABLE products ADD COLUMN IF NOT EXISTS bandwidth_limit bigint DEFAULT bandwidth;
ALTER TABLE products ADD COLUMN IF NOT EXISTS instance_limit bigint DEFAULT max_instances;
ALTER TABLE products ADD COLUMN IF NOT EXISTS status bigint DEFAULT is_enabled;

-- Insert products data
INSERT IGNORE INTO products (id, name, description, price, billing_cycle, cpu_limit, memory_limit, disk_limit, bandwidth_limit, instance_limit, features, status, sort_order, created_at, updated_at) VALUES
(1, '入门套餐', '适合个人用户的基础套餐，包含基本的虚拟化功能', 0.00, 'monthly', 1, 512, 10240, 100, 1, '{"cpu": "1核", "memory": "512MB", "disk": "10GB", "bandwidth": "100Mbps", "instances": "1个实例"}', 1, 1, NOW(), NOW()),
(2, '标准套餐', '适合小型团队的标准套餐，包含更多资源', 9.90, 'monthly', 2, 1024, 20480, 200, 3, '{"cpu": "2核", "memory": "1GB", "disk": "20GB", "bandwidth": "200Mbps", "instances": "3个实例"}', 1, 2, NOW(), NOW()),
(3, '专业套餐', '适合中型团队的专业套餐，包含完整功能', 29.90, 'monthly', 4, 2048, 40960, 500, 5, '{"cpu": "4核", "memory": "2GB", "disk": "40GB", "bandwidth": "500Mbps", "instances": "5个实例"}', 1, 3, NOW(), NOW()),
(4, '企业套餐', '适合大型团队的企业套餐，包含无限资源', 99.90, 'monthly', 8, 4096, 102400, 1000, 10, '{"cpu": "8核", "memory": "4GB", "disk": "100GB", "bandwidth": "1000Mbps", "instances": "10个实例"}', 1, 4, NOW(), NOW());
SQLEND
        echo "Default admin, user and system image data imported successfully"
    else
        echo "Users already exist, checking system image data..."
        # Check if system image data exists, import if not
        ANNOUNCEMENT_COUNT=$(mysql -h localhost -u root -e "USE oneclickvirt; SELECT COUNT(*) FROM announcements;" 2>/dev/null | tail -n 1 || echo 0)
        if [ "$ANNOUNCEMENT_COUNT" -eq 0 ]; then
            echo "Importing system image default data..."
            mysql -h localhost -u root <<SQLEND
USE oneclickvirt;

-- Create tables if they don't exist
CREATE TABLE IF NOT EXISTS announcements (
  id bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at datetime(3) NULL DEFAULT NULL,
  updated_at datetime(3) NULL DEFAULT NULL,
  deleted_at datetime(3) NULL DEFAULT NULL,
  title varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  content longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  content_html longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  type varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT 'homepage',
  priority bigint NULL DEFAULT 0,
  status bigint NULL DEFAULT 1,
  is_sticky tinyint(1) NULL DEFAULT 0,
  start_time datetime(3) NULL DEFAULT NULL,
  end_time datetime(3) NULL DEFAULT NULL,
  created_by bigint UNSIGNED NULL DEFAULT NULL,
  PRIMARY KEY (id) USING BTREE,
  INDEX idx_announcements_deleted_at(deleted_at ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

CREATE TABLE IF NOT EXISTS invite_codes (
  id bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  code varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  creator_id bigint UNSIGNED NOT NULL,
  creator_name varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  description varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  max_uses bigint NOT NULL DEFAULT 1,
  used_count bigint NOT NULL DEFAULT 0,
  expires_at datetime(3) NULL DEFAULT NULL,
  status bigint NOT NULL DEFAULT 1,
  created_at datetime(3) NULL DEFAULT NULL,
  updated_at datetime(3) NULL DEFAULT NULL,
  deleted_at datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (id) USING BTREE,
  UNIQUE INDEX idx_invite_codes_code(code ASC) USING BTREE,
  INDEX idx_invite_codes_creator_id(creator_id ASC) USING BTREE,
  INDEX idx_invite_codes_expires_at(expires_at ASC) USING BTREE,
  INDEX idx_invite_codes_status(status ASC) USING BTREE,
  INDEX idx_invite_codes_deleted_at(deleted_at ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

CREATE TABLE IF NOT EXISTS jwt_secrets (
  id bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  secret_key varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'JWT签名密钥',
  created_at datetime(3) NULL DEFAULT NULL,
  updated_at datetime(3) NULL DEFAULT NULL,
  deleted_at datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (id) USING BTREE,
  UNIQUE INDEX idx_jwt_secrets_secret_key(secret_key ASC) USING BTREE,
  INDEX idx_jwt_secrets_deleted_at(deleted_at ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- Insert data
INSERT INTO announcements (id, created_at, updated_at, deleted_at, title, content, content_html, type, priority, status, is_sticky, start_time, end_time, created_by) VALUES
(1, '2025-12-30 15:38:08.631', '2025-12-30 15:38:08.631', NULL, '欢迎使用虚拟化管理平台', '欢迎使用虚拟化管理平台，支持Docker、LXD、Incus、Proxmox VE等多种虚拟化技术。本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。', '<p>欢迎使用虚拟化管理平台，支持<strong>Docker</strong>、<strong>LXD</strong>、<strong>Incus</strong>、<strong>Proxmox VE</strong>等多种虚拟化技术。</p><p>本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。</p>', 'homepage', 10, 1, 1, NULL, NULL, NULL),
(2, '2025-12-30 15:38:08.633', '2025-12-30 15:38:08.633', NULL, '系统维护通知', '为了提供更好的服务质量，我们会定期进行系统维护。维护期间可能会影响部分功能的使用，请您谅解。', '<p>为了提供更好的服务质量，我们会定期进行系统维护。</p>', 'topbar', 5, 1, 0, NULL, NULL, NULL),
(3, '2025-12-30 15:38:08.644', '2025-12-30 15:38:08.644', NULL, '新手使用指南', '如果您是第一次使用本平台，建议先阅读使用文档。您可以在右上角的帮助菜单中找到详细的操作指南。', '<p>如果您是第一次使用本平台，建议先阅读使用文档。</p>', 'homepage', 8, 1, 0, NULL, NULL, NULL);

INSERT INTO invite_codes (id, code, creator_id, creator_name, description, max_uses, used_count, expires_at, status, created_at, updated_at, deleted_at) VALUES
(1, 'SC0Q19BW', 1, '', '', 1, 0, NULL, 1, '2025-12-31 10:59:55.167', '2025-12-31 10:59:55.167', NULL);

INSERT INTO jwt_secrets (id, secret_key, created_at, updated_at, deleted_at) VALUES
(1, 'b64dca17bf31d0e725285cccf00a6911a43b0e2c8d8d26ed458cdbf16e6a14b5', '2025-12-30 16:00:51.689', '2025-12-30 16:00:51.689', NULL);

-- Create products table if not exists
CREATE TABLE IF NOT EXISTS products (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品名称',
  description text COLLATE utf8mb4_unicode_ci COMMENT '产品描述',
  level bigint NOT NULL COMMENT '产品等级',
  price bigint NOT NULL COMMENT '价格(分)',
  period bigint NOT NULL COMMENT '周期(月), 0为永久',
  cpu bigint NOT NULL COMMENT 'CPU核心数',
  memory bigint NOT NULL COMMENT '内存(MB)',
  disk bigint NOT NULL COMMENT '磁盘(MB)',
  bandwidth bigint NOT NULL COMMENT '带宽(Mbps)',
  traffic bigint NOT NULL COMMENT '流量限制(MB)',
  max_instances bigint NOT NULL COMMENT '最大实例数',
  is_enabled bigint DEFAULT '1' COMMENT '是否启用(1:启用, 0:禁用)',
  sort_order bigint DEFAULT '0' COMMENT '排序',
  features text COLLATE utf8mb4_unicode_ci COMMENT '特性(JSON格式)',
  allow_repeat bigint DEFAULT '1' COMMENT '是否允许重复购买(1:允许, 0:不允许)',
  stock bigint DEFAULT '-1' COMMENT '库存(-1为无限)',
  sold_count bigint DEFAULT '0' COMMENT '已售数量',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create roles table if not exists
CREATE TABLE IF NOT EXISTS roles (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  name varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  code varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  description text COLLATE utf8mb4_unicode_ci,
  status bigint DEFAULT '1',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_roles_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create user_roles table if not exists
CREATE TABLE IF NOT EXISTS user_roles (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_id bigint unsigned NOT NULL,
  role_id bigint unsigned NOT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_user_roles_user_role (user_id,role_id),
  KEY idx_user_roles_role_id (role_id),
  CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create system_configs table if not exists
CREATE TABLE IF NOT EXISTS system_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text COLLATE utf8mb4_unicode_ci NOT NULL,
  description varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_system_configs_key (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create site_configs table if not exists
CREATE TABLE IF NOT EXISTS site_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `group` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  description varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY idx_site_configs_key (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create domain_configs table if not exists
CREATE TABLE IF NOT EXISTS domain_configs (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  max_domains_per_user bigint DEFAULT '3',
  max_domains_per_agent_user bigint DEFAULT '5',
  default_ttl bigint DEFAULT '300',
  auto_ssl bigint DEFAULT '0',
  allowed_suffixes text COLLATE utf8mb4_unicode_ci,
  dns_type varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT 'dnsmasq',
  dns_config_path varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '/etc/dnsmasq.d/oneclickvirt-hosts.conf',
  nginx_config_path varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '/etc/nginx/conf.d/oneclickvirt-domains',
  created_at datetime(3) DEFAULT NULL,
  updated_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert roles data
INSERT IGNORE INTO roles (name, code, description, status, created_at, updated_at) VALUES
('admin', 'admin', '系统管理员角色', 1, NOW(), NOW()),
('user', 'user', '普通用户角色', 1, NOW(), NOW());

-- Update users table to include necessary fields
ALTER TABLE users ADD COLUMN IF NOT EXISTS uuid varchar(36) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS nickname varchar(50) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone varchar(20) DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS level_expire_at datetime DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS user_type varchar(20) DEFAULT 'user';
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified tinyint(1) DEFAULT '0';
ALTER TABLE users ADD COLUMN IF NOT EXISTS real_name_verified tinyint(1) DEFAULT '0';
ALTER TABLE users ADD UNIQUE KEY IF NOT EXISTS idx_users_uuid (uuid);

-- Update existing users with missing fields
UPDATE users SET uuid = IFNULL(uuid, CONCAT('user-', id)), nickname = IFNULL(nickname, username), user_type = IFNULL(user_type, 'user') WHERE id IN (1, 2);

-- Insert user_roles data
INSERT IGNORE INTO user_roles (user_id, role_id, created_at, updated_at) VALUES
(1, 1, NOW(), NOW()),
(2, 2, NOW(), NOW());

-- Insert system_configs data
INSERT IGNORE INTO system_configs (`key`, `value`, description, created_at, updated_at) VALUES
('site_name', 'OneClickVirt', '网站名称', NOW(), NOW()),
('site_description', '虚拟化管理平台', '网站描述', NOW(), NOW()),
('site_keywords', '虚拟化,Docker,LXD,Incus,Proxmox', '网站关键词', NOW(), NOW()),
('enable_registration', 'true', '是否开启注册', NOW(), NOW()),
('enable_email_verify', 'false', '是否开启邮箱验证', NOW(), NOW()),
('default_user_level', '1', '默认用户等级', NOW(), NOW()),
('max_instances_per_user', '10', '每个用户最大实例数', NOW(), NOW()),
('default_instance_expiry_days', '30', '默认实例过期天数', NOW(), NOW()),
('enable_email_verification', 'false', '是否开启邮箱验证（注册后需验证邮箱）', NOW(), NOW()),
('email_activation_expire_hours', '24', '邮箱激活链接过期时间（小时）', NOW(), NOW()),
('enable_real_name', 'false', '是否开启实名认证', NOW(), NOW()),
('require_real_name', 'false', '是否强制实名认证后才能使用服务', NOW(), NOW()),
('enable_agent', 'true', '是否开启代理商功能', NOW(), NOW());

-- Insert site_configs data
INSERT IGNORE INTO site_configs (`key`, `value`, `type`, `group`, description, created_at, updated_at) VALUES
('site_name', 'OneClickVirt', 'string', 'basic', '网站名称', NOW(), NOW()),
('site_icon_url', '/favicon.ico', 'string', 'basic', '网站图标URL', NOW(), NOW()),
('site_logo_url', '/logo.png', 'string', 'basic', '网站Logo URL', NOW(), NOW()),
('footer_text', '© 2025 OneClickVirt. All rights reserved.', 'string', 'basic', '页脚文字', NOW(), NOW()),
('icp_number', '', 'string', 'basic', 'ICP备案号', NOW(), NOW()),
('police_number', '', 'string', 'basic', '公安备案号', NOW(), NOW());

-- Insert domain_configs data
INSERT IGNORE INTO domain_configs (max_domains_per_user, max_domains_per_agent_user, default_ttl, auto_ssl, allowed_suffixes, dns_type, dns_config_path, nginx_config_path, created_at, updated_at) VALUES
(3, 5, 300, 0, '', 'dnsmasq', '/etc/dnsmasq.d/oneclickvirt-hosts.conf', '/etc/nginx/conf.d/oneclickvirt-domains', NOW(), NOW());

-- Update products table to match init.sql structure
ALTER TABLE products ADD COLUMN IF NOT EXISTS billing_cycle varchar(20) DEFAULT 'monthly';
ALTER TABLE products ADD COLUMN IF NOT EXISTS cpu_limit bigint DEFAULT cpu;
ALTER TABLE products ADD COLUMN IF NOT EXISTS memory_limit bigint DEFAULT memory;
ALTER TABLE products ADD COLUMN IF NOT EXISTS disk_limit bigint DEFAULT disk;
ALTER TABLE products ADD COLUMN IF NOT EXISTS bandwidth_limit bigint DEFAULT bandwidth;
ALTER TABLE products ADD COLUMN IF NOT EXISTS instance_limit bigint DEFAULT max_instances;
ALTER TABLE products ADD COLUMN IF NOT EXISTS status bigint DEFAULT is_enabled;

-- Insert products data
INSERT IGNORE INTO products (id, name, description, price, billing_cycle, cpu_limit, memory_limit, disk_limit, bandwidth_limit, instance_limit, features, status, sort_order, created_at, updated_at) VALUES
(1, '入门套餐', '适合个人用户的基础套餐，包含基本的虚拟化功能', 0.00, 'monthly', 1, 512, 10240, 100, 1, '{"cpu": "1核", "memory": "512MB", "disk": "10GB", "bandwidth": "100Mbps", "instances": "1个实例"}', 1, 1, NOW(), NOW()),
(2, '标准套餐', '适合小型团队的标准套餐，包含更多资源', 9.90, 'monthly', 2, 1024, 20480, 200, 3, '{"cpu": "2核", "memory": "1GB", "disk": "20GB", "bandwidth": "200Mbps", "instances": "3个实例"}', 1, 2, NOW(), NOW()),
(3, '专业套餐', '适合中型团队的专业套餐，包含完整功能', 29.90, 'monthly', 4, 2048, 40960, 500, 5, '{"cpu": "4核", "memory": "2GB", "disk": "40GB", "bandwidth": "500Mbps", "instances": "5个实例"}', 1, 3, NOW(), NOW()),
(4, '企业套餐', '适合大型团队的企业套餐，包含无限资源', 99.90, 'monthly', 8, 4096, 102400, 1000, 10, '{"cpu": "8核", "memory": "4GB", "disk": "100GB", "bandwidth": "1000Mbps", "instances": "10个实例"}', 1, 4, NOW(), NOW());
SQLEND
            echo "System image default data imported successfully"
        else
            echo "System image data already exists, skipping import"
        fi
    fi
fi
EOF2

chmod +x /check_users.sh

# Start services first
echo "Starting services..."

# Run the user check script in the background
/bin/bash /check_users.sh &

exec supervisord -c /etc/supervisor/conf.d/supervisord.conf
