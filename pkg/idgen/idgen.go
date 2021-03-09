package idgen

import (
	"github.com/fztcjjl/zim/pkg/idgen/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	st := sonyflake.Settings{
		MachineID: getMachineId,
	}
	sf = sonyflake.NewSonyflake(st)
}

func getMachineId() (uint16, error) {
	// TODO
	return 1, nil
}

func Next() (id int64) {
	var i uint64
	if sf != nil {
		i, _ = sf.NextID()
		id = int64(i)
	}
	return
}

func GetOne() int64 {
	return Next()
}

func GetMulti(n int) (ids []int64) {
	for i := 0; i < n; i++ {
		ids = append(ids, Next())
	}
	return
}
