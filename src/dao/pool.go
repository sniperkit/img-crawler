package dao

import (
	"database/sql"
	"fmt"
	"img-crawler/src/conf"
	"img-crawler/src/utils"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db      *sql.DB
	tun     string // target unit name
	addr    string // target address
	slowlog time.Duration
}

func (db *DB) GetDB() *sql.DB {
	return db.db
}

type dbConfig struct {
	maxIdle     int
	maxOpen     int
	maxLifetime time.Duration
	slowlog     time.Duration
}

type Option func(*dbConfig)

type Pool struct {
	master    *DB
	slaves    []*DB
	nextSlave uint64 // Round-robin
}

func NewPool(rawurl string, options ...Option) *Pool {
	return &Pool{
		master:    newDB(rawurl, options...),
		slaves:    make([]*DB, 0),
		nextSlave: 0,
	}
}

func (p *Pool) Master() *DB {
	return p.master
}

func (p *Pool) Slave() *DB {
	slaveNum := uint64(len(p.slaves))
	if slaveNum == 0 {
		return p.master
	}
	return p.slaves[atomic.AddUint64(&p.nextSlave, 1)%slaveNum]
}

func (p *Pool) AddSlave(rawurl string, options ...Option) *Pool {
	p.slaves = append(p.slaves, newDB(rawurl, options...))
	return p
}

func (p *Pool) AddSlaves(rawurls []string) *Pool {
	for _, rawurl := range rawurls {
		p.slaves = append(p.slaves, newDB(rawurl))
	}
	return p
}

func newDB(rawurl string, opts ...Option) *DB {
	conf := &dbConfig{
		maxIdle:     5,
		maxOpen:     10,
		maxLifetime: 1800 * time.Second,
		slowlog:     1 * time.Second,
	}
	for _, o := range opts {
		o(conf)
	}

	u, err := url.Parse(rawurl)
	utils.CheckError(err)

	db, err := sql.Open("mysql", getDsn(u))
	utils.CheckError(err)

	db.SetMaxIdleConns(conf.maxIdle)
	db.SetMaxOpenConns(conf.maxOpen)
	db.SetConnMaxLifetime(conf.maxLifetime)

	return &DB{db: db, tun: getTargetInterface(u), addr: rawurl, slowlog: conf.slowlog}
}

func getDsn(u *url.URL) string {
	ret, _ := url.QueryUnescape(fmt.Sprintf("%s@tcp(%s)%s?%s", u.User, u.Host, u.Path, u.RawQuery))
	return ret
}

func getTargetInterface(u *url.URL) string {
	return "sql_" + strings.Replace(strings.Replace(u.Host, ".", "_", -1), ":", "_", -1)
}

var (
	Mpool *Pool
)

func init() {
	Mpool = NewPool(conf.Config.MySQL.Master)
	Mpool.AddSlaves(conf.Config.MySQL.Slaves)
}
