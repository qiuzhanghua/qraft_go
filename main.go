package main

import (
	"context"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/gommon/log"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"io"
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

	info, err := rdb.Info(ctx, "server").Result()
	idx := strings.Index(info, "\n")
	idx2 := strings.Index(info[idx+1:], "\n")
	idx3 := strings.Index(info[idx+1:], ":")
	redisVersion := ""
	if idx >= 0 && idx2 >= 0 && idx3 >= 0 && idx2 > idx3+2 {
		redisVersion = info[idx+idx3+2 : idx+idx2]
	}
	log.Infof("redis %s ready ...", redisVersion)

	//err = mq.WriteMessages(context.Background(), kafka.Message{
	//	Value: []byte("Hello, world!"),
	//})
	//log.Info(err)

	defer func(mq *kafka.Writer) {
		_ = mq.Close()
	}(mq)

	// test for mysql/pg only
	// for mssql     select @@version as v
	// for sqlite    SELECT sqlite_version() as v
	var ver string
	err = db.NewRaw("SELECT version() as v").Scan(ctx, &ver)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	defer func(db *bun.DB) {
		_ = db.Close()
	}(db)
	log.Infof("%s %s ready ...", dbtype, ver)

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

	router.GET("/rdb", func(w http.ResponseWriter, req bunrouter.Request) error {
		ans, err := rdb.Get(ctx, "hello").Result()

		if err == nil {
			_, _ = io.WriteString(w, ans)
		}
		return err
	})

	router.GET("/db", func(w http.ResponseWriter, req bunrouter.Request) error {

		var ver string
		err = db.NewRaw("SELECT version() as v").Scan(ctx, &ver)
		if err == nil {
			_, _ = io.WriteString(w, ver)
		}
		return err
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
	log.SetPrefix("qraft")
	format := strings.ToLower(os.Getenv("LOGGING_FORMAT"))
	if format != "json" {
		log.SetHeader(`${time_rfc3339_nano}, ${prefix}, ${level} ${short_file}(${line})`)
	}
	log.SetOutput(os.Stdout)
	level := strings.ToLower(os.Getenv("LOGGING_LEVEL"))
	x := levelOf(level)
	log.SetLevel(x)
	log.SetLevel(log.DEBUG)

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("config file not found!")
		} else {
			log.Errorf("config file error: %s", err.Error())
		}
	} else {
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Infof("Config file changed: %s, REBOOT please.", e.Name)
		})
		viper.WatchConfig()
	}
	viper.AutomaticEnv()
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
