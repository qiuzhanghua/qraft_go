package main

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"os"
	"runtime"
	"strings"
	"time"
)

var sqldb *sql.DB
var db *bun.DB
var dbtype string

const (
	DefaultDbType = "mysql"
	DefaultDbUrl  = "app:app@/app"
)

func DbInit() {
	viper.SetDefault("db.type", DefaultDbType)
	viper.SetDefault("db.url", DefaultDbUrl)
	dbtype = strings.ToLower(viper.GetString("db.type"))
	type2 := os.Getenv("DB_TYPE")
	if len(type2) >= 2 {
		dbtype = type2
	}
	switch dbtype {
	case "pg", "sqlite", "mysql", "mssql":
		break
	default:
		panic("db type not support")
	}
	url := viper.GetString("db.url")
	url2 := os.Getenv("DB_URL")
	if len(url2) >= 5 {
		url = url2
	}
	var err error
	switch dbtype {
	case "pg":
		sqldb = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(url)))
		db = bun.NewDB(sqldb, pgdialect.New())
		break
	case "mysql":
		sqldb, err = sql.Open("mysql", url)
		if err != nil {
			panic(err)
		}
		db = bun.NewDB(sqldb, mysqldialect.New())
		break
	case "mssql":
		sqldb, err = sql.Open("mssql", url)
		if err != nil {
			panic(err)
		}
		db = bun.NewDB(sqldb, mssqldialect.New())
		break
	case "sqlite":
		sqldb, err = sql.Open(sqliteshim.ShimName, url)
		if err != nil {
			panic(err)
		}
		db = bun.NewDB(sqldb, sqlitedialect.New())
		break
	default:
		panic("unreachable")
	}
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)
	sqldb.SetConnMaxIdleTime(2 * time.Second)
	sqldb.SetConnMaxLifetime(30 * time.Second)

	//switch db.Dialect().Name() {
	//case dialect.SQLite:
	//case dialect.PG:
	//case dialect.MySQL:
	//case dialect.MSSQL:
	//default:
	//	panic("not reached")
	//}

}
