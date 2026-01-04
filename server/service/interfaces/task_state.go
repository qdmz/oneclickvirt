package interfaces

// TaskStateManagerInterface 任务状态管理器接口
type TaskStateManagerInterface interface {
	// CompleteMainTask 完成主任务
	CompleteMainTask(taskID uint, success bool, errorMessage string, resultData map[string]interface{}) error

	// CancelMainTask 取消主任务
	CancelMainTask(taskID uint, reason string) error

	// CompleteConfigTask 完成配置任务
	CompleteConfigTask(taskID uint, success bool, errorMessage string, resultData map[string]interface{}) error

	// CancelConfigTask 取消配置任务
	CancelConfigTask(taskID uint, reason string) error

	// StartConfigTask 启动配置任务
	StartConfigTask(taskID uint) error

	// UpdateTaskProgress 更新任务进度
	UpdateTaskProgress(taskID uint, tableType string, progress int, message string) error
}
