package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

// 消息发件箱
type ImMsgSend struct {
	Id         int64                 `json:"id" gorm:"primaryKey;column:id;type:bigint(20) auto_increment"`          // 系统编号
	MsgId      int64                 `json:"msg_id" gorm:"column:msg_id;type:bigint(20);not null"`                   // 消息ID
	ConvType   int                   `json:"conv_type" gorm:"column:conv_type;type:tinyint(4);not null"`             // 会话类型[1:单聊;2:群聊]
	Content    string                `json:"content" gorm:"column:content;type:varchar(5000);not null"`              // 内容
	Extra      string                `json:"extra" gorm:"column:extra;type:varchar(1000);not null"`                  // 扩展
	Type       int                   `json:"type" gorm:"column:type;type:int(11);not null;default:0"`                // 消息类型[1:文本;2:图片消息;3:语音:4:视频;5:文件;6:地理位置;100:自定义]
	CreatedAt  time.Time             `json:"created_at" gorm:"column:created_at;type:datetime;not null"`             // 创建时间
	UpdatedAt  time.Time             `json:"updated_at" gorm:"column:updated_at;type:datetime;not null"`             // 更新时间
	DeletedAt  soft_delete.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:bigint(20);not null;default:0"` // 删除时间
	Sender     string                `json:"sender" gorm:"column:sender;type:varchar(50);not null"`                  // 发送者
	Target     string                `json:"target" gorm:"column:target;type:varchar(50);not null"`                  // 目标
	AtUserList string                `json:"at_user_list" gorm:"column:at_user_list;type:varchar(5000);not null"`    // @用户列表
	ReadTime   int64                 `json:"read_time" gorm:"column:read_time;type:bigint(20);not null;default:0"`   // 消息读取时间
}

func (_ *ImMsgSend) TableName() string {
	return "im_msg_send"
}
