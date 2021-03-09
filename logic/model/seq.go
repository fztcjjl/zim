package model

// 序号生成器
type Seq struct {
	Id      int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20) auto_increment"`
	ObjType int    `json:"obj_type" gorm:"column:obj_type;type:tinyint(4);not null;default:0"`
	ObjId   string `json:"obj_id" gorm:"column:obj_id;type:varchar(50);not null;default:0"`
	Seq     int64  `json:"seq" gorm:"column:seq;type:bigint(20);not null;default:0"`
}

func (_ *Seq) TableName() string {
	return "seq"
}