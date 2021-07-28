package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

// 任务成员
type ImGroupMember struct {
	Id        int64                 `json:"id" gorm:"primaryKey;column:id;type:bigint(20)"`                         // 系统编号
	GroupId   string                `json:"group_id" gorm:"column:group_id;type:varchar(50);not null"`              // 群ID
	Member    string                `json:"member" gorm:"column:member;type:varchar(50);not null"`                  // 成员ID
	CreatedAt time.Time             `json:"created_at" gorm:"column:created_at;type:datetime;not null"`             // 创建时间
	UpdatedAt time.Time             `json:"updated_at" gorm:"column:updated_at;type:datetime;not null"`             // 更新时间
	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:bigint(20);not null;default:0"` // 删除时间
}

func (_ *ImGroupMember) TableName() string {
	return "im_group_member"
}
