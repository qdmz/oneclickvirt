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

# Update nginx server_name if FRONTEND_URL is provided
if [ ! -z "$FRONTEND_URL" ]; then
    echo "Configuring frontend-url: $FRONTEND_URL"
    
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

echo "Setting initial permissions for MySQL directories..."
chown -R mysql:mysql /var/lib/mysql /var/run/mysqld /var/log/mysql
chmod -R 755 /var/lib/mysql
chmod 755 /var/run/mysqld
chmod 755 /var/log/mysql

echo "Permissions set successfully:" 
ls -la /var/lib/mysql
ls -la /var/run/mysqld
ls -la /var/log/mysql

# Check if database needs initialization
INIT_NEEDED=false
# Create database initialization flag file path (different from business init)
DB_INIT_FLAG="/app/.mysql_initialized"

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
    else
        echo "Database already initialized (flag exists), skipping initialization..."
    fi
fi

# 强制创建初始化标记，避免循环初始化
if [ "$INIT_NEEDED" = "true" ]; then
    echo "Creating initialization flag to prevent loop..."
    echo "$(date): Database initialization in progress" > "$DB_INIT_FLAG"
    chmod 644 "$DB_INIT_FLAG"
fi

# Always check and import default users if needed
CHECK_USERS=true

if [ "$INIT_NEEDED" = "true" ]; then
    # Stop any running database processes
    echo "Stopping any running database processes..."
    pkill -f "$DB_DAEMON" || true
    sleep 2
    
    # Remove old/corrupted data only when needed
    echo "Removing old database data..."
    rm -rf /var/lib/mysql/*
    
    # Set correct permissions
    echo "Setting permissions for MySQL directories..."
    chown -R mysql:mysql /var/lib/mysql /var/run/mysqld /var/log/mysql
    chmod -R 755 /var/lib/mysql
    chmod 755 /var/run/mysqld
    chmod 755 /var/log/mysql
    
    # Initialize database
    echo "Initializing MySQL database..."
    
    # 清理旧数据
    rm -rf /var/lib/mysql/*
    
    # 设置权限
    chown -R mysql:mysql /var/lib/mysql
    chmod -R 755 /var/lib/mysql
    
    # 初始化MySQL，使用--initialize-insecure
    echo "Running MySQL initialization..."
    # 捕获MySQL初始化的详细输出
    mysqld --initialize-insecure --user=mysql --datadir=/var/lib/mysql --verbose 2>&1
    
    if [ $? -ne 0 ]; then
        echo "MySQL initialization failed"
        echo "Error log:"
        cat /var/log/mysql/error.log || true
        # 不要立即退出，继续尝试
        echo "Continuing despite initialization error..."
    else
        echo "MySQL initialization completed successfully"
    fi
    
    # 检查初始化结果
    echo "Initialization result:"
    ls -la /var/lib/mysql/
    
    # 启动MySQL服务器
    echo "Starting MySQL server..."
    mysqld --user=mysql --pid-file=/var/run/mysqld/mysqld.pid --skip-networking &
    mysql_pid=$!
    
    # 等待MySQL启动
    echo "Waiting for MySQL to start..."
    for i in {1..60}; do
        if mysql -u root --skip-password -e "SELECT 1" >/dev/null 2>&1; then
            echo "MySQL started successfully"
            break
        fi
        echo "Waiting for MySQL... ($i/60)"
        sleep 1
    done
    
    # Check if MySQL is running
    if ! ps -p $mysql_pid > /dev/null; then
        echo "MySQL failed to start"
        cat /var/log/mysql/error.log || true
        exit 1
    fi
    
    # Create database
    echo "Creating database..."
    mysql -u root --skip-password -e "CREATE DATABASE IF NOT EXISTS ${MYSQL_DATABASE} CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
    
    if [ $? -ne 0 ]; then
        echo "Failed to create database"
        kill $mysql_pid 2>/dev/null || true
        exit 1
    fi
    
    # Set MySQL root user permissions to allow connections from all hosts
    echo "Setting MySQL root user permissions..."
    mysql -u root --skip-password -e "CREATE USER IF NOT EXISTS 'root'@'127.0.0.1' IDENTIFIED BY '';"
    mysql -u root --skip-password -e "CREATE USER IF NOT EXISTS 'root'@'localhost' IDENTIFIED BY '';"
    mysql -u root --skip-password -e "CREATE USER IF NOT EXISTS 'root'@'mysql' IDENTIFIED BY '';"
    mysql -u root --skip-password -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'127.0.0.1' WITH GRANT OPTION;"
    mysql -u root --skip-password -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost' WITH GRANT OPTION;"
    mysql -u root --skip-password -e "GRANT ALL PRIVILEGES ON *.* TO 'root'@'mysql' WITH GRANT OPTION;"
    mysql -u root --skip-password -e "FLUSH PRIVILEGES;"
    
    # Import default data
    echo "Importing default data..."
    mysql -u root --skip-password -D ${MYSQL_DATABASE} < /app/complete_init.sql
    
    if [ $? -ne 0 ]; then
        echo "Failed to import default data"
        kill $mysql_pid 2>/dev/null || true
        exit 1
    fi
    
    echo "Default data imported successfully"
    
    # Stop MySQL server
    echo "Stopping MySQL server..."
    kill $mysql_pid
    wait $mysql_pid 2>/dev/null || true
    
    echo "MySQL configuration completed."
    
    # Create database initialization flag
    echo "$(date): Database initialized successfully" > "$DB_INIT_FLAG"
    echo "Created initialization flag at $DB_INIT_FLAG"
    
    # 确保标记文件权限正确
    chmod 644 "$DB_INIT_FLAG"
else
    echo "Database already initialized, skipping..."
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
environment=DB_HOST="127.0.0.1",DB_PORT="3306",DB_USER="root",DB_PASSWORD="",DB_NAME="${MYSQL_DATABASE}"
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

# Check if users table exists and has data
TABLE_EXISTS=$(mysql -h "$DB_HOST" -u root -e "USE ${MYSQL_DATABASE}; SHOW TABLES LIKE 'users';" 2>/dev/null | tail -n 1 || echo "")
if [ -n "$TABLE_EXISTS" ]; then
    USER_COUNT=$(mysql -h "$DB_HOST" -u root -e "USE ${MYSQL_DATABASE}; SELECT COUNT(*) FROM users;" 2>/dev/null | tail -n 1 || echo 0)
    if [ "$USER_COUNT" -eq 0 ]; then
        echo "Users table exists but is empty, importing default data from complete_init.sql..."
        mysql -h "$DB_HOST" -u root < /app/complete_init.sql
        echo "Default data imported successfully"
    fi
fi
EOF2

chmod +x /check_users.sh

# Start the check_users.sh script in the background
nohup /check_users.sh > /var/log/check_users.log 2>&1 &

# Start supervisor
echo "Starting supervisor..."
exec /usr/bin/supervisord -n -c /etc/supervisor/conf.d/supervisord.conf