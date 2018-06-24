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

type ServerType struct {
	DisableRouterLog bool   `ini:"DISABLE_ROUTER_LOG"`
	StaticRootPath   string `ini:"STATIC_ROOT_PATH"`
	EnableGzip       bool   `ini:"ENABLE_GZIP"`
	RunMode          string `ini:"RUN_MODE"`
}

type SessionType struct {
	Provider        string `ini:"PROVIDER"`
	ProviderConfig  string `ini:"PROVIDER_CONFIG"`
	CookieName      string `ini:"COOKIE_NAME"`
	CookieSecure    bool   `ini:"COOKIE_SECURE"`
	EnableSetCookie bool   `ini:"ENABLE_SET_COOKIE"`
	GCIntervalTime  int    `ini:"GC_INTERVAL_TIME"`
	SessionLifeTime int    `ini:"SESSION_LIFE_TIME"`
	CSRFCookieName  string `ini:"CSRF_COOKIE_NAME"`
}

type DatabaseType struct {
	DbType  string `ini:"DB_TYPE"`
	Host    string `ini:"HOST"`
	Name    string `ini:"NAME"`
	User    string `ini:"USER"`
	Passwd  string `ini:"PASSWD"`
	SSLMode string `ini:"SSL_MODE"`
	Path    string `ini:"PATH"`
}

type LDAPType struct {
	Host     string `ini:"HOST"`
	Port     int    `ini:"PORT"`
	Base     string `ini:"BASE"`
	BindDn   string `ini:"BIND_DN"`
	Password string `ini:"PASSWORD"`
}

type LogType struct {
	RootPath  string `ini:"ROOT_PATH"`
	Mode      string `ini:"MODE"`
	BufferLen int    `ini:"BUFFER_LEN"`
	Level     string `ini:"LEVEL"`
}

type OtherType struct {
	ShowFooterTemplateLoadTime bool `ini:"SHOW_FOOTER_TEMPLATE_LOAD_TIME"`
}

type XormLogType struct {
	Rotate      bool  `ini:"ROTATE"`
	RotateDaily bool  `ini:"ROTATE_DAILY"`
	MaxSize     int64 `ini:"MAX_SIZE"`
	MaxDays     int64 `ini:"MAX_DAYS"`
}

var (
	// CfgPath is the path of configuration file
	CfgPath string = "config/app.ini"
	WebPort int    = APP_HTTPPORT

	// cfg is the file descriptor of configuration
	Cfg *ini.File

	// App settings
	AppVer         string = APP_VER
	AppName        string = APP_NAME
	AppURL         string
	AppSubURL      string
	AppSubURLDepth int // Number of slashes
	AppPath        string
	AppDataPath    string
	AppWorkDir     string

	ProdMode bool

	UseSQLite3    bool
	UseMySQL      bool
	UsePostgreSQL bool
	UseMSSQL      bool

	Server   ServerType
	Database DatabaseType
	Session  SessionType
	Ldap     LDAPType
	Log      LogType
	Other    OtherType
	XormLog  XormLogType
)

// LoadConfig load configuration settings
func LoadConfig(c *cli.Context) error {
	log.New(log.CONSOLE, log.ConsoleConfig{})

	AppWorkDir = util.ExecDir()

	CfgPath, _ = cfgAbsPath(c.GlobalString("config"))
	var err error
	Cfg, err = ini.Load(CfgPath)

	if err != nil {
		log.Fatal(2, "Fail to parse %s: %v", CfgPath, err)
	}

	if err = Cfg.Section("server").MapTo(&Server); err != nil {
		log.Fatal(2, "Fail to map server settings: %v", err)
	} else if err = Cfg.Section("session").MapTo(&Session); err != nil {
		log.Fatal(2, "Fail to map session settings: %v", err)
	} else if err = Cfg.Section("database").MapTo(&Database); err != nil {
		log.Fatal(2, "Fail to map database settings: %v", err)
	} else if err = Cfg.Section("ldap").MapTo(&Ldap); err != nil {
		log.Fatal(2, "Fail to map ldap settings: %v", err)
	} else if err = Cfg.Section("log").MapTo(&Log); err != nil {
		log.Fatal(2, "Fail to map log settings: %v", err)
	} else if err = Cfg.Section("log.xorm").MapTo(&XormLog); err != nil {
		log.Fatal(2, "Fail to map log.xorm settings: %v", err)
	} else if err = Cfg.Section("other").MapTo(&Other); err != nil {
		log.Fatal(2, "Fail to map other settings: %v", err)
	}

	if Log.RootPath == "" {
		Log.RootPath = path.Join(AppWorkDir, "log")
	}

	ProdMode = Server.RunMode == "prod"

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
			"App version":         AppVer,
			"App config filepath": CfgPath,
			"ProdMode":            ProdMode,
			"Server":              Server,
			"Session":             Session,
			"LDAP":                Ldap,
			"Database":            Database,
			"Log":                 Log,
			"XormLog":             XormLog,
			"Other":               Other,
		},
	)
}

func InitRoutes(m *macaron.Macaron) {
	m.Get("/setting", ConfigInfo)
}
