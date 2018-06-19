package setting

import (
	"net/http"
	"path"

	"github.com/urfave/cli"
	log "gopkg.in/clog.v1"
	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"

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

	// App settings
	AppVer         string
	AppName        string
	AppURL         string
	AppSubURL      string
	AppSubURLDepth int // Number of slashes
	AppPath        string
	AppDataPath    string

	// Database struct
	Database struct {
		DbType   string `ini:"DB_TYPE"`
		Host     string `ini:"HOST"`
		Name     string `ini:"NAME"`
		User     string `ini:"USER"`
		Password string `ini:"PASSWD"`
		SslMode  string `ini:"SSL_MODE"`
		Path     string `ini:"PATH"`
	}

	// Ldap struct
	Ldap struct {
		Host     string `ini:"HOST"`
		Port     int    `ini:"PORT"`
		Base     string `ini:"BASE"`
		BindDn   string `ini:"BIND_DN"`
		Password string `ini:"PASSWORD"`
	}

	// Other struct
	Other struct {
		ShowFooterTemplateLoadTime bool `ini:"SHOW_FOOTER_TEMPLATE_LOAD_TIME"`
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

	AppVer = APP_VER
	AppName = APP_NAME

	AppSubURL = ""

	if err = cfg.Section("database").MapTo(&Database); err != nil {
		log.Fatal(2, "Fail to map database settings: %v", err)
	} else if err = cfg.Section("ldap").MapTo(&Ldap); err != nil {
		log.Fatal(2, "Fail to map ldap settings: %v", err)
	} else if err = cfg.Section("other").MapTo(&Other); err != nil {
		log.Fatal(2, "Fail to map other settings: %v", err)
	}

	return nil
}

func cfgAbsPath(cf string) (string, error) {
	if !path.IsAbs(cf) {
		cf = path.Join(util.ExecDir(), cf)
	}

	return cf, nil
}

// ConfigInfo return application configuration info
func ConfigInfo(c *macaron.Context) {
	c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"Server DISABLE_ROURTER_LOG": DisableRouterLog,
			"Server STATIC_ROOT_PATH":    StaticRootPath,
			"Server ENABLE_GZIP":         EnableGzip,
			"App version":                AppVer,
			"App config filepath":        CfgPath,
			"LDAP":                       Ldap,
			"Database":                   Database,
			"Other":                      Other,
		},
	)
}

func InitRoutes(m *macaron.Macaron) {
	m.Get("/setting", ConfigInfo)
}
