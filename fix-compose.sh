#!/bin/bash

# Fix docker-compose.yaml: remove obsolete version attribute
# and ensure proper service configuration

docker-compose down

# Add missing entrypoint for web service
cat > docker-compose.conf << 'EOF'
version: '3.8'

services:
  oneclickvirt-mysql:
    image: mariadb:10.11
    container_name: oneclickvirt-mysql
    environment:
      MYSQL_ROOT_PASSWORD: root123
      MYSQL_DATABASE: oneclickvirt
      MYSQL_USER: oneclickvirt
      MYSQL_PASSWORD: oneclickvirt123
    ports:
      - "3307:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped
    networks:
      - oneclickvirt

  oneclickvirt-redis:
    image: redis:latest
    container_name: oneclickvirt-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped
    networks:
      - oneclickvirt

  oneclickvirt-init:
    build:
      context: .
      dockerfile: Dockerfile
      target: base
    image: oneclickvirt:fixed
    container_name: oneclickvirt-init
    depends_on:
      - oneclickvirt-mysql
    environment:
      DB_HOST: oneclickvirt-mysql
      DB_PORT: 3306
      MYSQL_ROOT_PASSWORD: root123
      MYSQL_DATABASE: oneclickvirt
    command: /autoinstall.sh
    restart: unless-stopped
    networks:
      - oneclickvirt

  oneclickvirt-api:
    build:
      context: .
      dockerfile: Dockerfile
    image: oneclickvirt:fixed
    container_name: oneclickvirt-api
    depends_on:
      - oneclickvirt-mysql
      - oneclickvirt-redis
    environment:
      DB_HOST: oneclickvirt-mysql
      DB_PORT: 3306
      REDIS_HOST: oneclickvirt-redis
    command: /app/main
    restart: unless-stopped
    ports:
      - "8890:8890"
    networks:
      - oneclickvirt

  oneclickvirt-web:
    build:
      context: .
      dockerfile: Dockerfile
      target: web
    image: oneclickvirt:fixed
    container_name: oneclickvirt-web
    depends_on:
      - oneclickvirt-api
    command: nginx -g "daemon off;"
    volumes:
      - nginx_logs:/var/log/nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    networks:
      - oneclickvirt

networks:
  oneclickvirt:
    driver: bridge

volumes:
  mysql_data:
  redis_data:
  nginx_logs:
EOF

docker-compose -f docker-compose.conf up -d --build
