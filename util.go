package main

import (
	"fmt"
	"github.com/yylt/gosqltool/config"
	"github.com/yylt/gosqltool/driver"
	"github.com/yylt/gosqltool/pkg/logger"
	"io"
	"os"
)

func getDriver(log logger.Mlogger) (*driver.Driver, error) {

	conf := &driver.Config{
		Net:    config.GetString("db.net"),
		User:   config.GetString("db.user"),
		Passwd: config.GetString("db.password"),
		Dbname: config.GetString("db.dbname"),
		Addr:   config.GetString("db.addr"),
	}
	return driver.ConnStorage(config.GetString("db.type"), conf, log)
}

func mustGetLogger() logger.Mlogger {
	var (
		logout io.Writer
	)
	lvl := config.GetString("log.level")
	if lvl == "" {
		lvl = "info"
	}
	jsonfmt := config.GetBool("log.jsonfmt")
	out := config.GetString("log.out")
	if out == "" {
		logout = os.Stdout
	} else {
		if s, err := os.Stat(out); err != nil {
			f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				panic(err)
			}
			logout = f
		} else {
			if s.IsDir() == true {
				panic(fmt.Errorf("%s is dir ,can not log!", out))
			}
		}
	}
	if jsonfmt {
		return logger.NewJsonLogger(logout, lvl)
	} else {
		return logger.NewJsonLogger(logout, lvl)
	}
}
