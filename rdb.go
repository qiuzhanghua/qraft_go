package main

import (
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.ClusterClient

const (
	DefaultRedisAddress = "localhost:7001,localhost:7002,localhost:7003,localhost:7004,localhost:7005,localhost:7006"
	DefaultPassword     = ""
	DefaultPoolSize     = 3
)

func init() {
	viper.SetDefault("redis.cluster", DefaultRedisAddress)
	viper.SetDefault("redis.password", DefaultPassword)
	viper.SetDefault("redis.poolsize", DefaultPoolSize)

	redisAddress := viper.GetString("redis.cluster")
	addr := os.Getenv("REDIS_CLUSTER")
	if len(addr) >= 5 {
		redisAddress = addr
	}

	password := viper.GetString("redis.password")
	p := os.Getenv("REDIS_PASSWORD")
	if len(p) >= 4 {
		password = p
	}

	poolsize := viper.GetInt("redis.poolsize")
	sizeStr := os.Getenv("REDIS_POOLSIZE")
	size, err := strconv.ParseInt(sizeStr, 0, 0)
	if err == nil {
		poolsize = int(size)
	}

	addrs := strings.Split(redisAddress, ",")
	rdb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
		PoolSize: poolsize,
	})

	//rdb := redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs: []string{":7001", ":7002", ":7003", ":7004", ":7005", ":7006"},
	//})
	//err := rdb.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
	//	return shard.Ping(ctx).Err()
	//})
	//if err != nil {
	//	println(err.Error())
	//}

	// rdb = redis.NewClient(&redis.Options{
	// 	Addr:     ":6379",
	// 	Password: "",
	// 	DB:       0,
	// })
}
