package model

// 序号生成器
type Seq struct {
	Id     int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20) auto_increment"`
	UserId string `json:"user_id" gorm:"column:user_id;type:varchar(50);not null;default:0"`
	Seq    int64  `json:"seq" gorm:"column:seq;type:bigint(20);not null;default:0"`
}

func (_ *Seq) TableName() string {
	return "seq"
}
