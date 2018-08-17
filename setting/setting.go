package setting

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/go-macaron/session"
	"github.com/urfave/cli"
	log "gopkg.in/clog.v1"
	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"
	libravatar "strk.kbt.io/projects/go/libravatar"

	"github.com/isymbo/pixpress/app/controllers/auth/ldap"
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

type CacheType struct {
	Adapter  string `ini:"ADAPTER"`
	Interval int    `ini:"INTERVAL"`
	Host     string `ini:"HOST"`
	Conn     string `ini:"-"`
}

type SessionType struct {
	Provider        string `ini:"PROVIDER"`
	ProviderConfig  string `ini:"PROVIDER_CONFIG"`
	CookieName      string `ini:"COOKIE_NAME"`
	CookieSecure    bool   `ini:"COOKIE_SECURE"`
	EnableSetCookie bool   `ini:"ENABLE_SET_COOKIE"`
	GCLifeTime      int64  `ini:"GC_LIFE_TIME"`
	SessionLifeTime int64  `ini:"SESSION_LIFE_TIME"`
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
	InstallLock              bool   `ini:"INSTALL_LOCK"`
	SecretKey                string `ini:"SECRET_KEY"`
	LoginRememberDays        int    `ini:"LOGIN_REMEMBER_DAYS"`
	CookieUserName           string `ini:"COOKIE_USERNAME"`
	CookieRememberName       string `ini:"COOKIE_REMEMBER_NAME"`
	CookieSecure             bool   `ini:"COOKIE_SECURE"`
	ReverseProxyAuthUser     string `ini:"REVERSE_PROXY_AUTHENTICATION_USER"`
	EnableLoginStatusCookie  bool   `ini:"ENABLE_LOGIN_STATUS_COOKIE"`
	LoginStatusCookieName    string `ini:"LOGIN_STATUS_COOKIE_NAME"`
	AccessControlAllowOrigin string `ini:"ACCESS_CONTROL_ALLOW_ORIGIN"`
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
	BufferLen int64  `ini:"BUFFER_LEN"`
	Level     string `ini:"LEVEL"`
}

type LogConsoleType struct {
	Level string `ini:"LEVEL"`
}

type LogFileType struct {
	Level        string `ini:"LEVEL"`
	LogRotate    bool   `ini:"LOG_ROTATE"`
	DailyRotate  bool   `ini:"DAILY_ROTATE"`
	MaxSizeShift int64  `ini:"MAX_SIZE_SHIFT"`
	MaxLines     int64  `ini:"MAX_LINES"`
	MaxDays      int64  `ini:"MAX_DAYS"`
}

type LogSlackType struct {
	Level string `ini:"LEVEL"`
	URL   string `ini:"URL"`
}

type LogDiscordType struct {
	Level    string `ini:"LEVEL"`
	URL      string `ini:"URL"`
	UserName string `ini:"USERNAME"`
}

type LogXormType struct {
	Rotate      bool  `ini:"ROTATE"`
	RotateDaily bool  `ini:"ROTATE_DAILY"`
	MaxSize     int64 `ini:"MAX_SIZE"`
	MaxDays     int64 `ini:"MAX_DAYS"`
}

// type LoginLDAPType ldap.Source

type LoginType struct {
	ID          int64  `ini:"ID"`
	Type        string `ini:"TYPE"`
	Name        string `ini:"NAME"`
	IsActivated bool   `ini:"IS_ACTIVATED"`
}

// UI settings
type UIType struct {
	ExplorePagingNum   int    `ini:"EXPLORE_PAGING_NUM"`
	IssuePagingNum     int    `ini:"ISSUE_PAGING_NUM"`
	FeedMaxCommitNum   int    `ini:"FEED_MAX_COMMIT_NUM"`
	ThemeColorMetaTag  string `ini:"THEME_COLOR_META_TAG"`
	MaxDisplayFileSize int64  `ini:"MAX_DISPLAY_FILE_SIZE"`

	Admin struct {
		UserPagingNum   int `ini:"USER_PAGING_NUM"`
		RepoPagingNum   int `ini:"REPO_PAGING_NUM"`
		NoticePagingNum int `ini:"NOTICE_PAGING_NUM"`
		OrgPagingNum    int `ini:"ORG_PAGING_NUM"`
	} `ini:"ui.admin"`

	User struct {
		RepoPagingNum     int `ini:"REPO_PAGING_NUM"`
		NewsFeedPagingNum int `ini:"NEWS_FEED_PAGING_NUM"`
		CommitsPagingNum  int `ini:"COMMITS_PAGING_NUM"`
		PostPagingNum     int `ini:"POST_PAGING_NUM"`
	} `ini:"ui.user"`
}

type AvatarType struct {
	AvatarUploadPath      string `ini:"AVATAR_UPLOAD_PATH"`
	GravatarSource        string `ini:"GRAVATAR_SOURCE"`
	DisableGravatar       bool   `ini:"DISABLE_GRAVATAR"`
	EnableFederatedAvatar bool   `ini:"ENABLE_FEDERATED_AVATAR"`
	// LibravatarService     *libravatar.Libravatar
}

type CoverType struct {
	Enabled      bool   `ini:"ENABLED"`
	Path         string `ini:"PATH"`
	AllowedTypes string `ini:"ALLOWED_TYPES"`
	MaxSize      int64  `ini:"MAX_SIZE"`
	MaxFiles     int    `ini:"MAX_FILES"`
}

type AttachmentType struct {
	Enabled      bool   `ini:"ENABLED"`
	Path         string `ini:"PATH"`
	AllowedTypes string `ini:"ALLOWED_TYPES"`
	MaxSize      int64  `ini:"MAX_SIZE"`
	MaxFiles     int    `ini:"MAX_FILES"`
}

type OtherType struct {
	ShowFooterTemplateLoadTime bool `ini:"SHOW_FOOTER_TEMPLATE_LOAD_TIME"`
	ShowFooterBranding         bool `ini:"SHOW_FOOTER_BRANDING"`
	ShowFooterVersion          bool `ini:"SHOW_FOOTER_VERSION"`
}

var (
	// Build information should only be set by -ldflags.
	BuildTime    string
	BuildGitHash string

	// Cfg is the file descriptor of configuration
	Cfg *ini.File

	// App settings
	AppVer          string = APP_VER
	AppName         string = APP_NAME
	AppCfgPath      string = APP_CONFIG
	AppAuthdCfgPath string = APP_AUTHD_CONFIG_PATH
	AppURL          string
	AppSubURL       string
	AppSubURLDepth  int // Number of slashes
	AppPath         string
	AppDataPath     string
	AppWorkDir      string

	ProdMode bool

	// Log settings
	LogRootPath string
	LogModes    []string
	LogConfigs  []interface{}

	UseSQLite3    bool
	UseMySQL      bool
	UsePostgreSQL bool
	UseMSSQL      bool

	SessionConfig     session.Options
	LibravatarService *libravatar.Libravatar

	Server       ServerType
	Database     DatabaseType
	Cache        CacheType
	Session      SessionType
	Ldap         LDAPType
	Security     SecurityType
	Time         TimeType
	Service      ServiceType
	Log          LogType
	LogConsole   LogConsoleType
	LogFile      LogFileType
	LogSlack     LogSlackType
	LogDiscord   LogDiscordType
	LogXorm      LogXormType
	LoginModes   []LoginType
	LoginSources []ldap.Source
	UI           UIType
	Avatar       AvatarType
	Cover        CoverType
	Attachment   AttachmentType
	Other        OtherType
)

func loadAppConfig(c *cli.Context) error {
	AppWorkDir = runtime.ExecDir()

	config := c.GlobalString("config")
	if !path.IsAbs(config) {
		AppCfgPath = path.Join(AppWorkDir, config)
	} else {
		AppCfgPath = config
	}
	log.Info("Use '%s' as the configuration file.", AppCfgPath)

	var err error
	Cfg, err = ini.Load(AppCfgPath)
	if err != nil {
		log.Fatal(2, "Fail to parse '%s': %s", AppCfgPath, err)
	}

	if err = Cfg.Section("server").MapTo(&Server); err != nil {
		log.Fatal(2, "Fail to map server settings: %s", err)
	} else if err = Cfg.Section("cache").MapTo(&Cache); err != nil {
		log.Fatal(2, "Fail to map cache settings: %s", err)
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
	} else if err = Cfg.Section("log.console").MapTo(&LogConsole); err != nil {
		log.Fatal(2, "Fail to map log.console settings: %s", err)
	} else if err = Cfg.Section("log.file").MapTo(&LogFile); err != nil {
		log.Fatal(2, "Fail to map log.file settings: %s", err)
	} else if err = Cfg.Section("log.slack").MapTo(&LogSlack); err != nil {
		log.Fatal(2, "Fail to map log.slack settings: %s", err)
	} else if err = Cfg.Section("log.discord").MapTo(&LogDiscord); err != nil {
		log.Fatal(2, "Fail to map log.discord settings: %s", err)
	} else if err = Cfg.Section("log.xorm").MapTo(&LogXorm); err != nil {
		log.Fatal(2, "Fail to map log.xorm settings: %s", err)
	} else if err = Cfg.Section("ui").MapTo(&UI); err != nil {
		log.Fatal(2, "Fail to map ui settings: %s", err)
	} else if err = Cfg.Section("avatar").MapTo(&Avatar); err != nil {
		log.Fatal(2, "Fail to map avatar settings: %s", err)
	} else if err = Cfg.Section("cover").MapTo(&Cover); err != nil {
		log.Fatal(2, "Fail to map cover settings: %s", err)
	} else if err = Cfg.Section("attachment").MapTo(&Attachment); err != nil {
		log.Fatal(2, "Fail to map attachment settings: %s", err)
	} else if err = Cfg.Section("other").MapTo(&Other); err != nil {
		log.Fatal(2, "Fail to map other settings: %s", err)
	}

	// Adapt to the configured LandingPageURL
	switch Server.LandingPageURL {
	case "explore":
		Server.LandingPageURL = LANDING_PAGE_EXPLORE
	default:
		Server.LandingPageURL = LANDING_PAGE_HOME
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

func LoadAuthSources(c *cli.Context) error {
	AppWorkDir = runtime.ExecDir()

	config := c.GlobalString("config")

	if !path.IsAbs(config) {
		AppAuthdCfgPath = path.Join(AppWorkDir, path.Dir(config), "auth.d")
	} else {
		AppAuthdCfgPath = path.Join(path.Dir(config), "auth.d")
	}
	log.Info("Use '%s' as the authd configuration path.", AppAuthdCfgPath)

	paths, err := com.GetFileListBySuffix(AppAuthdCfgPath, ".conf")
	if err != nil {
		log.Fatal(2, "Failed to list authentication sources: %v", err)
	}

	LoginModes = make([]LoginType, len(paths))
	LoginSources = make([]ldap.Source, len(paths))

	for i, p := range paths {
		authSource, err := ini.Load(p)
		if err != nil {
			log.Fatal(2, "Failed to load authentication source: %v, %v", p, err)
		}

		if err = authSource.Section("config").MapTo(&LoginSources[i]); err != nil {
			log.Fatal(2, "Fail to map authd config settings: %s", err)
		} else if err = authSource.MapTo(&LoginModes[i]); err != nil {
			log.Fatal(2, "Fail to map authd default settings: %s", err)
		}
	}

	return nil
}

// LoadConfig load configuration settings
func LoadConfig(c *cli.Context) error {
	log.New(log.CONSOLE, log.ConsoleConfig{})

	var err error
	if err = loadAppConfig(c); err != nil {
		return err
	}

	if err = LoadAuthSources(c); err != nil {
		return err
	}

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
			"Cache":           Cache,
			"Session":         Session,
			"Security":        Security,
			"Service":         Service,
			"LDAP":            Ldap,
			"LoginModes":      LoginModes,
			"LoginSources":    LoginSources,
			"Database":        Database,
			"LogConfigs":      LogConfigs,
			"Log":             Log,
			"LogConsole":      LogConsole,
			"LogFile":         LogFile,
			"LogSlack":        LogSlack,
			"LogDiscord":      LogDiscord,
			"XormLog":         LogXorm,
			"Other":           Other,
		},
	)
}

func InitRoutes(m *macaron.Macaron) {
	m.Get("/setting", ConfigInfo)
}

// Set Log configurations
func newLogService() {
	LogRootPath = Log.RootPath

	if len(BuildTime) > 0 {
		log.Trace("Build Time: %s", BuildTime)
		log.Trace("Build Git Hash: %s", BuildGitHash)
	}

	// Because we always create a console logger as primary logger before all settings are loaded,
	// thus if user doesn't set console logger, we should remove it after other loggers are created.
	hasConsole := false

	// Get and check log modes.
	LogModes = strings.Split(Log.Mode, ",")
	LogConfigs = make([]interface{}, len(LogModes))
	levelNames := map[string]log.LEVEL{
		"trace": log.TRACE,
		"info":  log.INFO,
		"warn":  log.WARN,
		"error": log.ERROR,
		"fatal": log.FATAL,
	}

	var level string
	for i, mode := range LogModes {
		mode = strings.ToLower(strings.TrimSpace(mode))
		validTypes := []string{"console", "file", "slack", "discord"}
		if !com.IsSliceContainsStr(validTypes, mode) {
			log.Fatal(2, "Unknown logger mode: %s", mode)
		}

		// Generate log configuration.
		switch log.MODE(mode) {
		case log.CONSOLE:
			hasConsole = true
			level = validateLogLevel(LogConsole.Level)
			LogConfigs[i] = log.ConsoleConfig{
				Level:      levelNames[level],
				BufferSize: Log.BufferLen,
			}

		case log.FILE:
			logPath := path.Join(LogRootPath, strings.ToLower(APP_NAME)+".log")
			if err := os.MkdirAll(path.Dir(logPath), os.ModePerm); err != nil {
				log.Fatal(2, "Fail to create log directory '%s': %v", path.Dir(logPath), err)
			}

			level = validateLogLevel(LogFile.Level)
			LogConfigs[i] = log.FileConfig{
				Level:      levelNames[level],
				BufferSize: Log.BufferLen,
				Filename:   logPath,
				FileRotationConfig: log.FileRotationConfig{
					Rotate:   LogFile.LogRotate,
					Daily:    LogFile.DailyRotate,
					MaxSize:  1 << uint(LogFile.MaxSizeShift),
					MaxLines: LogFile.MaxLines,
					MaxDays:  LogFile.MaxDays,
				},
			}

		case log.SLACK:
			level = validateLogLevel(LogSlack.Level)
			LogConfigs[i] = log.SlackConfig{
				Level:      levelNames[level],
				BufferSize: Log.BufferLen,
				URL:        LogSlack.URL,
			}

		case log.DISCORD:
			level = validateLogLevel(LogDiscord.Level)
			LogConfigs[i] = log.DiscordConfig{
				Level:      levelNames[level],
				BufferSize: Log.BufferLen,
				URL:        LogDiscord.URL,
				Username:   LogDiscord.UserName,
			}
		}

		log.New(log.MODE(mode), LogConfigs[i])
		log.Trace("Log Mode: %s (Level: %s)", strings.Title(mode), strings.Title(level))
	}

	// Make sure everyone gets version info printed.
	log.Info("%s %s", AppName, AppVer)
	if !hasConsole {
		log.Delete(log.CONSOLE)
	}
}

// Set cache related variables
func newCacheService() {
	switch Cache.Adapter {
	case "memory":
	case "redis", "memcache":
		Cache.Conn = strings.Trim(Cache.Host, "\" ")
	default:
		log.Fatal(2, "Unknown cache adapter: %s", Cache.Adapter)
	}

	log.Info("Cache Service Enabled")
}

// Set session config variables
func newSessionService() {
	SessionConfig.Provider = Session.Provider
	SessionConfig.ProviderConfig = strings.Trim(Session.ProviderConfig, "\" ")
	SessionConfig.CookieName = Session.CookieName
	SessionConfig.CookiePath = AppSubURL
	SessionConfig.Secure = Session.CookieSecure
	SessionConfig.Gclifetime = Session.GCLifeTime
	SessionConfig.Maxlifetime = Session.SessionLifeTime
	// CSRFCookieName = Session.CSRFCookieName

	log.Info("Session Service Enabled")
}

// Set auth related variables
func NewAuthdService() {

	log.Info("Authentication Service Enabled")
}

func NewServices() {
	newLogService()
	newCacheService()
	newSessionService()
}

func validateLogLevel(level string) string {
	// FIXME, TODO
	validLevels := []string{"trace", "info", "warn", "error", "fatal"}
	level = strings.ToLower(level)
	if com.IsSliceContainsStr(validLevels, level) {
		return level
	}

	return "trace"
}
