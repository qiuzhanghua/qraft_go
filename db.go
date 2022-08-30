package main

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"runtime"
	"time"
)

var db *bun.DB

func init() {
	sql, err := sql.Open("mysql", "app:app@/app")
	if err != nil {
		panic(err)
	}
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sql.SetMaxOpenConns(maxOpenConns)
	sql.SetMaxIdleConns(maxOpenConns)
	sql.SetConnMaxIdleTime(2 * time.Second)
	sql.SetConnMaxLifetime(30 * time.Second)

	db = bun.NewDB(sql, mysqldialect.New())

	switch db.Dialect().Name() {
	case dialect.SQLite:
	case dialect.PG:
	case dialect.MySQL:
	case dialect.MSSQL:
	default:
		panic("not reached")
	}

}
