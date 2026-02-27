#!/bin/bash

# 部署 one-click virtualization 项目到远程服务器
# 使用方法: 在远程服务器上运行此脚本

echo "开始部署 one-click virtualization 项目..."

# 创建项目目录
mkdir -p /opt/oneclickvirt
cd /opt/oneclickvirt

# 检查并停止已存在的容器
echo "检查并停止已存在的容器..."
docker stop oneclickvirt-api oneclickvirt-web oneclickvirt-mysql 2>/dev/null || true
docker rm oneclickvirt-api oneclickvirt-web oneclickvirt-mysql 2>/dev/null || true

# 如果之前有打包好的文件，可以解压（这里我们假设文件已经通过其他方式传输）
# tar -xzf oneclickvirt.tar.gz

# 或者从GitHub克隆最新代码
if [ ! -d "server" ] || [ ! -d "web" ]; then
    echo "克隆项目代码..."
    git clone https://github.com/luren2024/oneclickvirt.git tmp_project
    cp -r tmp_project/* .
    rm -rf tmp_project
fi

# 修改配置文件使用新端口
echo "修改配置文件..."
sed -i 's/addr: 8888/addr: 8890/g' server/config.yaml
sed -i 's/VITE_SERVER_PORT = 8888/VITE_SERVER_PORT = 8890/g' web/.env.development

# 构建后端镜像
echo "构建后端镜像..."
cd server
docker build -t oneclickvirt-server .
cd ..

# 构建前端镜像
echo "构建前端镜像..."
cd web
docker build -t oneclickvirt-web .
cd ..

# 启动服务
echo "启动服务..."
docker-compose -f docker-compose.yaml up -d

echo "部署完成!"
echo "前端访问地址: http://154.12.84.134:8080"
echo "后端API地址: http://154.12.84.134:8890"
echo "API文档地址: http://154.12.84.134:8890/swagger/index.html"