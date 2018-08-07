package admin

import (
	"github.com/isymbo/pixpress/app/controllers/context"
	"github.com/isymbo/pixpress/app/controllers/routes"
	"github.com/isymbo/pixpress/app/models"
	"github.com/isymbo/pixpress/setting"
)

const (
	PIXES   = "admin/pix/list"
	PIXINFO = "admin/pix/info"
)

func Pixes(c *context.Context) {
	c.Data["Title"] = "作品管理"
	c.Data["PageIsAdmin"] = true
	c.Data["PageIsAdminUsers"] = true
	c.Data["PageIsPixes"] = true

	routes.RenderPostSearch(c, &routes.PostSearchOptions{
		Type:     models.POST_TYPE_PIX,
		Counter:  models.CountPosts,
		Ranger:   models.Posts,
		PageSize: setting.UI.User.PostPagingNum,
		OrderBy:  "updated_unix DESC",
		TplName:  PIXES,
	})
}

func ShowPix(c *context.Context) {
	c.Data["Title"] = "作品信息"
	c.Data["PageIsAdmin"] = true
	c.Data["PageIsAdminUsers"] = true
	c.Data["PageIsShowPix"] = true

	post, err := models.GetPostByID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrPostNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetPostByID", err)
		}
		return
	}

	c.Data["title"] = post.Title
	c.Data["content"] = post.Content

	c.Success(PIXINFO)
}
