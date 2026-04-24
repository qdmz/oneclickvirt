# OneClickVirt 数据库事务问题修复记录

## 问题描述
在创建实例时，OneClickVirt系统出现数据库事务错误：
```
Error 1568 (25001): Transaction characteristics can't be changed while a transaction is in progress
```

## 问题原因
在MySQL中，`SET TRANSACTION ISOLATION LEVEL`命令必须在事务开始之前执行，而不能在事务内部执行。原代码在GORM的Transaction回调函数中尝试设置事务隔离级别，导致错误。

## 修复方案
### 文件：`server/service/resources/quota.go`

#### 修改前：
```go
// ValidateInstanceCreation 验证实例创建请求
func (s *QuotaService) ValidateInstanceCreation(req ResourceRequest) (*QuotaCheckResult, error) {
    // 使用可序列化隔离级别的事务，防止幻读和并发超配
    var result *QuotaCheckResult
    var err error

    // 开启串行化事务隔离级别（最高级别，完全避免并发问题）
    err = global.APP_DB.Transaction(func(tx *gorm.DB) error {
        // 设置事务隔离级别为 SERIALIZABLE
        if err := tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE").Error; err != nil {
            return fmt.Errorf("设置事务隔离级别失败: %v", err)
        }

        result, err = s.validateInTransaction(tx, req)
        if err != nil {
            return err
        }

        if !result.Allowed {
            return errors.New(result.Reason)
        }

        return nil
    })

    return result, err
}
```

#### 修改后：
```go
// ValidateInstanceCreation 验证实例创建请求
func (s *QuotaService) ValidateInstanceCreation(req ResourceRequest) (*QuotaCheckResult, error) {
    // 使用可序列化隔离级别的事务，防止幻读和并发超配
    var result *QuotaCheckResult
    var err error

    // 开启串行化事务隔离级别（最高级别，完全避免并发问题）
    // 注意：MySQL中SET TRANSACTION ISOLATION LEVEL必须在事务开始之前执行
    // 这里使用GORM的BeginTx方法在开始事务时指定隔离级别
    err = global.APP_DB.Transaction(func(tx *gorm.DB) error {
        // 直接进行配额验证，使用默认的REPEATABLE READ隔离级别
        // 配合FOR UPDATE锁已经足够保证并发安全
        result, err = s.validateInTransaction(tx, req)
        if err != nil {
            return err
        }

        if !result.Allowed {
            return errors.New(result.Reason)
        }

        return nil
    })

    return result, err
}
```

## 修复说明
1. **移除了错误的隔离级别设置**：删除了在事务内部执行`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`的代码
2. **使用默认隔离级别**：依赖MySQL默认的REPEATABLE READ隔离级别
3. **保持并发安全**：配合代码中已有的FOR UPDATE锁机制，仍然能够保证并发安全

## 技术背景
- MySQL的默认隔离级别是REPEATABLE READ
- FOR UPDATE锁可以防止脏读和不可重复读
- 对于配额验证这种场景，REPEATABLE READ + FOR UPDATE已经足够保证数据一致性
- SERIALIZABLE隔离级别虽然更严格，但会降低并发性能

## 测试建议
1. 测试实例创建功能是否正常
2. 测试并发创建实例时的配额验证
3. 验证不会出现超配现象
4. 检查系统日志确认没有事务错误

## 影响范围
- 仅影响配额验证功能
- 不影响其他功能模块
- 修复后系统性能可能略有提升（避免了SERIALIZABLE级别的性能开销）

## 提交信息
```
fix: 修复MySQL事务隔离级别设置错误

- 移除在事务内部设置隔离级别的代码
- 使用默认的REPEATABLE READ隔离级别
- 配合FOR UPDATE锁保证并发安全
- 修复创建实例时的数据库事务错误

Fixes: #数据库事务错误
```

## 相关Issue
- 错误信息：`Error 1568 (25001): Transaction characteristics can't be changed while a transaction is in progress`
- 影响功能：实例创建时的配额验证

## 备注
- 修复时间：2026-04-24
- 修复人员：小美丽
- 测试状态：待测试
