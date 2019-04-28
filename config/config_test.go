package config

import (
	"testing"
	"os"
	"fmt"
	"github.com/spf13/pflag"
)

func TestGetString(t *testing.T){
	var (
		err error
	)
	err = os.Setenv("DBNAME","ha")
	if err!=nil{
		t.Log(err)
	}

	var cs = Confs(map[string]*Conf{
		"db.name": &Conf{
			name:         "db.name",
			t:            STRING,
			e:            "dbname",
			defaultvalue: "default",
		},
	})
	fs:=pflag.NewFlagSet("test", pflag.ContinueOnError)
	cs.ApplyFlags(fs)

	fmt.Println(GetString("db.name"))
}
