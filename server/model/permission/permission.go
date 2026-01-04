package permission

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// UserPermission 用户权限组合模型
type UserPermission struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	UserID    uint           `json:"userId" gorm:"not null;index"`
	UserTypes string         `json:"userTypes" gorm:"size:100;not null"` // 逗号分隔的用户类型，如 "user,admin"
	Level     int            `json:"level" gorm:"default:1"`             // 该权限组合的有效等级
	IsActive  bool           `json:"isActive" gorm:"default:true"`       // 是否激活
	Remark    string         `json:"remark" gorm:"size:255"`             // 备注
}

// TableName 指定表名
func (UserPermission) TableName() string {
	return "user_permissions"
}

// GetUserTypes 获取用户类型列表
func (up *UserPermission) GetUserTypes() []string {
	if up.UserTypes == "" {
		return []string{}
	}
	types := strings.Split(up.UserTypes, ",")
	for i, t := range types {
		types[i] = strings.TrimSpace(t)
	}
	return types
}

// SetUserTypes 设置用户类型列表
func (up *UserPermission) SetUserTypes(types []string) {
	up.UserTypes = strings.Join(types, ",")
}

// GetEffectiveUserType 获取最高权限的用户类型
func (up *UserPermission) GetEffectiveUserType() string {
	types := up.GetUserTypes()

	// 权限优先级：admin > user
	for _, t := range types {
		if t == "admin" {
			return "admin"
		}
	}

	// 默认返回user类型
	if len(types) > 0 {
		return types[0]
	}
	return "user"
}
