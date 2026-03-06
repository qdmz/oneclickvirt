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
        # Generate self-signed certificate if not provided or empty
        if [ ! -f "/etc/nginx/ssl/cert.pem" ] || [ ! -f "/etc/nginx/ssl/key.pem" ] || [ ! -s "/etc/nginx/ssl/cert.pem" ] || [ ! -s "/etc/nginx/ssl/key.pem" ]; then
            echo "Generating self-signed SSL certificate for $DOMAIN..."
            openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/nginx/ssl/key.pem -out /etc/nginx/ssl/cert.pem -subj "/C=US/ST=State/L=City/O=Organization/CN=$DOMAIN"
            chmod 644 /etc/nginx/ssl/cert.pem
            chmod 600 /etc/nginx/ssl/key.pem
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
DROP USER IF EXISTS 'root'@'127.0.0.1';
DROP USER IF EXISTS 'root'@'%';
CREATE USER 'root'@'127.0.0.1' IDENTIFIED WITH mysql_native_password BY '';
CREATE USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '';
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
cat > /check_users.sh << 'EOF'
#!/bin/bash

# Wait for MySQL to start
echo "Checking if users exist..."
for i in {1..60}; do
    if mysql -h localhost -u root -e "SELECT 1" >/dev/null 2>&1; then
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
TABLE_EXISTS=$(mysql -h localhost -u root -e "USE oneclickvirt; SHOW TABLES LIKE 'users';" 2>/dev/null | tail -n 1 || echo "")
if [ -z "$TABLE_EXISTS" ]; then
    echo "Users table does not exist, skipping import"
else
    # Table exists, check if it has data
    USER_COUNT=$(mysql -h localhost -u root -e "USE oneclickvirt; SELECT COUNT(*) FROM users;" 2>/dev/null | tail -n 1 || echo 0)
    if [ "$USER_COUNT" -eq 0 ]; then
        echo "Users table exists but is empty, importing default admin and user data..."
        # Generate UUIDs for users
        ADMIN_UUID=$(cat /proc/sys/kernel/random/uuid)
        USER_UUID=$(cat /proc/sys/kernel/random/uuid)
        
        # Use simple password 'password' and let the system hash it later
        # For now, we'll use a placeholder and the system will handle the hashing
        mysql -h localhost -u root -e "USE oneclickvirt; INSERT INTO users (id, uuid, username, password, email, level, user_type, status, created_at, updated_at, max_instances, max_cpu, max_memory, max_disk) VALUES (1, '${ADMIN_UUID}', 'admin', 'password', 'admin@example.com', 5, 'admin', 1, NOW(), NOW(), 10, 8, 8192, 102400), (2, '${USER_UUID}', 'user', 'password', 'user@example.com', 1, 'user', 1, NOW(), NOW(), 1, 1, 512, 10240);"
        
        # Import system image default data
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
INSERT IGNORE INTO announcements (id, created_at, updated_at, deleted_at, title, content, content_html, type, priority, status, is_sticky, start_time, end_time, created_by) VALUES
(1, '2025-12-30 15:38:08.631', '2025-12-30 15:38:08.631', NULL, '欢迎使用虚拟化管理平台', '欢迎使用虚拟化管理平台，支持Docker、LXD、Incus、Proxmox VE等多种虚拟化技术。本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。', '<p>欢迎使用虚拟化管理平台，支持<strong>Docker</strong>、<strong>LXD</strong>、<strong>Incus</strong>、<strong>Proxmox VE</strong>等多种虚拟化技术。</p><p>本平台提供简单易用的Web界面，让您轻松管理各种虚拟化资源。</p>', 'homepage', 10, 1, 1, NULL, NULL, NULL),
(2, '2025-12-30 15:38:08.633', '2025-12-30 15:38:08.633', NULL, '系统维护通知', '为了提供更好的服务质量，我们会定期进行系统维护。维护期间可能会影响部分功能的使用，请您谅解。', '<p>为了提供更好的服务质量，我们会定期进行系统维护。</p>', 'topbar', 5, 1, 0, NULL, NULL, NULL),
(3, '2025-12-30 15:38:08.644', '2025-12-30 15:38:08.644', NULL, '新手使用指南', '如果您是第一次使用本平台，建议先阅读使用文档。您可以在右上角的帮助菜单中找到详细的操作指南。', '<p>如果您是第一次使用本平台，建议先阅读使用文档。</p>', 'homepage', 8, 1, 0, NULL, NULL, NULL);

INSERT IGNORE INTO invite_codes (id, code, creator_id, creator_name, description, max_uses, used_count, expires_at, status, created_at, updated_at, deleted_at) VALUES
(1, 'SC0Q19BW', 1, '', '', 1, 0, NULL, 1, '2025-12-31 10:59:55.167', '2025-12-31 10:59:55.167', NULL);

INSERT IGNORE INTO jwt_secrets (id, secret_key, created_at, updated_at, deleted_at) VALUES
(1, 'b64dca17bf31d0e725285cccf00a6911a43b0e2c8d8d26ed458cdbf16e6a14b5', '2025-12-30 16:00:51.689', '2025-12-30 16:00:51.689', NULL);
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
SQLEND
            echo "System image default data imported successfully"
        else
            echo "System image data already exists, skipping import"
        fi
    fi
fi
EOF

chmod +x /check_users.sh

# Start services first
echo "Starting services..."

# Run the user check script in the background
/bin/bash /check_users.sh &

exec supervisord -c /etc/supervisor/conf.d/supervisord.conf