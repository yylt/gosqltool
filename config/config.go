package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Conf struct {
	name         string
	t            T           // type
	e            string      // env
	defaultvalue interface{} // default
	desc         string
}

type T string
type Confs map[string]*Conf

var (
	INT    T = "int"
	BOOL   T = "bool"
	STRING T = "string"

	Prefix = "gosql"
)

var GlobalConfs = Confs(map[string]*Conf{
	"log.level": &Conf{
		name:         "log.level",
		t:            STRING,
		e:            "LOGLEVEL",
		defaultvalue: "info",
	},
	"log.jsonfmt": &Conf{
		name:         "log.jsonfmt",
		t:            BOOL,
		e:            "LOGJSONFMT",
		defaultvalue: true,
	},
	"log.out": &Conf{
		name:         "log.out",
		t:            STRING,
		e:            "LOGstdout",
		defaultvalue: "",
	},
})

var InitConfs = Confs(map[string]*Conf{
	"init.file": &Conf{
		name:         "init.file",
		t:            STRING,
		e:            "initfile",
		defaultvalue: "-",
		desc:         "sql file to init , default from stdin",
	},
	"init.ignoreErr": &Conf{
		name:         "init.ignoreErr",
		t:            BOOL,
		e:            "initignoreErr",
		defaultvalue: false,
		desc:         "ignore error when exec statement , but connErr will not ignore",
	},
})

var StorConfs = Confs(map[string]*Conf{
	"db.addr": &Conf{
		name:         "db.addr",
		t:            STRING,
		e:            "dbaddr",
		defaultvalue: "",
		desc:         "db addr, envname gosqldbaddr",
	},
	"db.net": &Conf{
		name:         "db.net",
		t:            STRING,
		e:            "dbnet",
		defaultvalue: "tcp",
		desc:         "db net, envname gosqldbnet",
	},
	"db.dbname": &Conf{
		name:         "db.dbname",
		t:            STRING,
		e:            "dbname",
		defaultvalue: "",
		desc:         "db dbname, envname gosqldbname",
	},
	"db.user": &Conf{
		name:         "db.user",
		t:            STRING,
		e:            "dbuser",
		defaultvalue: "admin",
		desc:         "db user, envname gosqldbuser",
	},
	"db.password": &Conf{
		name:         "db.password",
		t:            STRING,
		e:            "pw",
		defaultvalue: "admin",
		desc:         "db password, envname gosqlpw",
	},
	"db.type": &Conf{
		name:         "db.type",
		t:            STRING,
		e:            "type",
		defaultvalue: "mysql",
		desc:         "db type, now support mysql,  envname gosqltype",
	},
})

func (c Confs) ForEach(fn func(cf *Conf) error) error {
	for _, v := range c {
		err := fn(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Confs) ApplyConfFile(app string) error {
	if f, err := os.Open(app); err != nil {
		return err
	} else {
		return viper.ReadConfig(f)
	}
}

func (c Confs) Merge(app Confs) {
	for n, v := range app {
		c[n] = v
	}
}

func (c Confs) ApplyConf(cf *Conf) {
	if _, ok := c[cf.name]; ok {
		c[cf.name].defaultvalue = cf.defaultvalue
		c[cf.name].t = cf.t
		c[cf.name].desc = cf.desc
	} else {
		c[cf.name] = cf
	}
}

func (c Confs) ApplyFlags(fs *pflag.FlagSet) error {
	for _, v := range c {
		switch v.t {
		case BOOL:
			value, ok := v.defaultvalue.(bool)
			if ok {
				fs.Bool(v.name, value, v.desc)
				viper.BindEnv(v.name, v.e)
			}
		case INT:
			value, ok := v.defaultvalue.(int)
			if ok {
				fs.Int(v.name, value, v.desc)
				viper.BindEnv(v.name, v.e)
			}
		case STRING:
			value, ok := v.defaultvalue.(string)
			if ok {
				fs.String(v.name, value, v.desc)
				viper.BindEnv(v.name, v.e)
			}
		default:
			return fmt.Errorf("not correct type!")
		}
	}

	return viper.BindPFlags(fs)
}

func GetBool(name string) bool {
	return viper.GetBool(name)
}

func GetString(name string) string {
	return viper.GetString(name)
}

func GetInt(name string) int {
	return viper.GetInt(name)
}

func Lookup(name string) interface{} {
	return viper.Get(name)
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(Prefix)
}
