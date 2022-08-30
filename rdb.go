package main

import (
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.ClusterClient

const DefaultRedisAddress string = "localhost:7001,localhost:7002,localhost:7003,localhost:7004,localhost:7005,localhost:7006"

const DefaultPassword = ""

func init() {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	password := os.Getenv("REDIS_PASSWORD")
	if len(redisAddress) < 5 {
		redisAddress = DefaultRedisAddress
	}
	if len(password) < 5 {
		password = DefaultPassword
	}

	addrs := strings.Split(redisAddress, ",")
	rdb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
		PoolSize: 4,
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
