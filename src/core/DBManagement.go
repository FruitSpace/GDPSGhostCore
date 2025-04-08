package core

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

var DBTunnel *sqlx.DB

type MySQLConn struct {
	logger Logger
	DBName string
}

func (db *MySQLConn) ConnectBlob(config ConfigBlob) error {
	db.DBName = config.DBConfig.DBName

	return nil
}

func (db *MySQLConn) CloseDB() error {
	return nil
}

// PatchQuery replaces #DB# with the database name
func (db *MySQLConn) PatchQuery(query string) string {
	return strings.ReplaceAll(query, "#DB#", db.DBName)
}

func (db *MySQLConn) PrepareExec(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := DBTunnel.Prepare(db.PatchQuery(query))
	if err != nil {
		return nil, err
	}
	res, err1 := stmt.Exec(args...)
	return res, err1
}

func (db *MySQLConn) MustPrepareExec(query string, args ...interface{}) sql.Result {
	defer sentry.Recover()
	stmt, err := DBTunnel.Prepare(db.PatchQuery(query))
	if err != nil {
		db.logger.LogErr(db, err.Error())
	}
	fmt.Println(args...)
	res, err1 := stmt.Exec(args...)
	if err1 != nil {
		db.logger.LogErr(db, err1.Error())
	}
	return res
}

func (db *MySQLConn) MustQuery(query string, args ...interface{}) *sql.Rows {
	rows, err := DBTunnel.Query(db.PatchQuery(query), args...)
	if err != nil {
		db.logger.LogErr(db, err.Error())
	}
	return rows
}

func (db *MySQLConn) MustQueryRow(query string, args ...interface{}) *sql.Row {
	row := DBTunnel.QueryRow(db.PatchQuery(query), args...)
	if row.Err() != nil {
		db.logger.LogErr(db, row.Err().Error())
	}
	return row
}

func (db *MySQLConn) ShouldPrepareExec(query string, args ...interface{}) sql.Result {
	stmt, err := DBTunnel.Prepare(db.PatchQuery(query))
	if err != nil {
		db.logger.LogWarn(db, err.Error())
	}
	res, err1 := stmt.Exec(args...)
	if err1 != nil {
		db.logger.LogWarn(db, err1.Error())
	}
	return res
}

func (db *MySQLConn) ShouldQuery(query string, args ...interface{}) *sql.Rows {
	rows, err := DBTunnel.Query(db.PatchQuery(query), args...)
	if err != nil {
		db.logger.LogWarn(db, err.Error())
	}
	return rows
}

func (db *MySQLConn) ShouldQueryRow(query string, args ...interface{}) *sql.Row {
	row := DBTunnel.QueryRow(db.PatchQuery(query), args...)
	if row.Err() != nil {
		db.logger.LogWarn(db, row.Err().Error())
	}
	return row
}

func (db *MySQLConn) ShouldExec(query string, args ...interface{}) {
	_, err := DBTunnel.Exec(db.PatchQuery(query), args...)
	if err != nil {
		db.logger.LogErr(db, err.Error())
	}
}

type RedisConn struct {
	context context.Context
	DB      *redis.Client
}

func (rdb *RedisConn) ConnectBlob(config GlobalConfig) error {
	rdb.context = context.Background()
	rdb.DB = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})
	return rdb.DB.Ping(rdb.context).Err()
}
