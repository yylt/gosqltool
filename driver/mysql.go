package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"time"
)

type mysqlDb struct {
	conntimeout string
	wtimeout    string
	rtimeout    string

	multiStat string

	*sql.DB
}

type mysqlOpt func(*mysqlDb)

func MysqlTimeoutOpt(duration time.Duration) mysqlOpt {
	return func(o *mysqlDb) {
		seconds := int(math.Ceil(duration.Seconds()))
		o.conntimeout = fmt.Sprintf("timeout=%ds", seconds)
	}
}

func MysqlMultiStatOpt() mysqlOpt {
	return func(o *mysqlDb) {
		o.multiStat = fmt.Sprintf("multiStatement=true")
	}
}

func mysqlConn(dsn string) (*sql.DB, error) {
	var (
		opts = []mysqlOpt{MysqlTimeoutOpt(time.Second)}
		db   = new(mysqlDb)
	)

	db.ApplyOpts(opts...)
	return db.conn(dsn)
}

func (m *mysqlDb) ApplyOpts(opts ...mysqlOpt) {
	for _, v := range opts {
		v(m)
	}
}

func (m *mysqlDb) conn(dsn string) (*sql.DB, error) {
	var (
		err error
	)
	dsn = addMark(dsn)
	if m.multiStat != "" {
		dsn = fmt.Sprintf("%s%s", dsn, m.multiStat)
	}
	if m.conntimeout != "" {
		dsn = fmt.Sprintf("%s%s", dsn, m.conntimeout)
	}
	m.DB, err = sql.Open("mysql", dsn)
	return m.DB, err
}
