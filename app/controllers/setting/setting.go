package setting

import (
	"net/http"

	macaron "gopkg.in/macaron.v1"

	"github.com/isymbo/pixpress/setting"
)

// ConfigInfo return application configuration info
func ConfigInfo(c *macaron.Context) {
	c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"App version":                setting.APP_VER,
			"App config filepath":        setting.CfgPath,
			"Database host":              setting.Database.Host,
			"Database port":              setting.Database.Port,
			"Ldap host":                  setting.Ldap.Host,
			"Ldap port":                  setting.Ldap.Port,
			"Ldap bind_dn":               setting.Ldap.BindDn,
			"Ldap user":                  setting.Ldap.User,
			"Ldap password":              setting.Ldap.Password,
			"Server DISABLE_ROURTER_LOG": setting.DisableRouterLog,
			"Server STATIC_ROOT_PATH":    setting.StaticRootPath,
			"Server ENABLE_GZIP":         setting.EnableGzip,
		},
	)
}
