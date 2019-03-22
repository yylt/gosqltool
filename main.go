package main

import (
	"fmt"
	"io"
	"os"

	"github.com/yylt/gosqltool/config"
	"github.com/yylt/gosqltool/pkg/logger"

	"github.com/spf13/cobra"
)

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

func main() {
	newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "gosqltool",
		Short:        "sql tool ,such as init/shell ",
		SilenceUsage: true,
	}
	flags := cmd.PersistentFlags()

	err := config.GlobalConfs.ApplyFlags(flags)
	if err != nil {
		panic(err)
	}
	err = config.StorConfs.ApplyFlags(flags)
	if err != nil {
		panic(err)
	}
	cmd.AddCommand(
		//  commands
		newInitcmd(),
	)

	return cmd
}
