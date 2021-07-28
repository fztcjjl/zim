module github.com/fztcjjl/zim

go 1.15

replace (
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20210204162551-dae29bb719dd
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20210204162551-dae29bb719dd
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)

require (
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/fztcjjl/tiger v0.0.0-20210213024013-a5b016d9a95a
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis/v8 v8.7.1
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/nats-io/nats.go v1.10.0
	github.com/panjf2000/gnet v1.4.1
	github.com/spf13/cast v1.3.0
	github.com/spf13/viper v1.7.1
	github.com/zentures/cityhash v0.0.0-20131128155616-cdd6a94144ab
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.25.0
	gorm.io/driver/mysql v1.0.4
	gorm.io/gorm v1.21.7
	gorm.io/plugin/soft_delete v1.0.2 // indirect
)
