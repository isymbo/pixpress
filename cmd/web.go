package cmd

import (
	"path"
	"strconv"
	"strings"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/gzip"
	"github.com/go-macaron/session"
	"github.com/go-macaron/toolbox"
	"github.com/urfave/cli"
	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"

	"github.com/isymbo/pixpress/app/controllers/admin"
	"github.com/isymbo/pixpress/app/controllers/context"
	"github.com/isymbo/pixpress/app/controllers/post"
	"github.com/isymbo/pixpress/app/controllers/routes"
	"github.com/isymbo/pixpress/app/controllers/template"
	"github.com/isymbo/pixpress/app/controllers/user"
	"github.com/isymbo/pixpress/app/models"
	"github.com/isymbo/pixpress/setting"
)

// Web start web server
var Web = cli.Command{
	Name:        "web",
	Usage:       "Start web server",
	Description: `PixPress web server is the only thing you need to run, and it takes care of all the other things for you.`,
	Action:      runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", strconv.Itoa(setting.APP_HTTPPORT), "Specify network port to serve"),
	},
}

// newMacaron initializes Macaron instance.
func newMacaron() *macaron.Macaron {
	m := macaron.New()
	if !setting.Server.DisableRouterLog {
		m.Use(macaron.Logger())
	}
	m.Use(macaron.Recovery())
	if setting.Server.EnableGzip {
		m.Use(gzip.Gziper())
	}

	m.Use(macaron.Static(
		path.Join(setting.Server.StaticRootPath, "public"),
		macaron.StaticOptions{
			SkipLogging: setting.Server.DisableRouterLog,
		},
	))

	m.Use(macaron.Static(
		setting.Avatar.AvatarUploadPath,
		macaron.StaticOptions{
			Prefix:      "avatars",
			SkipLogging: setting.Server.DisableRouterLog,
		},
	))

	funcMap := template.NewFuncMap()
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:  path.Join(setting.Server.StaticRootPath, "app/views/templates"),
		Funcs:      funcMap,
		IndentJSON: macaron.Env != macaron.PROD,
	}))

	m.Use(cache.Cacher(cache.Options{
		Adapter:       setting.Cache.Adapter,
		AdapterConfig: setting.Cache.Conn,
		Interval:      setting.Cache.Interval,
	}))

	m.Use(session.Sessioner(setting.SessionConfig))
	m.Use(csrf.Csrfer(csrf.Options{
		Secret:     setting.Security.SecretKey,
		Cookie:     setting.Session.CSRFCookieName,
		SetCookie:  true,
		Header:     "X-Csrf-Token",
		CookiePath: setting.AppSubURL,
	}))

	m.Use(toolbox.Toolboxer(m, toolbox.Options{
		HealthCheckFuncs: []*toolbox.HealthCheckFuncDesc{
			&toolbox.HealthCheckFuncDesc{
				Desc: "Database connection",
				Func: models.Ping,
			},
		},
	}))

	m.Use(context.Contexter())
	return m

}

// GlobalInit use loaded configurations to set internal variables and init services
func GlobalInit() {
	// Set database type
	models.InitDBType()

	// Initialize database engine
	if err := models.NewEngine(); err != nil {
		log.Fatal(2, "Fail to initialize ORM engine: %v", err)
	}
	models.HasEngine = true

	// Kept for development/test usage only
	if models.EnableSQLite3 {
		log.Info("SQLite3 Supported")
	}
	checkRunMode()

}

// Set authd cofigurations into database here
func NewAuthdService() {
	length := len(setting.LoginModes)

	loginTypes := map[string]models.LoginType{
		"ldap_simple_auth": models.LOGIN_LDAP,
		"ldap_bind_dn":     models.LOGIN_DLDAP,
	}

	for i := 0; i < length; i++ {
		// log.Info("LoginMode: %d (LoginType: %+v)", i, models.LoginType(setting.LoginModes[i].Type))

		ls := models.LoginSource{
			ID:        setting.LoginModes[i].ID,
			Type:      loginTypes[setting.LoginModes[i].Type],
			Name:      setting.LoginModes[i].Name,
			IsActived: setting.LoginModes[i].IsActivated,
			Cfg: &models.LDAPConfig{
				Source: &setting.LoginSources[i],
			},
		}

		if err := models.CreateLoginSource(&ls); err != nil {
			if models.IsErrLoginSourceAlreadyExist(err) {
				log.Info("LoginSource exists: %s", err)
				if err = models.UpdateLoginSource(&ls); err != nil {
					log.Error(2, "Update LoginSource Error: %+v, %s ", setting.LoginSources[i], err)
				} else {
					log.Trace("Update LoginSource successfully: %+v", setting.LoginSources[i])
				}
			} else {
				log.Error(2, "Create LoginSource Error: %+v, %s ", setting.LoginSources[i], err)
			}
		} else {
			log.Trace("Add new LoginSource successfully: %+v", setting.LoginSources[i])
		}
	}

	// comment out for now, added for debug only
	// models.LoadAuthSources(path.Join(setting.AppWorkDir, "config", "auth.d"))

	log.Info("Authentication Service Started")

}

func runWeb(c *cli.Context) error {
	// override HTTPPort if it is set by command run option -p / --port
	if c.IsSet("port") {
		setting.Server.HTTPPort = c.Int("port")
	}

	// Get all services right
	setting.NewServices()

	// Get backend database running
	GlobalInit()

	// TODO, FIXME
	// During this initial development phase, use configuration files to set authd.
	// As the authd configuration stores data in the database, so it has to be called after database service is ready.
	// Later on, implementation of creating/updating/deleting authd service via admin web pages.
	NewAuthdService()

	m := newMacaron()
	m.SetAutoHead(true)
	initRoutes(m)

	m.Run(setting.Server.HTTPPort)

	return nil
}

func initRoutes(m *macaron.Macaron) {

	// Not found handler.
	m.NotFound(routes.NotFound)

	m.Get("/", context.IgnSignIn, routes.Home)
	m.Group("/explore", func() {
		m.Get("", func(c *context.Context) {
			c.Redirect(setting.AppSubURL + "/explore/works")
		})
		m.Get("/works", routes.ExploreWorks)
		m.Get("/works/:pixid", post.AnonViewPix)
	}, context.IgnSignIn)

	setting.InitRoutes(m)
	admin.InitRoutes(m)
	user.InitRoutes(m)
	post.InitRoutes(m)

}

func checkRunMode() {
	if setting.ProdMode {
		macaron.Env = macaron.PROD
		macaron.ColorLog = false
	}
	log.Info("Run Mode: %s", strings.Title(macaron.Env))
}
