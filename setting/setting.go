package setting

import (
	"path"

	"github.com/urfave/cli"
	log "gopkg.in/clog.v1"
	"gopkg.in/ini.v1"

	"github.com/isymbo/pixpress/util"
)

var (
	// CfgPath is the path of configuration file
	CfgPath     = "config/app.ini"
	WebPort int = APP_HTTPPORT

	// cfg is the file descriptor of configuration
	cfg *ini.File

	// Server settings
	DisableRouterLog bool
	StaticRootPath   string
	EnableGzip       bool

	// Database struct
	Database struct {
		Host string `ini:"HOST"`
		Port int    `ini:"PORT"`
	}

	// Ldap struct
	Ldap struct {
		Host     string `ini:"HOST"`
		Port     int    `ini:"PORT"`
		BindDn   string `ini:"BIND_DN"`
		User     string `ini:"USER"`
		Password string `ini:"PASSWORD"`
	}
)

// LoadConfig load configuration settings
func LoadConfig(c *cli.Context) error {
	log.New(log.CONSOLE, log.ConsoleConfig{})

	CfgPath, _ = cfgAbsPath(c.GlobalString("config"))
	var err error
	cfg, err = ini.Load(CfgPath)
	if err != nil {
		log.Fatal(2, "Fail to parse %s: %v", CfgPath, err)
	}

	workDir := util.ExecDir()

	sec := cfg.Section("server")
	DisableRouterLog = sec.Key("DISABLE_ROUTER_LOG").MustBool(false)
	StaticRootPath = sec.Key("STATIC_ROOT_PATH").MustString(workDir)
	EnableGzip = sec.Key("ENABLE_GZIP").MustBool()

	if err = cfg.Section("database").MapTo(&Database); err != nil {
		log.Fatal(2, "Fail to map database settings: %v", err)
	} else if err = cfg.Section("ldap").MapTo(&Ldap); err != nil {
		log.Fatal(2, "Fail to map ldap settings: %v", err)
	}

	return nil
}

func cfgAbsPath(cf string) (string, error) {
	if !path.IsAbs(cf) {
		cf = path.Join(util.ExecDir(), cf)
	}

	return cf, nil
}
