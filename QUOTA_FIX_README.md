# OneClickVirt 数据库事务问题修复补丁

## 修复文件
- `server/service/resources/quota.go`

## 修复内容
移除在事务内部设置MySQL事务隔离级别的代码，避免出现"Transaction characteristics can't be changed while a transaction is in progress"错误。

## 应用方法

### 方法1：手动修改
1. 打开文件 `server/service/resources/quota.go`
2. 找到 `ValidateInstanceCreation` 函数
3. 移除以下代码：
   ```go
   // 设置事务隔离级别为 SERIALIZABLE
   if err := tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE").Error; err != nil {
       return fmt.Errorf("设置事务隔离级别失败: %v", err)
   }
   ```
4. 保存文件

### 方法2：使用git apply
```bash
cd /path/to/oneclickvirt
git apply < quota_fix.patch
```

### 方法3：直接替换文件
```bash
cp quota.go.fixed server/service/resources/quota.go
```

## 重新编译
```bash
cd /path/to/oneclickvirt
go build -o main server/main.go
```

## 重新构建Docker镜像
```bash
docker build -t oneclickvirt:fixed .
```

## 测试验证
1. 启动修复后的容器
2. 尝试创建实例
3. 检查是否还有事务错误
4. 验证配额验证功能正常

## 回滚方法
如果修复后出现问题，可以使用git回滚：
```bash
git checkout HEAD -- server/service/resources/quota.go
```

## 注意事项
1. 修复后使用默认的REPEATABLE READ隔离级别
2. 配合FOR UPDATE锁仍然能保证并发安全
3. 性能可能略有提升（避免了SERIALIZABLE级别的开销）
4. 建议在测试环境充分测试后再部署到生产环境
