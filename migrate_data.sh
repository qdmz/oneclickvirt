#!/bin/bash

# 数据迁移脚本：从v1.2迁移到v1.3
# 该脚本将旧数据目录中的MySQL数据复制到新的数据结构中

# 定义变量
OLD_DATA_DIR="/root/oneclickvirt-new1/data/mysql"
NEW_DATA_DIR="/root/oneclickvirt-new/data/mysql"
CONTAINER_NAME="oneclickvirt-new"

# 检查旧数据目录是否存在
if [ ! -d "$OLD_DATA_DIR" ]; then
    echo "错误：旧数据目录 $OLD_DATA_DIR 不存在！"
    exit 1
fi

# 停止当前运行的容器
if docker ps -a | grep -q "$CONTAINER_NAME"; then
    echo "停止并移除现有容器..."
    docker stop "$CONTAINER_NAME"
    docker rm "$CONTAINER_NAME"
fi

# 备份新数据目录（如果存在）
if [ -d "$NEW_DATA_DIR" ]; then
    echo "备份新数据目录..."
    BACKUP_DIR="$NEW_DATA_DIR.bak_$(date +%Y%m%d_%H%M%S)"
    cp -r "$NEW_DATA_DIR" "$BACKUP_DIR"
    echo "新数据已备份到 $BACKUP_DIR"
    
    # 清空新数据目录
    echo "清空新数据目录..."
    rm -rf "$NEW_DATA_DIR"/*
else
    # 创建新数据目录
    echo "创建新数据目录..."
    mkdir -p "$NEW_DATA_DIR"
fi

# 复制旧数据到新数据目录
echo "复制旧数据到新数据目录..."
cp -r "$OLD_DATA_DIR"/* "$NEW_DATA_DIR/"

# 设置正确的权限
echo "设置数据目录权限..."
chown -R nobody:nogroup "$NEW_DATA_DIR"

# 重新构建并启动容器
echo "重新构建并启动容器..."
cd /root/oneclickvirt-new && docker build -t oneclickvirt . && docker run -d --name oneclickvirt-new -p 80:80 -p 443:443 -v /root/oneclickvirt-new/data/app:/app/storage -v /root/oneclickvirt-new/data/mysql:/var/lib/mysql -v /root/oneclickvirt-new/data/config:/app/config -e FRONTEND_URL=https://heyun.ypvps.com oneclickvirt

# 等待容器启动
echo "等待容器启动..."
sleep 10

# 检查容器状态
echo "检查容器状态..."
docker ps | grep "$CONTAINER_NAME"

# 提示迁移完成
echo "\n数据迁移完成！"
echo "旧数据已成功迁移到新的数据结构中。"
echo "您可以登录系统检查数据是否完整。"
