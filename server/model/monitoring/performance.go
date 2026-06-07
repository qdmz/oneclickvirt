package monitoring

import (
	"time"

	"gorm.io/gorm"
)

// PerformanceMetric 性能指标历史记录
type PerformanceMetric struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 时间戳
	Timestamp time.Time `gorm:"index:idx_timestamp;not null" json:"timestamp"`

	// 系统基础指标
	GoroutineCount int `json:"goroutineCount"`
	CPUCount       int `json:"cpuCount"`

	// 内存指标 (MB)
	MemoryAlloc      uint64 `json:"memoryAlloc"`       // 当前分配的内存
	MemoryTotalAlloc uint64 `json:"memoryTotalAlloc"` // 累计分配的内存
	MemorySys        uint64 `json:"memorySys"`         // 从系统获取的内存
	MemoryHeapAlloc  uint64 `json:"memoryHeapAlloc"`  // 堆上分配的内存
	MemoryHeapSys    uint64 `json:"memoryHeapSys"`    // 堆从系统获取的内存
	MemoryStackInuse uint64 `json:"memoryStackInuse"` // 栈使用的内存

	// GC 指标
	GCCount      uint32 `json:"gcCount"`       // GC次数
	GCPauseTotal uint64 `json:"gcPauseTotal"` // GC总暂停时间(ns)
	GCPauseAvg   uint64 `json:"gcPauseAvg"`   // GC平均暂停时间(ns)
	GCLastPause  uint64 `json:"gcLastPause"`  // 上次GC暂停时间(ns)
	NextGC       uint64 `json:"nextGC"`        // 下次GC触发阈值

	// 数据库连接池状态
	DBMaxOpenConnections int   `json:"dbMaxOpenConnections"` // 最大连接数
	DBOpenConnections    int   `json:"dbOpenConnections"`     // 当前打开的连接数
	DBInUse              int   `json:"dbInUse"`               // 正在使用的连接数
	DBIdle               int   `json:"dbIdle"`                 // 空闲连接数
	DBWaitCount          int64 `json:"dbWaitCount"`           // 等待连接的总次数
	DBWaitDuration       int64 `json:"dbWaitDuration"`        // 等待连接的总时间(ns)
	DBMaxIdleClosed      int64 `json:"dbMaxIdleClosed"`      // 因超过最大空闲数而关闭的连接数
	DBMaxLifetimeClosed  int64 `json:"dbMaxLifetimeClosed"`  // 因超过最大生命周期而关闭的连接数

	// SSH连接池状态
	SSHTotalConnections     int     `json:"sshTotalConnections"`     // SSH总连接数
	SSHHealthyConnections   int     `json:"sshHealthyConnections"`   // SSH健康连接数
	SSHUnhealthyConnections int     `json:"sshUnhealthyConnections"` // SSH不健康连接数
	SSHIdleConnections      int     `json:"sshIdleConnections"`      // SSH空闲连接数
	SSHActiveConnections    int     `json:"sshActiveConnections"`    // SSH活跃连接数
	SSHMaxConnections       int     `json:"sshMaxConnections"`       // SSH最大连接数限制
	SSHUtilization          float64 `json:"sshUtilization"`           // SSH连接池利用率(%)
	SSHOldestConnectionAge  int64   `json:"sshOldestConnectionAge"` // 最老连接年龄(秒)
	SSHNewestConnectionAge  int64   `json:"sshNewestConnectionAge"` // 最新连接年龄(秒)
	SSHAvgConnectionAge     int64   `json:"sshAvgConnectionAge"`    // 平均连接年龄(秒)

	// 任务系统状态
	TaskRunningContexts int `json:"taskRunningContexts"` // 运行中的任务上下文数量
	TaskProviderPools   int `json:"taskProviderPools"`   // Provider工作池数量
	TaskTotalQueueSize  int `json:"taskTotalQueueSize"` // 总队列大小
}

// TableName 指定表名
func (PerformanceMetric) TableName() string {
	return "performance_metrics"
}
