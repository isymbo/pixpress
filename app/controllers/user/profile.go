package user

import (
	"path"
	"strings"

	"github.com/isymbo/pixpress/app/controllers/context"
	"github.com/isymbo/pixpress/app/models"
	"github.com/isymbo/pixpress/app/models/errors"
	"github.com/isymbo/pixpress/setting"
)

func GetUserByName(c *context.Context, name string) *models.User {
	user, err := models.GetUserByName(name)
	if err != nil {
		c.NotFoundOrServerError("GetUserByName", errors.IsUserNotExist, err)
		return nil
	}
	return user
}

// // GetUserByParams returns user whose name is presented in URL paramenter.
// func GetUserByParams(c *context.Context) *models.User {
// 	return GetUserByName(c, c.Params(":username"))
// }

func GetUserByID(c *context.Context, id int64) *models.User {
	user, err := models.GetUserByID(id)
	if err != nil {
		c.NotFoundOrServerError("GetUserByName", errors.IsUserNotExist, err)
		return nil
	}
	return user
}

func Profile(c *context.Context) {
	uname := c.Params(":loginname")
	//log.Trace("loginname: %+v", uname)
	// Special handle for FireFox requests favicon.ico.
	if uname == "favicon.ico" {
		c.ServeFile(path.Join(setting.Server.StaticRootPath, "public/img/favicon.png"))
		return
	} else if strings.HasSuffix(uname, ".png") {
		c.Error(404)
		return
	}

	ctxUser := GetUserByName(c, strings.TrimSuffix(uname, ".keys"))
	//log.Trace("ctxUser: %+v", ctxUser)
	if c.Written() {
		return
	}

	// if ctxUser.IsOrganization() {
	// 	showOrgProfile(c)
	// 	return
	// }

	c.Data["Title"] = ctxUser.DisplayName
	c.Data["PageIsUserProfile"] = true
	//c.Data["Owner"] = ctxUser

	// orgs, err := models.GetOrgsByUserID(ctxUser.ID, c.IsLogged && (c.User.IsAdmin || c.User.ID == ctxUser.ID))
	// if err != nil {
	// 	c.Handle(500, "GetOrgsByUserIDDesc", err)
	// 	return
	// }

	// c.Data["Orgs"] = orgs

	// tab := c.Query("tab")
	// c.Data["TabName"] = tab
	// switch tab {
	// case "activity":
	// 	retrieveFeeds(c, ctxUser, -1, true)
	// 	if c.Written() {
	// 		return
	// 	}
	// default:
	// 	page := c.QueryInt("page")
	// 	if page <= 0 {
	// 		page = 1
	// 	}

	// 	showPrivate := c.IsLogged && (ctxUser.ID == c.User.ID || c.User.IsAdmin)
	// 	c.Data["Repos"], err = models.GetUserRepositories(&models.UserRepoOptions{
	// 		UserID:   ctxUser.ID,
	// 		Private:  showPrivate,
	// 		Page:     page,
	// 		PageSize: setting.UI.User.RepoPagingNum,
	// 	})
	// 	if err != nil {
	// 		c.Handle(500, "GetRepositories", err)
	// 		return
	// 	}

	// 	count := models.CountUserRepositories(ctxUser.ID, showPrivate)
	// 	c.Data["Page"] = paginater.New(int(count), setting.UI.User.RepoPagingNum, page, 5)
	// }

	//c.HTML(200, PROFILE)
	c.Success(PROFILE)
}

func UProfile(c *context.Context) {
	uid := c.ParamsInt64(":uid")

	ctxUser := GetUserByID(c, uid)
	if c.Written() {
		return
	}

	c.Data["Title"] = ctxUser.DisplayName
	c.Data["PageIsUserUProfile"] = true

	c.Success(UPROFILE)
}
