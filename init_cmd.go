package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"context"

	"github.com/yylt/gosqltool/driver"
	"github.com/yylt/gosqltool/pkg/logger"
	"github.com/yylt/gosqltool/config"

	"github.com/spf13/cobra"

)

var (
	bufPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		}}

	defaultTimeout = time.Second * 3
)

func getDriver(log logger.Mlogger) (*driver.Driver, error) {

	conf, err := driver.ParseDSN(config.GetString("store.dsn"))
	if err != nil {
		return nil, err
	}
	return driver.ConnStorage(config.GetString("store.type"), conf, log)
}

func newInitcmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init storage ,such as create db/table/user ",

		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			log := mustGetLogger()
			db, err := getDriver(log)
			if err != nil {
				log.Error("msg", fmt.Sprintf("conn failed: %s", err.Error()))
				return err
			}
			return run(db, log)
		},
	}
	flags := cmd.Flags()
	err := config.InitConfs.ApplyFlags(flags)
	if err != nil {
		panic(err)
	}
	return cmd
}

func run(db *driver.Driver, log logger.Mlogger) error {
	var (
		initsql   io.ReadCloser
		ignoreErr = config.GetBool("init.ignoreErr")
		bs        string
		err       error
	)
	fpath := config.GetString("init.file")
	if fpath == "-" {
		initsql = os.Stdin
	} else {
		f, err := os.Open(fpath)
		if err != nil {
			log.Error("msg", "open sqlfile error")
			return err
		}
		initsql = f
	}
	defer initsql.Close()

	buf := bufio.NewReader(initsql)
	querybuf := bufPool.Get().(*bytes.Buffer)
	querybuf.Reset()
	ctx,canclefunc := context.WithTimeout(context.Background(),defaultTimeout)
	for {
		bs, err = buf.ReadString('\n')
		bs = strings.TrimSpace(bs)

		if err == io.EOF {
			if bs != "" {
				querybuf.WriteString(bs)
				log.Debug("msg", fmt.Sprintf("msg: %s", querybuf.String()))
				_, err = db.ExecContext(ctx,querybuf.String())
				break
			}
			err = nil
			break
		} else if err != nil {
			break
		} else {
			if bs == "" {
				continue
			} else {
				if strings.HasSuffix(bs, ";") != true {
					querybuf.WriteString(bs)
					continue
				} else {
					querybuf.WriteString(bs)
					log.Debug("msg", fmt.Sprintf("msg: %s", querybuf.String()))
					_, err = db.ExecContext(ctx,querybuf.String())
					querybuf.Reset()
					if err != nil {
						break
					}
				}
			}
		}
	}
	canclefunc()
	if err != nil {
		log.Error("msg", fmt.Sprintf("stmt exec failed: %s", err.Error()))
	}
	if ignoreErr {
		return nil
	} else {
		return err
	}
}
