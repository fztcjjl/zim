package util

import "time"

func GetTimeMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
