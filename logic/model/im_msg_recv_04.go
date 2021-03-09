package model

import (
	"time"
)

// 消息收件箱
type ImMsgRecv04 struct {
	Id         int64     `json:"id" gorm:"primaryKey;column:id;type:bigint(20)"`                           // 消息ID
	ConvType   int       `json:"conv_type" gorm:"column:conv_type;type:tinyint(4);not null"`               // 会话类型[1:单聊;2:群聊]
	Type       int       `json:"type" gorm:"column:type;type:int(11);not null;default:0"`                  // 消息类型[1:文本;2:图片消息;3:语音:4:视频;5:文件;6:地理位置;100:自定义]
	Content    string    `json:"content" gorm:"column:content;type:varchar(5000);not null"`                // 内容
	Extra      string    `json:"extra" gorm:"column:extra;type:varchar(1000);not null"`                    // 额外内容
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;type:datetime;not null"`               // 创建时间
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime;not null"`               // 更新时间
	From       string    `json:"from" gorm:"column:from;type:varchar(50);not null"`                        // 发送者
	To         string    `json:"to" gorm:"column:to;type:varchar(50);not null"`                            // 接收者
	ReadTime   time.Time `json:"read_time" gorm:"column:read_time;type:datetime"`                          // 阅读时间
	Delivered  int       `json:"delivered" gorm:"column:delivered;type:tinyint(4);not null;default:0"`     // 送达状态[0:未送达;1:已送达]
	Target     string    `json:"target" gorm:"column:target;type:varchar(50);not null"`                    // 目标
	Seq        int64     `json:"seq" gorm:"column:seq;type:bigint(20);not null;default:0"`                 // 消息序号
	ClientTime int64     `json:"client_time" gorm:"column:client_time;type:bigint(20);not null;default:0"` // 客户端发送时间
}

func (_ *ImMsgRecv04) TableName() string {
	return "im_msg_recv_04"
}
