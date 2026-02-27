#!/bin/bash

# 在远程服务器上部署 one-click virtualization 项目
echo "开始部署 one-click virtualization 项目..."

# 创建项目目录
mkdir -p /opt/oneclickvirt
cd /opt/oneclickvirt

# 从GitHub克隆项目
git clone https://github.com/luren2024/oneclickvirt.git temp_project

# 复制文件
cp -r temp_project/server ./
cp -r temp_project/web ./
cp temp_project/docker-compose.yaml ./

# 修改配置文件使用新端口
sed -i 's/addr: 8888/addr: 8890/g' server/config.yaml
sed -i 's/VITE_SERVER_PORT = 8888/VITE_SERVER_PORT = 8890/g' web/.env.development

# 修改docker-compose文件端口映射
sed -i 's/- "8889:8888"/- "8890:8888"/g' docker-compose.yaml

# 删除临时目录
rm -rf temp_project

# 构建Docker镜像并启动服务
echo "构建并启动服务..."
docker-compose up -d

echo "部署完成!"
echo "前端访问地址: http://154.12.84.134:8099"
echo "后端API地址: http://154.12.84.134:8890"
echo "API文档地址: http://154.12.84.134:8890/swagger/index.html"

# 检查服务状态
docker-compose ps