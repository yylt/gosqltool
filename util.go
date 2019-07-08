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
	jsonfmt := config.GetBool("log.jsonfmt")
	logout = os.Stdout
	if jsonfmt {
		return logger.NewJsonLogger(logout, lvl)
	} else {
		return logger.NewStrLogger(logout, lvl)
	}
}
