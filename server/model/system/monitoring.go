package system

import "time"

// SystemStats 系统统计信息
type SystemStats struct {
	CPU       CPUStats      `json:"cpu"`
	Memory    MemoryStats   `json:"memory"`
	Disk      DiskStats     `json:"disk"`
	Network   NetworkStats  `json:"network"`
	Database  DatabaseStats `json:"database"`
	Runtime   RuntimeStats  `json:"runtime"`
	Timestamp time.Time     `json:"timestamp"`
}

type CPUStats struct {
	Usage     float64 `json:"usage"`       // CPU使用率
	Cores     int     `json:"cores"`       // CPU核心数
	LoadAvg1  float64 `json:"load_avg_1"`  // 1分钟负载平均值
	LoadAvg5  float64 `json:"load_avg_5"`  // 5分钟负载平均值
	LoadAvg15 float64 `json:"load_avg_15"` // 15分钟负载平均值
}

type MemoryStats struct {
	Total     uint64  `json:"total"`      // 总内存
	Used      uint64  `json:"used"`       // 已使用内存
	Free      uint64  `json:"free"`       // 空闲内存
	Usage     float64 `json:"usage"`      // 内存使用率
	SwapTotal uint64  `json:"swap_total"` // 交换分区总大小
	SwapUsed  uint64  `json:"swap_used"`  // 交换分区已使用
}

type DiskStats struct {
	Total uint64  `json:"total"` // 磁盘总大小
	Used  uint64  `json:"used"`  // 磁盘已使用
	Free  uint64  `json:"free"`  // 磁盘空闲
	Usage float64 `json:"usage"` // 磁盘使用率
}

type NetworkStats struct {
	BytesReceived uint64 `json:"bytes_received"` // 接收字节数
	BytesSent     uint64 `json:"bytes_sent"`     // 发送字节数
	PacketsRecv   uint64 `json:"packets_recv"`   // 接收包数
	PacketsSent   uint64 `json:"packets_sent"`   // 发送包数
}

type DatabaseStats struct {
	Connections    int    `json:"connections"`     // 当前连接数
	MaxConnections int    `json:"max_connections"` // 最大连接数
	QueriesTotal   uint64 `json:"queries_total"`   // 总查询数
	SlowQueries    uint64 `json:"slow_queries"`    // 慢查询数
	Uptime         string `json:"uptime"`          // 运行时间
}

type RuntimeStats struct {
	Goroutines int       `json:"goroutines"` // 协程数量
	HeapAlloc  uint64    `json:"heap_alloc"` // 堆内存分配
	HeapSys    uint64    `json:"heap_sys"`   // 堆内存系统
	HeapIdle   uint64    `json:"heap_idle"`  // 堆内存空闲
	HeapInuse  uint64    `json:"heap_inuse"` // 堆内存使用中
	GCCycles   uint32    `json:"gc_cycles"`  // GC循环次数
	LastGC     time.Time `json:"last_gc"`    // 最后GC时间
	Uptime     string    `json:"uptime"`     // 运行时间
}
