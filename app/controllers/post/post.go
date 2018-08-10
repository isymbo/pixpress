package post

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"

	"github.com/isymbo/pixpress/app/controllers/context"
	"github.com/isymbo/pixpress/app/controllers/form"
	"github.com/isymbo/pixpress/app/controllers/routes"
	"github.com/isymbo/pixpress/app/models"
	"github.com/isymbo/pixpress/setting"
)

const (
	PIXNEW    = "pix/new"
	PIXEDIT   = "pix/edit"
	PIXDELETE = "pix/delete"
	PIXES     = "pix/list"
	PIXHOME   = "/pix"
)

func InitRoutes(m *macaron.Macaron) {

	reqSignIn := context.ReqSignIn
	bindIgnErr := binding.BindIgnErr

	m.Group("/pix", func() {
		m.Get("", reqSignIn, ListPix)
		m.Combo("/:pixid").
			Get(reqSignIn, EditPix).
			Post(bindIgnErr(form.CreatePost{}), EditPixPost)
		m.Get("/:pixid/delete", reqSignIn, DeletePix)
		m.Combo("/new").
			Get(reqSignIn, NewPix).
			Post(bindIgnErr(form.CreatePost{}), NewPixPost)
	})
}

func NewPix(c *context.Context) {
	c.Data["Title"] = "新建作品"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsNewPix"] = true
	renderAttachmentSettings(c)

	if c.Written() {
		return
	}

	c.Success(PIXNEW)
}

func EditPix(c *context.Context) {
	c.Data["Title"] = "作品信息"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsEditPix"] = true
	renderAttachmentSettings(c)

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

	c.Success(PIXEDIT)
}

func EditPixPost(c *context.Context, f form.CreatePost) {
	c.Data["Title"] = "作品信息"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsEditPixPost"] = true

	if c.HasError() {
		c.Flash.Error(c.Data["ErrorMsg"].(string))
		c.HTML(200, PIXEDIT)
		return
	}

	post, err := models.GetPostByID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrPostNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetPostByID", err)
		}
		return
	}

	post.Title = f.Title
	post.Content = f.Content

	if err = models.UpdatePost(post); err != nil {
		c.Handle(500, "UpdatePost", err)
		return
	}

	c.SubURLRedirect(PIXHOME)
}

func ListPix(c *context.Context) {
	c.Data["Title"] = "作品列表"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsListPix"] = true

	routes.RenderPostSearch(c, &routes.PostSearchOptions{
		Type:     models.POST_TYPE_PIX,
		Counter:  models.CountPosts,
		Ranger:   models.Posts,
		PageSize: setting.UI.User.PostPagingNum,
		OrderBy:  "updated_unix DESC",
		TplName:  PIXES,
	})
}

func NewPixPost(c *context.Context, f form.CreatePost) {
	c.Data["Title"] = "新建作品"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsNewPix"] = true
	renderAttachmentSettings(c)

	if c.HasError() {
		c.Success(PIXNEW)
		return
	}

	post := &models.Post{
		AuthorID: c.User.ID,
		Author:   c.User,
		Title:    f.Title,
		Content:  f.Content,
		PostType: models.POST_TYPE_PIX,
	}

	if err := models.NewPost(post); err != nil {
		c.Handle(500, "NewPost", err)
		return
	}

	c.SubURLRedirect(PIXHOME)
}

func DeletePix(c *context.Context) {
	post, err := models.GetPostByID(c.ParamsInt64(":pixid"))
	if err != nil {
		if models.IsErrPostNotExist(err) {
			c.Handle(404, "", nil)
		} else {
			c.Handle(500, "GetPostByID", err)
		}
		return
	}

	if err = models.DeletePost(post); err != nil {
		c.Handle(500, "DeletePost", err)
		return
	}

	c.SubURLRedirect(PIXHOME)
}

// func renderAttachmentSettings(c *context.Context) {
// 	c.Data["RequireDropzone"] = true
// 	c.Data["IsAttachmentEnabled"] = setting.Attachment.Enabled
// 	c.Data["AttachmentAllowedTypes"] = setting.Attachment.AllowedTypes
// 	c.Data["AttachmentMaxSize"] = setting.Attachment.MaxSize
// 	c.Data["AttachmentMaxFiles"] = setting.Attachment.MaxFiles
// }
