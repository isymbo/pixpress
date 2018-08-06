package post

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"

	"github.com/isymbo/pixpress/app/controllers/context"
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

type Post struct {
	Title    string `form:"Title" binding:"Required"`
	Content  string `form:"Content" binding:"Required"`
	CoverImg string
}

func InitRoutes(m *macaron.Macaron) {

	reqSignIn := context.ReqSignIn
	// bindIgnErr := binding.BindIgnErr

	m.Group("/pix", func() {
		m.Get("", reqSignIn, ListPix)
		m.Combo("/:pixid").
			Get(reqSignIn, EditPix).
			Post(binding.Bind(Post{}), EditPixPost)
		m.Get("/:pixid/delete", reqSignIn, DeletePix)
		m.Combo("/new").
			Get(reqSignIn, NewPix).
			Post(binding.Bind(Post{}), NewPixPost)
	})
}

func NewPix(c *context.Context) {
	c.Data["Title"] = "新建作品"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsNewPix"] = true

	c.Success(PIXNEW)
}

func EditPix(c *context.Context) {
	c.Data["Title"] = "作品列表"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsEditPix"] = true

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

func EditPixPost(c *context.Context, f Post) {
	c.Data["Title"] = "作品列表"
	c.Data["PageIsPixes"] = true
	c.Data["PageIsEditPixPost"] = true

	if c.HasError() {
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
	c.Data["PageIsListPix"] = true

	routes.RenderPostSearch(c, &routes.PostSearchOptions{
		Type:     models.POST_TYPE_PIX,
		Counter:  models.CountPosts,
		Ranger:   models.Posts,
		PageSize: setting.UI.Admin.UserPagingNum,
		OrderBy:  "updated_unix DESC",
		TplName:  PIXES,
	})
}

func NewPixPost(c *context.Context, f Post) {
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
