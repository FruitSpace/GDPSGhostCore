package core

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLConn struct {
	DB *sql.DB
	logger Logger
}

func (db *MySQLConn) ConnectBlob(config ConfigBlob) error {
	db.DB, _ =sql.Open("mysql",fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.DBConfig.User,config.DBConfig.Password,config.DBConfig.Host,config.DBConfig.Port,config.DBConfig.DBName))
	err:= db.DB.Ping()
	if err!=nil{
		db.logger.LogErr(db,err.Error())
	}
	return err
}

func (db *MySQLConn) PrepareExec(query string, args ...interface{}) (sql.Result, error) {
	stmt, err:=db.DB.Prepare(query)
	if err!=nil {return nil, err}
	res, err1:= stmt.Exec(args...)
	return res, err1
}

func (db *MySQLConn) MustPrepareExec(query string, args ...interface{}) sql.Result {
	stmt, err:=db.DB.Prepare(query)
	if err!=nil {db.logger.LogErr(db,err.Error())}
	res, err1:= stmt.Exec(args...)
	if err1!=nil {db.logger.LogErr(db,err1.Error())}
	return res
}

func (db *MySQLConn) MustQuery(query string, args ...interface{}) *sql.Rows {
	rows, err:= db.DB.Query(query, args...)
	if err!=nil {db.logger.LogErr(db,err.Error())}
	return rows
}

func (db *MySQLConn) MustQueryRow(query string, args ...interface{}) *sql.Row {
	row:=db.DB.QueryRow(query, args...)
	if row.Err()!=nil {db.logger.LogErr(db,row.Err().Error())}
	return row
}

func (db *MySQLConn) ShouldPrepareExec(query string, args ...interface{}) sql.Result {
	stmt, err:=db.DB.Prepare(query)
	if err!=nil {db.logger.LogWarn(db,err.Error())}
	res, err1:= stmt.Exec(args...)
	if err1!=nil {db.logger.LogWarn(db,err1.Error())}
	return res
}

func (db *MySQLConn) ShouldQuery(query string, args ...interface{}) *sql.Rows {
	rows, err:= db.DB.Query(query, args...)
	if err!=nil {db.logger.LogWarn(db,err.Error())}
	return rows
}

func (db *MySQLConn) ShouldQueryRow(query string, args ...interface{}) *sql.Row {
	row:=db.DB.QueryRow(query, args...)
	if row.Err()!=nil {db.logger.LogWarn(db,row.Err().Error())}
	return row
}


type RedisConn struct {
	context context.Context
	DB *redis.Client
}

func (rdb *RedisConn) ConnectBlob(config GlobalConfig) error {
	rdb.context=context.Background()
	rdb.DB = redis.NewClient(&redis.Options{
		Addr: config.RedisHost+":"+config.RedisPort,
		Password: config.RedisPassword,
		DB: config.RedisDB,
	})
	return rdb.DB.Ping(rdb.context).Err()
}
