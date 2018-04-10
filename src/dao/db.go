package dao

import (
	"img-crawler/src/conf"

	_ "github.com/go-sql-driver/mysql"
)

type Pool struct {
	master    *DB
	slaves    []*DB
	nextSlave uint64 // Round-robin
}

var (
	MysqlPool *mysql.Pool
)

func init() {
	Pool = mysql.NewPool(conf.Config.MySQL.Master)
	Pool.AddSlaves(conf.Config.MySQL.Slaves)
}
