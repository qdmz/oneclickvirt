package initialize

import (
	"strings"

	"oneclickvirt/global"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// schemaSyncResult 汇总一次增量结构同步做了什么，便于日志和测试。
type schemaSyncResult struct {
	CreatedTables  []string
	AddedColumns   []string // table.column
	CreatedIndexes []string // table.index
	Errors         int
}

// ensureAdditiveSchema 对已初始化（已有数据）的数据库执行“只增不改”的结构同步：
//   - 缺失的表：CreateTable
//   - 已有表缺失的列：AddColumn
//   - 已有表缺失的索引：CreateIndex
//
// 它从不 ALTER/DROP 已存在的列、索引或表，因此对已有数据是安全的，
// 且每次启动都可重复执行（幂等）。与 AutoMigrate 不同，它不会因为列类型/长度/
// 可空性的细微差异去修改已有列，也避免了 MariaDB 上 AutoMigrate 挂起的问题。
// 已有列的类型修正仍由 FixSchemaColumns 负责。
func ensureAdditiveSchema(db *gorm.DB, models []interface{}) schemaSyncResult {
	var res schemaSyncResult
	if db == nil {
		return res
	}
	existing := loadExistingSchema(db)
	m := db.Migrator()
	addedDesiredState := false

	for _, mdl := range models {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(mdl); err != nil {
			global.APP_LOG.Warn("解析模型结构失败，跳过增量同步", zap.Error(err))
			res.Errors++
			continue
		}
		table := stmt.Schema.Table
		key := strings.ToLower(table)

		tableExists := false
		if existing != nil {
			_, tableExists = existing.columns[key]
		} else {
			tableExists = m.HasTable(mdl)
		}
		if !tableExists {
			global.APP_LOG.Info("增量同步：创建缺失的表", zap.String("table", table))
			if err := m.CreateTable(mdl); err != nil {
				global.APP_LOG.Warn("增量同步：创建表失败", zap.String("table", table), zap.Error(err))
				res.Errors++
				continue
			}
			res.CreatedTables = append(res.CreatedTables, table)
			continue
		}

		for _, f := range stmt.Schema.Fields {
			if f.DBName == "" || f.IgnoreMigration {
				continue
			}
			var has bool
			if existing != nil {
				has = existing.columns[key][strings.ToLower(f.DBName)]
			} else {
				has = m.HasColumn(mdl, f.DBName)
			}
			if has {
				continue
			}
			global.APP_LOG.Info("增量同步：添加缺失的列", zap.String("table", table), zap.String("column", f.DBName))
			if err := m.AddColumn(mdl, f.Name); err != nil {
				global.APP_LOG.Warn("增量同步：添加列失败", zap.String("table", table), zap.String("column", f.DBName), zap.Error(err))
				res.Errors++
				continue
			}
			res.AddedColumns = append(res.AddedColumns, table+"."+f.DBName)
			if table == "instances" && f.DBName == "desired_state" {
				addedDesiredState = true
			}
		}

		for _, idx := range stmt.Schema.ParseIndexes() {
			if idx == nil || idx.Name == "" {
				continue
			}
			var has bool
			if existing != nil {
				has = existing.indexes[key][strings.ToLower(idx.Name)]
			} else {
				has = m.HasIndex(mdl, idx.Name)
			}
			if has {
				continue
			}
			global.APP_LOG.Info("增量同步：创建缺失的索引", zap.String("table", table), zap.String("index", idx.Name))
			if err := m.CreateIndex(mdl, idx.Name); err != nil {
				global.APP_LOG.Warn("增量同步：创建索引失败", zap.String("table", table), zap.String("index", idx.Name), zap.Error(err))
				res.Errors++
				continue
			}
			res.CreatedIndexes = append(res.CreatedIndexes, table+"."+idx.Name)
		}
	}

	// 新增 desired_state 列时，与全量迁移一致地为历史实例回填运行意图。
	// 仅在本次刚添加该列时执行，避免每次启动改写管理员设置的意图。
	if addedDesiredState {
		backfillInstanceDesiredStates(db)
	}

	if len(res.CreatedTables)+len(res.AddedColumns)+len(res.CreatedIndexes) > 0 || res.Errors > 0 {
		global.APP_LOG.Info("增量表结构同步完成",
			zap.Strings("createdTables", res.CreatedTables),
			zap.Strings("addedColumns", res.AddedColumns),
			zap.Strings("createdIndexes", res.CreatedIndexes),
			zap.Int("errors", res.Errors))
	} else {
		global.APP_LOG.Info("增量表结构同步完成：表结构已是最新")
	}
	return res
}

type existingSchema struct {
	columns map[string]map[string]bool // table -> column set (lowercase)
	indexes map[string]map[string]bool // table -> index set (lowercase)
}

// loadExistingSchema 在 MySQL/MariaDB 上用两次 information_schema 查询一次性读取
// 当前库的全部列和索引，避免逐列查询在慢机器上拖慢启动。其他方言返回 nil，
// 由调用方回退到 Migrator 的逐项检查。
func loadExistingSchema(db *gorm.DB) *existingSchema {
	if db.Dialector == nil || db.Dialector.Name() != "mysql" {
		return nil
	}
	type colRow struct {
		TableName  string `gorm:"column:TABLE_NAME"`
		ColumnName string `gorm:"column:COLUMN_NAME"`
	}
	type idxRow struct {
		TableName string `gorm:"column:TABLE_NAME"`
		IndexName string `gorm:"column:INDEX_NAME"`
	}
	var cols []colRow
	if err := db.Raw("SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE()").Scan(&cols).Error; err != nil {
		global.APP_LOG.Warn("读取现有列信息失败，回退到逐项检查", zap.Error(err))
		return nil
	}
	var idxs []idxRow
	if err := db.Raw("SELECT DISTINCT TABLE_NAME, INDEX_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE()").Scan(&idxs).Error; err != nil {
		global.APP_LOG.Warn("读取现有索引信息失败，回退到逐项检查", zap.Error(err))
		return nil
	}
	s := &existingSchema{columns: map[string]map[string]bool{}, indexes: map[string]map[string]bool{}}
	for _, c := range cols {
		t := strings.ToLower(c.TableName)
		if s.columns[t] == nil {
			s.columns[t] = map[string]bool{}
		}
		s.columns[t][strings.ToLower(c.ColumnName)] = true
	}
	for _, i := range idxs {
		t := strings.ToLower(i.TableName)
		if s.indexes[t] == nil {
			s.indexes[t] = map[string]bool{}
		}
		s.indexes[t][strings.ToLower(i.IndexName)] = true
	}
	return s
}
