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
}

func (db *MySQLConn) ConnectBlob(config ConfigBlob) error {
	db.DB, _ =sql.Open("mysql",fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.DBConfig.User,config.DBConfig.Password,config.DBConfig.Host,config.DBConfig.Port,config.DBConfig.DBName))
	err:= db.DB.Ping()
	return err
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
