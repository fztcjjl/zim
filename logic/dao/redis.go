package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

const (
	prefixUin = "uin:%s"
)

const (
	Online     = 1
	PushOnline = 2
	Offline    = 3
)

const (
	PushOnlineKeepDays = 7 // 推送在线状态保持天数
)

type ConnInfo struct {
	ConnId         string `json:"conn_id"`
	Platform       string `json:"platform"`
	Device         string `json:"device"`
	Server         string `json:"server"`
	LoginTime      int64  `json:"login_time"`
	DisconnectTime int64  `json:"disconnect_time"`
	Status         int    `json:"status"`
}

func keyUin(uin string) string {
	return fmt.Sprintf(prefixUin, uin)
}

var (
	rs        *Redis
	onceRedis sync.Once
)

type Redis struct {
	client *redis.Client
}

func getRedis() *Redis {
	onceRedis.Do(func() {
		rs = new(Redis)
	})
	return rs
}

func SetupRedis(addr, password string, db int) {
	r := getRedis()
	r.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}

func GetRedisClient() *redis.Client {
	return getRedis().client
	//onceRedisClient.Do(func() {
	//	redisClient = redis.NewClient(&redis.Options{
	//		Addr:     "127.0.0.1:6379",
	//		Password: "123456",
	//		DB:       0,
	//	})
	//})
	//
	//return redisClient
}

//func GetRedisClient() *redis.Client {
//	onceRedisClient.Do(func() {
//		redisClient = redis.NewClient(&redis.Options{
//			Addr:     "127.0.0.1:6379",
//			Password: "",
//			DB:       1,
//		})
//	})
//
//	return redisClient
//}

func AddConn(ctx context.Context, uin string, info *ConnInfo) (err error) {
	client := GetRedisClient()
	b, err := json.Marshal(info)
	if err != nil {
		return
	}
	_, err = client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		key := keyUin(uin)
		if err := pipe.HSet(ctx, key, info.Platform, string(b)).Err(); err != nil {
			return err
		}
		if err := pipe.Expire(ctx, key, time.Duration(PushOnlineKeepDays*24)*time.Hour).Err(); err != nil {
			return err
		}
		return nil
	})
	return
}

func DelConn(ctx context.Context, uin, platform string) (err error) {
	client := GetRedisClient()
	_, err = client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if err := client.HDel(ctx, keyUin(uin), platform).Err(); err != nil {
			return err
		}
		return nil
	})

	return
}

func ExpireConn(ctx context.Context, uin string) (err error) {
	client := GetRedisClient()
	if err = client.Expire(ctx, keyUin(uin), time.Duration(PushOnlineKeepDays*24)*time.Hour).Err(); err != nil {
		return
	}

	return
}

func GetConnByPlatform(ctx context.Context, uin, platform string) *ConnInfo {
	client := GetRedisClient()

	key := keyUin(uin)

	if b, err := client.HGet(ctx, key, platform).Bytes(); err != nil {
		return nil
	} else {
		info := &ConnInfo{}
		if err := json.Unmarshal(b, info); err != nil {
			return nil
		}
		return info
	}
}

func GetConnByUin(ctx context.Context, uin string) (conns map[string][]*ConnInfo, err error) {
	conns = make(map[string][]*ConnInfo)
	client := GetRedisClient()
	r := client.HGetAll(ctx, keyUin(uin))
	if err = r.Err(); err != nil {
		return
	}

	for _, v := range r.Val() {
		info := ConnInfo{}
		if err := json.Unmarshal([]byte(v), &info); err != nil {
			continue
		}
		conns[info.Server] = append(conns[info.Server], &info)
	}

	return
}
