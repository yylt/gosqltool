package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yylt/gosqltool/config"
)

func main() {
	err := newRootCmd().Execute()
	if err != nil {
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
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

	cmd.AddCommand(
		//  commands
		newInitcmd(),
	)

	return cmd
}
