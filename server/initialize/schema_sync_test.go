package initialize

import (
	"fmt"
	"testing"
	"time"

	"oneclickvirt/global"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type schemaSyncLegacy struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func (schemaSyncLegacy) TableName() string { return "schema_sync_items" }

type schemaSyncCurrent struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:64"`
	State string `gorm:"size:16;index"`
}

func (schemaSyncCurrent) TableName() string { return "schema_sync_items" }

type schemaSyncNewTable struct {
	ID    uint `gorm:"primaryKey"`
	Value string
}

func TestEnsureAdditiveSchemaAddsMissingTablesColumnsAndIndexes(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s_%d?mode=memory&cache=shared", t.Name(), time.Now().UnixNano())), &gorm.Config{})
	if err == nil {
		err = db.Exec("SELECT 1").Error
	}
	if err != nil {
		t.Skipf("sqlite unavailable (needs cgo): %v", err)
	}
	sqlDB, _ := db.DB()
	t.Cleanup(func() { _ = sqlDB.Close() })
	previousLog := global.APP_LOG
	global.APP_LOG = zap.NewNop()
	t.Cleanup(func() { global.APP_LOG = previousLog })

	if err := db.AutoMigrate(&schemaSyncLegacy{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&schemaSyncLegacy{Name: "keep-me"}).Error; err != nil {
		t.Fatal(err)
	}

	res := ensureAdditiveSchema(db, []interface{}{&schemaSyncCurrent{}, &schemaSyncNewTable{}})
	if res.Errors != 0 {
		t.Fatalf("unexpected errors: %+v", res)
	}
	m := db.Migrator()
	if !m.HasColumn(&schemaSyncCurrent{}, "state") {
		t.Fatal("missing column was not added")
	}
	if !m.HasIndex(&schemaSyncCurrent{}, "idx_schema_sync_items_state") {
		t.Fatal("missing index was not created")
	}
	if !m.HasTable(&schemaSyncNewTable{}) {
		t.Fatal("missing table was not created")
	}
	var n int64
	db.Table("schema_sync_items").Where("name = ?", "keep-me").Count(&n)
	if n != 1 {
		t.Fatalf("existing data lost, count=%d", n)
	}

	// Idempotent: a second run changes nothing.
	res = ensureAdditiveSchema(db, []interface{}{&schemaSyncCurrent{}, &schemaSyncNewTable{}})
	if len(res.CreatedTables)+len(res.AddedColumns)+len(res.CreatedIndexes)+res.Errors != 0 {
		t.Fatalf("second run should be a no-op: %+v", res)
	}
}
