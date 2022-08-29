package main

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

var db *bun.DB

func init() {
	sql, err := sql.Open("mysql", "app:app@/app")
	if err != nil {
		panic(err)
	}
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
