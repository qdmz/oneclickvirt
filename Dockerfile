# OneClickVirt All-in-One Container

FROM node:22-slim AS frontend-builder
ARG TARGETARCH
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci --include=optional
RUN if [ "$TARGETARCH" = "amd64" ]; then \
        npm install --no-save @rollup/rollup-linux-x64-gnu; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
        npm install --no-save @rollup/rollup-linux-arm64-gnu; \
    fi
COPY web/ ./
RUN npm run build


FROM golang:1.24-alpine AS backend-builder
ARG TARGETARCH
WORKDIR /app/server
RUN apk add --no-cache git ca-certificates
COPY server/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -installsuffix cgo -ldflags "-w -s" -o main .

FROM debian:12-slim
ARG TARGETARCH

# Install database and other services based on architecture
RUN apt-get update --allow-insecure-repositories && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --allow-unauthenticated \
        gnupg2 wget lsb-release procps nginx supervisor ca-certificates openssl python3 python3-venv python3-pip && \
    if [ "$TARGETARCH" = "amd64" ]; then \
        echo "Installing MySQL for AMD64..." && \
        gpg --keyserver keyserver.ubuntu.com --recv-keys B7B3B788A8D3785C && \
        gpg --export B7B3B788A8D3785C > /usr/share/keyrings/mysql.gpg && \
        echo "deb [signed-by=/usr/share/keyrings/mysql.gpg] http://repo.mysql.com/apt/debian bookworm mysql-8.0" > /etc/apt/sources.list.d/mysql.list && \
        apt-get update --allow-insecure-repositories && \
        DEBIAN_FRONTEND=noninteractive apt-get install -y --allow-unauthenticated mysql-server mysql-client; \
    else \
        echo "Installing MariaDB for ARM64..." && \
        DEBIAN_FRONTEND=noninteractive apt-get install -y --allow-unauthenticated mariadb-server mariadb-client; \
    fi && \
    # Install Certbot for SSL certificate management
    apt-get update --allow-insecure-repositories && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --allow-unauthenticated certbot python3-certbot-nginx && \
    apt-get clean

ENV TZ=Asia/Shanghai
WORKDIR /app
RUN mkdir -p /var/lib/mysql /var/log/mysql /var/run/mysqld /var/log/supervisor \
    && mkdir -p /app/storage/{cache,certs,configs,exports,logs,temp,uploads} \
    && mkdir -p /etc/mysql/conf.d \
    && mkdir -p /etc/nginx/ssl

# Copy SSL certificate files if provided (empty files for now)
RUN touch /etc/nginx/ssl/cert.pem \
    && touch /etc/nginx/ssl/key.pem

COPY --from=backend-builder /app/server/main ./main
COPY --from=backend-builder /app/server/config.yaml ./config.yaml.default
RUN if [ ! -f /app/config.yaml ]; then mv /app/config.yaml.default /app/config.yaml; else rm /app/config.yaml.default; fi
# Create config directory for persistence
RUN mkdir -p /app/config
# Use config from persistence directory if exists
RUN if [ -f /app/config/config.yaml ]; then cp /app/config/config.yaml /app/config.yaml; fi
COPY --from=frontend-builder /app/web/dist /var/www/html

# Copy initialization scripts
COPY scripts/init.sql /app/scripts/init.sql
COPY scripts/init.sh /app/scripts/init.sh
COPY complete_init.sql /app/complete_init.sql
RUN chmod +x /app/scripts/init.sh

RUN mkdir -p /var/run/mysqld && \
    chown -R mysql:mysql /var/lib/mysql /var/log/mysql /var/run/mysqld && \
    chown -R www-data:www-data /var/www/html && \
    chmod -R 755 /var/www/html && \
    chmod 755 /app/main && \
    chmod 666 /app/config.yaml && \
    chmod 750 /app/storage && \
    chmod -R 750 /app/storage/*

# Create database configuration based on architecture
RUN if [ "$TARGETARCH" = "amd64" ]; then \
        echo '[mysqld]' > /etc/mysql/conf.d/custom.cnf && \
        echo 'datadir=/var/lib/mysql' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'socket=/var/run/mysqld/mysqld.sock' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'user=mysql' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'pid-file=/var/run/mysqld/mysqld.pid' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'bind-address=0.0.0.0' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'port=3306' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'character-set-server=utf8mb4' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'collation-server=utf8mb4_unicode_ci' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'authentication_policy=mysql_native_password' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'max_connections=200' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'skip-name-resolve' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'secure-file-priv=""' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'innodb_buffer_pool_size=256M' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'innodb_redo_log_capacity=67108864' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'innodb_force_recovery=0' >> /etc/mysql/conf.d/custom.cnf; \
    else \
        echo '[mysqld]' > /etc/mysql/conf.d/custom.cnf && \
        echo 'datadir=/var/lib/mysql' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'socket=/var/run/mysqld/mysqld.sock' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'user=mysql' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'pid-file=/var/run/mysqld/mysqld.pid' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'bind-address=0.0.0.0' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'port=3306' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'character-set-server=utf8mb4' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'collation-server=utf8mb4_unicode_ci' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'max_connections=200' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'skip-name-resolve' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'secure-file-priv=""' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'innodb_buffer_pool_size=256M' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'innodb_log_file_size=64M' >> /etc/mysql/conf.d/custom.cnf && \
        echo 'innodb_force_recovery=0' >> /etc/mysql/conf.d/custom.cnf; \
    fi

RUN echo 'user www-data;' > /etc/nginx/nginx.conf && \
    echo 'worker_processes auto;' >> /etc/nginx/nginx.conf && \
    echo 'error_log /var/log/nginx/error.log;' >> /etc/nginx/nginx.conf && \
    echo 'pid /run/nginx.pid;' >> /etc/nginx/nginx.conf && \
    echo 'events { worker_connections 1024; }' >> /etc/nginx/nginx.conf && \
    echo 'http {' >> /etc/nginx/nginx.conf && \
    echo '    include /etc/nginx/mime.types;' >> /etc/nginx/nginx.conf && \
    echo '    default_type application/octet-stream;' >> /etc/nginx/nginx.conf && \
    echo '    sendfile on;' >> /etc/nginx/nginx.conf && \
    echo '    keepalive_timeout 65;' >> /etc/nginx/nginx.conf && \
    echo '    gzip on;' >> /etc/nginx/nginx.conf && \
    echo '    ' >> /etc/nginx/nginx.conf && \
    echo '    # WebSocket upgrade mapping' >> /etc/nginx/nginx.conf && \
    echo '    map $http_upgrade $connection_upgrade {' >> /etc/nginx/nginx.conf && \
    echo '        default upgrade;' >> /etc/nginx/nginx.conf && \
    echo '        "" close;' >> /etc/nginx/nginx.conf && \
    echo '    }' >> /etc/nginx/nginx.conf && \
    echo '    ' >> /etc/nginx/nginx.conf && \
    echo '    # HTTP server' >> /etc/nginx/nginx.conf && \
    echo '    server {' >> /etc/nginx/nginx.conf && \
    echo '        listen 80;' >> /etc/nginx/nginx.conf && \
    echo '        server_name localhost;' >> /etc/nginx/nginx.conf && \
    echo '        root /var/www/html;' >> /etc/nginx/nginx.conf && \
    echo '        index index.html;' >> /etc/nginx/nginx.conf && \
    echo '        client_max_body_size 10M;' >> /etc/nginx/nginx.conf && \
    echo '        ' >> /etc/nginx/nginx.conf && \
    echo '        location /api/ {' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_pass http://127.0.0.1:8890;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_set_header Host $host;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_set_header X-Real-IP $remote_addr;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_set_header REMOTE-HOST $remote_addr;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_set_header X-Forwarded-Proto $scheme;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_set_header X-Forwarded-Port $server_port;' >> /etc/nginx/nginx.conf && \
	    echo '            ' >> /etc/nginx/nginx.conf && \
	    echo '            # WebSocket support' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_set_header Upgrade $http_upgrade;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_set_header Connection $connection_upgrade;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_http_version 1.1;' >> /etc/nginx/nginx.conf && \
	    echo '            ' >> /etc/nginx/nginx.conf && \
	    echo '            # SSL settings' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_ssl_server_name off;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_ssl_name $proxy_host;' >> /etc/nginx/nginx.conf && \
	    echo '            ' >> /etc/nginx/nginx.conf && \
	    echo '            # Timeout settings for SSH connections' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_connect_timeout 60s;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_send_timeout 600s;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_read_timeout 600s;' >> /etc/nginx/nginx.conf && \
	    echo '            ' >> /etc/nginx/nginx.conf && \
	    echo '            # Disable buffering for real-time data' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_buffering off;' >> /etc/nginx/nginx.conf && \
	    echo '            add_header X-Cache $upstream_cache_status;' >> /etc/nginx/nginx.conf && \
	    echo '            add_header Cache-Control no-cache;' >> /etc/nginx/nginx.conf && \
	    echo '        }' >> /etc/nginx/nginx.conf && \
	    echo '        ' >> /etc/nginx/nginx.conf && \
	    echo '        location /swagger/ {' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_pass http://127.0.0.1:8890;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_set_header Host $host;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_set_header X-Real-IP $remote_addr;' >> /etc/nginx/nginx.conf && \
	    echo '            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;' >> /etc/nginx/nginx.conf && \
	    echo '        }' >> /etc/nginx/nginx.conf && \
    echo '        ' >> /etc/nginx/nginx.conf && \
    echo '        # WebSocket endpoints for SSH connections' >> /etc/nginx/nginx.conf && \
    echo '        location /v1/ {' >> /etc/nginx/nginx.conf && \
    echo '            proxy_pass http://127.0.0.1:8890;' >> /etc/nginx/nginx.conf && \
    echo '            proxy_set_header Host $host;' >> /etc/nginx/nginx.conf && \
    echo '            proxy_set_header X-Real-IP $remote_addr;' >> /etc/nginx/nginx.conf && \
    echo '            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;' >> /etc/nginx/nginx.conf && \
    echo '            proxy_set_header REMOTE-HOST $remote_addr;' >> /etc/nginx/nginx.conf && \
    echo '            proxy_set_header X-Forwarded-Proto $scheme;' >> /etc/nginx/nginx.conf && \
    echo '            proxy_set_header X-Forwarded-Port $server_port;' >> /etc/nginx/nginx.conf && \
    echo '            ' >> /etc/nginx/nginx.conf && \
    echo '            # WebSocket support' >> /etc/nginx/nginx.conf && \
    echo '            proxy_set_header Upgrade $http_upgrade;' >> /etc/nginx/nginx.conf && \
    echo '            proxy_set_header Connection $connection_upgrade;' >> /etc/nginx/nginx.conf && \
    echo '            proxy_http_version 1.1;' >> /etc/nginx/nginx.conf && \
    echo '            ' >> /etc/nginx/nginx.conf && \
    echo '            # Timeout settings for SSH connections' >> /etc/nginx/nginx.conf && \
    echo '            proxy_connect_timeout 60s;' >> /etc/nginx/nginx.conf && \
    echo '            proxy_send_timeout 600s;' >> /etc/nginx/nginx.conf && \
    echo '            proxy_read_timeout 600s;' >> /etc/nginx/nginx.conf && \
    echo '            ' >> /etc/nginx/nginx.conf && \
    echo '            # Disable buffering for real-time data' >> /etc/nginx/nginx.conf && \
    echo '            proxy_buffering off;' >> /etc/nginx/nginx.conf && \
    echo '            add_header X-Cache $upstream_cache_status;' >> /etc/nginx/nginx.conf && \
    echo '            add_header Cache-Control no-cache;' >> /etc/nginx/nginx.conf && \
    echo '        }' >> /etc/nginx/nginx.conf && \
    echo '        ' >> /etc/nginx/nginx.conf && \
    echo '        location / {' >> /etc/nginx/nginx.conf && \
    echo '            try_files $uri $uri/ /index.html;' >> /etc/nginx/nginx.conf && \
    echo '        }' >> /etc/nginx/nginx.conf && \
    echo '    }' >> /etc/nginx/nginx.conf && \
    echo '    ' >> /etc/nginx/nginx.conf && \
    # echo '    # HTTPS server (SSL)' >> /etc/nginx/nginx.conf && \
    # echo '    server {' >> /etc/nginx/nginx.conf && \
    # echo '        listen 443 ssl;' >> /etc/nginx/nginx.conf && \
    # echo '        server_name localhost;' >> /etc/nginx/nginx.conf && \
    # echo '        ssl_certificate /etc/nginx/ssl/cert.pem;' >> /etc/nginx/nginx.conf && \
    # echo '        ssl_certificate_key /etc/nginx/ssl/key.pem;' >> /etc/nginx/nginx.conf && \
    # echo '        ssl_protocols TLSv1.2 TLSv1.3;' >> /etc/nginx/nginx.conf && \
    # echo '        ssl_ciphers HIGH:!aNULL:!MD5;' >> /etc/nginx/nginx.conf && \
    # echo '        root /var/www/html;' >> /etc/nginx/nginx.conf && \
    # echo '        index index.html;' >> /etc/nginx/nginx.conf && \
    # echo '        client_max_body_size 10M;' >> /etc/nginx/nginx.conf && \
    # echo '        ' >> /etc/nginx/nginx.conf && \
    # echo '        location /api/ {' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_pass http://127.0.0.1:8890;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header Host $host;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header X-Real-IP $remote_addr;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header REMOTE-HOST $remote_addr;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header X-Forwarded-Proto https;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header X-Forwarded-Port 443;' >> /etc/nginx/nginx.conf && \
    # echo '            ' >> /etc/nginx/nginx.conf && \
    # echo '            # WebSocket support' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header Upgrade $http_upgrade;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header Connection $connection_upgrade;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_http_version 1.1;' >> /etc/nginx/nginx.conf && \
    # echo '            ' >> /etc/nginx/nginx.conf && \
    # echo '            # Timeout settings for SSH connections' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_connect_timeout 60s;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_send_timeout 600s;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_read_timeout 600s;' >> /etc/nginx/nginx.conf && \
    # echo '            ' >> /etc/nginx/nginx.conf && \
    # echo '            # Disable buffering for real-time data' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_buffering off;' >> /etc/nginx/nginx.conf && \
    # echo '            add_header X-Cache $upstream_cache_status;' >> /etc/nginx/nginx.conf && \
    # echo '            add_header Cache-Control no-cache;' >> /etc/nginx/nginx.conf && \
    # echo '        }' >> /etc/nginx/nginx.conf && \
    # echo '        ' >> /etc/nginx/nginx.conf && \
    # echo '        location /swagger/ {' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_pass http://127.0.0.1:8890;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header Host $host;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header X-Real-IP $remote_addr;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;' >> /etc/nginx/nginx.conf && \
    # echo '        }' >> /etc/nginx/nginx.conf && \
    # echo '        ' >> /etc/nginx/nginx.conf && \
    # echo '        # WebSocket endpoints for SSH connections' >> /etc/nginx/nginx.conf && \
    # echo '        location /v1/ {' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_pass http://127.0.0.1:8890;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header Host $host;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header X-Real-IP $remote_addr;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header REMOTE-HOST $remote_addr;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header X-Forwarded-Proto https;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header X-Forwarded-Port 443;' >> /etc/nginx/nginx.conf && \
    # echo '            ' >> /etc/nginx/nginx.conf && \
    # echo '            # WebSocket support' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header Upgrade $http_upgrade;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_set_header Connection $connection_upgrade;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_http_version 1.1;' >> /etc/nginx/nginx.conf && \
    # echo '            ' >> /etc/nginx/nginx.conf && \
    # echo '            # Timeout settings for SSH connections' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_connect_timeout 60s;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_send_timeout 600s;' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_read_timeout 600s;' >> /etc/nginx/nginx.conf && \
    # echo '            ' >> /etc/nginx/nginx.conf && \
    # echo '            # Disable buffering for real-time data' >> /etc/nginx/nginx.conf && \
    # echo '            proxy_buffering off;' >> /etc/nginx/nginx.conf && \
    # echo '            add_header X-Cache $upstream_cache_status;' >> /etc/nginx/nginx.conf && \
    # echo '            add_header Cache-Control no-cache;' >> /etc/nginx/nginx.conf && \
    # echo '        }' >> /etc/nginx/nginx.conf && \
    # echo '        ' >> /etc/nginx/nginx.conf && \
    # echo '        location / {' >> /etc/nginx/nginx.conf && \
    # echo '            try_files $uri $uri/ /index.html;' >> /etc/nginx/nginx.conf && \
    # echo '        }' >> /etc/nginx/nginx.conf && \
    # echo '    }' >> /etc/nginx/nginx.conf && \
    echo '}' >> /etc/nginx/nginx.conf

# Create base supervisor directory
RUN mkdir -p /etc/supervisor/conf.d

# Copy autoinstall.sh script
COPY autoinstall.sh /autoinstall.sh
RUN chmod +x /autoinstall.sh

# Expose ports
EXPOSE 80 443

# Start the application using autoinstall.sh
CMD ["/autoinstall.sh"]