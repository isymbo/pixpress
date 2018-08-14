package user

import (
	"github.com/isymbo/pixpress/app/controllers/context"
)

func Settings(c *context.Context) {
	// uname := c.Params(":loginname")
	// log.Trace("loginname: %+v", uname)
	// // Special handle for FireFox requests favicon.ico.
	// if uname == "favicon.ico" {
	// 	c.ServeFile(path.Join(setting.Server.StaticRootPath, "public/img/favicon.png"))
	// 	return
	// } else if strings.HasSuffix(uname, ".png") {
	// 	c.Error(404)
	// 	return
	// }

	// ctxUser := GetUserByName(c, strings.TrimSuffix(uname, ".keys"))
	// log.Trace("ctxUser: %+v", ctxUser)
	// if c.Written() {
	// 	return
	// }

	// c.Data["Title"] = ctxUser.DisplayName
	// c.Data["PageIsUserProfile"] = true
	// c.Data["Owner"] = ctxUser

	// c.Success(PROFILE)

	c.Title("用户头像")
	c.PageIs("SettingsAvatar")
	c.Success(SETTINGS_AVATAR)

}
