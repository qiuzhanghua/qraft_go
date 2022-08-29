package main

import (
	"context"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"net/http"
	"os"
	"strings"
)

var ctx = context.Background()

func main() {
	// check redis ready
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	defer func(rdb *redis.ClusterClient) {
		_ = rdb.Close()
	}(rdb)

	log.Info("Redis ready...")

	//err = mq.WriteMessages(context.Background(), kafka.Message{
	//	Value: []byte("Hello, world!"),
	//})
	//log.Info(err)

	defer func(mq *kafka.Writer) {
		_ = mq.Close()
	}(mq)

	var ver string
	err = db.NewRaw("SELECT version() as v").Scan(ctx, &ver)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	defer func(db *bun.DB) {
		_ = db.Close()
	}(db)
	log.Infof("MySQL %s ready...", ver)

	router := bunrouter.New(
		bunrouter.WithMiddleware(reqlog.NewMiddleware(
			reqlog.FromEnv("BUNDEBUG"),
		)),
	)

	router.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		// req embeds *http.Request and has all the same fields and methods
		fmt.Println("Beijing Health Kit")
		return nil
	})

	port := os.Getenv("HT_PORT")

	if len(port) < 1 {
		port = ":9999"
	}
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Error("Can't Open Web", err)
	}
}

func init() {
	log.SetPrefix("ht")
	format := strings.ToLower(os.Getenv("LOGGING_FORMAT"))
	if format != "json" {
		log.SetHeader(`${time_rfc3339_nano}, ${prefix}, ${level} ${short_file}(${line})`)
	}
	log.SetOutput(os.Stdout)
	level := strings.ToLower(os.Getenv("LOGGING_LEVEL"))
	x := levelOf(level)
	log.SetLevel(x)
	log.SetLevel(log.DEBUG)

}

func levelOf(s string) log.Lvl {
	switch s {
	case "off":
		return log.OFF
	case "debug":
		return log.DEBUG
	case "info":
		return log.INFO
	case "warn":
		return log.WARN
	case "error":
		return log.ERROR
	default:
		return log.ERROR
	}
}
