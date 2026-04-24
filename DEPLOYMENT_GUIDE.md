# OneClickVirt 修复部署指南

## 当前状态
✅ 代码修复完成
✅ Git提交完成
⏳ 需要重新构建和部署

## 部署方法

### 方法1：重新构建Docker镜像（推荐）

#### 1. 准备修复后的代码
```bash
cd /root/.openclaw/workspace/oneclickvirt

# 确认修复已应用
git diff HEAD~1 server/service/resources/quota.go
```

#### 2. 构建Docker镜像
```bash
# 使用修复后的代码构建
docker build -t oneclickvirt:fixed .

# 或者指定Dockerfile
docker build -f Dockerfile -t oneclickvirt:fixed .
```

#### 3. 停止旧容器
```bash
# 停止运行中的容器
docker stop oneclickvirt

# 删除旧容器
docker rm oneclickvirt
```

#### 4. 启动新容器
```bash
# 启动修复后的容器
docker run -d --name oneclickvirt \
  -p 8080:80 \
  -p 8443:443 \
  -v /path/to/data:/app/storage \
  oneclickvirt:fixed

# 或者使用docker-compose
docker-compose up -d
```

#### 5. 验证部署
```bash
# 检查容器状态
docker ps | grep oneclickvirt

# 查看容器日志
docker logs oneclickvirt --tail 50

# 测试API访问
curl -I https://oneclickvirt.ypvps.com/api/v1/health
```

### 方法2：直接替换二进制文件

#### 1. 重新编译
```bash
cd /root/.openclaw/workspace/oneclickvirt

# 安装Go依赖
go mod download

# 编译修复后的代码
go build -o main_fixed server/main.go
```

#### 2. 替换容器中的二进制文件
```bash
# 复制新编译的二进制文件到容器
docker cp main_fixed oneclickvirt:/app/main

# 重启容器
docker restart oneclickvirt
```

#### 3. 验证修复
```bash
# 查看启动日志
docker logs oneclickvirt --tail 20

# 测试创建实例
curl -X POST https://oneclickvirt.ypvps.com/api/v1/admin/instances \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"test-fixed","provider":"1","image":"ubuntu:20.04","cpu":1,"memory":512,"disk":10,"bandwidth":100}'
```

### 方法3：使用Docker Compose

#### 1. 更新docker-compose.yml
```yaml
version: '3.8'

services:
  oneclickvirt:
    image: oneclickvirt:fixed
    container_name: oneclickvirt
    ports:
      - "8080:80"
      - "8443:443"
    volumes:
      - ./data:/app/storage
    restart: unless-stopped
```

#### 2. 重新构建和启动
```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

## 验证测试

### 1. 基础功能测试
```bash
# 获取管理员token
TOKEN=$(curl -s -X POST https://oneclickvirt.ypvps.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"Xk9#mP2$vL8@qR5!"}' | \
  jq -r '.data.token')

# 测试创建实例
curl -X POST https://oneclickvirt.ypvps.com/api/v1/admin/instances \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-fixed-container",
    "provider": "1",
    "image": "ubuntu:20.04",
    "cpu": 1,
    "memory": 512,
    "disk": 10,
    "bandwidth": 100
  }'

# 预期结果：创建成功，无事务错误
```

### 2. 检查日志
```bash
# 查看错误日志
docker exec oneclickvirt cat /app/storage/logs/$(date +%Y-%m-%d)/error.log | tail -10

# 查看信息日志
docker exec oneclickvirt cat /app/storage/logs/$(date +%Y-%m-%d)/info.log | tail -10

# 检查是否有事务错误
docker exec oneclickvirt grep -i "事务\|transaction" /app/storage/logs/$(date +%Y-%m-%d)/error.log
```

### 3. 性能测试
```bash
# 测试创建时间
time curl -X POST https://oneclickvirt.ypvps.com/api/v1/admin/instances \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"test-perf","provider":"1","image":"ubuntu:20.04","cpu":1,"memory":512,"disk":10,"bandwidth":100}'

# 预期结果：创建时间 < 5秒
```

## 回滚方案

### 如果修复后出现问题

#### 1. 快速回滚
```bash
# 停止新容器
docker stop oneclickvirt
docker rm oneclickvirt

# 启动旧容器（如果有备份）
docker run -d --name oneclickvirt \
  -p 8080:80 \
  -p 8443:443 \
  -v /path/to/data:/app/storage \
  oneclickvirt:previous

# 或者使用Git回滚
cd /root/.openclaw/workspace/oneclickvirt
git checkout HEAD~1 -- server/service/resources/quota.go
# 重新构建和部署
```

#### 2. 数据备份
```bash
# 备份数据
docker cp oneclickvirt:/app/storage /path/to/backup/storage_$(date +%Y%m%d_%H%M%S)

# 备份配置
docker cp oneclickvirt:/app/config.yaml /path/to/backup/config_$(date +%Y%m%d_%H%M%S).yaml
```

## 监控和日志

### 实时监控
```bash
# 实时查看日志
docker logs -f oneclickvirt

# 监控容器资源使用
docker stats oneclickvirt

# 检查容器健康状态
docker inspect oneclickvirt | jq '.[0].State.Health'
```

### 日志分析
```bash
# 统计错误数量
docker exec oneclickvirt grep -c "error" /app/storage/logs/$(date +%Y-%m-%d)/error.log

# 查看最近的错误
docker exec oneclickvirt tail -20 /app/storage/logs/$(date +%Y-%m-%d)/error.log

# 搜索特定错误
docker exec oneclickvirt grep "事务隔离级别" /app/storage/logs/$(date +%Y-%m-%d)/error.log
```

## 性能优化

### 资源限制
```bash
# 启动时设置资源限制
docker run -d --name oneclickvirt \
  -p 8080:80 \
  -p 8443:443 \
  --memory="2g" \
  --cpus="2.0" \
  -v /path/to/data:/app/storage \
  oneclickvirt:fixed
```

### 数据库优化
```bash
# 检查数据库连接数
docker exec oneclickvirt mysql -u root -p123456 -e "SHOW PROCESSLIST;"

# 优化数据库配置
docker exec oneclickvirt sed -i 's/max-open-conns: 100/max-open-conns: 200/g' /app/config.yaml
docker restart oneclickvirt
```

## 故障排除

### 问题1：容器无法启动
**症状**：`docker ps`看不到容器
**解决**：
```bash
# 查看容器日志
docker logs oneclickvirt

# 检查端口占用
netstat -tlnp | grep 8080

# 重新启动
docker restart oneclickvirt
```

### 问题2：API无法访问
**症状**：curl请求超时
**解决**：
```bash
# 检查容器网络
docker exec oneclickvirt ping -c 3 8.8.8.8

# 检查端口映射
docker port oneclickvirt

# 检查防火墙
iptables -L -n | grep 8080
```

### 问题3：数据库连接失败
**症状**：日志显示数据库错误
**解决**：
```bash
# 检查数据库状态
docker exec oneclickvirt mysql -u root -p123456 -e "SELECT 1;"

# 检查数据库配置
docker exec oneclickvirt cat /app/config.yaml | grep -A 10 "mysql:"

# 重启数据库
docker exec oneclickvirt supervisorctl restart mysql
```

## 部署检查清单

### 部署前
- [ ] 代码修复已验证
- [ ] 数据已备份
- [ ] 配置文件已检查
- [ ] 依赖已安装

### 部署中
- [ ] 镜像构建成功
- [ ] 容器启动成功
- [ ] 端口映射正确
- [ ] 数据卷挂载正确

### 部署后
- [ ] API访问正常
- [ ] 日志无错误
- [ ] 功能测试通过
- [ ] 性能测试通过
- [ ] 监控配置完成

## 自动化部署脚本

### 部署脚本
```bash
#!/bin/bash

# OneClickVirt修复部署脚本
set -e

echo "开始部署OneClickVirt修复..."

# 1. 备份当前版本
echo "备份当前版本..."
docker cp oneclickvirt:/app/main /tmp/main_backup_$(date +%Y%m%d_%H%M%S)
docker cp oneclickvirt:/app/config.yaml /tmp/config_backup_$(date +%Y%m%d_%H%M%S).yaml

# 2. 停止旧容器
echo "停止旧容器..."
docker stop oneclickvirt || true
docker rm oneclickvirt || true

# 3. 构建新镜像
echo "构建新镜像..."
cd /root/.openclaw/workspace/oneclickvirt
docker build -t oneclickvirt:fixed .

# 4. 启动新容器
echo "启动新容器..."
docker run -d --name oneclickvirt \
  -p 8080:80 \
  -p 8443:443 \
  -v /path/to/data:/app/storage \
  oneclickvirt:fixed

# 5. 等待容器启动
echo "等待容器启动..."
sleep 10

# 6. 验证部署
echo "验证部署..."
if curl -f https://oneclickvirt.ypvps.com/api/v1/health > /dev/null 2>&1; then
    echo "✅ 部署成功！"
else
    echo "❌ 部署失败，开始回滚..."
    docker stop oneclickvirt
    docker rm oneclickvirt
    docker run -d --name oneclickvirt \
      -p 8080:80 \
      -p 8443:443 \
      -v /path/to/data:/app/storage \
      oneclickvirt:previous
    echo "回滚完成"
    exit 1
fi

echo "部署完成！"
```

### 使用方法
```bash
chmod +x deploy_fix.sh
./deploy_fix.sh
```

## 后续维护

### 定期更新
```bash
# 拉取最新代码
cd /root/.openclaw/workspace/oneclickvirt
git pull origin main

# 重新构建和部署
docker build -t oneclickvirt:latest .
docker stop oneclickvirt
docker rm oneclickvirt
docker run -d --name oneclickvirt \
  -p 8080:80 \
  -p 8443:443 \
  -v /path/to/data:/app/storage \
  oneclickvirt:latest
```

### 监控告警
建议配置以下监控：
- 容器健康状态
- API响应时间
- 错误日志数量
- 资源使用情况

---

**最后更新**：2026-04-24
**维护人员**：小美丽
