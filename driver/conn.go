package driver

import (
	"database/sql"
	"fmt"
	"github.com/yylt/gosqltool/pkg/logger"
	"strings"
	"time"

	"context"
)

var (
	errInvalidDsn  = fmt.Errorf("dsn is invalid!")
	errInvalidKind = fmt.Errorf("db type can not use!")
	defaultTimeout = time.Second
)

type Driver struct {
	*sql.DB
	kind string
	c    *Config
	log  logger.Mlogger
}

type Config struct {
	Net string

	User   string
	Passwd string
	Addr   string

	Dbname string
}

func ping(db *sql.DB) error {
	ctx, clefunc := context.WithTimeout(context.Background(), defaultTimeout)
	defer clefunc()
	err := db.PingContext(ctx)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func ConnStorage(kind string, c *Config, l logger.Mlogger) (*Driver, error) {
	var (
		dsn string
		db  *sql.DB
		err error
	)
	dsn = fmt.Sprintf("%s:%s@%s(%s)/%s", c.User, c.Passwd, c.Net, c.Addr, c.Dbname)
	l.Debug("msg", fmt.Sprintf("start connecting: %s", dsn))
	switch strings.ToLower(kind) {
	case "mysql":
		db, err = mysqlConn(dsn)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errInvalidKind
	}
	if err := ping(db); err != nil {
		return nil, err
	}
	return &Driver{
		DB:   db,
		kind: kind,
		c:    c,
		log:  l,
	}, nil
}

// 执行事务
func (d *Driver) Transaction(txFunc func(tx *sql.Tx) error) (err error) {
	tx, err := d.Begin()
	if err != nil {
		return
	}
	defer func() {
		tx.Rollback()
	}()
	err = txFunc(tx)
	return
}

//执行statement
func (d *Driver) Statement(prepare string, stmtFunc func(stmt *sql.Stmt) error) (err error) {
	var st *sql.Stmt
	if prepare == "" {
		return fmt.Errorf("prepare should not be none")
	}
	st, err = d.Prepare(prepare)

	if err != nil {
		return
	}
	defer func() {
		st.Close()
	}()
	err = stmtFunc(st)
	return
}
