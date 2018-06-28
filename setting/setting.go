package setting

import (
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/go-macaron/session"
	"github.com/urfave/cli"
	log "gopkg.in/clog.v1"
	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"

	"github.com/isymbo/pixpress/util/runtime"
)

type LandingPage string

type ServerType struct {
	Protocol         string      `ini:"PROTOCOL"`
	Domain           string      `ini:"DOMAIN"`
	HTTPAddr         string      `ini:"HTTP_ADDR"`
	HTTPPort         int         `ini:"HTTP_PORT"`
	RootURL          string      `ini:"ROOT_URL"`
	DisableRouterLog bool        `ini:"DISABLE_ROUTER_LOG"`
	StaticRootPath   string      `ini:"STATIC_ROOT_PATH"`
	EnableGzip       bool        `ini:"ENABLE_GZIP"`
	RunMode          string      `ini:"RUN_MODE"`
	LandingPageURL   LandingPage `ini:"LANDING_PAGE"`
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

type SecurityType struct {
	InstallLock             bool   `ini:"INSTALL_LOCK"`
	SecretKey               string `ini:"SECRET_KEY"`
	LoginRememberDays       int    `ini:"LOGIN_REMEMBER_DAYS"`
	CookieUserName          string `ini:"COOKIE_USERNAME"`
	CookieRememberName      string `ini:"COOKIE_REMEMBER_NAME"`
	CookieSecure            bool   `ini:"COOKIE_SECURE"`
	ReverseProxyAuthUser    string `ini:"REVERSE_PROXY_AUTHENTICATION_USER"`
	EnableLoginStatusCookie bool   `ini:"ENABLE_LOGIN_STATUS_COOKIE"`
	LoginStatusCookieName   string `ini:"LOGIN_STATUS_COOKIE_NAME"`
}

type ServiceType struct {
	ActiveCodeLives                int  `ini:"ACTIVE_CODE_LIVE_MINUTES"`
	ResetPwdCodeLives              int  `ini:"RESET_PASSWD_CODE_LIVE_MINUTES"`
	RegisterEmailConfirm           bool `ini:"REGISTER_EMAIL_CONFIRM"`
	DisableRegistration            bool `ini:"DISABLE_REGISTRATION"`
	ShowRegistrationButton         bool `ini:"SHOW_REGISTRATION_BUTTON"`
	RequireSignInView              bool `ini:"REQUIRE_SIGNIN_VIEW"`
	EnableNotifyMail               bool `ini:"ENABLE_NOTIFY_MAIL"`
	EnableReverseProxyAuth         bool `ini:"ENABLE_REVERSE_PROXY_AUTHENTICATION"`
	EnableReverseProxyAutoRegister bool `ini:"ENABLE_REVERSE_PROXY_AUTO_REGISTRATION"`
	EnableCaptcha                  bool `ini:"ENABLE_CAPTCHA"`
}

type TimeType struct {
	Format string `ini:"FORMAT"`
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
	// Cfg is the file descriptor of configuration
	Cfg *ini.File

	// App settings
	AppVer         string = APP_VER
	AppName        string = APP_NAME
	AppCfgPath     string = APP_CONFIG
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

	SessionConfig session.Options

	Server   ServerType
	Database DatabaseType
	Session  SessionType
	Ldap     LDAPType
	Security SecurityType
	Time     TimeType
	Service  ServiceType
	Log      LogType
	Other    OtherType
	XormLog  XormLogType
)

// LoadConfig load configuration settings
func LoadConfig(c *cli.Context) error {
	log.New(log.CONSOLE, log.ConsoleConfig{})

	AppWorkDir = runtime.ExecDir()

	config := c.GlobalString("config")
	if !path.IsAbs(config) {
		AppCfgPath = path.Join(AppWorkDir, config)
	}
	log.Info("Use '%s' as the configuration file.", AppCfgPath)

	var err error
	Cfg, err = ini.Load(AppCfgPath)
	if err != nil {
		log.Fatal(2, "Fail to parse '%s': %s", AppCfgPath, err)
	}

	if err = Cfg.Section("server").MapTo(&Server); err != nil {
		log.Fatal(2, "Fail to map server settings: %s", err)
	} else if err = Cfg.Section("session").MapTo(&Session); err != nil {
		log.Fatal(2, "Fail to map session settings: %s", err)
	} else if err = Cfg.Section("database").MapTo(&Database); err != nil {
		log.Fatal(2, "Fail to map database settings: %s", err)
	} else if err = Cfg.Section("ldap").MapTo(&Ldap); err != nil {
		log.Fatal(2, "Fail to map ldap settings: %s", err)
	} else if err = Cfg.Section("security").MapTo(&Security); err != nil {
		log.Fatal(2, "Fail to map security settings: %s", err)
	} else if err = Cfg.Section("time").MapTo(&Time); err != nil {
		log.Fatal(2, "Fail to map time settings: %s", err)
	} else if err = Cfg.Section("service").MapTo(&Service); err != nil {
		log.Fatal(2, "Fail to map service settings: %s", err)
	} else if err = Cfg.Section("log").MapTo(&Log); err != nil {
		log.Fatal(2, "Fail to map log settings: %s", err)
	} else if err = Cfg.Section("log.xorm").MapTo(&XormLog); err != nil {
		log.Fatal(2, "Fail to map log.xorm settings: %s", err)
	} else if err = Cfg.Section("other").MapTo(&Other); err != nil {
		log.Fatal(2, "Fail to map other settings: %s", err)
	}

	AppURL = Server.RootURL
	if AppURL[len(AppURL)-1] != '/' {
		AppURL += "/"
	}

	// Check if has app suburl.
	url, err := url.Parse(AppURL)
	if err != nil {
		log.Fatal(2, "Invalid ROOT_URL '%s': %s", AppURL, err)
	}
	// Suburl should start with '/' and end without '/', such as '/{subpath}'.
	// This value is empty if site does not have sub-url.
	AppSubURL = strings.TrimSuffix(url.Path, "/")
	AppSubURLDepth = strings.Count(AppSubURL, "/")

	if Log.RootPath == "" {
		Log.RootPath = path.Join(AppWorkDir, "log")
	}

	ProdMode = Server.RunMode == "prod"

	Time.Format = map[string]string{
		"ANSIC":       time.ANSIC,
		"UnixDate":    time.UnixDate,
		"RubyDate":    time.RubyDate,
		"RFC822":      time.RFC822,
		"RFC822Z":     time.RFC822Z,
		"RFC850":      time.RFC850,
		"RFC1123":     time.RFC1123,
		"RFC1123Z":    time.RFC1123Z,
		"RFC3339":     time.RFC3339,
		"RFC3339Nano": time.RFC3339Nano,
		"Kitchen":     time.Kitchen,
		"Stamp":       time.Stamp,
		"StampMilli":  time.StampMilli,
		"StampMicro":  time.StampMicro,
		"StampNano":   time.StampNano,
	}[Time.Format]

	return nil
}

// ConfigInfo return application configuration info
func ConfigInfo(c *macaron.Context) {
	c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"App Version":     AppVer,
			"App Config Path": AppCfgPath,
			"App Root URL":    AppURL,
			"App Sub URL":     AppSubURL,
			"ProdMode":        ProdMode,
			"Server":          Server,
			"Session":         Session,
			"LDAP":            Ldap,
			"Database":        Database,
			"Log":             Log,
			"XormLog":         XormLog,
			"Other":           Other,
		},
	)
}

func InitRoutes(m *macaron.Macaron) {
	m.Get("/setting", ConfigInfo)
}

// func newSessionService() {
// 	SessionConfig.Provider = Cfg.Section("session").Key("PROVIDER").In("memory",
// 		[]string{"memory", "file", "redis", "mysql"})
// 	SessionConfig.ProviderConfig = strings.Trim(Cfg.Section("session").Key("PROVIDER_CONFIG").String(), "\" ")
// 	SessionConfig.CookieName = Cfg.Section("session").Key("COOKIE_NAME").MustString("i_like_gogs")
// 	SessionConfig.CookiePath = AppSubURL
// 	SessionConfig.Secure = Cfg.Section("session").Key("COOKIE_SECURE").MustBool()
// 	SessionConfig.Gclifetime = Cfg.Section("session").Key("GC_INTERVAL_TIME").MustInt64(3600)
// 	SessionConfig.Maxlifetime = Cfg.Section("session").Key("SESSION_LIFE_TIME").MustInt64(86400)
// 	CSRFCookieName = Cfg.Section("session").Key("CSRF_COOKIE_NAME").MustString("_csrf")

// 	log.Info("Session Service Enabled")
// }
