package driver

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"github.com/yylt/gosqltool/pkg/logger"

	"context"
)

var (
	errInvalidDsn  = fmt.Errorf("dsn is invalid!")
	errInvalidKind = fmt.Errorf("db type can not use!")
	defaultTimeout= time.Second
)

type Driver struct {
	*sql.DB
	kind string
	c    *ConnConfig
	log  logger.Mlogger
}

type ConnConfig struct {
	Net string

	User   string
	Passwd string
	Addr   string

	Dbname string
}

func ParseDSN(dsn string) (*ConnConfig, error) {
	cfg := new(ConnConfig)
	// [user[:password]@][net[(addr)]]/dbname
	found := false
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' {
			found = true
			var j, k int
			// left part is empty if i <= 0
			if i > 0 {
				for j = i; j >= 0; j-- {
					if dsn[j] == '@' {
						// username[:password]
						for k = 0; k < j; k++ {
							if dsn[k] == ':' {
								cfg.Passwd = dsn[k+1 : j]
								break
							}
						}
						cfg.User = dsn[:k]
						break
					}
				}

				// [protocol[(address)]]
				for k = j + 1; k < i; k++ {
					if dsn[k] == '(' {
						if dsn[i-1] != ')' {
							return nil, errInvalidDsn
						}
						cfg.Addr = dsn[k+1 : i-1]
						break
					}
				}
				cfg.Net = dsn[j+1 : k]
			}
			cfg.Dbname = dsn[i+1 : len(dsn)]
			break
		}
	}

	if !found && len(dsn) > 0 {
		return nil, errInvalidDsn
	}

	return cfg, nil

}

func ping(db *sql.DB) (error){
	ctx,clefunc:=context.WithTimeout(context.Background(),defaultTimeout)
	defer clefunc()
	err:=db.PingContext(ctx)
	if err!=nil{
		return err
	}else{
		return nil
	}
}

func ConnStorage(kind string, c *ConnConfig, l logger.Mlogger) (*Driver, error) {
	var (
		dsn string
		db  *sql.DB
		err error
	)
	dsn = fmt.Sprintf("%s:%s@%s(%s)/%s", c.User, c.Passwd, c.Net, c.Addr, c.Dbname)
	l.Debug("msg",fmt.Sprintf("start connecting: %s://%s",kind,dsn))
	switch strings.ToLower(kind) {
	case "mysql":
		dsn = fmt.Sprintf("%s?timeout=1s",dsn)
		db, err =  sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}
		if err:=ping(db);err!=nil{
			return nil,err
		}
	default:
		return nil, errInvalidKind
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
