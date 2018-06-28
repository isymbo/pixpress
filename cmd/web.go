package cmd

import (
	"path"
	"strconv"
	"strings"

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/gzip"
	"github.com/go-macaron/session"
	"github.com/go-macaron/toolbox"
	"github.com/urfave/cli"
	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"

	"github.com/isymbo/pixpress/app/controllers/context"
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

	// m.Use(macaron.Static(
	// 	setting.AvatarUploadPath,
	// 	macaron.StaticOptions{
	// 		Prefix:      "avatars",
	// 		SkipLogging: setting.DisableRouterLog,
	// 	},
	// ))

	funcMap := template.NewFuncMap()
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:  path.Join(setting.Server.StaticRootPath, "app/views/templates"),
		Funcs:      funcMap,
		IndentJSON: macaron.Env != macaron.PROD,
	}))

	// m.Use(cache.Cacher(cache.Options{
	// 	Adapter:       setting.CacheAdapter,
	// 	AdapterConfig: setting.CacheConn,
	// 	Interval:      setting.CacheInterval,
	// }))

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

func runWeb(c *cli.Context) error {

	if c.IsSet("port") {
		setting.Server.HTTPPort = c.Int("port")
	}

	m := newMacaron()

	models.LoadConfigs()
	if err := models.NewEngine(); err != nil {
		log.Fatal(2, "Fail to initialize ORM engine: %v", err)
	}
	models.HasEngine = true

	m.SetAutoHead(true)
	initRoutes(m)
	m.Run(setting.Server.HTTPPort)

	return nil
}

func initRoutes(m *macaron.Macaron) {

	setting.InitRoutes(m)
	user.InitRoutes(m)

}

func checkRunMode() {
	if setting.ProdMode {
		macaron.Env = macaron.PROD
		macaron.ColorLog = false
	}
	log.Info("Run Mode: %s", strings.Title(macaron.Env))
}
